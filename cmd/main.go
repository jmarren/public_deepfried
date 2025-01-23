package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jmarren/deepfried/run"
)

func main() {
	environment := flag.String("env", "dev", "Environment to use: dev or prod")
	authenticated := flag.String("auth", "none", "Authenticated: true or false")
	logLevel := flag.String("log_level", "low", "Log Level: low or high")
	resetDb := flag.String("reset_db", "false", "whether to clear db schema")
	uploadStatic := flag.String("upload_static", "false", "whether to upload js and css files to s3")
	flag.Parse()
	os.Setenv("auth", *authenticated)
	os.Setenv("env", *environment)
	os.Setenv("loglevel", *logLevel)
	os.Setenv("resetdb", *resetDb)
	os.Setenv("uploadstatic", *uploadStatic)

	fmt.Println("--------------- args -----------")
	fmt.Printf("\tauth: %s\n\tenv: %s\n\tlog_level: %s\n\tresetdb: %s\n\tuploadstatic: %s\n\n", *authenticated, *environment, *logLevel, *resetDb, *uploadStatic)

	ctx := context.Background()
	// defer cancel()
	s, err := run.Run(ctx, os.Stdout, os.Stdin, os.Getenv, os.Args, *environment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer s.Close()
}
