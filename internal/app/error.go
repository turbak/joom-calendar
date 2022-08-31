package app

type CodableError struct {
	StatusCode int   `json:"code"`
	Err        error `json:"error"`
}

func (e CodableError) Error() string {
	return e.Err.Error()
}

func (e CodableError) Code() int {
	return e.StatusCode
}
