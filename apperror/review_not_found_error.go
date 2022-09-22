package apperror

type ReviewNotFoundError struct{}

func (_ *ReviewNotFoundError) Error() string {
	return "review not found"
}
