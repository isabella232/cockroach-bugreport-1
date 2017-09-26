package main

import (
	"context"
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func testInsertsWithoutTransactions(ctx context.Context, db *sql.DB) {
	// Insert two rows into the "accounts" table.
	failures := 0
	inserts := 0
	for i := 1000; i < 2000; i++ {
		_, err := db.ExecContext(ctx, "INSERT INTO accounts (id, balance) VALUES ($1, 1000)", i)
		if err != nil {
			failures = failures + 1
		}
		inserts = inserts + 1
	}

	log.Printf("RESULT: %d / %d inserts without transactions failed due to transaction errors\n", failures, inserts)
}

func testInsertsWithTransactions(ctx context.Context, db *sql.DB) {
	// Insert two rows into the "accounts" table.
	failures := 0
	inserts := 0
	for i := 0; i < 1000; i++ {
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Begin: %v", err)
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO accounts (id, balance) VALUES ($1, 1000)", i)
		if err != nil {
			log.Fatalf("Exec: %v", err)
		}

		err = tx.Commit()
		if err != nil {
			failures = failures + 1
		}
		inserts = inserts + 1
	}

	log.Printf("RESULT: %d / %d inserts with transactions failed due to transaction errors\n", failures, inserts)
}

func main() {
	// Connect to the "bank" database.
	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:26257/bank?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	// Drop the "accounts" table.
	if _, err := db.Exec("DROP TABLE IF EXISTS accounts"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)"); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			// Repeatedly select ids
			row := db.QueryRowContext(ctx, "SELECT id FROM accounts ORDER BY id DESC limit 1")

			var id int
			err = row.Scan(&id)
			if err == context.Canceled {
				return
			}

			if err != nil {
				if strings.Contains(err.Error(), "no rows in result set") {
					continue
				}
				log.Fatal(err)
			}

			_ = id
		}
	}()

	testInsertsWithTransactions(ctx, db)
	testInsertsWithoutTransactions(ctx, db)

	cancel()
}
