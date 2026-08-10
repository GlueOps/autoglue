package handlers

import "testing"

func TestOriginFromState(t *testing.T) {
	cases := map[string]string{
		// What AuthStart actually builds.
		"a1b2|mode=spa|origin=http%3A%2F%2F127.0.0.1%3A8080": "http://127.0.0.1:8080",
		"a1b2|mode=spa|origin=http%3A%2F%2Flocalhost%3A8080": "http://localhost:8080",
		// Non-SPA state carries no origin.
		"a1b2":          "",
		"a1b2|mode=spa": "",
		"":              "",
		"origin=%zz":    "", // undecodable
		"origin=plain":  "plain",
	}
	for state, want := range cases {
		if got := originFromState(state); got != want {
			t.Errorf("originFromState(%q) = %q, want %q", state, got, want)
		}
	}
}
