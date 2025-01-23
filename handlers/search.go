package handlers

import (
	// "context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"github.com/jmarren/deepfried/services"
	"github.com/jmarren/deepfried/util"
	"net/http"
	"strconv"
)

type Search struct {
	JHandle
	searchService SearchService
}

type ReqParse struct {
	filters  *services.SearchAudioFilters
	nextPage string
	prevPage string
}

type SearchProps struct {
	keyword string
	*services.AudioSearchResults
}

func NewSearch(j JHandle) *Search {
	s := new(Search)
	s.JHandle = j

	return s
}

func (s *Search) ParseRequest() *ReqParse {
	pageStr := s.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	util.EMsg(err, "converting page query to int")

	nextPage, err := util.UpdateQueryParam(s.URL.RequestURI(), "page", fmt.Sprintf("%d", page+1))
	prevPage, err := util.UpdateQueryParam(s.URL.RequestURI(), "page", fmt.Sprintf("%d", page-1))

	filtersStrs := services.SearchAudioFilters{
		Page:                   page,
		Keyword:                s.URL.Query().Get("keyword"),
		BpmRadioStr:            s.FormValue("bpm-radio"),
		ExactBpmStr:            s.FormValue("exact-bpm"),
		MaxBpmStr:              s.FormValue("max-bpm"),
		MinBpmStr:              s.FormValue("min-bpm"),
		MusicalKeyStr:          s.FormValue("musical-key"),
		MajorMinorStr:          s.FormValue("major-minor"),
		MusicalKeySignatureStr: s.FormValue("musical-key-signature"),
		IncludesStemsOnlyStr:   s.FormValue("includes-stems-only"),
	}
	return &ReqParse{
		&filtersStrs,
		nextPage,
		prevPage,
	}

}

func (s *Search) GetComponent() *templ.Component {
	// get data from request
	searchService := services.NewSearchService()
	reqParse := s.ParseRequest()

	keyword := reqParse.filters.Keyword

	// get tagbar component
	tagService := services.NewTagService(s.Context())
	tagBar := NewTagBar(s.Context(), keyword, tagService).GetComponent()

	// get user search component
	userSearch := NewUserSearch(s.Context(), searchService, reqParse.filters.Keyword).GetComponent()

	// get audio search table component
	audioSearch := NewAudioSearch(s.Context(), reqParse.prevPage, reqParse.nextPage, reqParse.filters, searchService).GetComponent()

	// get full component
	component := components.SearchPageBody(reqParse.filters.Keyword, *tagBar, *userSearch, *audioSearch)
	return &component
}

func (s *Search) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := s.GetComponent()
	(*component).Render(r.Context(), w)
}
