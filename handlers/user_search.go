package handlers

import (
	"context"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
)

type UserSearch struct {
	searchService SearchService
	keyword       string
	ctx           context.Context
}

func NewUserSearch(ctx context.Context, searchService SearchService, keyword string) *UserSearch {
	u := new(UserSearch)
	u.ctx = ctx
	u.keyword = keyword
	u.searchService = searchService
	return u
}

func (u *UserSearch) GetComponent() *templ.Component {
	searchResults := u.searchService.GetUserSearchResults(u.ctx, &u.keyword)
	component := components.UserSearchSectionBody(searchResults)
	return &component
}

func (u *UserSearch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := u.GetComponent()
	(*component).Render(r.Context(), w)
}
