package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/glueops/autoglue/internal/api"
	"github.com/glueops/autoglue/internal/app"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/config"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long: "Start the HTTP API. This process does not run background jobs: it " +
		"only enqueues them. Run `autoglue worker` for the other half.",
	RunE: func(_ *cobra.Command, _ []string) error {
		rt := app.NewRuntime()
		defer rt.Close()

		cfg := rt.Cfg

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Insert-only River client. It registers no workers and no queues, so
		// this process never fetches work.
		jobs, err := bg.NewInsertClient(rt.Pool)
		if err != nil {
			return fmt.Errorf("init river client: %w", err)
		}

		_ = auth.Refresh(rt.DB, rt.Cfg.JWTPrivateEncKey)
		go func() {
			t := time.NewTicker(60 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					_ = auth.Refresh(rt.DB, rt.Cfg.JWTPrivateEncKey)
				case <-ctx.Done():
					return
				}
			}
		}()

		opts := api.RouterOpts{}

		riverUI, err := api.MountRiverUI(ctx, jobs)
		if err != nil {
			return fmt.Errorf("init river ui: %w", err)
		}
		opts.RiverUI = riverUI
		log.Printf("river dashboard mounted at %s/", api.RiverUIPrefix)

		if cfg.DBStudioEnabled {
			dbURL := cfg.DbURLRO
			if dbURL == "" {
				dbURL = cfg.DbURL
			}

			studio, err := api.MountDbStudio(dbURL, "db-studio", false)
			if err != nil {
				return fmt.Errorf("init db studio: %w", err)
			}
			opts.Studio = studio
			log.Printf("pgweb mounted at /db-studio/")
		}

		r := api.NewRouter(rt.DB, jobs, cfg, opts)

		addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

		srv := &http.Server{
			Addr:         addr,
			Handler:      TimeoutExceptUpgrades(r, 60*time.Second, "request timed out"), // global safety
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		go func() {
			fmt.Printf("🚀 API listening on %s (ui.dev=%v)\n", addr, cfg.UIDev)
			fmt.Printf("   Open %s\n", appURL(cfg, addr))
			for _, w := range originWarnings(cfg, addr) {
				fmt.Printf("   ⚠️  %s\n", w)
			}
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("server error: %v", err)
			}
		}()

		<-ctx.Done()
		fmt.Println("\n⏳ Shutting down API...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// appURL is the URL a browser should actually use. OAUTH_REDIRECT_BASE is the
// authority here, not the bind address: the OAuth popup posts tokens back to
// that origin and the session cookie is scoped to that host, so opening the app
// on any other spelling of the same address logs you nowhere.
func appURL(cfg config.Config, addr string) string {
	if o := strings.TrimSuffix(strings.TrimSpace(cfg.OAuthRedirectBase), "/"); o != "" {
		return o
	}
	return "http://" + addr
}

// originWarnings flags configuration that will half-work in a browser: the app
// answers, but login silently fails because the origins do not match.
func originWarnings(cfg config.Config, addr string) []string {
	base := strings.TrimSpace(cfg.OAuthRedirectBase)
	if base == "" {
		return []string{"OAUTH_REDIRECT_BASE is unset; OAuth popup login will not complete"}
	}

	u, err := url.Parse(base)
	if err != nil || u.Host == "" {
		return []string{fmt.Sprintf("OAUTH_REDIRECT_BASE %q is not a valid URL", base)}
	}

	var out []string

	origin := (&url.URL{Scheme: u.Scheme, Host: u.Host}).String()
	if !config.IsAllowedOrigin(origin) {
		out = append(out, fmt.Sprintf(
			"%s is not in CORS_ALLOWED_ORIGINS; browser calls from it will be blocked", origin))
	}

	// Binding 127.0.0.1 while OAUTH_REDIRECT_BASE says localhost still works,
	// but only if the browser uses the localhost spelling. Say so, because the
	// failure is silent: postMessage is dropped and the cookie is not sent.
	if host, _, splitErr := net.SplitHostPort(addr); splitErr == nil {
		if !strings.EqualFold(host, u.Hostname()) && isLoopback(host) && isLoopback(u.Hostname()) {
			out = append(out, fmt.Sprintf(
				"bound to %s but OAUTH_REDIRECT_BASE uses %s — open the %s URL above, "+
					"the other spelling is a different browser origin and login will fail silently",
				host, u.Hostname(), u.Hostname()))
		}
	}
	return out
}

func isLoopback(host string) bool {
	h := strings.ToLower(host)
	return h == "localhost" || h == "127.0.0.1" || h == "::1"
}

func TimeoutExceptUpgrades(next http.Handler, d time.Duration, msg string) http.Handler {
	timeout := http.TimeoutHandler(next, d, msg)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If this is an upgrade (e.g., websocket), don't wrap.
		if isUpgrade(r) {
			next.ServeHTTP(w, r)
			return
		}
		timeout.ServeHTTP(w, r)
	})
}

func isUpgrade(r *http.Request) bool {
	// Connection: Upgrade, Upgrade: websocket
	if strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") {
		return true
	}
	return false
}
