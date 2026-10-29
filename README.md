# Freshdesk CLI

An independent Go command-line client for Freshdesk, built for people and coding agents. Not affiliated with Freshworks.

Named commands cover **266 documented API operation variants**, including ticketing, contacts, companies, administration, knowledge base translations, community forums, custom-object records, imports/exports, automations, SLAs, email, and collaboration. Table, JSON, YAML, and CSV output; named profiles; shell completion; one cross-platform binary.

[![CI](https://github.com/piyush-gambhir/freshdesk-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/piyush-gambhir/freshdesk-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/piyush-gambhir/freshdesk-cli)](https://github.com/piyush-gambhir/freshdesk-cli/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[Public documentation](https://projects.piyushgambhir.com/freshdesk-cli/docs) · [Releases](https://github.com/piyush-gambhir/freshdesk-cli/releases)

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/piyush-gambhir/freshdesk-cli/main/install.sh | sh
```

Installs `freshdesk` to `~/.local/bin`, verifying SHA-256 checksums first. Override with `INSTALL_DIR`, or pin a version with `VERSION=v0.1.3`. Prebuilt archives support macOS, Linux, and Windows on amd64 and arm64. Download Windows ZIPs from [Releases](https://github.com/piyush-gambhir/freshdesk-cli/releases).

From source, with Go 1.27.3 or later:

```sh
git clone https://github.com/piyush-gambhir/freshdesk-cli.git
cd freshdesk-cli
make build      # cli-go/bin/freshdesk
make install    # $(go env GOPATH)/bin/freshdesk
```

`freshdesk update --check` checks the latest GitHub release. `freshdesk update --yes` verifies its archive and replaces the binary. Source installs in Go bin directories print the source update command. There are no background update checks.

## Quick start

Get your API key from Freshdesk: profile picture > Profile settings > API key.

```sh
export FRESHDESK_DOMAIN=yourcompany
freshdesk auth login --profile support    # hidden API key prompt; verifies before saving

freshdesk tickets list
freshdesk tickets search --query 'status:2 AND priority:3' -o json
freshdesk tickets get 123 --include stats
freshdesk tickets conversations 123 --all -o json
freshdesk contacts search --query "email:'customer@example.com'"

# No saved profile is required for environment-based authentication.
export FRESHDESK_API_KEY='your-key'
freshdesk --read-only agents me -o json
```

A bare account slug becomes `https://yourcompany.freshdesk.com`. You may also provide an HTTPS hostname. Freshdesk documents using the account's Freshdesk domain, rather than a custom portal CNAME. Credentials are never forwarded through redirects.

## Writes and previews

Every API write requires confirmation. In scripts, pass `--yes`. Preview with `--dry-run` first: it needs a domain but no API key, does not contact Freshdesk, and works in read-only mode.

```sh
freshdesk tickets create \
  --subject 'Cannot sign in' \
  --description '<p>The account returns an error.</p>' \
  --email customer@example.com \
  --status open --priority high \
  --set 'tags=["login","urgent"]' \
  --set 'custom_fields.cf_team=Support' \
  --attach ./evidence.png \
  --dry-run -o json

# Use the same arguments and the preview's plan.plan_hash to bind execution.
freshdesk tickets update 123 --status resolved --yes
freshdesk tickets reply 123 --body-file ./reply.html --yes
freshdesk tickets note 123 --body 'Investigating' --set private=true --yes
freshdesk contacts update 42 --input-json @contact.json --dry-run -o json
freshdesk automations create 1 --input-json @rule.json --dry-run -o json
freshdesk contacts imports create --upload import_file=contacts.csv \
  --input-json @import-options.json --dry-run -o json
```

`--input-json` accepts an inline object, `@FILE`, or `-` for stdin. `--set FIELD=VALUE` accepts JSON values or plain text; dotted names build nested objects. Explicit convenience flags override those fields. Standard ticket status and priority names map to Freshdesk's numbers; account-specific statuses can be supplied numerically. `--body-file` sends its text as-is, including HTML; it does not convert Markdown.

`--attach FILE` uploads `attachments[]`. `--upload FIELD=FILE` supports named upload fields, such as import CSVs and avatars, using the field name in the endpoint reference. Nested JSON fields are converted to bracketed multipart fields. Uploads are limited to 20 MiB in total. Attachment bytes, filenames, target URL, payload, API registry, CLI version, and optional `--if-match` are bound into the plan hash. `--expect-plan HASH` rejects a different request before sending it. API-specific support for an If-Match precondition is not assumed.

`--read-only`, `FRESHDESK_READ_ONLY=1`, or a saved profile restriction blocks remote writes, local profile changes, and self-update. Environment restrictions cannot be undone with `--read-only=false` or `--no-input=false`. Cloning a ticket form is classified as a write even though Freshdesk documents it as GET.

No API requests are automatically retried. Throttling returns `retry_after`; an ambiguous mutation failure reports `outcome: unknown`. Bulk, import, and export responses can represent a queued job: inspect the returned ID using the corresponding get command.

## Commands

| Area | Command groups |
| --- | --- |
| Ticketing | `tickets`, `conversations`, `attachments`, `outbound-messages` |
| Customers | `contacts`, `contact-fields`, `companies`, `company-fields` |
| Administration | `agents`, `groups`, `admin-groups`, `skills`, `roles`, `ticket-fields`, `ticket-forms` |
| Knowledge base | `solutions categories`, `solutions folders`, `solutions articles`, including `translations` |
| Community | `discussions categories`, `discussions forums`, `discussions topics`, `discussions comments` |
| Custom objects | `custom-objects`, `custom-objects records` |
| Workflows | `automations`, `scenario-automations`, `sla-policies`, `business-hours`, `activities` |
| Support operations | `time-entries`, `service`, `csat`, `surveys`, `products` |
| Configuration | `email`, `email-configs`, `settings`, `account`, `jobs`, `audit-log` |
| Collaboration | `threads`, `threads messages` |
| Local tools | `auth`, `config`, `doctor`, `commands`, `api`, `completion`, `version`, `update` |

See the generated [command reference](docs/commands.md) and [API coverage map](docs/api-coverage.md) for exact routes, positional arguments, effects, and upstream discrepancies.

## Pagination and output

```sh
freshdesk tickets list --updated-since 2026-01-01 --all --max-pages 100 -o json
freshdesk contacts list --page 2 --per-page 100 -o csv > contacts-page-2.csv
freshdesk custom-objects records list 30 --param page_size=100 -o json
freshdesk api request GET /api/v2/agents/me -o json
```

`--param KEY=VALUE` supplies any native query parameter. Standard list commands support `--page`, `--per-page` (1-100), and bounded `--all`. New CSAT APIs use `--next-token`/`--prev-token` and can also collect pages with `--all`. Search uses 30 results per page and at most 10 pages. Ticket lists have a 300-page ceiling and default to tickets created in the last 30 days unless `--updated-since` is supplied. Custom-object pagination has a different protocol and remains manual through `--param` and native response metadata. Raw calls use manual pagination.

JSON and YAML API results use `{contract_version, command_id, data, meta}`. Single-page `data` preserves the vendor response. With `--all`, `data` is the combined result array. `meta.complete` describes pagination from the selected start page, `meta.next_page_url` gives the next page when one remains, and `meta.requests` reports requests made. Non-paginated endpoints use `complete: null`, without a completeness claim. A multi-page failure exits with no partial success output. CSV and table formats render result rows; use JSON for full context. Command data goes to stdout and diagnostics to stderr. Exit codes are 0 for success, 1 for failure.

## Profiles and credentials

```sh
freshdesk auth login --domain company-a --profile company-a
freshdesk auth login --domain company-b --profile company-b
freshdesk auth list
freshdesk auth use company-b
freshdesk auth status             # local; does not load the secret
freshdesk auth status --verify    # one GET /agents/me request
freshdesk auth logout             # removes only local credentials
```

API key precedence: `--api-key`, then `FRESHDESK_API_KEY`, then the selected saved profile. Domain precedence: `--domain`, then `FRESHDESK_DOMAIN`, then the profile. Profile precedence: `--profile`, then `FRESHDESK_PROFILE`, then the saved current profile. Environment credentials override a selected profile; unset them to use saved account credentials.

Profiles live in `~/.config/freshdesk-cli/config.yaml`, respecting `XDG_CONFIG_HOME`. `FRESHDESK_CONFIG` overrides the full config path. API keys are stored in the OS keychain by default. Without a keychain, login fails unless you explicitly choose `--insecure-storage` (a separate plaintext `secrets.yaml`, mode 0600 on Unix) or use an environment key. Windows uses the user's filesystem ACLs. Config writes are locked and atomic. `--save-read-only` persists a restriction in a new profile.

| Environment variable | Purpose |
| --- | --- |
| `FRESHDESK_DOMAIN` | Account slug or hostname |
| `FRESHDESK_API_KEY` | API key |
| `FRESHDESK_PROFILE` | Selected saved profile |
| `FRESHDESK_CONFIG` | Explicit config path |
| `FRESHDESK_READ_ONLY`, `FRESHDESK_NO_INPUT` | Sticky restrictions (`1` or `true`) |
| `FRESHDESK_VERBOSE`, `FRESHDESK_QUIET` | Diagnostic controls |

## Coding agents

The operational skill is [freshdesk/SKILL.md](freshdesk/SKILL.md).

```sh
freshdesk commands -o json
freshdesk commands describe freshdesk.tickets.create -o json
freshdesk commands resolve -- tickets reply 123 --body-file reply.html
freshdesk api methods -o json
```

Discovery and resolution work offline, without loading configuration or touching the keychain. The resolver uses the real command parser. Unknown raw routes are classified as writes, including GET. Credentials are redacted from resolved flags. Request previews and plan hashes support review before execution; authorization remains the caller's responsibility.

Resolver classifications include explicit `secret_output` and `interactive` booleans for harnesses such as Reckon. `secret_output` is false; `auth login` is marked interactive unless `--key-stdin` is selected.

## Coverage and verification limits

The pinned inventory comes from the [official REST reference](https://developers.freshdesk.com/api/), retrieved October 29, 2026. All 266 documented operation variants have named commands and fixture tests. This count includes translation and field-service variants that share routes. It does not imply that every Freshdesk UI feature has an API.

The initial v0.1.0 release used local fixtures. Subsequent [live checks](docs/live-verification.md) verified many read operations and identified compatibility fixes for admin ticket fields and CSAT pagination. An explicitly approved one-ticket test verified creation, priority updates, private-note creation/edit/deletion, attachment upload/deletion, and ticket soft-deletion, with ten successful mutations and thirteen read-back assertions. Other writes remain fixture-tested. Automated tests cover request construction, effect guards, pagination, uploads, authentication, errors, and release verification. Account entitlements and custom business rules still require tenant-specific verification. Conflicting route examples are listed in the [coverage report](docs/api-coverage.md#reference-discrepancies).

The reference exposes no mutation endpoints for roles, products, business hours, or helpdesk settings, and no custom-object schema-creation endpoint. Analytics dashboards, billing, portal styling, Freddy configuration, and FDK app development are outside this inventory. Attachment URLs remain in ticket/conversation responses for download with an HTTPS client.

## Development

```sh
make test       # race-enabled; no Freshdesk account or real keychain
make vet
make build
make docs       # regenerates command and coverage references
```

See [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md), and [AGENTS.md](AGENTS.md). MIT licensed.
