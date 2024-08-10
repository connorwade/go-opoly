package main

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Spaces [40]Space
}

type Space interface {
	getName() string
	Action(p *Player)
}

type NeutralSpace struct {
	Name string `json:"name,omitempty"`
}

func (ns *NeutralSpace) Action(p *Player) {}

func (ns *NeutralSpace) getName() string {
	return ns.Name
}

type TaxSpace struct {
	Name   string `json:"name,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

func (ts *TaxSpace) Action(p *Player) {
	p.payBank(ts.Amount)
}

func (ts *TaxSpace) getName() string {
	return ts.Name
}

type GoToJailSpace struct {
	Name string `json:"name,omitempty"`
}

func (g *GoToJailSpace) Action(p *Player) {
	p.Position = 10
	p.InJail = true
}

func (g *GoToJailSpace) getName() string {
	return g.Name
}

type Property struct {
	Name  string  `json:"name,omitempty"`
	Price int     `json:"price,omitempty"`
	Owner *Player `json:"owner"` // nil if unown,omitemptyed
	Rent  int     `json:"rent,omitempty"`
	Group string  `json:"group,omitempty"`
}

func (p *Property) Action(pl *Player) {
	if p.Group == "utility" {
		pl.payPlayer(p.Owner, 60) // adjust to dice roll
		return
	}
	if p.Owner == nil {
		pl.BuyProperty(p)
		return
	}
	// mortgage properties if necessary
	for pl.Money < p.Rent && len(pl.Properties) > 0 {
		prop := pl.PickPropertyToMortgage()
		pl.MortgageProperty(prop)
	}
	if pl.Money < p.Rent {
		pl.Bankrupt = true
		return
	}
	pl.payPlayer(p.Owner, p.Rent)
}

func (p *Property) getName() string {
	return p.Name
}

type CardDrawSpace struct {
	Name string `json:"name,omitempty"`
	CardDeck
}

// func pop(slice []Card) ([]Card, Card) {
// 	return slice[:len(slice)-1], slice[len(slice)-1]
// }

func (c *CardDeck) DrawCard() Card {
	card := c.Cards[rand.Int()%len(c.Cards)]
	return card
}

func (c *CardDrawSpace) Action(p *Player) {
	fmt.Printf("Deck: %v\n", c.CardDeck)
	fmt.Printf("%s drew a card\n", p.Name)
	card := c.DrawCard()
	card.Action(p)
}

func (c *CardDrawSpace) getName() string {
	return c.Name
}
