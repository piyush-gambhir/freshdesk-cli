package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/output"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/registry"
	"github.com/spf13/cobra"
)

const maxInput = 4 << 20

type requestFlags struct {
	input, bodyFile, ifMatch          string
	set, params, attachments, uploads []string
	fields                            map[string]*string
	page, perPage, maxPages           int
	all, multipart                    bool
}

func (a *app) addOperation(root *cobra.Command, op registry.Operation) {
	parent := root
	for _, name := range op.Command[:len(op.Command)-1] {
		var found *cobra.Command
		for _, c := range parent.Commands() {
			if c.Name() == name {
				found = c
				break
			}
		}
		if found == nil {
			found = &cobra.Command{Use: name, Short: "Manage " + strings.ReplaceAll(name, "-", " ")}
			parent.AddCommand(found)
		}
		parent = found
	}
	use := op.Command[len(op.Command)-1]
	for _, p := range op.PathParameters {
		use += " " + strings.ToUpper(strings.ReplaceAll(p, "_", "-"))
	}
	short := strings.ToUpper(op.Method) + " " + op.Path
	c := &cobra.Command{Use: use, Short: short, Long: short + "\n\nReference: https://developers.freshdesk.com/api/#" + op.SourceAnchor,
		Args: cobra.ExactArgs(len(op.PathParameters)), Annotations: map[string]string{"effect": op.Effect, "operation": op.ID}}
	if op.Note != "" {
		c.Long += "\n\n" + op.Note
	}
	flags := a.requestFlags(c, op.Effect == "write", op.Paginated)
	if op.Pagination == "cursor" {
		_ = c.Flags().MarkHidden("page")
	}
	c.RunE = func(cmd *cobra.Command, args []string) error {
		path := op.Path
		for i, param := range op.PathParameters {
			if err := validSegment(args[i]); err != nil {
				return fmt.Errorf("%s: %w", param, err)
			}
			path = strings.ReplaceAll(path, "{"+param+"}", url.PathEscape(args[i]))
		}
		return a.request(cmd, op, path, flags)
	}
	parent.AddCommand(c)
}

func validSegment(s string) error {
	if s == "" || s == "." || s == ".." || len(s) > 256 || strings.ContainsAny(s, "/\\?#%{}\r\n\x00\t ") {
		return fmt.Errorf("expected a nonempty resource ID or language code, without URL delimiters or whitespace")
	}
	return nil
}

func (a *app) requestFlags(c *cobra.Command, write, paginated bool) *requestFlags {
	r := &requestFlags{fields: map[string]*string{}}
	f := c.Flags()
	f.StringArrayVar(&r.params, "param", nil, "Query parameter KEY=VALUE; repeat for multiple values")
	for _, key := range []string{"query", "term", "include", "filter", "updated_since", "order_by", "order_type", "email", "requester_id", "company_id", "since", "until", "next_token", "prev_token"} {
		if write && (key == "email" || key == "requester_id" || key == "company_id") {
			continue
		}
		flag := strings.ReplaceAll(key, "_", "-")
		r.fields["query:"+key] = f.String(flag, "", "Query parameter "+key)
	}
	if write {
		f.StringVar(&r.input, "input-json", "", "JSON object, @FILE, or - for stdin")
		f.StringArrayVar(&r.set, "set", nil, "Body FIELD=VALUE (JSON value or text); dotted fields build nested objects")
		f.StringVar(&r.bodyFile, "body-file", "", "Read the body field as UTF-8 text or HTML from FILE (or -)")
		f.StringArrayVar(&r.attachments, "attach", nil, "Attach a local file as attachments[]; repeat for more files")
		f.StringArrayVar(&r.uploads, "upload", nil, "Upload FIELD=FILE for import CSVs, avatars, or other named file fields")
		f.BoolVar(&r.multipart, "multipart", false, "Send multipart form data, even without attachments")
		f.StringVar(&r.ifMatch, "if-match", "", "Send an If-Match precondition (when supported by the endpoint)")
		for _, key := range []string{"subject", "description", "body", "email", "name", "phone", "mobile", "title", "content_html", "status", "priority", "source", "requester_id", "responder_id", "group_id", "company_id", "type", "user_id", "private", "billable", "note", "time_spent"} {
			flag := strings.ReplaceAll(key, "_", "-")
			if c.Flags().Lookup(flag) != nil {
				continue
			}
			r.fields["body:"+key] = f.String(flag, "", "Body field "+key)
		}
	}
	if paginated {
		f.IntVar(&r.page, "page", 1, "Page number (starts at 1)")
		f.IntVar(&r.perPage, "per-page", 30, "Page size, 1-100 (search APIs use 30)")
		f.BoolVar(&r.all, "all", false, "Collect pages, stopping at --max-pages or the API limit")
		f.IntVar(&r.maxPages, "max-pages", 100, "Maximum pages to collect with --all, 1-500")
	}
	return r
}

