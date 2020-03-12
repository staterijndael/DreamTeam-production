package errors

import (
	"dt/logwrap"
	"github.com/semrush/zenrpc"
)

func New(n int, err error, data interface{}) *zenrpc.Error {
	if n == 1 && err != nil {
		logwrap.Debug("internal err: %s", err.Error())
	}

	return &zenrpc.Error{
		Code:    n,
		Message: errMap[n],
		Data:    data,
		Err:     err,
	}
}
