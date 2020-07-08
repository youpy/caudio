package main

import (
	"math"

	"github.com/youpy/caudio"
)

func main() {
	fn := func(t float64, i int) float64 {
		return math.Sin(800 * t * math.Pi * math.Sin(float64((i/6000)%16)))
	}

	osc := caudio.New(fn)
	osc.Start()

	var wait chan bool
	<-wait
}
