package generator

import (
	"context"
	"sync"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name       string
		start, end int
		want       []int
	}{
		{"simple range", 1, 5, []int{1, 2, 3, 4}},
		{"single value", 0, 1, []int{0}},
		{"empty when start equals end", 3, 3, nil},
		{"empty when start after end", 5, 2, nil},
		{"negative start", -2, 2, []int{-2, -1, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var got []int
			for v := range Count(ctx, tt.start, tt.end) {
				got = append(got, v)
			}
			if !equal(got, tt.want) {
				t.Errorf("Count(%d, %d) = %v, want %v", tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"non-empty", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"single", []string{"x"}, []string{"x"}},
		{"empty", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			var got []string
			for v := range FromSlice(ctx, tt.in) {
				got = append(got, v)
			}
			if !equal(got, tt.want) {
				t.Errorf("FromSlice(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	ctx := context.Background()
	src := Count(ctx, 1, 5) // 1,2,3,4
	var got []int
	for v := range Map(ctx, src, func(n int) int { return n * n }) {
		got = append(got, v)
	}
	want := []int{1, 4, 9, 16}
	if !equal(got, want) {
		t.Errorf("Map squares = %v, want %v", got, want)
	}
}

func TestMapChangesType(t *testing.T) {
	ctx := context.Background()
	src := FromSlice(ctx, []int{1, 22, 333})
	var got []int
	for n := range Map(ctx, src, func(n int) int { return len(itoa(n)) }) {
		got = append(got, n)
	}
	want := []int{1, 2, 3}
	if !equal(got, want) {
		t.Errorf("Map lengths = %v, want %v", got, want)
	}
}

func TestTake(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{"take fewer than available", 3, []int{0, 1, 2}},
		{"take zero", 0, nil},
		{"take negative", -1, nil},
		{"take more than available", 100, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Each Take gets its own cancellable source so the upstream
			// generator can exit when Take stops early — no goroutine leak.
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			src := Count(ctx, 0, 10)
			var got []int
			for v := range Take(ctx, src, tt.n) {
				got = append(got, v)
			}
			if !equal(got, tt.want) {
				t.Errorf("Take(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

// TestCancellationStopsProducer verifies that cancelling the context unblocks a
// producer that is blocked on send, so its goroutine does not leak. We
// synchronize entirely with channels — no time.Sleep.
func TestCancellationStopsProducer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// done is closed by the producer when it returns; if cancellation does
	// not unblock the send, this test deadlocks (caught by the test timeout)
	// and the producer goroutine would leak.
	done := make(chan struct{})
	out := Generate(ctx, func(yield func(int) bool) {
		defer close(done)
		// Infinite producer: only stops when yield reports cancellation.
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	})

	// Receive a couple of values to prove the stream is live...
	if v := <-out; v != 0 {
		t.Fatalf("first value = %d, want 0", v)
	}
	if v := <-out; v != 1 {
		t.Fatalf("second value = %d, want 1", v)
	}

	// ...then abandon it and cancel. The producer must observe cancellation
	// (either on its blocked send or its next yield) and return.
	cancel()

	// Drain any in-flight value so the producer's pending send completes,
	// allowing it to reach the cancellation branch and close out.
	for range out { //nolint:revive // intentionally draining to completion
	}
	<-done // blocks until the producer goroutine has returned
}

// TestConcurrentConsumersDrainExactly checks that many independent generators
// run concurrently and each emits its full, correct sequence under -race.
func TestConcurrentConsumersDrainExactly(t *testing.T) {
	const (
		gens   = 50
		perGen = 200
	)
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(gens)
	sums := make([]int, gens)
	for g := 0; g < gens; g++ {
		go func(idx int) {
			defer wg.Done()
			sum := 0
			for v := range Count(ctx, 0, perGen) {
				sum += v
			}
			sums[idx] = sum // each goroutine writes a distinct index: race-free
		}(g)
	}
	wg.Wait()

	want := perGen * (perGen - 1) / 2
	for g, got := range sums {
		if got != want {
			t.Errorf("generator %d sum = %d, want %d", g, got, want)
		}
	}
}

// TestPipelineComposition exercises Count -> Map -> Take composed together,
// the canonical generator-pipeline usage, with shared cancellation.
func TestPipelineComposition(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nums := Count(ctx, 1, 1_000) // large source
	doubled := Map(ctx, nums, func(n int) int { return n * 2 })
	first5 := Take(ctx, doubled, 5)

	var got []int
	for v := range first5 {
		got = append(got, v)
	}
	want := []int{2, 4, 6, 8, 10}
	if !equal(got, want) {
		t.Errorf("pipeline = %v, want %v", got, want)
	}
	// cancel() (deferred) lets the unconsumed Count/Map producers exit.
}

// --- helpers ---

func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
