package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
)

type Player struct {
	JHandle
}

type PlayerReqParse struct {
	artworkSrc string
	title      string
	creator    string
}

func NewPlayer(j JHandle) *Player {
	p := new(Player)
	p.JHandle = j
	return p
}

func (p *Player) GetComponent() *templ.Component {
	fmt.Printf("url: %v\n ", p.URL.Query())

	newCurrent := p.URL.Query().Get("playing")
	newQueue := p.URL.Query().Get("queue")

	playableService := services.NewPlayableService(p.Context())

	fmt.Printf("audioQueue: %v\n", newQueue)

	queueArr := strings.Split(newQueue, ",")

	data := playableService.GetPlayerData(newCurrent, queueArr)
	fmt.Printf("playerData: %v\n", *data.Current)
	component := components.Player(data)
	return &component
}

func (p *Player) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// playerPresent := r.Header.Get("X-Player-Present")

	// if playerPresent == "true" {
	// 	w.Header().Set("HX-Reswap", "morphdom")
	// }

	component := p.GetComponent()
	(*component).Render(p.Context(), w)
}
