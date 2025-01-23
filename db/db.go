package db

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmarren/deepfried/awssdk"
	"github.com/jmarren/deepfried/sqlc"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// func Init(ctx context.Context) {
// }
//
//go:embed sql/schema.sql
var initQueryProd string

// /go:embed sql/restart_db.sql
var restartDBQuery string

var (
	Dbtx  *pgxpool.Pool
	Query *sqlc.Queries
)

func Init(ctx context.Context, environment string) error {
	fmt.Printf("initializing %s environment\n", environment)

	if environment == "dev" {
		return initDev(ctx)
	}
	if environment == "prod" {
		return initProd(ctx)
	}
	return fmt.Errorf("environment not specified")
}

func initDev(ctx context.Context) error {
	err := godotenv.Load("/home/john-marren/Projects/deepfried/.env")
	if err != nil {
		return err
	}
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Printf("psqlInfo:%s\n", psqlInfo)

	Dbtx, err = pgxpool.New(ctx, psqlInfo)
	if err != nil {
		return err
	}

	if os.Getenv("resetdb") == "true" {
		resetDB, err := os.ReadFile("/home/john-marren/Projects/deepfried/db/sql/restart.sql")
		if err != nil {
			return fmt.Errorf("error reading restart.sql file: %s", err)
		}
		_, err = Dbtx.Exec(ctx, string(resetDB))
		if err != nil {
			log.Fatalf("failed to execute reset script: %v", err)
		}
	}

	initQuery, err := os.ReadFile("/home/john-marren/Projects/deepfried/db/sql/schema.sql")
	if err != nil {
		return fmt.Errorf("error reading schema.sql file: %s", err)
	}

	status, err := Dbtx.Exec(ctx, string(initQuery))
	if err != nil {
		log.Fatalf("failed to execute init script: %v", err)
	}
	fmt.Printf("status of init exec: %s\n", status)

	Query = sqlc.New(Dbtx)

	if Query == nil {
		return fmt.Errorf("error: nil db.Query !")
	}
	fmt.Printf("db.Query: %v\n", *Query)

	// fmt.Printf("status of testDataQuery exec: %s\n", status)
	fmt.Printf("database initialized successfully\n")
	return nil
}

func initProd(ctx context.Context) error {
	fmt.Println("running production init (initProd)")
	// TESTING *****
	// initQuery, err := os.ReadFile("/Users/johnmarren//db/sql/schema.sql")
	// *********
	// initQuery, err := os.ReadFile("~/app/sql/schema.sql")
	// if err != nil {
	// 	return fmt.Errorf("error reading schema.sql file: %s", err)
	// }

	var err error
	Dbtx, err = pgxpool.New(ctx, awssdk.DbDsn)
	if err != nil {
		return fmt.Errorf("error connection to prod db: %s", err)
	}
	fmt.Printf("Dbtx: %v\n", Dbtx)
	status, err := Dbtx.Exec(ctx, initQueryProd)
	if err != nil {
		return fmt.Errorf("failed to execute init script: %v", err)
	}
	fmt.Printf("status of init exec: %s\n", status)

	Query = sqlc.New(Dbtx)
	if Query == nil {
		return fmt.Errorf("error: nil db.Query !")
	}
	fmt.Printf("db.Query: %v\n", *Query)
	if Query == nil {
		return fmt.Errorf("error: nil db.Query !")
	}
	fmt.Printf("db.Query: %v\n", *Query)

	return nil
}

func ErrorCode(err error) (string, bool) {
	pgerr, ok := err.(*pgconn.PgError)
	if !ok {
		return "", false
	}
	return pgerr.Code, true
}
