package sms

func IsApiErr(err error) bool {
	_, ok1 := err.(ApiErr)
	_, ok2 := err.(*ApiErr)
	return ok1 || ok2
}

func IsManagerErr(err error) bool {
	_, ok1 := err.(ManagerErr)
	_, ok2 := err.(*ManagerErr)
	return ok1 || ok2
}

type ManagerErr struct {
	*ManagerErr
	err error
}

func (err ManagerErr) Error() string {
	msg := "<nil>"
	if err.err != nil {
		msg = err.err.Error()
	}

	return msg
}

type ApiErr struct {
	*ApiErr
	err error
}

func (err ApiErr) Error() string {
	msg := "<nil>"
	if err.err != nil {
		msg = err.err.Error()
	}

	return msg
}
