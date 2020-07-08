package caudio

// New returns a new Audio
func New(fn Callback) *Audio {
	return _new(fn)
}

// Start starts making sound
func (o *Audio) Start() error {
	return o._start()
}
