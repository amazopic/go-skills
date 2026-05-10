---
name: concurrency-bounded-parallelism
description: Bounded parallelism — process N items with at most K workers concurrently using a worker-pool goroutine + channels. Use to cap concurrency on an unbounded input list.
category: concurrency
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/concurrency/bounded_parallelism.md
example: examples/concurrency/bounded-parallelism/
---

# Bounded Parallelism Pattern

[Bounded parallelism](https://blog.golang.org/pipelines#TOC_9.) is similar to [parallelism](parallelism.md), but allows limits to be placed on allocation.

# Implementation and Example

An example showing implementation and usage can be found in [bounded_parallelism.go](bounded_parallelism.go).