func readLimited(r io.Reader, limit int64) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("input exceeds %d bytes", limit)
	}
	return b, nil
}

func (a *app) readFile(path string) ([]byte, error) {
	if path == "-" {
		return readLimited(a.in, maxInput)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		return nil, fmt.Errorf("input must be a regular file")
	}
	return readLimited(f, maxInput)
}

func decodeObject(raw []byte) (map[string]any, error) {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var value map[string]any
	if err := d.Decode(&value); err != nil || value == nil {
		return nil, fmt.Errorf("input must be one JSON object")
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("input must contain exactly one JSON object")
	}
	return value, nil
}

func setField(m map[string]any, key string, value any) error {
	parts := strings.Split(key, ".")
	for i, p := range parts {
		if p == "" || strings.ContainsAny(p, "\r\n\x00") {
			return fmt.Errorf("invalid body field name")
		}
		if i == len(parts)-1 {
			m[p] = value
			return nil
		}
		if m[p] == nil {
			m[p] = map[string]any{}
		}
		next, ok := m[p].(map[string]any)
		if !ok {
			return fmt.Errorf("body field %s conflicts with an existing scalar", key)
		}
		m = next
	}
	return nil
}

func value(raw string) any {
	d := json.NewDecoder(strings.NewReader(raw))
	d.UseNumber()
	var v any
	if d.Decode(&v) == nil {
		var extra any
		if d.Decode(&extra) == io.EOF {
			return v
		}
	}
	return raw
}

