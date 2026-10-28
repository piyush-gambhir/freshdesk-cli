package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/config"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/registry"
	"github.com/zalando/go-keyring"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func isolate(t *testing.T) {
	t.Helper()
	for _, key := range []string{"FRESHDESK_API_KEY", "FRESHDESK_DOMAIN", "FRESHDESK_PROFILE", "FRESHDESK_READ_ONLY", "FRESHDESK_NO_INPUT", "FRESHDESK_VERBOSE"} {
		t.Setenv(key, "")
	}
	t.Setenv("FRESHDESK_CONFIG", filepath.Join(t.TempDir(), "config.yaml"))
}

func call(t *testing.T, tr http.RoundTripper, input string, args ...string) (int, string, string) {
	t.Helper()
	var out, stderr bytes.Buffer
	code := execute(context.Background(), &app{in: strings.NewReader(input), out: &out, errOut: &stderr, transport: tr}, args)
	return code, out.String(), stderr.String()
}

func response(status int, body string, headers http.Header) *http.Response {
	if headers == nil {
		headers = http.Header{}
	}
	return &http.Response{StatusCode: status, Header: headers, Body: io.NopCloser(strings.NewReader(body))}
}

func TestEveryOperationRoutesThroughTheClient(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_API_KEY", "test-api-key")
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	for _, op := range registry.Operations {
		t.Run(op.ID, func(t *testing.T) {
			args := append([]string{}, op.Command...)
			expected := op.Path
			for _, param := range op.PathParameters {
				args = append(args, "42")
				expected = strings.ReplaceAll(expected, "{"+param+"}", "42")
			}
			for _, key := range op.RequiredQuery {
				args = append(args, "--param", key+"=example")
			}
			args = append(args, "--yes", "--no-input", "-o", "json")
			calls := 0
			tr := transportFunc(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Method != op.Method || r.URL.Path != expected {
					t.Errorf("got %s %s, expected %s %s", r.Method, r.URL.Path, op.Method, expected)
				}
				if key, _, ok := r.BasicAuth(); !ok || key != "test-api-key" {
					t.Error("missing authentication")
				}
				return response(200, `{"id":42}`, nil), nil
			})
			code, out, stderr := call(t, tr, "", args...)
			if code != 0 || calls != 1 {
				t.Fatalf("exit=%d calls=%d stderr=%s", code, calls, stderr)
			}
			var v map[string]any
			if err := json.Unmarshal([]byte(out), &v); err != nil {
				t.Fatalf("invalid command JSON: %v %s", err, out)
			}
			if v["command_id"] != op.ID {
				t.Fatalf("wrong command ID: %v", v["command_id"])
			}
		})
	}
}

func TestEveryWriteBlockedAndPreviewableWithoutKey(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_READ_ONLY", "true")
	tr := transportFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("safety checks called the network")
		return nil, nil
	})
	for _, op := range registry.Operations {
		if op.Effect != "write" {
			continue
		}
		t.Run(op.ID, func(t *testing.T) {
			args := append([]string{}, op.Command...)
			for range op.PathParameters {
				args = append(args, "42")
			}
			for _, k := range op.RequiredQuery {
				args = append(args, "--param", k+"=test")
			}
			args = append(args, "--read-only=false", "--yes", "-o", "json")
			code, _, stderr := call(t, tr, "", args...)
			if code == 0 || !strings.Contains(stderr, "read_only") {
				t.Fatalf("write bypassed read-only: %s", stderr)
			}
			code, out, stderr := call(t, tr, "", append(args, "--dry-run")...)
			if code != 0 || !strings.Contains(out, `"plan_hash"`) {
				t.Fatalf("preview failed: %s %s", out, stderr)
			}
		})
	}
}

func TestPreviewHashMatchesExecutionAndBindsAttachments(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	file := filepath.Join(t.TempDir(), "evidence.txt")
	_ = os.WriteFile(file, []byte("original"), 0600)
	args := []string{"tickets", "create", "--subject", "A ticket", "--email", "support@example.com", "--set", "custom_fields.cf_team=Support", "--attach", file, "-o", "json"}
	code, out, stderr := call(t, nil, "", append(args, "--dry-run")...)
	if code != 0 {
		t.Fatal(stderr)
	}
	var preview struct {
		Plan struct {
			Hash string `json:"plan_hash"`
		}
	}
	if err := json.Unmarshal([]byte(out), &preview); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FRESHDESK_API_KEY", "test-key")
	calls := 0
	tr := transportFunc(func(r *http.Request) (*http.Response, error) { calls++; return response(201, `{"id":123}`, nil), nil })
	execArgs := append(append([]string{}, args...), "--expect-plan", preview.Plan.Hash, "--yes")
	code, out, stderr = call(t, tr, "", execArgs...)
	if code != 0 || calls != 1 || !strings.Contains(out, preview.Plan.Hash) {
		t.Fatalf("execution disagrees with plan: %s %s", out, stderr)
	}
	_ = os.WriteFile(file, []byte("changed"), 0600)
	code, _, stderr = call(t, tr, "", execArgs...)
	if code == 0 || calls != 1 || !strings.Contains(stderr, "plan_mismatch") {
		t.Fatalf("changed attachment was sent: %s", stderr)
	}
}

