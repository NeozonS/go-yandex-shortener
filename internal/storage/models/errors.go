package models

type ErrURLConflict struct {
	ExistingURL string
	Err         error
}

func (e *ErrURLConflict) Unwrap() error {
	return e.Err
}

func (e *ErrURLConflict) Error() string {
	return "URL already exists"
}
