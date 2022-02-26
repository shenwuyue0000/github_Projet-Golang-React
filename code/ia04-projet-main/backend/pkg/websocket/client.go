package websocket

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID      int
	Conn    *websocket.Conn
	Pool    *Pool
	oppActs []string
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type CMPair struct {
	client *Client
	msg    Message
}

func NewClient(id int, pool *Pool, conn *websocket.Conn) *Client {
	sl := make([]string, 0)
	return &Client{id, conn, pool, sl}
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	switch {
	case c.ID == 1:
		nowTour := 1
		for nowTour <= c.Pool.MaxTour {
			cmpair := c.agent_taire()
			c.Pool.Choix <- cmpair
			fmt.Printf("Agent1 Tour%d Choix : %+v\n", nowTour, cmpair.msg.Body)
			prevTour := nowTour
			nowTour = c.Pool.Tour
			for prevTour >= nowTour {
				time.Sleep(time.Second)
				nowTour = c.Pool.Tour
			}
		}
	case c.ID == 2:
		nowTour := 1
		for nowTour <= c.Pool.MaxTour {
			cmpair := c.agent_trahir()
			c.Pool.Choix <- cmpair
			fmt.Printf("Agent2 Tour%d Choix : %+v\n", nowTour, cmpair.msg.Body)
			prevTour := nowTour
			nowTour = c.Pool.Tour
			for prevTour >= nowTour {
				time.Sleep(time.Second)
				nowTour = c.Pool.Tour
			}
		}
	case c.ID == 3:
		nowTour := 1
		for nowTour <= c.Pool.MaxTour {
			cmpair := c.agent_oeil_pour_oeil(nowTour, c.oppActs)
			c.Pool.Choix <- cmpair
			fmt.Printf("Agent3 Tour%d Choix: %s\n", nowTour, cmpair.msg.Body)
			c.collectOppAct(1)
			prevTour := nowTour
			nowTour = c.Pool.Tour
			for prevTour >= nowTour {
				time.Sleep(time.Second)
				nowTour = c.Pool.Tour
			}
		}
	case c.ID == 4:
		nowTour := 1
		for nowTour <= c.Pool.MaxTour {
			cmpair := c.agent_aleatoire()
			c.Pool.Choix <- cmpair
			fmt.Printf("Agent4 Tour%d  Choix : %+v\n", nowTour, cmpair.msg.Body)
			prevTour := nowTour
			nowTour = c.Pool.Tour
			for prevTour >= nowTour {
				time.Sleep(time.Second)
				nowTour = c.Pool.Tour
			}
		}
	case c.ID == 5:
		nowTour := 1
		taille := 5
		for nowTour <= c.Pool.MaxTour {
			cmpair := c.agent_i_d(nowTour, c.oppActs)
			c.Pool.Choix <- cmpair
			fmt.Printf("Agent5 Tour%d  Choix : %+v\n", nowTour, cmpair.msg.Body)
			c.collectOppAct(taille)
			prevTour := nowTour
			nowTour = c.Pool.Tour
			for prevTour >= nowTour {
				time.Sleep(time.Second)
				nowTour = c.Pool.Tour
			}
		}
	case c.ID == 6:
		nowTour := 1
		taille := 5
		for nowTour <= c.Pool.MaxTour {
			cmpair := c.agent_poid_recent(nowTour, c.oppActs)
			c.Pool.Choix <- cmpair
			fmt.Printf("Agent6 Tour%d  Choix : %+v\n", nowTour, cmpair.msg.Body)
			c.collectOppAct(taille)
			prevTour := nowTour
			nowTour = c.Pool.Tour
			for prevTour >= nowTour {
				time.Sleep(time.Second)
				nowTour = c.Pool.Tour
			}
		}
	}
}

func (c *Client) collectOppAct(taille int) {
	if len(c.oppActs) >= taille {
		c.oppActs = c.oppActs[1:]
	}
	for len(c.Pool.Msgs) < 2 {
		time.Sleep(time.Microsecond)
	}
	for _, cm := range c.Pool.Msgs {
		if cm.client.ID != c.ID {
			c.oppActs = append(c.oppActs, cm.msg.Body)
			break
		}
	}
}

func (c *Client) agent_taire() (pair CMPair) {
	str := "se tait"
	message := Message{Type: 1, Body: str}
	pair = CMPair{client: c, msg: message}
	return
}

