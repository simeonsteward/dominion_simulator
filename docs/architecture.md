# Architecture

Core gameplay is implemented in `internal/core`, strategy parsing/evaluation in `internal/strategy`, and simulation orchestration in `internal/sim`.

## Packages
- `internal/core`: cards, state, player turn flow, kingdom/supply, game loop.
- `internal/strategy`: strategy JSON schema and condition evaluators.
- `internal/io`: strategy and snapshot file I/O.
- `internal/sim`: multi-game execution, benchmarking, tournaments.
- `pkg/api`: stable facade for future UIs/services.
- `cmd/engine`: CLI simulation runner.
- `cmd/strategy-cli`: strategy create/validate/migrate tool.
- `cmd/webui`: lightweight HTTP UI backend.
