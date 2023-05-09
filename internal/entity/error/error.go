package error

import "errors"

type Error struct {
	Code    int
	Message string
	Error   error
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message, Error: errors.New(message)}
}
func ExposeError(err error, es ...*Error) *Error {
	for _, e := range es {
		if errors.Is(err, e.Error) {
			return e
		}
	}
	return Unknown
}

var (
	Unknown                         = NewError(1, "Unknown error")
	FieldRequired                   = NewError(2, "Field required")
	UserAccountNotFound             = NewError(3, "User account not found")
	UserAccountBanned               = NewError(13, "User account not found")
	PermissionDenied                = NewError(4, "Permission denied")
	PaymentMethodNotSupported       = NewError(5, "Payment method not supported")
	PaymentNotFound                 = NewError(6, "Payment not found")
	PaymentTransactionNotFound      = NewError(7, "Payment transaction not found")
	PaymentDone                     = NewError(8, "Payment is done")
	PaymentAmountNotMatch           = NewError(9, "Payment amount not match")
	PaymentAmountNotReachedMinimum  = NewError(10, "Payment amount has not been reached minimum")
	UserAccountWrongPassword        = NewError(11, "Wrong password")
	UserAccountSaveFailed           = NewError(12, "Account save failed")
	UserAccountResetPasswordExpired = NewError(14, "Reset password request has expired")
	EmailInvalidate                 = NewError(15, "Email invalidate")
	PhoneNumberInvalidate           = NewError(16, "Phone number invalidate")
	EmailAlreadyUse                 = NewError(17, "Email already in use")
	PhoneNumberAlreadyUse           = NewError(18, "Phone number already in use")
	NameAlreadyUse                  = NewError(19, "Name already in use")
	SomethingWrong                  = NewError(20, "Something wrong")
)
