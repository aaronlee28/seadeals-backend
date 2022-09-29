package apperror

type VoucherNotFoundError struct{}

func (_ *VoucherNotFoundError) Error() string {
	return "voucher not found"
}
