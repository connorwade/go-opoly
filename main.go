package main

func main() {
	g := Game{
		Players: []*Player{
			{Name: "Player 1"},
			{Name: "Player 2"},
		},
	}

	g.Start()
}
