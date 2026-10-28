---
name: freshdesk
description: Read and manage a Freshdesk helpdesk through the freshdesk CLI, including tickets, contacts, administration, knowledge base, and workflows. Use for Freshdesk operations, not Freshworks FDK app development.
---

# Freshdesk

Use the `freshdesk` binary. It is an independent REST API client. Native response fields and account-specific business rules come from Freshdesk; the CLI's documented operation inventory is not a promise that the account enables every feature.

Discover capabilities offline:

```sh
freshdesk commands -o json
freshdesk commands describe freshdesk.tickets.reply -o json
freshdesk commands resolve -- tickets reply 123 --body-file reply.html
freshdesk api methods -o json
```

Discovery uses the real command parser without reading config or accessing credentials. Prefer stable command IDs and explicit effect metadata over guesses from verb names. In particular, ticket-form cloning is a write over GET; unknown raw routes are also writes.

For reads, use `--read-only --no-input -o json`. `auth status` inspects the active domain/profile and credential source without loading the key; `auth status --verify` makes one API call. Keys resolve from `--api-key`, then `FRESHDESK_API_KEY`, then the saved profile. Domain and profile overrides can change the target, so check the resolved account against the user's request.

API JSON results contain `data` and `meta`. Single-page data preserves the vendor response; `--all` returns a combined array. Inspect `meta.complete`, `meta.next_page_url`, `meta.start_page`, and warnings before treating results as complete. Ticket lists default to recent tickets; use `--updated-since` for older records. Search has a 300-result cap. Custom-object cursor/page-size handling is manual through native `--param` fields.

For an authorized write, build the request with native JSON or documented flags:

```sh
freshdesk tickets update 123 --input-json @update.json --dry-run -o json
freshdesk tickets update 123 --input-json @update.json --expect-plan HASH --yes --no-input -o json
```

The dry-run needs a domain but no key and sends nothing. Review the target, request body, attachment metadata, and plan hash. Preserve the same input and hash for execution. If input changes, generate a new plan. `--body-file` sends text or HTML verbatim, without Markdown conversion. `--set custom_fields.cf_team=Support` builds nested body fields; use `--param KEY=VALUE` for query parameters.

Writes require `--yes` in non-interactive runs. Sticky read-only/no-input restrictions cannot be removed by false-valued flags. There are no mutation retries. An `outcome: unknown` requires inspecting the remote state before considering a retry. A job or export/import ID means queued work: check it with the corresponding get command.

Read [the command reference](../docs/commands.md) for exact positional arguments and flags, or [the coverage report](../docs/api-coverage.md) when assessing unsupported features or conflicting upstream route examples. The initial release has fixture coverage without live tenant verification. API keys, signed URLs, and customer records should not be copied into public reports.
