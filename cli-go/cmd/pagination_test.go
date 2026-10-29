package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestCSATCursorPaginationPreservesFilters(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	calls := 0
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Query().Get("state") != "ACTIVE" || r.URL.Query().Get("page") != "" {
			t.Errorf("cursor query %s", r.URL.RawQuery)
		}
		if calls == 1 {
			return response(200, `{"data":[{"id":"first"}],"paging":{"next":"https://acme.freshdesk.com/api/v2/customer-satisfaction/surveys?next_token=next123&per_page=30"}}`, nil), nil
		}
		if r.URL.Query().Get("next_token") != "next123" {
			t.Errorf("lost cursor: %s", r.URL.RawQuery)
		}
		return response(200, `{"data":[{"id":"second"}],"paging":{"next":null}}`, nil), nil
	})
	code, out, stderr := call(t, tr, "", "csat", "surveys", "list", "--param", "state=ACTIVE", "--all", "-o", "json")
	if code != 0 || calls != 2 || !strings.Contains(out, `"complete": true`) || !strings.Contains(out, "second") {
		t.Fatalf("%s %s calls=%d", out, stderr, calls)
	}
}

func TestCSATHTTPLinkKeepsRequestsOnHTTPS(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	calls := 0
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Scheme != "https" || r.URL.Host != "acme.freshdesk.com" {
			t.Fatalf("cursor downgraded request security: %s", r.URL)
		}
		if calls == 1 {
			return response(200, `{"data":[{"id":"first"}],"paging":{"next":"http://acme.freshdesk.com/api/v2/customer-satisfaction/surveys?next_token=next123&per_page=1"}}`, nil), nil
		}
		if r.URL.Query().Get("next_token") != "next123" {
			t.Error("cursor token was not preserved")
		}
		return response(200, `{"data":[{"id":"second"}],"paging":{"next":null}}`, nil), nil
	})
	code, out, stderr := call(t, tr, "", "csat", "surveys", "list", "--per-page", "1", "--all", "-o", "json")
	if code != 0 || calls != 2 || !strings.Contains(out, `"complete": true`) {
		t.Fatalf("%s %s calls=%d", out, stderr, calls)
	}
}

func TestCSATCursorRejectsCrossOriginAndLoops(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	for _, next := range []string{"https://evil.invalid/api/v2/customer-satisfaction/surveys?next_token=leak", "https://acme.freshdesk.com/api/v2/customer-satisfaction/surveys?next_token=again"} {
		calls := 0
		tr := transportFunc(func(r *http.Request) (*http.Response, error) {
			calls++
			return response(200, fmt.Sprintf(`{"data":[],"paging":{"next":%q}}`, next), nil), nil
		})
		code, _, stderr := call(t, tr, "", "csat", "surveys", "list", "--next-token", "again", "--all", "-o", "json")
		if code == 0 || calls != 1 || !strings.Contains(stderr, "pagination") {
			t.Fatalf("unsafe cursor accepted: %s calls=%d", stderr, calls)
		}
	}
}

func TestAutomationRuleEnvelopeAndRequestedStartPage(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		return response(200, `{"rules":[{"id":1}],"meta":{"count":1}}`, nil), nil
	})
	code, out, stderr := call(t, tr, "", "automations", "list", "1", "--all", "--param", "page=3", "-o", "json")
	if code != 0 || !strings.Contains(out, `"start_page": 3`) || !strings.Contains(out, `"id": 1`) {
		t.Fatalf("%s %s", out, stderr)
	}
}
