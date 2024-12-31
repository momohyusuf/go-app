package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/momoh-yusuf/note-app/generated_sql"
)

type DatabaseConnection struct {
	connection *generated_sql.Queries
}

func Db_Query() *generated_sql.Queries {
	ctx := context.Background()
	DB_URL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(ctx, DB_URL)

	if err != nil {
		log.Fatalf("Error occurred while connecting: %v \n", err)
	}

	fmt.Println("Database connection successfull")

	queries := DatabaseConnection{
		connection: generated_sql.New(conn),
	}
	return queries.connection
}
