package models

type ErrURLConflict struct {
	ExistingURL string
}

func (e ErrURLConflict) Error() string {
	return "URL already exists"
}
