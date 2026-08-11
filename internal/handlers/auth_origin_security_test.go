package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/glueops/autoglue/internal/handlers/dto"
)

// The property under test:
//
//	OAuth tokens are only ever delivered to an allowlisted browser origin.
//
// AuthStart copies `origin` from the query string into the OAuth state without
// validating it, so by the time AuthCallback reads it back it is
// attacker-controlled. postMessage's second argument is what stops a page from
// reading the message, so if that value came from state unchecked, a page the
// attacker controls could open the callback in a popup — making itself the
// opener — and be handed a live access and refresh token.
//
// These tests run against postMessageOrigin, the function that makes the
// decision. Deleting the allowlist check inside it turns the first test red.

const trustedBase = "https://autoglue.glueopshosted.rocks"

// spaState builds the state string exactly as AuthStart does.
func spaState(origin string) string {
	return "b1a2c3d4|mode=spa|origin=" + url.QueryEscape(origin)
}

// postMessageTarget pulls the literal second argument out of the rendered
// postMessage call and undoes the template's JS string escaping, so assertions
// can talk about the origin rather than about escaping.
func postMessageTarget(t *testing.T, body string) string {
	t.Helper()

	m := regexp.MustCompile(`postMessage\(\s*\{[^}]*\},\s*"([^"]*)"`).FindStringSubmatch(body)
	if m == nil {
		t.Fatalf("could not find a postMessage target in the rendered page:\n%s", body)
	}
	return strings.ReplaceAll(m[1], `\/`, "/")
}

func TestPostMessageOrigin_RejectsAttackerControlledOrigin(t *testing.T) {
	// The attack: an origin the operator never allowlisted must never be
	// returned, because returning it is what delivers the tokens.
	t.Setenv("CORS_ALLOWED_ORIGINS", trustedBase)

	attacks := map[string]string{
		"plain attacker origin":      "https://evil.example",
		"http downgrade of trusted":  "http://autoglue.glueopshosted.rocks",
		"subdomain of attacker":      "https://autoglue.glueopshosted.rocks.evil.example",
		"attacker with trusted path": "https://evil.example/autoglue.glueopshosted.rocks",
		"userinfo smuggling":         "https://autoglue.glueopshosted.rocks@evil.example",
		"port confusion":             "https://autoglue.glueopshosted.rocks:8443",
		"uppercase scheme":           "HTTPS://EVIL.EXAMPLE",
		"protocol relative":          "//evil.example",
		"whitespace padded":          " https://evil.example ",
		"embedded null":              "https://evil.example\x00.glueopshosted.rocks",
		"newline injection":          "https://evil.example\nhttps://autoglue.glueopshosted.rocks",
	}

	for name, attacker := range attacks {
		t.Run(name, func(t *testing.T) {
			got := postMessageOrigin(spaState(attacker), trustedBase)

			if got != trustedBase {
				t.Fatalf("origin = %q, want the configured base %q — tokens would be delivered to an unvetted origin",
					got, trustedBase)
			}
			// Belt and braces: the attacker's host must not appear at all, in
			// case a future change starts composing rather than replacing.
			if strings.Contains(strings.ToLower(got), "evil.example") {
				t.Fatalf("origin = %q still references the attacker host", got)
			}
		})
	}
}

func TestPostMessageOrigin_AllowsTheLegitimateOrigin(t *testing.T) {
	// Without this, "return the base unconditionally" would pass the suite and
	// silently break every real login, because the browser drops a postMessage
	// whose target does not match the opener exactly.
	t.Setenv("CORS_ALLOWED_ORIGINS", trustedBase+",http://localhost:5173")

	for _, legit := range []string{trustedBase, "http://localhost:5173"} {
		if got := postMessageOrigin(spaState(legit), trustedBase); got != legit {
			t.Errorf("origin = %q, want %q — a legitimate SPA would never receive its tokens", got, legit)
		}
	}
}

