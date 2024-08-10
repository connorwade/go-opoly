package main

import "fmt"

type Player struct {
	Name            string
	Money           int
	Position        int
	PositionHistory []int
	InJail          bool
	JailTurns       int
	Properties      []*Property
	Bankrupt        bool
}

func (p *Player) Move(spaces int) int {
	if p.Position+spaces > 39 {
		p.Position = p.Position + spaces - MAX_SPACES
		p.recieveMoneyFromBank(MONEY_FOR_PASSING_GO)
	} else {
		p.Position += spaces
		p.PositionHistory = append(p.PositionHistory, p.Position)
	}

	return p.Position
}

func (p *Player) payPlayer(p2 *Player, m int) {
	t := Transaction{
		From:   p,
		To:     p2,
		Amount: m,
	}

	t.Execute()
}

func (p *Player) recieveMoneyFromBank(m int) {
	t := Transaction{
		From:   nil,
		To:     p,
		Amount: m,
	}

	t.Execute()
}

func (p *Player) payBank(m int) {
	t := Transaction{
		From:   p,
		To:     nil,
		Amount: m,
	}

	t.Execute()
}

func (p *Player) increaseMoney(m int) {
	p.Money += m
}

func (p *Player) decreaseMoney(m int) {
	p.Money -= m
}

func (p *Player) PickPropertyToMortgage() *Property {
	minRent := 10000
	var minProp *Property
	for _, prop := range p.Properties {
		if prop.Rent < minRent {
			minRent = prop.Rent
			minProp = prop
		}
	}
	return minProp
}

func (p *Player) MortgageProperty(prop *Property) {
	fmt.Println(p.Name + " mortgaged " + prop.Name)
	p.recieveMoneyFromBank(prop.Price)
	prop.Owner = nil
}

func (p *Player) BuyProperty(prop *Property) {
	fmt.Println(p.Name + " bought " + prop.Name)
	p.payBank(prop.Price)
	prop.Owner = p
}
