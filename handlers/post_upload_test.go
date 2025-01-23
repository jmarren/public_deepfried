package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func NewPostUpload(ctx context.Context, us UserService) *PostUpload {

func Test_NewPostUpload(t *testing.T) {
	mockUserService := NewMockUserService(t)
	ctx := context.Background()
	postUpload := NewPostUpload(ctx, mockUserService)

}
