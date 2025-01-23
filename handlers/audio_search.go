package handlers

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"

	"net/http"
)

type AudioSearch struct {
	ctx           context.Context
	searchService SearchService
	filters       *services.SearchAudioFilters
	prevPage      string
	nextPage      string
}

func NewAudioSearch(ctx context.Context, prevPage string, nextPage string, filters *services.SearchAudioFilters, searchService SearchService) *AudioSearch {
	s := new(AudioSearch)
	s.ctx = ctx
	s.filters = filters
	s.searchService = searchService
	s.prevPage = prevPage
	s.nextPage = nextPage
	return s
}

func (a *AudioSearch) GetComponent() *templ.Component {
	results := a.searchService.GetAudioSearchResults(a.ctx, a.filters)
	hasMore := true

	if len(*results.AudioFiles) < 20 {
		fmt.Printf("len of files: %d\n", len(*results.AudioFiles))
		hasMore = false
	}

	component := components.AudioSearchTableBody(a.filters.Page, a.prevPage, a.nextPage, hasMore, results)

	return &component
}

func (a *AudioSearch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := a.GetComponent()
	(*component).Render(r.Context(), w)
}
