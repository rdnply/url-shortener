package baseconv

import (
	"errors"
	"math"
)

const Charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type BaseConv struct {
	base        uint
	alphabetMap map[rune]uint
}

func NewBaseConv(base uint) (*BaseConv, error) {
	if base == 0 || base > 62 {
		return nil, errors.New("invalid base")
	}

	runes := []rune(Charset[0:base])
	runeMap := make(map[rune]uint)

	var i uint
	for i = 0; i < base; i++ {
		runeMap[runes[i]] = i
	}

	return &BaseConv{
		base:        base,
		alphabetMap: runeMap,
	}, nil
}

func (c *BaseConv) Encode(number uint) string {
	remainders := make([]rune, 0)

	n := number
	for n > 0 {
		r := rune(Charset[n%c.base])
		remainders = append(remainders, r)
		n /= c.base
	}

	return reverse(string(remainders))
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func (c *BaseConv) Decode(code string) uint {
	var value float64
	value = 0

	for pos, r := range reverse(code) {
		value += float64(c.alphabetMap[r]) * math.Pow(float64(c.base), float64(pos))
	}

	return uint(value)
}
