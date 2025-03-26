package models

type ErrURLConflict struct {
	ExistingURL string
}

func (e *ErrURLConflict) Error() string {
	return "URL already exists"
}
func (e *ErrURLConflict) Unwrap() error {
	return nil // или возвращайте вложенную ошибку, если есть
}
