package innkeep

type innkeepError struct {
	htttpStatus int
	err         error
}

func (e innkeepError) StatusCode() int {
	if e.htttpStatus <= 0 {
		return 500
	} else {
		return e.htttpStatus
	}
}

func (e innkeepError) Error() string {
	return e.err.Error()
}

func NewError(err error) error {
	return innkeepError{err: err}
}

func NewStatusError(err error, status int) error {
	return innkeepError{
		err:         err,
		htttpStatus: status,
	}
}
