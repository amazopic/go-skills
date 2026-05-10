---
name: concurrency-parallelism
description: Parallel scatter-gather — fan a workload across goroutines and collect results via WaitGroup or channel. Use for embarrassingly parallel CPU-bound work.
category: concurrency
go-version-min: "1.18"
sources:
  - go-old-pattern/go-patterns-1/concurrency/parallelism.md
example: examples/concurrency/parallelism/
---

# Parallelism Pattern

[Parallelism](https://blog.golang.org/pipelines#TOC_8.) allows multiple "jobs" or tasks to be run concurrently and asynchronously.

# Implementation and Example

An example showing implementation and usage can be found in [parallelism.go](parallelism.go).
