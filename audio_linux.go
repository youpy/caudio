package caudio

import "errors"

// Audio is responsible for producing the sound using a given callback
type Audio struct {
}

func _new(fn Callback) *Audio {
	return &Audio{}
}

func (o *Audio) _start() error {
	return errors.New("platform not supported: linux")
}
