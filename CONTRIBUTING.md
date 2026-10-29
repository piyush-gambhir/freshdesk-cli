# Contributing

The Go/Cobra implementation is in `cli-go/`. Use Go 1.27.3 or later. Tests use fake transports, temporary config directories, and the keyring mock; they never need a live helpdesk.

```sh
make test
make vet
make build
make docs
```

Run targeted tests first when changing a subsystem. CI also runs staticcheck, govulncheck, cross-platform tests, a source install, installer fixtures, and GoReleaser snapshot packaging.

## API changes

`cli-go/internal/registry/operations.json` is the reviewed operation registry. It drives named commands, effects, raw-route classification, and generated coverage documentation. The separate source-facts snapshot is in `cli-go/internal/registry/testdata/freshdesk-api.json`.

To refresh the source inventory:

```sh
python3 scripts/snapshot_api.py --output .scratch/freshdesk-api.json
```

Review the diff before replacing the pinned snapshot. The snapshot keeps endpoint facts, not the source page or prose. The extractor includes nested operations and normalizes the reference's `update`/`clone` labels to their documented HTTP verbs. `scripts/build_registry.py` assigns commands and records route corrections. Its one-to-one mapping intentionally fails when the source inventory size changes: update the mapping and review effects before generating a new registry.

Preserve the source `anchor`, `method`, and `path` in the snapshot. Corrections belong in the registry builder with an explanation. The snapshot's `source_sha256` identifies the original downloaded HTML. Never infer read/write effects solely from HTTP verbs: ticket-form cloning uses GET and is a write.

After updating the reviewed mapping:

```sh
python3 scripts/build_registry.py
make docs test vet
```

JSON input stays open to native and account-specific fields. Do not advertise offline business-rule validation that is not implemented. Avoid automatic retries for writes, account-dependent guesses, or undocumented endpoint claims.

## Releases

Bump `cli-go/VERSION` in a pull request. Required CI and CodeQL checks must pass before merging to `main`. The Release workflow tags that commit, tests it, creates six platform archives and SBOMs, attests their checksums, and then publishes an immutable release. Do not push release tags by hand. Recover a failed release by rerunning its original workflow run; a published version is skipped. See [Maintaining the repository](docs/maintaining.md) for the repository safeguards and website deployment.

A source checkout can use `make install`; end users can use the checksummed installer or `freshdesk update`.

Keep tests focused on behavior, request construction, and safety boundaries. Document user-visible flag or behavior changes and regenerate `docs/` in the same change.

## Website

The public documentation site is in `web/`. Run `pnpm install --frozen-lockfile`, `pnpm types:check`, `pnpm build:cloudflare`, and `pnpm test:search` there. Generated reference pages come from `docs/`; update their source rather than the generated MDX.
