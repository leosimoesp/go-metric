package errorwrapper

type ErrorWrapper struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
	Err     error  `json:"-"`
}

func NewErrorWrapper(message string, code int, err error) error {
	return ErrorWrapper{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (err ErrorWrapper) Error() string {
	if err.Err == nil {
		return ""
	}
	return err.Err.Error()
}

func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func (err ErrorWrapper) Unwrap() error {
	return err.Err
}
