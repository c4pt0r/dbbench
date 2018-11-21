package databases

import (
	"database/sql"
	"log"
	"os"
)

// SQLite implements the bencher implementation.
type SQLite struct {
	db *sql.DB
}

// NewSQLite a new mysql bencher.
func NewSQLite() *SQLite {
	// TODO: filename as flag
	db, err := sql.Open("sqlite3", "./dbbench.sqlite?cache=shared")
	if err != nil {
		log.Fatalf("failed to open connection: %v\n", err)
	}

	db.SetMaxOpenConns(1)
	p := &SQLite{db: db}
	return p
}

// Benchmarks returns the individual benchmark functions for the mysql db.
func (m *SQLite) Benchmarks() []Benchmark {
	return []Benchmark{
		{"inserts", m.inserts},
		{"updates", m.updates},
		{"selects", m.selects},
		{"deletes", m.deletes},
	}
}

// Setup initializes the database for the benchmark.
func (m *SQLite) Setup(...string) {
	if _, err := m.db.Exec("CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance DECIMAL);"); err != nil {
		log.Fatalf("failed to create table: %v\n", err)
	}
}

// Cleanup removes all remaining benchmarking data.
func (m *SQLite) Cleanup() {
	if _, err := m.db.Exec("DROP TABLE accounts"); err != nil {
		log.Printf("failed to drop table: %v\n", err)
	}
	if err := m.db.Close(); err != nil {
		log.Printf("failed to close connection: %v", err)
	}
	if err := os.Remove("./dbbench.sqlite"); err != nil {
		log.Printf("not able to delete old database file: %v\n", err)
	}
}

func (m *SQLite) inserts(from, to int) {
	const q = "INSERT INTO accounts VALUES(?, ?);"
	for i := from; i < to; i++ {
		if _, err := m.db.Exec(q, i, i); err != nil {
			log.Fatalf("failed to insert: %v\n", err)
		}
	}
}

func (m *SQLite) selects(from, to int) {
	const q = "SELECT * FROM accounts WHERE id = ?;"
	for i := from; i < to; i++ {
		if _, err := m.db.Exec(q, i); err != nil {
			log.Fatalf("failed to select: %v\n", err)
		}
	}
}

func (m *SQLite) updates(from, to int) {
	const q = "UPDATE accounts SET balance = ? WHERE id = ?;"
	for i := from; i < to; i++ {
		if _, err := m.db.Exec(q, i, i); err != nil {
			log.Fatalf("failed to update: %v\n", err)
		}
	}
}

func (m *SQLite) deletes(from, to int) {
	const q = "DELETE FROM accounts WHERE id = ?"
	for i := from; i < to; i++ {
		if _, err := m.db.Exec(q, i); err != nil {
			log.Fatalf("failed to delete: %v\n", err)
		}
	}
}