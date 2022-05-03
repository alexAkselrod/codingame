package card

import (
	"strconv"
)

type Card int

const (
	Two Card = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	J
	Q
	K
	A
)

func (c Card) String() string {
	switch c {
	case J, Q, K, A:
		return "image " + strconv.Itoa(int(c))
	default:
		return "digit " + strconv.Itoa(int(c))
	}
}

func NewCard (s string) Card {
	s = s[:len(s) - 1]
	switch s {
	case "J":
		return J
	case "Q":
		return Q
	case "K":
		return K
	case "A":
		return A
	default:
		d, _ := strconv.Atoi(s)
		return Card(d - 2)
	}
}
