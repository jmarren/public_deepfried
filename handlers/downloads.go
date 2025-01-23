package handlers

import (
	"fmt"
	"github.com/jmarren/deepfried/services"
	"net/http"
)

type DownloadHandler struct {
	JHandle
}

func NewDownloadHandler(j JHandle) *DownloadHandler {
	d := new(DownloadHandler)
	d.JHandle = j
	return d
}

func (d *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("download handler")
	user, auth := d.GetFromCtx()
	if !auth {
		w.WriteHeader(http.StatusOK)
		return
	}
	downloadService := services.NewDownloadsService(d.Context())
	downloadId := d.PathValue("id")
	downloadService.AddUserDownload(user.ID, downloadId)
	fmt.Printf("downloadId: %s\n", downloadId)
	w.WriteHeader(http.StatusOK)
}
