package main

import (
	"fmt"
	"io"
	"bufio"
	"strconv"
	"os"
	"errors"
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

func (c *Card) String() string {
	switch *c {
	case J, Q, K, A:
		s := []string{"J", "Q", "K", "A"}
		return s[*c - J]
	default:
		return strconv.Itoa(int(*c) + 2)
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

type Cards []Card

func (cards *Cards) pop(c int) Cards {
	res := []Card(*cards)[:c]
	*cards = []Card(*cards)[c:]
	
	return res
}

type Player struct {
	cards Cards
	name string
}

func (player *Player) String() string {
	return player.name
}

func NewPlayer (cards Cards, name string) Player {
	return Player {cards, name}
}

func (player *Player) MakeMove() Card {
	res := player.cards[0]
	player.cards = player.cards[1:]
	return res
}

func (player *Player) AddCards(cards []Card) {
	player.cards = append (player.cards, cards...)
}

func (player *Player) SetCards (cards Cards) {
	player.cards = cards
}

type War struct {
	A Player
	B Player
	pileA []Card
	pileB []Card
}

func (war *War) AddCards (cardsA []Card, cardsB []Card) {
	war.pileA, war.pileB = append (war.pileA, cardsA...), append (war.pileB, cardsB...) 
}

func (war *War) PurgeCards (player *Player) {
	player.AddCards (war.pileA)
	player.AddCards (war.pileB)

	war.pileA, war.pileB = make ([]Card, 0), make ([]Card, 0)
}

func (war *War) MakePile () error {
	if (len(war.A.cards) < 3 || len(war.B.cards) < 3) {
		return errors.New("PAT")
	} 
	war.AddCards (war.A.cards.pop(3), war.B.cards.pop(3))
	return nil
}

func (war *War) FightWar () (*Player, error) {
	err := war.MakePile()

	if (err != nil) {
		return nil, err
	}
	
	return war.MakeMove()
}

func (war *War) IsWar () bool {
	if (len(war.pileA) > 2) {
		return true
	} else {
		return false
	}
}

func (war *War) FightBattle (cardA Card, cardB Card) (*Player, error) {
	var p *Player
	if (cardA > cardB) {
		p = &war.A
	} else {
		p = &war.B
	}
	war.PurgeCards(p)
	return p, nil
}

func (war *War) GetWinner () *Player {
	var p *Player
	if (len(war.A.cards) > 0 && len (war.B.cards) == 0) {
		p = &war.A
	}
	if (len(war.A.cards) == 0 && len (war.B.cards) > 0) {
		p = &war.B
	}

	return p
} 

func (war *War) MakeMove() (*Player, error) {
	p := war.GetWinner()
	if (p != nil) {
		if (war.IsWar()) {
			war.PurgeCards(p)
			return p, errors.New("PAT")
		} else {
			return p, nil
		}
	}
	cardA, cardB := war.A.MakeMove(), war.B.MakeMove()

	war.AddCards([]Card{cardA}, []Card{cardB})
	
	var winner *Player
	var err error
	
	if (cardA == cardB) {
		winner, err = war.FightWar()
	} else {
		winner, err = war.FightBattle(cardA, cardB)
	}
	return winner, err
}

func readCards (scanner *bufio.Scanner) []Card {
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	res := make ([]Card, n)
	for i := 0; i < n && scanner.Scan(); i++ {
		res[i] = NewCard(scanner.Text())
	}

	return res
}

func readFile (r io.Reader) War {
	scanner := bufio.NewScanner (r)
	scanner.Split(bufio.ScanWords)
	
	return War {NewPlayer(readCards(scanner), "1"), NewPlayer(readCards(scanner), "2"), make([]Card, 0), make([]Card, 0)}
}

func main() {
	f, _ := os.Open ("input.txt")
	//f := os.Stdin
	war := readFile(bufio.NewReader (f))

	player := war.GetWinner()
	i := 0
	for player == nil {
		_, err := war.MakeMove()
		if (err != nil) {
			fmt.Println(err)
			return
		}
		player = war.GetWinner()
		i += 1
	}
	fmt.Println (player, i)
}
