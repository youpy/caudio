package caudio

// Callback is a function that the user registers with the `Audio` instance to make a sound
type Callback func(t float64, i int) float64
