package handlers

import (
	"fmt"
	"net/http"
)

type Notifications struct {
	JHandle
}

func NewNotifications(j JHandle) *Notifications {
	n := new(Notifications)
	n.JHandle = j
	return n
}

func (n *Notifications) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Notifications")
	if n.Method == "PATCH" {
		n.MarkFollowNotficationSeen()
	}

	fmt.Fprintf(w, "")
}
