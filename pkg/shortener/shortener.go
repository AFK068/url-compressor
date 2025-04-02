package shortener

import "strings"

type Shortener struct {
	Alphabet string
	Base     uint64
	Length   uint64
}

func NewShortener(alphabet string, length uint64) (*Shortener, error) {
	if length == 0 {
		return nil, ErrInvalidLength
	}

	if alphabet == "" {
		return nil, ErrInvalidLengthAlphabet
	}

	return &Shortener{
		Alphabet: alphabet,
		Base:     uint64(len(alphabet)),
		Length:   length,
	}, nil
}

func (s *Shortener) Encode(num uint64) (string, error) {
	maxValue := uint64(1)
	for i := uint64(0); i < s.Length; i++ {
		maxValue *= s.Base
	}

	if num >= maxValue {
		return "", ErrNumberOverflow
	}

	data := make([]byte, s.Length)
	for i := s.Length; i > 0; i-- {
		data[i-1] = s.Alphabet[num%s.Base]
		num /= s.Base
	}

	return string(data), nil
}

func (s *Shortener) Decode(str string) (uint64, error) {
	if uint64(len(str)) != s.Length {
		return 0, ErrInvalidDecoderLength
	}

	var num uint64

	for i := 0; i < len(str); i++ {
		index := strings.IndexByte(s.Alphabet, str[i])

		if index == -1 {
			return 0, ErrInvalidCharacter
		}

		num = num*s.Base + uint64(index) //nolint
	}

	return num, nil
}
