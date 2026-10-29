# Security

Report a suspected vulnerability privately through [GitHub private vulnerability reporting](https://github.com/piyush-gambhir/freshdesk-cli/security/advisories/new). Do not include API keys, customer data, or signed attachment URLs in public issues.

Keys use the OS keychain by default. Plaintext storage requires an explicit `auth login --insecure-storage`; that file uses mode 0600 on Unix. Windows relies on the user's filesystem ACLs. Environment-based authentication avoids credential persistence. Profile metadata is stored separately from secrets.

API credentials travel in an HTTPS Basic authorization header. API redirects are not followed. Raw requests stay on the configured helpdesk origin, and traversal or absolute paths are rejected. Verbose output logs method, path, and status rather than payloads or query values. Error and discovery output redact API keys.

Read-only controls cover remote writes, local profile changes, and self-update. A dry-run freezes and hashes its payload and attachment contents without sending a request. `--expect-plan` checks that rebuilt request before execution. Unknown raw routes are classified conservatively as writes. These controls do not grant authorization or replace an operator's approval process.

Writes are never retried automatically. Network failures, invalid success responses, and server errors can leave a mutation's outcome unknown. Inspect server state before deciding whether to repeat it.

Release archives include SHA-256 checksums and GitHub build provenance. Version 0.1.2 and later also include an SBOM for each archive and are immutable after publication. Checksums detect a mismatch against a release's checksum file; they do not independently establish who built that release. Verify provenance with:

```sh
gh attestation verify freshdesk-cli_darwin_arm64.tar.gz --repo piyush-gambhir/freshdesk-cli
```
