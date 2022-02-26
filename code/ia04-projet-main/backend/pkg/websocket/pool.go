package websocket

import (
	"fmt"
	"math/rand"
	"time"
)

type notes struct {
	id   int
	note float32
}

type Pool struct {
	MaxTour    int
	Tour       int
	Gains      [2]notes
	Msgs       []CMPair
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Choix      chan CMPair
	AgtIDName  map[int]string
}

func NewPool() *Pool {
	AgtIDName := map[int]string{
		1: "prisonnier muet",
		2: "prisonnier méchant",
		3: "prisonnier oeil pour oeil",
		4: "prisonnier aléatoire",
		5: "prisonnier intelligent",
		6: "prisonnier plus intelligent",
	}
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Choix:      make(chan CMPair),
		Msgs:       make([]CMPair, 0),
		AgtIDName:  AgtIDName,
	}
}

func (pl *Pool) Handle() {
	pl.Tour++
	for pl.Tour <= pl.MaxTour {
		// afficher le tour actuel
		fmt.Printf("POLICE: tour%d\n", pl.Tour)
		smallerID, biggerID := pl.rankID()
		pl.setGainsID(smallerID, biggerID)
		ligneMsg, colonneMsg := pl.rankMsg()
		pl.Count(ligneMsg, colonneMsg)
		// afficher les notes de chaque client
		fmt.Printf("%s action: %s, %s action: %s\n", pl.AgtIDName[smallerID], ligneMsg, pl.AgtIDName[biggerID], colonneMsg)
		fmt.Printf("%s note: %f, %s note: %f\n", pl.AgtIDName[smallerID], pl.Gains[0].note, pl.AgtIDName[biggerID], pl.Gains[1].note)
		actMsg := fmt.Sprintf("%s action: %s; %s action: %s", pl.AgtIDName[smallerID], ligneMsg, pl.AgtIDName[biggerID], colonneMsg)
		message := Message{Type: 1, Body: actMsg}
		fmt.Println("Sending message to all clients in Pool")
		for client := range pl.Clients {
			if err := client.Conn.WriteJSON(message); err != nil {
				fmt.Println(err)
				return
			}
		}
		if pl.Tour == pl.MaxTour {
			pl.Msgs = pl.Msgs[:0]
			pl.showWinMsg()
			break
		}
		time.Sleep(time.Second)
		pl.Tour++
		pl.Msgs = pl.Msgs[:0]
		for len(pl.Msgs) < 2 {
			choix := <-pl.Choix
			fmt.Printf("Tour%d receive a message %q from agent%d",pl.Tour, choix.msg.Body, choix.client.ID)
			if len(pl.Msgs) == 0 {
				pl.Msgs = append(pl.Msgs, choix)
			} else if len(pl.Msgs) == 1 {
				if choix.client.ID != pl.Msgs[0].client.ID {
					pl.Msgs = append(pl.Msgs, choix)
				}
			}
		}
	}
}

func (pl *Pool) Start() {
	pl.setMaxTour()
	for {
		select {
		case client := <-pl.Register:
			fmt.Println("Taille de la Connection Pool: ", len(pl.Clients))
			if len(pl.Clients) < 1 {
				pl.Clients[client] = true
				client.Conn.WriteJSON(Message{Type: 1, Body: "Attendez un autre prisionnier..."})
			} else if len(pl.Clients) == 1 {
				pl.Clients[client] = true
				for client := range pl.Clients {
					client.Conn.WriteJSON(Message{Type: 1, Body: "Débat commence..."})
				}
			}
		case client := <-pl.Unregister:
			delete(pl.Clients, client)
			fmt.Println("Taille de la Connection Pool: ", len(pl.Clients))
			for client := range pl.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "Un prisonnier part..."})
			}
		case choix := <-pl.Choix:
			fmt.Printf("Tour%d receive a message %q from agent%d",1, choix.msg.Body, choix.client.ID)
			if len(pl.Msgs) == 0 {
				pl.Msgs = append(pl.Msgs, choix)
			} else if len(pl.Msgs) == 1 {
				if choix.client.ID != pl.Msgs[0].client.ID {
					pl.Msgs = append(pl.Msgs, choix)
				}
			}
			if len(pl.Msgs) == 2 {
				pl.Handle()
			}
		}
	}
}

func (srv *Pool) showWinMsg() {
	fmt.Println("Sending a winner message to all clients")
	avg1 := (0 - srv.Gains[0].note) / float32(srv.MaxTour)
	avg2 := (0 - srv.Gains[1].note) / float32(srv.MaxTour)
	agt1 := srv.AgtIDName[srv.Gains[0].id]
	agt2 := srv.AgtIDName[srv.Gains[1].id]
	winMsg := fmt.Sprintf("%s est condamné à %f ans de prison; %s est condamné à %f ans de prison", agt1, avg1, agt2, avg2)
	message := Message{Type: 1, Body: winMsg}
	for client := range srv.Clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (pl *Pool) rankID() (smallerID int, biggerID int) {
	if pl.Msgs[0].client.ID < pl.Msgs[1].client.ID {
		return pl.Msgs[0].client.ID, pl.Msgs[1].client.ID
	}
	return pl.Msgs[1].client.ID, pl.Msgs[0].client.ID
}

func (pl *Pool) rankMsg() (ligneMsg string, colonneMsg string) {
	if pl.Msgs[0].client.ID == pl.Gains[0].id {
		return pl.Msgs[0].msg.Body, pl.Msgs[1].msg.Body
	}
	return pl.Msgs[1].msg.Body, pl.Msgs[0].msg.Body
}

func (pl *Pool) Count(ligneMsg string, colonneMsg string) {
	switch {
	case ligneMsg == "se tait" && colonneMsg == "se tait":
		pl.Gains[0].note += -0.5
		pl.Gains[1].note += -0.5
	case ligneMsg == "se tait" && colonneMsg == "trahit":
		pl.Gains[0].note += -10
	case ligneMsg == "trahit" && colonneMsg == "se tait":
		pl.Gains[1].note += -10
	case ligneMsg == "trahit" && colonneMsg == "trahit":
		pl.Gains[0].note += -5
		pl.Gains[1].note += -5
	}
}

func (pl *Pool) setGainsID(id1 int, id2 int) {
	pl.Gains[0].id = id1
	pl.Gains[1].id = id2
}

func (pl *Pool) setMaxTour() {
	rand.Seed(time.Now().UnixNano())
	pl.MaxTour = rand.Intn(10) + 10 // entier aléatoire dans [10, 20)
	fmt.Printf("MaxTour: %d\n", pl.MaxTour)
}
