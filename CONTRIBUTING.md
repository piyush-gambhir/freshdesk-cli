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

Update `cli-go/VERSION`, run checks, and push a version tag matching that file, such as `v0.1.0`. The Release workflow tests the tagged source and publishes six platform archives, checksums, and build provenance. A source checkout can use `make install`; end users can use the checksummed installer or `freshdesk update`.

Keep tests focused on behavior, request construction, and safety boundaries. Document user-visible flag or behavior changes and regenerate `docs/` in the same change.
