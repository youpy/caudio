package caudio

import "C"
import (
	"sync"
)

//export baudio_callback
func baudio_callback(index C.int, t C.float) C.float {
	fn := lookup(int(index))

	return C.float(fn(float64(t)))
}

type Callback func(float64) float64

// based on https://github.com/golang/go/wiki/cgo#function-variables

var mu sync.Mutex
var index int
var fns = make(map[int]Callback)

func register(fn Callback) int {
	mu.Lock()
	defer mu.Unlock()
	index++
	for fns[index] != nil {
		index++
	}
	fns[index] = fn
	return index
}

func lookup(i int) Callback {
	mu.Lock()
	defer mu.Unlock()
	return fns[i]
}

func unregister(i int) {
	mu.Lock()
	defer mu.Unlock()
	delete(fns, i)
}