func (a *app) prepare(cmd *cobra.Command, op registry.Operation, path string, r *requestFlags) (*client.Plan, error) {
	if r.input == "-" && r.bodyFile == "-" {
		return nil, fmt.Errorf("stdin can supply either --input-json or --body-file, not both")
	}
	if _, err := client.Path(path); err != nil {
		return nil, err
	}
	origin, err := a.target()
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	for k, v := range op.QueryDefaults {
		q.Set(k, v)
	}
	for _, p := range r.params {
		k, v, ok := strings.Cut(p, "=")
		if !ok || k == "" || strings.ContainsAny(k, "\r\n\x00") {
			return nil, fmt.Errorf("--param requires KEY=VALUE")
		}
		q.Add(k, v)
	}
	for k, p := range r.fields {
		if strings.HasPrefix(k, "query:") {
			key := strings.TrimPrefix(k, "query:")
			name := strings.ReplaceAll(key, "_", "-")
			if cmd.Flags().Changed(name) {
				q.Set(key, *p)
			}
		}
	}
	for _, k := range op.RequiredQuery {
		if q.Get(k) == "" {
			return nil, fmt.Errorf("this operation requires --%s or --param %s=VALUE", strings.ReplaceAll(k, "_", "-"), k)
		}
	}
	if strings.Contains(path, "/search/") && q.Get("query") != "" {
		query := q.Get("query")
		if !strings.HasPrefix(query, `"`) {
			query = `"` + query + `"`
		}
		q.Set("query", query)
	}
	if op.Paginated {
		if r.page < 1 || r.page > 500 || r.perPage < 1 || r.perPage > 100 || r.maxPages < 1 || r.maxPages > 500 {
			return nil, fmt.Errorf("page/max-pages must be 1-500 and per-page must be 1-100")
		}
		if op.Pagination == "cursor" && (cmd.Flags().Changed("page") || q.Get("page") != "") {
			return nil, fmt.Errorf("this endpoint uses --next-token or --prev-token instead of --page")
		}
		if op.Pagination != "cursor" && (q.Get("page") == "" || cmd.Flags().Changed("page")) {
			q.Set("page", strconv.Itoa(r.page))
		}
		if strings.HasPrefix(path, "/api/v2/search/") {
			if cmd.Flags().Changed("per-page") && r.perPage != 30 {
				return nil, fmt.Errorf("search APIs use 30 results per page")
			}
		} else if q.Get("per_page") == "" || cmd.Flags().Changed("per-page") {
			q.Set("per_page", strconv.Itoa(r.perPage))
		}
	}
	for _, key := range []string{"page", "per_page"} {
		if s := q.Get(key); s != "" {
			n, err := strconv.Atoi(s)
			max := 500
			if key == "per_page" {
				max = 100
			}
			if err != nil || n < 1 || n > max || len(q[key]) != 1 {
				return nil, fmt.Errorf("%s must be a single integer from 1 to %d", key, max)
			}
		}
	}
	var body map[string]any
	if r.input != "" {
		raw := []byte(r.input)
		if r.input == "-" {
			raw, err = a.readFile("-")
		} else if strings.HasPrefix(r.input, "@") {
			raw, err = a.readFile(r.input[1:])
		}
		if err != nil {
			return nil, err
		}
		if len(raw) > maxInput {
			return nil, fmt.Errorf("JSON input exceeds 4 MiB")
		}
		body, err = decodeObject(raw)
		if err != nil {
			return nil, err
		}
	}
	if len(r.set) > 0 || r.bodyFile != "" {
		if body == nil {
			body = map[string]any{}
		}
	}
	for _, s := range r.set {
		k, v, ok := strings.Cut(s, "=")
		if !ok {
			return nil, fmt.Errorf("--set requires FIELD=VALUE")
		}
		if err := setField(body, k, value(v)); err != nil {
			return nil, err
		}
	}
	for k, p := range r.fields {
		if strings.HasPrefix(k, "body:") {
			key := strings.TrimPrefix(k, "body:")
			name := strings.ReplaceAll(key, "_", "-")
			if cmd.Flags().Changed(name) {
				if body == nil {
					body = map[string]any{}
				}
				if err := setField(body, key, bodyValue(key, *p)); err != nil {
					return nil, err
				}
			}
		}
	}
	if r.bodyFile != "" {
		raw, err := a.readFile(r.bodyFile)
		if err != nil {
			return nil, err
		}
		body["body"] = string(raw)
	}
	for k, v := range op.BodyDefaults {
		if body == nil {
			body = map[string]any{}
		}
		if current, ok := body[k]; ok && current != v {
			return nil, fmt.Errorf("%s requires %s=%v", strings.Join(op.Command, " "), k, v)
		}
		body[k] = v
	}
	if strings.HasPrefix(op.ID, "freshdesk.tickets.") || strings.HasPrefix(op.ID, "freshdesk.service.tasks.") {
		if raw, ok := body["status"].(string); ok {
			if n, found := map[string]int{"open": 2, "pending": 3, "resolved": 4, "closed": 5}[strings.ToLower(raw)]; found {
				body["status"] = n
			}
		}
		if raw, ok := body["priority"].(string); ok {
			if n, found := map[string]int{"low": 1, "medium": 2, "high": 3, "urgent": 4}[strings.ToLower(raw)]; found {
				body["priority"] = n
			}
		}
	}
	p := &client.Plan{Contract: "agent-cli/v1", CommandID: op.ID, CatalogHash: registry.Hash(), Version: build.Version, Effect: op.Effect, Method: op.Method, URL: origin + path, Multipart: r.multipart, Precondition: r.ifMatch}
	if body != nil {
		p.Body = body
	}
	if len(q) > 0 {
		p.URL += "?" + q.Encode()
	}
	for _, f := range r.attachments {
		r.uploads = append(r.uploads, "attachments[]="+f)
	}
	total := 0
	for _, upload := range r.uploads {
		field, file, ok := strings.Cut(upload, "=")
		if !ok {
			return nil, fmt.Errorf("--upload requires FIELD=FILE")
		}
		attachment, err := client.ReadAttachment(field, file)
		if err != nil {
			return nil, err
		}
		total += attachment.Size
		if total > client.MaxAttachments {
			return nil, fmt.Errorf("total attachment size exceeds 20 MiB")
		}
		p.Attachments = append(p.Attachments, attachment)
	}
	if len(p.Attachments) > 0 {
		p.Multipart = true
	}
	if p.Precondition != "" && strings.ContainsAny(p.Precondition, "\r\n\x00") {
		return nil, fmt.Errorf("invalid If-Match value")
	}
	if err := p.Seal(); err != nil {
		return nil, err
	}
	return p, nil
}

