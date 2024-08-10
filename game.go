package main

import (
	"fmt"
	"math/rand"
)

type Game struct {
	Players       []*Player
	Board         Board
	TurnData      []TurnData
	Chest         *CardDeck
	Chance        *CardDeck
	CurrentPlayer *Player
	PlayerOrder   []*Player
	Finished      bool
}

func rollDie() int {
	return rand.Int()%6 + 1
}

func rollDice() (int, bool) {
	r1 := rollDie()
	r2 := rollDie()
	d := r1 == r2
	return r1 + r2, d
}

func (g *Game) Turn(p *Player) {
	fmt.Printf("%s's turn; they have $%d\n", p.Name, p.Money)
	if p.Money < 0 {
		fmt.Println("Player is bankrupt")
		p.Bankrupt = true
		g.RemovePlayer(p)
	}
	if g.Finished {
		return
	}
	if !p.InJail {
		p.JailTurns++
		if p.JailTurns == MAX_TURNS_IN_JAIL {
			p.payBank(BAIL)
			p.InJail = false
			p.JailTurns = 0
			fmt.Println("Player has been in jail for 3 turns and must pay $50 to get out")
			g.Turn(p)
			return
		}

		// 50% chance the player rolls or pays
		if rand.Int()%2 == 0 {
			_, doubles := rollDice()
			if doubles {
				p.InJail = false
				p.JailTurns = 0
				fmt.Println("Player rolled doubles and got out of jail")
				g.Turn(p)
				return
			}
		} else {
			p.payBank(BAIL)
			p.InJail = false
			p.JailTurns = 0
			g.Turn(p)
			return
		}

		// fmt.Println("You are in jail")
		// if p.Money >= 50 {
		// 	fmt.Println("Would you like to pay $50 to get out of jail? (y/n)")
		// 	var input string
		// 	fmt.Scanln(&input)
		// 	if input == "y" {
		// 		p.payBank(50)
		// 		p.InJail = false
		// 	}
		// } else {
		// 	fmt.Println("You do not have enough money to pay the $50 fine")
		// }
	}
	diceRoll, doubles := rollDice()
	pos := p.Move(diceRoll)
	space := g.Board.Spaces[pos]
	turnData := TurnData{
		CurrentPlayer: p,
		DiceRoll:      diceRoll,
		Space:         space,
	}
	fmt.Printf("%s rolled %d and landed on %s\n", p.Name, diceRoll, space.getName())
	space.Action(p)
	for _, p := range g.Players {
		if p.Bankrupt {
			g.RemovePlayer(p)
		}
	}

	g.TurnData = append(g.TurnData, turnData)

	if g.Finished {
		return
	}

	if !doubles {
		g.CurrentPlayer = g.GetNextPlayer()
	}
	g.Turn(g.CurrentPlayer)
}

func (g *Game) GetNextPlayer() *Player {
	for i, p := range g.PlayerOrder {
		if p == g.CurrentPlayer {
			if i == len(g.PlayerOrder)-1 {
				return g.PlayerOrder[0]
			}
			return g.PlayerOrder[i+1]
		}
	}
	return nil
}

func (g *Game) RemovePlayer(p *Player) {
	for i, player := range g.Players {
		if player == p {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
		}
	}
	if len(g.Players) == 1 {
		g.EndGame(g.Players[0])
	}
}

func (g *Game) EndGame(p *Player) {
	fmt.Println(p.Name + " wins!")
	fmt.Println(g.TurnData)
	g.Finished = true
}

func (g *Game) Start() {

	for _, p := range g.Players {
		p.recieveMoneyFromBank(1500)
	}

	g.Chest = &CardDeck{
		Cards: []Card{
			{
				Text: "Advance to Go",
				Action: func(p *Player) {
					p.Position = 0
					p.recieveMoneyFromBank(200)
				},
			},
			{
				Text: "Bank error in your favor",
				Action: func(p *Player) {
					p.recieveMoneyFromBank(200)
				},
			},
			{
				Text: "Doctor's fees",
				Action: func(p *Player) {
					p.payBank(50)
				},
			},
		},
	}

	g.Chance = &CardDeck{
		Cards: []Card{
			{
				Text: "Advance to Go",
				Action: func(p *Player) {
					p.Position = 0
					p.recieveMoneyFromBank(200)
				},
			},
			{
				Text: "Bank error in your favor",
				Action: func(p *Player) {
					p.recieveMoneyFromBank(200)
				},
			},
			{
				Text: "Doctor's fees",
				Action: func(p *Player) {
					p.payBank(50)
				},
			},
		},
	}
	g.SetupBoardSpaces()

	g.CurrentPlayer = g.Players[0]
	g.PlayerOrder = g.Players
	g.Turn(g.CurrentPlayer)
}