func TestPostMessageOrigin_AllowsLoopbackAlias(t *testing.T) {
	// localhost and 127.0.0.1 are different browser origins but the same host
	// to a developer. AllowedOrigins expands the alias deliberately; pin it, so
	// tightening the allowlist later does not silently break local login.
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:8080")

	if got := postMessageOrigin(spaState("http://127.0.0.1:8080"), trustedBase); got != "http://127.0.0.1:8080" {
		t.Errorf("origin = %q, want the 127.0.0.1 alias to be accepted", got)
	}
}

func TestPostMessageOrigin_FallsBackWhenStateCarriesNoOrigin(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", trustedBase)

	cases := map[string]string{
		"no origin key":     "b1a2c3d4|mode=spa",
		"empty origin":      "b1a2c3d4|mode=spa|origin=",
		"bare state":        "b1a2c3d4",
		"empty state":       "",
		"undecodable value": "b1a2c3d4|mode=spa|origin=%zz",
		"not a url":         "b1a2c3d4|mode=spa|origin=" + url.QueryEscape("not-a-url"),
	}

	for name, state := range cases {
		t.Run(name, func(t *testing.T) {
			if got := postMessageOrigin(state, trustedBase); got != trustedBase {
				t.Errorf("origin = %q, want the configured base %q", got, trustedBase)
			}
		})
	}
}

func TestPostMessageOrigin_NeverReturnsWildcard(t *testing.T) {
	// postMessage(payload, "*") delivers the token pair to whatever opened the
	// window. It must be unreachable by any input, including the degenerate
	// case where nothing is configured at all.
	for _, allowed := range []string{trustedBase, ""} {
		t.Setenv("CORS_ALLOWED_ORIGINS", allowed)
		for _, state := range []string{spaState("*"), spaState("https://evil.example"), ""} {
			if got := postMessageOrigin(state, trustedBase); got == "*" {
				t.Fatalf("origin = %q for state %q — a wildcard target hands tokens to any opener", got, state)
			}
		}
	}
}

func TestWritePostMessageHTML_TargetsExactlyOneOriginAndNotAWildcard(t *testing.T) {
	// The rendered page is the last step: whatever origin was chosen has to
	// reach postMessage as a literal target, never "*".
	rr := httptest.NewRecorder()
	writePostMessageHTML(rr, trustedBase, dto.TokenPair{
		AccessToken:  "access-token-value",
		RefreshToken: "refresh-token-value",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()

	// html/template escapes "/" as "\/" inside the JS string context, so compare
	// the extracted target rather than substring-matching the raw origin.
	target := postMessageTarget(t, body)
	if target != trustedBase {
		t.Errorf("postMessage target = %q, want %q; the SPA would never receive the tokens", target, trustedBase)
	}
	if target == "*" {
		t.Errorf("rendered page posts to a wildcard target:\n%s", body)
	}
	// The tokens are in the page, so the target origin is the only thing
	// restricting who can read them. Assert they are actually there, otherwise
	// the wildcard check above is vacuously satisfied by an empty page.
	if !strings.Contains(body, "PayloadB64") && !strings.Contains(body, "atob(") {
		t.Errorf("rendered page does not look like the postMessage bridge:\n%s", body)
	}
}

func TestWritePostMessageHTML_EscapesAnOriginContainingQuotes(t *testing.T) {
	// Defence in depth: postMessageOrigin should never let this value through,
	// but if it ever did, the template must not let it break out of the string
	// literal and rewrite the script.
	rr := httptest.NewRecorder()
	writePostMessageHTML(rr, `https://evil.example" ,"*"); //`, dto.TokenPair{AccessToken: "a"})

	target := postMessageTarget(t, rr.Body.String())
	if target == "*" {
		t.Errorf("a crafted origin broke out of the string literal and became a wildcard target")
	}
	if !strings.Contains(target, "evil.example") {
		t.Errorf("target = %q; expected the crafted value to be contained as data, not interpreted", target)
	}
}
