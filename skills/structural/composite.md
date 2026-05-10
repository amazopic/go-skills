---
name: structural-composite
description: Use when you need to treat individual objects and compositions of objects uniformly. Classic for tree structures: file systems, UI widget hierarchies, scene graphs, org charts.
category: structural
go-version-min: "1.21"
sources:
  - go-old-pattern/go-patterns-1/SUMMARY.md (taxonomy)
  - go-old-pattern/go-patterns-2/Structural/Composite/
example: examples/structural/composite/
---

# Composite

## Intent

Composite lets clients treat individual objects (leaves) and groups of objects (composites/branches) through a single unified interface. The tree can be traversed without the caller knowing — or caring — whether it is talking to a leaf or a subtree. Operations on the tree are naturally recursive.

## Context

Whenever your domain has a part-whole hierarchy — where a "whole" contains zero or more "parts" that are themselves valid wholes — Composite is the right shape. File systems are the canonical example: a directory is a collection of files and directories, and both respond to `Print`, `Size`, or `Walk`. UI toolkits, scene graphs, expression trees, bill-of-materials, and org charts all exhibit the same pattern.

The key question: "Do leaves and composites need to satisfy the same interface?" If yes, Composite is appropriate. If leaves and composites have fundamentally different operations, keep them as separate types.

## Implementation in Go

`examples/structural/composite/` models a file-system tree. Both `Directory` (branch) and `File` (leaf) implement the `Component` interface. The composite (`Directory`) holds a slice of `Component`, so it can contain both files and subdirectories.

```go
// Component is the uniform interface for both leaves and branches.
type Component interface {
    Add(child Component)
    Name() string
    Child() []Component
    Print(prefix string) string
}

// Directory is a composite — it holds other Components.
type Directory struct {
    name   string
    childs []Component
}

func (d *Directory) Add(child Component) { d.childs = append(d.childs, child) }
func (d *Directory) Name() string        { return d.name }
func (d *Directory) Child() []Component  { return d.childs }

func (d *Directory) Print(prefix string) string {
    result := prefix + "/" + d.Name() + "\n"
    for _, c := range d.Child() {
        result += c.Print(prefix + "/" + d.Name()) // recursive
    }
    return result
}

// File is a leaf — Add is a no-op, Child returns empty.
type File struct{ name string }

func (f *File) Add(_ Component)      {}
func (f *File) Name() string         { return f.name }
func (f *File) Child() []Component   { return []Component{} }
func (f *File) Print(prefix string) string {
    return prefix + "/" + f.Name() + "\n"
}
```

### Go-specific consideration: the leaf `Add` no-op

In many Composite implementations the leaf's `Add` silently does nothing. That is acceptable for traversal-only use cases but can hide bugs. If callers should never call `Add` on a leaf, consider making `Add` a method only on the composite and narrowing the interface:

```go
// Narrower approach: split the interface
type Node interface {
    Name() string
    Print(prefix string) string
}

type Container interface {
    Node
    Add(child Node)
    Children() []Node
}
```

This is more idiomatic Go (prefer small interfaces) at the cost of losing the perfect uniformity the GoF shape provides. Choose based on whether callers must treat leaves and composites identically.

## When to use

- The domain is a recursive part-whole hierarchy.
- Client code should not distinguish between individual items and collections of items.
- You want operations (size, print, search) to propagate automatically through the tree without special-casing.
- Adding new leaf or composite types should not require changes to client traversal code.

## When NOT to use

- The hierarchy is shallow and fixed — a simple slice or map is clearer.
- Leaves and composites have meaningfully different operations; forcing them into one interface produces awkward no-ops.
- You need fine-grained control over traversal order with complex context — consider a separate Visitor instead of recursion baked into the components.
- Cycles are possible in the graph — Composite's recursive traversal will hang. Use an explicit visited set or choose a different traversal strategy.

## Gotchas

- **Unbounded recursion.** Deep trees can exhaust the call stack. For trees deeper than a few thousand levels, replace recursion with an explicit stack (`[]Component`).
- **Silently ignoring `Add` on leaves.** The no-op is convenient but masks programming errors. At minimum log a warning; in strict APIs return an error.
- **Shared children.** If the same `Component` is added to two parents, mutations (rename, delete) will affect both subtrees. Enforce single-parent ownership or make components immutable.
- **Nil children.** Guard `Add` against nil arguments; a nil element in the child slice causes a panic on the next `Print` call.
- **Thread safety.** Concurrent `Add` and `Print` on the same `Directory` is a data race. Protect the slice with a `sync.RWMutex` if the tree is modified after construction.

## See also

- skills/structural/decorator.md — wraps a single object; Composite aggregates many
- skills/structural/proxy.md — delegates to one object; Composite delegates to N
- examples/structural/composite/
