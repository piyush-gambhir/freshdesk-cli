// Package client sends bounded, authenticated Freshdesk requests without retries.
package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const MaxBody = 32 << 20
const MaxAttachments = 20 << 20

var slug = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
var hostPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9.-]*[a-zA-Z0-9])?$`)

// Domain accepts a Freshdesk account slug or an explicit HTTPS hostname. Paths,
// credentials, queries, and ports are rejected rather than silently discarded.
func Domain(value string) (string, error) {
	value = strings.TrimSpace(value)
	if slug.MatchString(value) {
		value += ".freshdesk.com"
	}
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}
	u, err := url.Parse(value)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") || u.Port() != "" || !hostPattern.MatchString(u.Hostname()) || !strings.Contains(u.Hostname(), ".") {
		return "", fmt.Errorf("domain must be an account slug, hostname, or HTTPS origin without a path, port, or credentials")
	}
	for _, part := range strings.Split(u.Hostname(), ".") {
		if !slug.MatchString(part) {
			return "", fmt.Errorf("domain contains an invalid hostname label")
		}
	}
	return "https://" + strings.ToLower(u.Hostname()), nil
}

// Path confines raw calls and pagination to helpdesk API routes.
func Path(value string) (string, error) {
	if !strings.HasPrefix(value, "/api/v2/") && value != "/scenario_automations.json" {
		return "", fmt.Errorf("path must start with /api/v2/ (or be /scenario_automations.json)")
	}
	u, err := url.Parse(value)
	decoded, derr := url.PathUnescape(value)
	if err != nil || derr != nil || u.Host != "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || strings.ContainsAny(decoded, "%\\\r\n\x00?#{}") || strings.Contains(decoded, "//") || strings.Count(decoded, "/") != strings.Count(value, "/") {
		return "", fmt.Errorf("path must be a relative API path without a query, fragment, or placeholders; use --param for query values")
	}
	for _, part := range strings.Split(decoded, "/") {
		if part == ".." || part == "." {
			return "", fmt.Errorf("path traversal is not allowed")
		}
	}
	return value, nil
}

type Attachment struct {
	Field  string `json:"field"`
	Name   string `json:"name"`
	Size   int    `json:"size"`
	SHA256 string `json:"sha256"`
	Data   []byte `json:"-"`
}

// ReadAttachment freezes content before preview and execution so the hash binds
// the bytes sent, rather than an on-disk filename that can change later.
func ReadAttachment(field, path string) (Attachment, error) {
	if field == "" || strings.ContainsAny(field, "\r\n\x00") {
		return Attachment{}, fmt.Errorf("invalid upload field")
	}
	f, err := os.Open(path)
	if err != nil {
		return Attachment{}, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		return Attachment{}, fmt.Errorf("attachment must be a regular file")
	}
	data, err := io.ReadAll(io.LimitReader(f, MaxAttachments+1))
	if err != nil {
		return Attachment{}, err
	}
	if len(data) > MaxAttachments {
		return Attachment{}, fmt.Errorf("attachment exceeds 20 MiB")
	}
	sum := sha256.Sum256(data)
	return Attachment{Field: field, Name: filepath.Base(path), Size: len(data), SHA256: hex.EncodeToString(sum[:]), Data: data}, nil
}

type Plan struct {
	Contract     string       `json:"contract_version"`
	CommandID    string       `json:"command_id"`
	CatalogHash  string       `json:"catalog_hash"`
	Version      string       `json:"cli_version"`
	Effect       string       `json:"effect"`
	Method       string       `json:"method"`
	URL          string       `json:"url"`
	Body         any          `json:"body,omitempty"`
	Multipart    bool         `json:"multipart,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
	Precondition string       `json:"if_match,omitempty"`
	Hash         string       `json:"plan_hash"`
}

func (p *Plan) Seal() error {
	p.Hash = ""
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	p.Hash = hex.EncodeToString(sum[:])
	return nil
}

type APIError struct {
	Kind       string `json:"kind"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"http_status,omitempty"`
	RetryAfter string `json:"retry_after,omitempty"`
	Details    any    `json:"details,omitempty"`
	Outcome    string `json:"outcome"`
}

func (e *APIError) Error() string { return e.Message }

type Response struct {
	Data   any
	Header http.Header
	Status int
}

type Client struct {
	Origin    string
	APIKey    string
	Version   string
	Timeout   time.Duration
	Transport http.RoundTripper
	Log       io.Writer
}

func (c Client) Scrub(text string) string {
	if c.APIKey != "" {
		text = strings.ReplaceAll(text, c.APIKey, "[REDACTED]")
	}
	return text
}

