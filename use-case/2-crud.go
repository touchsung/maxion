package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
)

func connectDB(dbDriver string, dbSource string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}
	
	// Create users table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			age INTEGER NOT NULL
		)`
	
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

func insertUser(db *sql.DB, name string, age int) error {
	query := "INSERT INTO users (name, age) VALUES ($1, $2)"
	_, err := db.Exec(query, name, age)
	return err
}

func getUser(db *sql.DB, name string) (int, string, int, error) {
	var id, age int
	var userName string
	query := "SELECT id, name, age FROM users WHERE name = $1"
	err := db.QueryRow(query, name).Scan(&id, &userName, &age)
	return id, userName, age, err
}

func updateUserAge(db *sql.DB, name string, newAge int) error {
	query := "UPDATE users SET age = $1 WHERE name = $2"
	_, err := db.Exec(query, newAge, name)
	return err
}

func deleteUser(db *sql.DB, name string) error {
	query := "DELETE FROM users WHERE name = $1"
	_, err := db.Exec(query, name)
	return err
}

func main() {
	db, err := connectDB(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := insertUser(db, "Alice", 25); err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}

	id, name, age, err := getUser(db, "Alice")
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Printf("User: ID=%d, Name=%s, Age=%d\n", id, name, age)

	if err := updateUserAge(db, "Alice", 26); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	id, name, age, err = getUser(db, "Alice")
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Printf("Updated User: ID=%d, Name=%s, Age=%d\n", id, name, age)

	if err := deleteUser(db, "Alice"); err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
}
