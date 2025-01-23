package handlers

import (
	"context"
	// "fmt"
	"github.com/a-h/templ"
	"github.com/jmarren/deepfried/components"
	"net/http"
	"slices"
)

type TagBar struct {
	ctx        context.Context
	tagService TagService
	keyword    string
}

func NewTagBar(ctx context.Context, keyword string, tagService TagService) *TagBar {
	t := new(TagBar)
	t.ctx = ctx
	t.tagService = tagService
	t.keyword = keyword
	return t
}

func (t *TagBar) GetComponent() *templ.Component {
	tags := t.tagService.GetMostPopular()

	keywordTag := ""

	for i := 0; i < len(tags); i++ {
		if tags[i] == t.keyword {
			keywordTag = tags[i]
			tags = slices.Delete(tags, i, i+1)
			break
		}
	}

	component := components.TagBarBody(tags, keywordTag)
	return &component
}

func (t *TagBar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	component := t.GetComponent()
	(*component).Render(r.Context(), w)
}
