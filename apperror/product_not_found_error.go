package apperror

type ProductNotFoundError struct{}

func (_ *ProductNotFoundError) Error() string {
	return "product not found"
}
