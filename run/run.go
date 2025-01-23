package run

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/cache"
	"github.com/jmarren/deepfried/db"
	"github.com/jmarren/deepfried/handlers"
	"github.com/jmarren/deepfried/util"
)

func printRoutineCount() {
	for {
		fmt.Printf("routines: %d\n", runtime.NumGoroutine())
		time.Sleep(time.Second * 5)
	}
}

//go:embed upload-worker.js
var uploadWorkerScript string

func Run(ctx context.Context, w io.Writer, r io.Reader, envfunc func(string) string, args []string, environment string) (*http.Server, error) {
	_, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cache.Init()
	awssdk.InitAWS()
	err := db.Init(ctx, environment)
	util.EMsg(err, "initializing db")

	if environment == "dev" {
		err := db.InitTestData(ctx)
		util.EMsg(err, "initializing test data")
	}

	if os.Getenv("uploadstatic") == "true" {
		awssdk.UploadStaticAssets()
	}

	mux := http.NewServeMux()
	var staticDir http.Dir
	if environment == "dev" {
		os.Setenv("env", "dev")
		staticDir = http.Dir("/home/john-marren/Projects/deepfried/static/assets/")
		os.Setenv("STATIC_DOMAIN", "/static/assets/")
		fs := http.FileServer(staticDir)
		mux.Handle("GET /static/assets/", http.StripPrefix("/static/assets/", fs))
	}
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /workers/upload-worker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s", uploadWorkerScript)
	})

	mux.Handle("/", http.StripPrefix("", handlers.NewDefaultHandler()))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if os.Getenv("loglevel") == "high" {
		go printRoutineCount()
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		return nil, errors.New("failed to listen")
	}

	return s, nil
}
