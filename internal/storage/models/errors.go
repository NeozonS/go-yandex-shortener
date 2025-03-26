package models

type ErrURLConflict struct {
	ExistingToken string
}

func (e ErrURLConflict) Error() string {
	return "URL already exists"
}
