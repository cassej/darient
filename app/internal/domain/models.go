package domain
import "errors"

var (
	ErrInvalidInput = errors.New("invalid input parameters")
	ErrNotFound     = errors.New("resource not found")
    ErrAlreadyExists = errors.New("already exists")
    ErrForeignKey    = errors.New("foreign key violation")
)