func bodyValue(key, raw string) any {
	// Free text must stay text even if it happens to look like JSON or a number.
	switch key {
	case "status", "priority", "source", "requester_id", "responder_id", "group_id", "company_id", "user_id", "private", "billable":
		return value(raw)
	}
	return raw
}

func (a *app) request(cmd *cobra.Command, op registry.Operation, path string, r *requestFlags) error {
	p, err := a.prepare(cmd, op, path, r)
	if err != nil {
		return err
	}
	if a.expectPlan != "" && a.expectPlan != p.Hash {
		return &client.APIError{Kind: "plan_mismatch", Message: "request differs from --expect-plan; run --dry-run again", Outcome: "failed"}
	}
	if a.dryRun {
		return a.print(map[string]any{"dry_run": true, "plan": p})
	}
	if p.Effect == "write" {
		ro, err := a.profileReadOnly()
		if err != nil {
			return err
		}
		if ro {
			return &client.APIError{Kind: "read_only", Message: "operation is blocked by read-only mode", Outcome: "failed"}
		}
		if err := a.confirm(p.Method + " " + path); err != nil {
			return err
		}
	}
	c, err := a.client()
	if err != nil {
		return err
	}
	return a.send(cmd, c, *p, op, r)
}

var nextLink = regexp.MustCompile(`<([^>]+)>\s*;\s*rel="?next"?`)

func nextPage(p client.Plan, res *client.Response, op registry.Operation) (string, error) {
	if !op.Paginated {
		return "", nil
	}
	current, _ := url.Parse(p.URL)
	if op.Pagination == "cursor" {
		if m, ok := res.Data.(map[string]any); ok {
			if paging, ok := m["paging"].(map[string]any); ok {
				if value, ok := paging["next"].(string); ok && value != "" {
					next, err := current.Parse(value)
					if err != nil || next.Scheme != current.Scheme || next.Host != current.Host || next.Path != current.Path || next.User != nil || next.Fragment != "" {
						return "", fmt.Errorf("unsafe cursor pagination link from Freshdesk")
					}
					token := next.Query().Get("next_token")
					if token == "" || token == current.Query().Get("next_token") {
						return "", fmt.Errorf("cursor pagination did not advance")
					}
					q := current.Query()
					q.Del("prev_token")
					q.Set("next_token", token)
					current.RawQuery = q.Encode()
					return current.String(), nil
				}
			}
		}
		return "", nil
	}
	if match := nextLink.FindStringSubmatch(res.Header.Get("Link")); match != nil {
		next, err := current.Parse(match[1])
		if err != nil || next.Scheme != current.Scheme || next.Host != current.Host || next.Path != current.Path || next.User != nil || next.Fragment != "" {
			return "", fmt.Errorf("unsafe pagination link from Freshdesk")
		}
		n, err := strconv.Atoi(next.Query().Get("page"))
		previous, _ := strconv.Atoi(current.Query().Get("page"))
		if err != nil || n != previous+1 {
			return "", fmt.Errorf("pagination did not advance by one page")
		}
		// Reuse the original filters, never API-provided arbitrary query values.
		q := current.Query()
		q.Set("page", strconv.Itoa(n))
		current.RawQuery = q.Encode()
		return current.String(), nil
	}
	if m, ok := res.Data.(map[string]any); ok {
		if rows, ok := m["results"].([]any); ok {
			total, _ := strconv.Atoi(fmt.Sprint(m["total"]))
			page, _ := strconv.Atoi(current.Query().Get("page"))
			if len(rows) > 0 && page*30 < total {
				q := current.Query()
				q.Set("page", strconv.Itoa(page+1))
				current.RawQuery = q.Encode()
				return current.String(), nil
			}
		}
	}
	return "", nil
}