func TestPaginationAndBounds(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	for _, bound := range []int{1, 2} {
		calls := 0
		tr := transportFunc(func(r *http.Request) (*http.Response, error) {
			calls++
			page := r.URL.Query().Get("page")
			if page == "1" {
				return response(200, `[{"id":1}]`, http.Header{"Link": []string{`<https://acme.freshdesk.com/api/v2/tickets?page=2>; rel="next"`}}), nil
			}
			return response(200, `[{"id":2}]`, nil), nil
		})
		code, out, stderr := call(t, tr, "", "tickets", "list", "--all", "--max-pages", fmt.Sprint(bound), "-o", "json")
		if code != 0 || calls != bound {
			t.Fatalf("pagination: %d %s", calls, stderr)
		}
		var envelope struct {
			Data []any
			Meta struct {
				Complete bool
				Requests int
			}
		}
		if err := json.Unmarshal([]byte(out), &envelope); err != nil {
			t.Fatal(err)
		}
		if envelope.Meta.Complete != (bound == 2) || len(envelope.Data) != bound {
			t.Fatalf("misleading completion: %s", out)
		}
	}
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		return response(200, `[]`, http.Header{"Link": []string{`<https://attacker.invalid/api/v2/tickets?page=2>; rel=next`}}), nil
	})
	code, _, stderr := call(t, tr, "", "tickets", "list", "--all", "-o", "json")
	if code == 0 || !strings.Contains(stderr, "unsafe pagination") {
		t.Fatalf("unsafe link accepted: %s", stderr)
	}
}

func TestSearchQuotesAndFixedPaging(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	calls := 0
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Query().Get("query") != `"status:2"` {
			t.Errorf("query=%s", r.URL.RawQuery)
		}
		if calls == 1 {
			return response(200, `{"total":31,"results":[{"id":1}]}`, nil), nil
		}
		return response(200, `{"total":31,"results":[{"id":2}]}`, nil), nil
	})
	code, out, stderr := call(t, tr, "", "tickets", "search", "--query", "status:2", "--all", "-o", "json")
	if code != 0 || calls != 2 || !strings.Contains(out, `"complete": true`) {
		t.Fatalf("%s %s", out, stderr)
	}
}

func TestLocalValidationBeforeNetwork(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatal("invalid input reached network")
		return nil, nil
	})
	for _, args := range [][]string{{"tickets", "get", "../secret"}, {"tickets", "list", "--page", "0"}, {"tickets", "list", "--param", "page=0"}, {"tickets", "search"}, {"tickets", "search", "--query", "x", "--per-page", "100"}, {"tickets", "create", "--input-json", "{} {}", "--yes"}, {"tickets", "create", "--input-json", "[]", "--yes"}, {"tickets", "create", "--upload", "bad", "--yes"}, {"tickets", "delete", "42"}, {"api", "request", "GET", "https://attacker.invalid/api/v2/tickets"}} {
		code, _, stderr := call(t, tr, "", append(args, "-o", "json")...)
		if code == 0 || !json.Valid([]byte(stderr)) {
			t.Fatalf("invalid input accepted or unstructured error %v: %s", args, stderr)
		}
	}
}

func TestOfflineDiscoveryNoConfigNetworkOrSecretOutput(t *testing.T) {
	isolate(t)
	path := os.Getenv("FRESHDESK_CONFIG")
	_ = os.WriteFile(path, []byte("invalid: [yaml"), 0600)
	tr := transportFunc(func(*http.Request) (*http.Response, error) { t.Fatal("discovery made a request"); return nil, nil })
	for _, args := range [][]string{{"commands", "-o", "json"}, {"commands", "describe", "freshdesk.tickets.create", "-o", "json"}, {"commands", "resolve", "--", "api", "request", "GET", "/api/v2/ticket-forms/42/clone", "--api-key", "private-key"}, {"version", "-o", "json"}} {
		code, out, stderr := call(t, tr, "", args...)
		if code != 0 || !json.Valid([]byte(out)) || strings.Contains(out, "private-key") {
			t.Fatalf("offline discovery failed: %v %s %s", args, out, stderr)
		}
		if args[1] == "resolve" && !strings.Contains(out, `"effect": "write"`) {
			t.Fatal("GET clone was classified as read")
		}
	}
	for _, test := range []struct{ method, path, effect string }{{"GET", "/api/v2/tickets/42", "read"}, {"GET", "/api/v2/ticket-forms/42/clone", "write"}, {"GET", "/api/v2/unknown", "write"}, {"POST", "/api/v2/tickets", "write"}} {
		if op := classify(test.method, test.path); op.Effect != test.effect {
			t.Fatalf("%+v resolved as %+v", test, op)
		}
	}
}

