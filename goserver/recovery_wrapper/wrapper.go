package recoveryWrapper

import (
	"errors"
)

type Wrapper struct {
	Error error
}

var ErrPanic = errors.New("panic")

func (w *Wrapper) Clear() *Wrapper {
	w.Error = nil
	return w
}

func (w *Wrapper) Do(processor func() error) (wr *Wrapper) {
	wr = w
	if w.Error != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			w.Error = ErrPanic
		}
	}()

	w.Error = processor()
	return
}
