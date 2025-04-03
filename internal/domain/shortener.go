package domain

type Shortener interface {
	Encode(num uint64) (string, error)
	Decode(str string) (uint64, error)
}
