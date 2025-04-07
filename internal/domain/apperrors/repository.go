package apperrors

type ErrRepositoryIsFull struct {
	Message string
}

func (e *ErrRepositoryIsFull) Error() string { return e.Message }

type ErrURLNotFound struct {
	Message string
}

func (e *ErrURLNotFound) Error() string { return e.Message }
