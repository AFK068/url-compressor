package shortener_test

import (
	"testing"

	"github.com/AFK068/compressor/pkg/shortener"
	"github.com/gookit/goutil/testutil/assert"
)

func Test_NewShortener_InvalidLength_Failure(t *testing.T) {
	_, err := shortener.NewShortener("abc", 0)
	assert.Error(t, err)
}

func Test_NewShortener_InvalidLengthAlphabet_Failure(t *testing.T) {
	_, err := shortener.NewShortener("", 5)
	assert.Error(t, err)
}

func Test_NewShortener_Success(t *testing.T) {
	shortener, err := shortener.NewShortener("abc", 5)
	assert.Nil(t, err)

	assert.Equal(t, "abc", shortener.Alphabet)
	assert.Equal(t, uint64(3), shortener.Base)
	assert.Equal(t, uint64(5), shortener.Length)
}

func Test_Encode_InvalidNumberOverflow_Failure(t *testing.T) {
	shortener, _ := shortener.NewShortener("abc", 5)

	_, err := shortener.Encode(444)
	assert.Error(t, err)
}

func Test_Encode_Success(t *testing.T) {
	shortener, _ := shortener.NewShortener("abc", 5)

	result, err := shortener.Encode(4)
	assert.Nil(t, err)

	assert.Equal(t, "aaabb", result)
}

func Test_Decode_InvalidStringLength_Failure(t *testing.T) {
	shortener, _ := shortener.NewShortener("abc", 5)

	_, err := shortener.Decode("aaaaaa")
	assert.Error(t, err)
}

func Test_Decode_InvalidCharacter_Failure(t *testing.T) {
	shortener, _ := shortener.NewShortener("abc", 5)

	_, err := shortener.Decode("aaaaa1")
	assert.Error(t, err)
}

func Test_Decode_Success(t *testing.T) {
	shortener, _ := shortener.NewShortener("abc", 5)

	result, err := shortener.Decode("aaabb")
	assert.Nil(t, err)

	assert.Equal(t, uint64(4), result)
}