func TestAuthProfilesAndSecrets(t *testing.T) {
	isolate(t)
	keyring.MockInit()
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	code, out, stderr := call(t, nil, "my-secret-key\n", "auth", "login", "--key-stdin", "--no-verify", "--profile", "support", "-o", "json")
	if code != 0 || strings.Contains(out, "my-secret-key") {
		t.Fatalf("%s %s", out, stderr)
	}
	path := os.Getenv("FRESHDESK_CONFIG")
	raw, _ := os.ReadFile(path)
	if bytes.Contains(raw, []byte("my-secret-key")) {
		t.Fatal("credential leaked into config")
	}
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		key, _, _ := r.BasicAuth()
		if key != "my-secret-key" {
			t.Error("did not load profile key")
		}
		return response(200, `{"id":1}`, nil), nil
	})
	code, _, stderr = call(t, tr, "", "auth", "status", "--verify", "-o", "json")
	if code != 0 {
		t.Fatal(stderr)
	}
	code, _, stderr = call(t, nil, "", "auth", "logout", "-o", "json")
	if code != 0 {
		t.Fatal(stderr)
	}
	code, out, stderr = call(t, nil, "", "auth", "status", "-o", "json")
	if code != 0 || !strings.Contains(out, `"configured": false`) {
		t.Fatalf("%s %s", out, stderr)
	}
}

func TestStickyNoInputAndProfileReadOnly(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_NO_INPUT", "1")
	code, _, stderr := call(t, nil, "yes\n", "tickets", "delete", "1", "--no-input=false", "-o", "json")
	if code == 0 || !strings.Contains(stderr, "--yes") {
		t.Fatalf("no-input bypass: %s", stderr)
	}
	path := os.Getenv("FRESHDESK_CONFIG")
	err := config.Update(context.Background(), path, func(c *config.Config) error {
		c.CurrentProfile = "locked"
		c.Profiles["locked"] = config.Profile{Domain: "https://acme.freshdesk.com", ReadOnly: true}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"tickets", "delete", "1", "--yes"}, {"auth", "logout"}, {"config", "use-profile", "another"}, {"auth", "login", "--no-verify"}, {"update", "--yes"}} {
		code, _, stderr := call(t, nil, "", append(args, "-o", "json")...)
		if code == 0 || !strings.Contains(stderr, "read-only") {
			t.Fatalf("profile restriction bypass %v: %s", args, stderr)
		}
	}
}

func TestErrorsNeverExposeCredential(t *testing.T) {
	isolate(t)
	secret := "secret-quote-\"-key"
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", secret)
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(map[string]any{"description": secret})
		return response(400, string(body), nil), nil
	})
	code, out, stderr := call(t, tr, "", "tickets", "list", "-o", "json")
	if code != 1 || out != "" || strings.Contains(stderr, "secret-quote") || !json.Valid([]byte(stderr)) {
		t.Fatalf("credential leak or malformed error: %s", stderr)
	}
}

func TestRawReadOnlyDoesNotSendBody(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	code, _, stderr := call(t, nil, "", "api", "request", "GET", "/api/v2/tickets", "--input-json", `{"subject":"x"}`, "--dry-run", "-o", "json")
	if code == 0 || !strings.Contains(stderr, "do not accept") {
		t.Fatal(stderr)
	}
}

func TestEmptyPaginationAndExplicitEmptyBody(t *testing.T) {
	isolate(t)
	t.Setenv("FRESHDESK_DOMAIN", "acme")
	t.Setenv("FRESHDESK_API_KEY", "key")
	tr := transportFunc(func(r *http.Request) (*http.Response, error) {
		return response(200, `[]`, nil), nil
	})
	code, out, stderr := call(t, tr, "", "contacts", "list", "--all", "-o", "json")
	if code != 0 || !strings.Contains(out, `"data": []`) {
		t.Fatalf("empty results must remain an array: %s %s", out, stderr)
	}
	code, out, stderr = call(t, nil, "", "tickets", "create", "--input-json", "{}", "--dry-run", "-o", "json")
	if code != 0 || !strings.Contains(out, `"body": {}`) {
		t.Fatalf("explicit empty JSON was dropped: %s %s", out, stderr)
	}
	code, out, stderr = call(t, nil, "", "tickets", "update", "42", "--status", "resolved", "--priority", "high", "--dry-run", "-o", "json")
	if code != 0 || !strings.Contains(out, `"status": 4`) || !strings.Contains(out, `"priority": 3`) {
		t.Fatalf("named ticket values were not mapped: %s %s", out, stderr)
	}
}