func (c *Client) agent_trahir() (pair CMPair) {
	str := "trahit"
	message := Message{Type: 1, Body: str}
	pair = CMPair{client: c, msg: message}
	return
}

func (c *Client) agent_oeil_pour_oeil(tour int, choix_histoire_autre []string) CMPair {
	str := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if tour == 1 {
		arbitraire := r.Intn(100)
		if 0 <= arbitraire && arbitraire < 50 {
			str = "se tait"
			message := Message{Type: 1, Body: str}
			return CMPair{client: c, msg: message}
		} else {
			str = "trahit"
			message := Message{Type: 1, Body: str}
			return CMPair{client: c, msg: message}
		}
	}
	if choix_histoire_autre[0] == "se tait" {
		str = "se tait"
	} else if choix_histoire_autre[0] == "trahit" {
		str = "trahit"
	}
	message := Message{Type: 1, Body: str}
	return CMPair{client: c, msg: message}
}

func (c *Client) agent_aleatoire() CMPair {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	arbitraire := r.Intn(100)
	str := ""
	if 0 <= arbitraire && arbitraire < 50 {
		str = "se tait"
		message := Message{Type: 1, Body: str}
		return CMPair{client: c, msg: message}
	} else {
		str = "trahit"
		message := Message{Type: 1, Body: str}
		return CMPair{client: c, msg: message}
	}
}

func (c *Client) agent_i_d(tour int, choix_histoire_autre []string) CMPair {
	str := ""
	message := Message{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if tour == 1 {
		arbitraire := r.Intn(100)
		if 0 <= arbitraire && arbitraire < 50 {
			str = "se tait"
			message = Message{Type: 1, Body: str}
			return CMPair{client: c, msg: message}
		} else {
			str = "trahit"
			message = Message{Type: 1, Body: str}
			return CMPair{client: c, msg: message}
		}
	}
	if tour == 2 {
		str = choix_histoire_autre[0]
		message = Message{Type: 1, Body: str}
		return CMPair{client: c, msg: message}
	}

	arbitraire := r.Intn(101)

	len_histoire := len(choix_histoire_autre)
	interval := int((1 / len_histoire) * 100)
	for i := 0; i < len_histoire; i++ {
		if i*interval <= arbitraire && arbitraire <= (i+1)*interval {
			str = choix_histoire_autre[i]
			message = Message{Type: 1, Body: str}
			break
		}
		if i == len_histoire-1 {
			if i*interval <= arbitraire && arbitraire <= 100 {
				str = choix_histoire_autre[i]
				message = Message{Type: 1, Body: str}
				break
			}
		}
	}
	return CMPair{client: c, msg: message}
}

func (c *Client) agent_poid_recent(tour int, choix_histoire_autre []string) CMPair {
	str := ""
	message := Message{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if tour == 1 {
		arbitraire := r.Intn(100)
		if 0 <= arbitraire && arbitraire < 50 {
			str = "se tait"
			message = Message{Type: 1, Body: str}
			return CMPair{client: c, msg: message}
		} else {
			str = "trahit"
			message := Message{Type: 1, Body: str}
			return CMPair{client: c, msg: message}
		}
	}
	if tour == 2 {
		str = choix_histoire_autre[0]
		message = Message{Type: 1, Body: str}
		return CMPair{client: c, msg: message}
	}
	arbitraire := r.Intn(101)
	len_histoire := len(choix_histoire_autre)

	var poids = make([]int, len_histoire)
	max := 101
	temp := 0
	for i := 0; i < len_histoire; i++ {
		if i == len_histoire-1 {
			for j := 0; j < i; j++ {
				temp = temp + poids[j]
			}
			poids[i] = 100 - temp
			break
		}
		poids[i] = r.Intn(max)
		max = max - poids[i]
	}
	poids = append(poids, 0)
	sort.Ints(poids)
	for i := 0; i < len_histoire+1; i++ {
		if i >= 2 {
			poids[i] = poids[i] + poids[i-1]
		}
	}
	for i := 0; i < len_histoire; i++ {
		if poids[i] <= arbitraire && arbitraire <= poids[i+1] {
			str = choix_histoire_autre[i]
			message = Message{Type: 1, Body: str}
			break
		}
	}
	return CMPair{client: c, msg: message}
}
