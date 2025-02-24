package errors

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrNotFound   = &DomainError{Code: "NOT_FOUND", Message: "Entity not found"}
	ErrValidation = &DomainError{Code: "VALIDATION_ERROR", Message: "Validation error"}
	ErrRepository = &DomainError{Code: "REPOSITORY_ERROR", Message: "Repository error"}
)
