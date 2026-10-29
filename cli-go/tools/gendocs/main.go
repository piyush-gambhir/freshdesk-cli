package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/cmd"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/registry"
	"github.com/spf13/cobra"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "coverage" {
		coverage()
		return
	}
	root := cmd.NewRoot(strings.NewReader(""), io.Discard, io.Discard)
	fmt.Println("# Freshdesk command reference\n\nGenerated with `make docs`. Resource IDs appear as positional arguments in route order.\nUse `freshdesk COMMAND --help` for individual flags.")
	fmt.Printf("\n## Global flags\n\n```text\n%s```\n", root.PersistentFlags().FlagUsages())
	for _, args := range [][]string{{"tickets", "list"}, {"tickets", "create"}} {
		c, _, _ := root.Find(args)
		fmt.Printf("\n## Shared flags: %s\n\nAll corresponding API commands accept these flags. Fields that are not relevant to an endpoint should be supplied only when its reference documents them.\n\n```text\n%s```\n", strings.Join(args, " "), c.LocalNonPersistentFlags().FlagUsages())
	}
	fmt.Println("\n## API commands\n\nJSON input accepts the endpoint's native fields, including custom fields and nested structures.\n`--param KEY=VALUE` supplies query parameters not listed as shortcuts.\nEvery write supports `--dry-run`, `--expect-plan HASH`, and `--yes`.\n\n| Command | HTTP | Effect | Automatic pagination |\n| --- | --- | --- | --- |")
	for _, op := range registry.Operations {
		c, _, _ := root.Find(op.Command)
		fmt.Printf("| `%s` | `%s %s` | %s | %t |\n", c.UseLine(), op.Method, op.Path, op.Effect, op.Paginated)
	}
	fmt.Println("\n## Local and discovery commands")
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		if c.Runnable() && c.Annotations["operation"] == "" && c.Annotations["effect"] != "" {
			fmt.Printf("\n### %s\n\n%s\n\n```text\n%s\n```\n", c.CommandPath(), c.Short, c.UseLine())
			if c.Long != "" {
				fmt.Printf("\n%s\n", c.Long)
			}
			if c.LocalNonPersistentFlags().HasAvailableFlags() {
				fmt.Printf("\n```text\n%s```\n", c.LocalNonPersistentFlags().FlagUsages())
			}
		}
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(root)
}

func coverage() {
	fmt.Printf("# Freshdesk API coverage\n\nThe reviewed registry covers **%d documented operation variants**, including translations and field-service variants that reuse an endpoint. This is a documentation coverage count, not a count of unique routes or live-verified features.\n\nSource: [Freshdesk API reference](https://developers.freshdesk.com/api/), retrieved 2026-10-29.\nThe endpoint-facts snapshot is in `cli-go/internal/registry/testdata/freshdesk-api.json`.\nRegistry SHA-256: `%s`.\n\n", len(registry.Operations), registry.Hash())
	fmt.Println("All rows have named commands and pass request-routing fixture tests. See [live verification notes](live-verification.md) for read checks and the approved one-ticket lifecycle test. Other writes remain fixture-tested; plan entitlements and agent permissions still apply.\n\n| Command | Verb | Path | Effect | Reference |\n| --- | --- | --- | --- | --- |")
	for _, op := range registry.Operations {
		fmt.Printf("| `%s` | %s | `%s` | %s | [reference](https://developers.freshdesk.com/api/#%s) |\n", strings.Join(op.Command, " "), op.Method, op.Path, op.Effect, op.SourceAnchor)
	}
	fmt.Println("\n## Reference discrepancies\n\nThese are deliberate choices from the public reference. They require tenant verification before claiming live compatibility.")
	for _, op := range registry.Operations {
		if op.Note != "" {
			fmt.Printf("\n- `%s`: %s\n", strings.Join(op.Command, " "), op.Note)
		}
	}
	fmt.Println("\n## Boundaries\n\n- The reference documents creating a custom-object schema in the admin UI, but provides no schema-creation REST endpoint. Record CRUD, schema reads, and record filtering are included.\n- Roles, products, business hours, and helpdesk settings expose read-only operations in this reference.\n- Analytics dashboards, portal themes, billing, Freddy configuration, and app deployment are outside this REST API inventory. FDK remains the app-development tool.\n- Custom-object APIs use their own page-size/cursor protocol. Use `--param` and the returned native metadata; `--all` is not offered for these operations.\n- The new CSAT APIs use `/customer-satisfaction`; the legacy survey APIs remain separate. Some contact/company search APIs are marked beta upstream.\n- JSON and `--set` can express endpoint payloads; the server validates account-specific fields and business rules. There is no claim of complete offline validation for every payload.\n- Bulk and export/import requests can return asynchronous jobs. Inspect the returned job/export/import ID with the corresponding get command.\n- Attachment upload is included through `--attach` and `--upload`. Download URLs are native fields in ticket/conversation responses; download them with an HTTPS client. The CLI does not save remote attachment files automatically.")
}
