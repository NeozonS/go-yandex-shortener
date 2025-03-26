package models

type ErrURLConflict struct {
	ExistingURL string
	Err         error
}

func (e *ErrURLConflict) Error() string {
	return "URL already exists"
}
