package apperror

const (
	DuplicateEmailErrorText = "Duplicate email"
)

type DuplicateEmailError struct{}

func (*DuplicateEmailError) Error() string {
	return DuplicateEmailErrorText
}
