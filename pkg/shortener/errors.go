package shortener

import "errors"

var (
	ErrInvalidLength         = errors.New("length must be greater than 0")
	ErrInvalidDecoderLength  = errors.New("length must be equal to the specified value")
	ErrInvalidLengthAlphabet = errors.New("alphabet must not be empty")
	ErrInvalidStringLength   = errors.New("invalid string length")
	ErrInvalidCharacter      = errors.New("invalid character in string")
	ErrNumberOverflow        = errors.New("number overflow, too large to encode")
)
