# Migration Parity Spec (Phase 0)

## Scope for first parity release
- Implement and validate the currently active kingdom card set used in Rust init:
  - Treasure/Victory: Copper, Silver, Gold, Estate, Duchy, Province
  - Kingdom: Smithy, Village, Council Room, Chapel
- Preserve turn order semantics:
  - Optional second-player opening turn when player one is not first.
  - `action_phase -> buy_phase -> cleanup`
- Preserve game end:
  - Province pile empty OR 3 supply piles empty.
- Preserve result comparison:
  - Higher VP wins.
  - Tie in VP: player with fewer turns wins.
  - Equal VP and turns is tie.

## Strategy schema compatibility
- Strategy JSON format remains the same as current Rust files.
- Condition operators supported:
  - `True`, `EqualTo`, `NotEqualTo`, `GreaterThan`, `GreaterThanOrEqualTo`, `LessThan`, `LessThanOrEqualTo`
- Condition values supported:
  - `Int`, `CountInDeck`, `CountTypeInDeck`, `CountInSupply`, `CountAllCardsInDeck`, `CountVp`, `CountOpponentVp`, `CountCoin`

## Acceptance criteria
- Existing files in `strategies/` load without manual edits.
- Seeded simulation runs are reproducible.
- Engine can run >= 100k games between two strategies.
- Baseline parity harness can compare Rust vs Go aggregate win/tie/loss distribution.