func (c Client) Do(ctx context.Context, p Plan) (*Response, error) {
	u, err := url.Parse(p.URL)
	if err != nil || u.Scheme+"://"+u.Host != c.Origin || u.User != nil || u.Fragment != "" {
		return nil, fmt.Errorf("request target must match the configured helpdesk origin")
	}
	if _, err := Path(u.EscapedPath()); err != nil {
		return nil, err
	}
	body, contentType, err := encodeBody(p)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, p.Method, p.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %s", c.Scrub(err.Error()))
	}
	req.SetBasicAuth(c.APIKey, "X")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "freshdesk-cli/"+c.Version)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if p.Precondition != "" {
		req.Header.Set("If-Match", p.Precondition)
	}
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	h := &http.Client{Transport: c.Transport, Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	if c.Log != nil {
		// Query values can contain customer data. Log only the method and path.
		fmt.Fprintf(c.Log, "%s %s\n", p.Method, c.Scrub(u.EscapedPath()))
	}
	res, err := h.Do(req)
	if err != nil {
		outcome := "failed"
		if p.Effect == "write" {
			outcome = "unknown"
		}
		return nil, &APIError{Kind: "network", Message: c.Scrub(err.Error()), Outcome: outcome}
	}
	defer res.Body.Close()
	if c.Log != nil {
		fmt.Fprintf(c.Log, "HTTP %d\n", res.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(res.Body, MaxBody+1))
	if err != nil || len(data) > MaxBody {
		outcome := "failed"
		if p.Effect == "write" {
			outcome = "unknown"
		}
		return nil, &APIError{Kind: "response", Message: "response is incomplete or exceeds 32 MiB", HTTPStatus: res.StatusCode, Outcome: outcome}
	}
	var value any
	if len(bytes.TrimSpace(data)) != 0 {
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		err = decoder.Decode(&value)
		var extra any
		if err == nil && !errors.Is(decoder.Decode(&extra), io.EOF) {
			err = fmt.Errorf("multiple JSON values")
		}
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		kind := "api"
		switch res.StatusCode {
		case 401:
			kind = "authentication"
		case 403:
			kind = "permission"
		case 404:
			kind = "not_found"
		case 409, 412:
			kind = "conflict"
		case 429:
			kind = "rate_limit"
		case 400, 422:
			kind = "validation"
		}
		outcome := "failed"
		if p.Effect == "write" && res.StatusCode >= 500 {
			outcome = "unknown"
		}
		message := fmt.Sprintf("Freshdesk returned HTTP %d", res.StatusCode)
		if res.StatusCode >= 300 && res.StatusCode < 400 {
			message += "; redirects are not followed, check the Freshdesk account domain"
		}
		if err != nil {
			value = nil // HTML errors can contain credentials and should not reach terminals.
		}
		return nil, &APIError{Kind: kind, Message: message, HTTPStatus: res.StatusCode, Details: value, RetryAfter: res.Header.Get("Retry-After"), Outcome: outcome}
	}
	if err != nil {
		outcome := "failed"
		if p.Effect == "write" {
			outcome = "unknown"
		}
		return nil, &APIError{Kind: "response", Message: "Freshdesk returned invalid JSON", HTTPStatus: res.StatusCode, Outcome: outcome}
	}
	return &Response{Data: value, Header: res.Header, Status: res.StatusCode}, nil
}

func encodeBody(p Plan) ([]byte, string, error) {
	if !p.Multipart && len(p.Attachments) == 0 {
		if p.Body == nil {
			return nil, "", nil
		}
		body, err := json.Marshal(p.Body)
		return body, "application/json", err
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if p.Body != nil {
		m, ok := p.Body.(map[string]any)
		if !ok {
			return nil, "", fmt.Errorf("multipart input must be a JSON object")
		}
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if err := formField(w, k, m[k]); err != nil {
				return nil, "", err
			}
		}
	}
	total := 0
	for _, f := range p.Attachments {
		total += len(f.Data)
		if total > MaxAttachments {
			return nil, "", fmt.Errorf("total attachment size exceeds 20 MiB")
		}
		h := make(textproto.MIMEHeader)
		quote := strings.NewReplacer("\\", "\\\\", "\"", "\\\"", "\r", "", "\n", "")
		h.Set("Content-Disposition", `form-data; name="`+quote.Replace(f.Field)+`"; filename="`+quote.Replace(f.Name)+`"`)
		h.Set("Content-Type", http.DetectContentType(f.Data))
		part, err := w.CreatePart(h)
		if err != nil {
			return nil, "", err
		}
		if _, err = part.Write(f.Data); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return b.Bytes(), w.FormDataContentType(), nil
}

func formField(w *multipart.Writer, key string, v any) error {
	switch x := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if err := formField(w, key+"["+k+"]", x[k]); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range x {
			if err := formField(w, key+"[]", item); err != nil {
				return err
			}
		}
	case nil:
		return w.WriteField(key, "")
	case string:
		return w.WriteField(key, x)
	default:
		return w.WriteField(key, fmt.Sprint(v))
	}
	return nil
}
