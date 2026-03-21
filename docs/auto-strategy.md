# Auto Strategy Pipeline v1

Implemented baseline evolutionary scaffolding in `internal/autostrategy`:
- initialize candidate population from baseline priorities
- evaluate via short simulations
- rank by score
- keep top half and clone loop hook

Next upgrades:
- real mutation operators for buy/action/trash priorities
- kingdom-specific and cross-kingdom fitness
- result persistence and resumable runs
