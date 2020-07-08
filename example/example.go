package main

import (
	"log"
	"math"
	"time"

	"github.com/youpy/caudio"
)

func main() {
	n := 0.0
	fn := func(t float64, stepCount int) float64 {
		log.Print(stepCount)
		x := math.Sin(t*262 + math.Sin(n))
		n += math.Sin(t)

		return x
	}

	osc := caudio.New(fn)
	osc.Start()

	time.Sleep(5000 * time.Millisecond)
}