func (a *app) send(cmd *cobra.Command, c *client.Client, p client.Plan, op registry.Operation, r *requestFlags) error {
	requests := 0
	initialURL, _ := url.Parse(p.URL)
	initialPage, _ := strconv.Atoi(initialURL.Query().Get("page"))
	var data any
	rows := []any{}
	var next string
	var status int
	warnings := []string{}
	visited := map[string]bool{p.URL: true}
	if op.Note != "" {
		warnings = append(warnings, op.Note)
	}
	if op.ID == "freshdesk.tickets.list" && initialURL.Query().Get("updated_since") == "" {
		warnings = append(warnings, "Freshdesk lists tickets created in the last 30 days by default; use --updated-since for older tickets. The API returns at most 300 pages.")
	}
	for {
		res, err := c.Do(cmd.Context(), p)
		if err != nil {
			return err
		}
		requests++
		status = res.Status
		if !r.all {
			data = res.Data
		} else {
			pageRows, ok := extractRows(res.Data)
			if !ok {
				return fmt.Errorf("this response has no recognized result array; use a single-page request")
			}
			rows = append(rows, pageRows...)
			data = rows
		}
		next, err = nextPage(p, res, op)
		if err != nil {
			return err
		}
		if next == "" || !r.all {
			break
		}
		if visited[next] {
			return fmt.Errorf("pagination repeated an earlier page")
		}
		visited[next] = true
		u, _ := url.Parse(next)
		n, _ := strconv.Atoi(u.Query().Get("page"))
		limit := 500
		if op.ID == "freshdesk.tickets.list" {
			limit = 300
		}
		if strings.HasPrefix(u.Path, "/api/v2/search/") {
			limit = 10
		}
		if requests >= r.maxPages || n > limit {
			warnings = append(warnings, "Stopped at the requested page bound or Freshdesk's pagination limit; results are incomplete.")
			break
		}
		p.URL = next
	}
	meta := map[string]any{"requests": requests, "http_status": status, "warnings": warnings}
	if op.Paginated {
		if op.Pagination == "cursor" {
			meta["pagination"] = "cursor"
		} else {
			meta["start_page"] = initialPage
		}
		meta["complete"] = next == ""
		if next != "" {
			meta["next_page_url"] = next
		}
	} else {
		meta["complete"] = nil
	}
	if p.Effect == "write" {
		meta["outcome"] = "applied"
		meta["plan_hash"] = p.Hash
		if status == 202 {
			meta["outcome"] = "accepted"
		}
		if object, ok := data.(map[string]any); ok && object["job_id"] != nil {
			meta["outcome"] = "accepted"
		}
	}
	if !a.quiet && (a.format == "table" || a.format == "csv") {
		for _, warning := range warnings {
			fmt.Fprintln(a.errOut, "Note:", warning)
		}
	}
	full := map[string]any{"contract_version": "agent-cli/v1", "command_id": op.ID, "data": data, "meta": meta}
	var tableRows []map[string]any
	if list, ok := extractRows(data); ok {
		tableRows = make([]map[string]any, 0, len(list))
		for _, v := range list {
			if row, ok := v.(map[string]any); ok {
				tableRows = append(tableRows, row)
			} else {
				tableRows = append(tableRows, map[string]any{"value": v})
			}
		}
	}
	if (a.format == "table" || a.format == "csv") && tableRows == nil {
		if data == nil && p.Effect == "write" {
			return output.Print(a.out, a.format, map[string]any{"http_status": status, "outcome": meta["outcome"]})
		}
		return output.Print(a.out, a.format, data)
	}
	return output.View(a.out, a.format, full, tableRows)
}

func extractRows(data any) ([]any, bool) {
	if list, ok := data.([]any); ok {
		return list, true
	}
	if m, ok := data.(map[string]any); ok {
		for _, key := range []string{"results", "records", "data", "surveys", "responses", "custom_objects", "schemas", "agents", "groups", "ticket_fields", "ticket_forms", "threads", "messages", "imports", "tickets", "rules", "sections", "watcher_ids", "user_ids"} {
			if list, ok := m[key].([]any); ok {
				return list, true
			}
		}
	}
	return nil, false
}
