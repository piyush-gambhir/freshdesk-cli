# Freshdesk CLI agent instructions

- This repository is a Go/Cobra CLI in `cli-go/`, module `github.com/piyush-gambhir/freshdesk-cli/cli-go`.
- Follow the sibling CLI suite's stdout/stderr, profiles, output formats, completion, and checksummed release conventions.
- Every API operation comes from `internal/registry/operations.json`. Effects are explicit. GET ticket-form cloning writes state. Unknown raw routes are writes.
- Build previews and execution through the same request builder. Preserve attachment-content hashing and plan mismatch checks.
- Read-only/no-input environment settings are sticky. Offline discovery must not load config, access the keychain, use the network, or write files.
- Default credential persistence is the keychain. Plaintext secrets require explicit opt-in. Never include credentials in logs, status, errors, previews, or resolved flags.
- Test with fake transports and temporary config directories. Use `keyring.MockInit()` for credential tests, never a real keychain or helpdesk.
- Run focused tests, then `make test vet build docs`. Generated docs must be committed with registry or flag changes.
- Preserve native JSON fields and integer precision. Report incomplete pages and ambiguous mutation outcomes accurately.
- Keep verification claims precise: `docs/live-verification.md` records read checks and the approved one-ticket lifecycle test. Other writes remain fixture-tested. Keep upstream discrepancies visible in `docs/api-coverage.md`.
- No em dashes. No deployment, production config changes, or edits to sibling repositories as part of ordinary work here.

- The public docs site is in `web/`, using the approved CLI-suite template and coral accent. Follow `web/AGENTS.md`; run type checks, the static build, search checks, and rendered desktop/mobile checks for website changes.
- Release by changing `cli-go/VERSION` in a pull request. Merging to protected `main` publishes the version automatically. Never push release tags manually.
