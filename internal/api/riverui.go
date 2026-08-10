package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/glueops/autoglue/internal/bg"
	"riverqueue.com/riverui"
)

// RiverUIPrefix is where the bundled River dashboard is mounted. It replaces the
// hand-rolled /admin/archer endpoints and the SPA page that consumed them.
const RiverUIPrefix = "/admin/river"

// MountRiverUI builds the River dashboard handler. The returned handler is
// stateful: it runs background caching services, so it must be started with the
// returned start func before it serves requests, and it stops when ctx is done.
func MountRiverUI(ctx context.Context, client *bg.Client) (http.Handler, error) {
	handler, err := riverui.NewHandler(&riverui.HandlerOpts{
		Endpoints: riverui.NewEndpoints(client, nil),
		Logger:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Prefix:    RiverUIPrefix,
	})
	if err != nil {
		return nil, err
	}

	if err := handler.Start(ctx); err != nil {
		return nil, err
	}
	return handler, nil
}
