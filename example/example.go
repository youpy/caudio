package main

import (
	"math"
	"time"

	"github.com/youpy/caudio"
)

func main() {
	n := 0.0
	fn := func(t float64, stepCount int) float64 {
		x := math.Sin(t*262 + math.Sin(n))
		n += math.Sin(t)

		return x
	}

	osc := caudio.New(fn)
	osc.Start()

	time.Sleep(5000 * time.Millisecond)
}
