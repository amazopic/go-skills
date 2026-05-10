# Known issues — examples

The following packages have failing tests inherited from upstream go-patterns-2 archive. They are tracked here; bodies will be reviewed in a future iteration.

| Package | Failure |
|---|---|
| `idiom/specification` | `TestSpecification`: `InCollectionSpecification.IsSatisfiedBy` uses `!elm.IsSent` (semantics inverted — field `IsSent: false` means not yet sent to collection, but the implementation treats it as "already in collection"); expected `true`, got `false`. Test error message is also inverted. Upstream bug — non-trivial semantic fix required. |
