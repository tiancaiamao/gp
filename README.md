# A simple goroutine pool 

## Rational

Goroutine is cheap, but not free. Especially when the goroutine trigger `runtime.morestack`, the cost become high.

This package is mainly aimed to handle that. Put back the stack-growed goroutine to a pool, and reuse that goroutine can eliminate the `runtime.morestack` cost.

See a blog post (Chinese) https://www.zenlife.tk/goroutine-pool.md

## Usage

Just replace your `go f()` call with `gp.Go(f)`:

```
import "github.com/tiancaiamao/gp"

var gP = gp.New(N, time.Duration)

gp.Go(func() {
	// ...
})
```
