package main

import "fmt"



type TurnData struct {
	CurrentPlayer    *Player
	DiceRoll         int
	Space            Space
	Transactions     []Transaction
	PropertyPurchase *Property
	PropertyMortgage *Property
}

type Transaction struct {
	From   *Player // nil if bank
	To     *Player // nil if bank
	Amount int
}

func (t *Transaction) Execute() {
	if t.From != nil {
		t.From.decreaseMoney(t.Amount)
	}
	if t.To != nil {
		t.To.increaseMoney(t.Amount)
	}
	toName := "bank"
	fromName := "bank"
	if t.To != nil {
		toName = t.To.Name
	}
	if t.From != nil {
		fromName = t.From.Name
	}

	fmt.Printf("%s paid %s $%d\n", fromName, toName, t.Amount)
}

type CardDeck struct {
	Cards []Card
}

type Card struct {
	Text   string
	Action func(p *Player)
}
