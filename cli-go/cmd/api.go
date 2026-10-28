package cmd

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/registry"
	"github.com/spf13/cobra"
)

func (a *app) apiCmd() *cobra.Command {
	c := &cobra.Command{Use: "api", Short: "Inspect API coverage or send a raw request"}
	c.AddCommand(&cobra.Command{Use: "methods", Short: "List the reviewed API operations offline", Args: cobra.NoArgs, Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error { return a.print(registry.Operations) }})
	raw := &cobra.Command{Use: "request METHOD PATH", Short: "Call an API route with the same preview and write controls", Long: "Known routes use the registry's semantic effect, including GET /ticket-forms/ID/clone.\nUnknown routes are treated as writes, including GET. Paths stay on the configured helpdesk origin.\nPagination is manual for raw calls; use --param page=N. Named commands support --all where documented.", Args: cobra.ExactArgs(2), Annotations: effect("conditional")}
	f := a.requestFlags(raw, true, false)
	raw.RunE = func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		switch method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			return fmt.Errorf("method must be GET, POST, PUT, PATCH, or DELETE")
		}
		op := classify(method, args[1])
		if op.Effect == "read" && (f.input != "" || len(f.set) > 0 || len(f.attachments) > 0 || len(f.uploads) > 0 || f.bodyFile != "" || f.multipart) {
			return fmt.Errorf("read operations do not accept a request body or uploads")
		}
		// Raw requests do not acquire named-operation positional/default rules.
		op.QueryDefaults = nil
		op.RequiredQuery = nil
		op.BodyDefaults = nil
		op.Paginated = false
		return a.request(cmd, op, args[1], f)
	}
	c.AddCommand(raw)
	return c
}

var placeholder = regexp.MustCompile(`\{[^}]+\}`)

func classify(method, path string) registry.Operation {
	best := registry.Operation{ID: "freshdesk.api.request", Method: method, Path: path, Effect: "write", Confirmation: true}
	if _, err := client.Path(path); err != nil {
		return best
	}
	bestScore := -1
	for _, op := range registry.Operations {
		if op.Method != method {
			continue
		}
		pattern := regexp.QuoteMeta(op.Path)
		for _, param := range op.PathParameters {
			pattern = strings.ReplaceAll(pattern, regexp.QuoteMeta("{"+param+"}"), `[^/]+`)
		}
		match, _ := regexp.MatchString("^"+pattern+"/?$", path)
		if match {
			score := len(placeholder.ReplaceAllString(op.Path, ""))
			if score > bestScore || score == bestScore && op.Effect == "write" {
				best = op
				bestScore = score
			}
		}
	}
	return best
}
