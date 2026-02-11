package ucerror

type UCError struct {
	details error // error's details, used for internal log to debug and trace error
	public  error // internal wrapped error, used for external display
}

func New(details error, err error) *UCError {
	return &UCError{details: details, public: err}
}

func (e *UCError) Error() string {
	return e.public.Error()
}

func (e *UCError) Unwrap() error {
	return e.details
}

func (e *UCError) Details() error {
	return e.details
}

func (e *UCError) Public() error {
	return e.public
}
