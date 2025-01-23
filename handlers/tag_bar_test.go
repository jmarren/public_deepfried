package handlers

import (
	"context"
	"fmt"
	"testing"
)

func Test_NewTagBar(t *testing.T) {
	ctx := context.Background()

	keywords := []string{"hi", "pop", "rap"}

	mockTagService := NewMockTagService(t)

	mockTagService.On("GetMostPopular").Return([]string{"hi", "hello"})

	for _, t_keyword := range keywords {
		got := NewTagBar(ctx, t_keyword, mockTagService)

		if got == nil {
			t.Fail()
		}
		if got.GetComponent() == nil {
			t.Fail()
		}
		component := got.GetComponent()
		if fmt.Sprintf("%T", component) != "*templ.Component" {
			t.Fail()
		}
	}

}
