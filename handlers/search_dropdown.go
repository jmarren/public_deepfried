package handlers

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

type SearchDropdown struct {
	JHandle
}

func NewSearchDropdown(j JHandle) *SearchDropdown {
	s := new(SearchDropdown)
	s.JHandle = j
	return s
}

func (s *SearchDropdown) ParseReq() string {
	keyword := s.URL.Query().Get("keyword")
	return keyword
}

func (s *SearchDropdown) GetComponent() *templ.Component {
	fmt.Println("getting dropdown items")

	// parse request
	keyword := s.ParseReq()

	searchService := services.NewSearchService()

	// Get data from search service
	data := searchService.GetDropdownItems(s.Context(), &keyword)

	for _, item := range data {
		fmt.Printf("data.ArtworkSrc: %s\n", item.ArtworkSrc)
	}

	// get full component
	component := components.SearchDropdown(data)
	return &component
}

func (s *SearchDropdown) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := s.GetComponent()
	(*component).Render(r.Context(), w)
}
