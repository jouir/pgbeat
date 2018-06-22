package base

import (
	"database/sql"
	"fmt"
	// Force database/sql to use libpq
	_ "github.com/lib/pq"
)

// Db centralizes connection to the database
type Db struct {
	dsn  string
	conn *sql.DB
}

// NewDb creates a Db object
func NewDb(dsn string) *Db {
	return &Db{
		dsn: dsn,
	}
}

// Connect connects to the instance and ping it to ensure connection is working
func (db *Db) Connect() {
	conn, err := sql.Open("postgres", db.dsn)
	Panic(err)

	err = conn.Ping()
	Panic(err)

	db.conn = conn
}

// Disconnect ends connection cleanly
func (db *Db) Disconnect() {
	err := db.conn.Close()
	Panic(err)
}

// TableExists checks if a table exists
// Not using SQL "create table if not exists" statement because some users
// don't have DDL privileges and table could already exists
func (db *Db) TableExists(table Table) bool {
	var exists bool
	query := `select true as exists from information_schema.tables where table_schema = $1 and table_name = $2 limit 1;`
	err := db.conn.QueryRow(query, table.Schema, table.Name).Scan(&exists)
	if err == sql.ErrNoRows {
		exists = false
	} else {
		Panic(err)
	}
	return exists
}

// CreateTable initializes table structure on instance
func (db *Db) CreateTable(table Table) {
	if !db.TableExists(table) {
		query := fmt.Sprintf("create table %s (id bigint primary key, ts timestamptz not null);", table)
		_, err := db.conn.Exec(query)
		Panic(err)
	}
}

// BeatExists checks for beat existance
func (db *Db) BeatExists(table Table, serverID int) bool {
	var exists bool
	query := fmt.Sprintf("select true as result from %s where id = $1 limit 1;", table)
	err := db.conn.QueryRow(query, serverID).Scan(&exists)
	if err == sql.ErrNoRows {
		exists = false
	} else {
		Panic(err)
	}
	return exists
}

// InsertBeat insert a beat into the table
func (db *Db) InsertBeat(table Table, serverID int) {
	query := fmt.Sprintf("insert into %s (id, ts) values ($1, now());", table)
	_, err := db.conn.Exec(query, serverID)
	Panic(err)
}

// UpdateBeat updates an already existing beat in the table
func (db *Db) UpdateBeat(table Table, serverID int) {
	query := fmt.Sprintf("update %s set ts = now() where id = $1;", table)
	_, err := db.conn.Exec(query, serverID)
	Panic(err)
}

// InRecovery checks if instance is in recovery mode (read-only)
func (db *Db) InRecovery() (recovery bool) {
	query := "select pg_is_in_recovery();"
	err := db.conn.QueryRow(query).Scan(&recovery)
	Panic(err)
	return recovery
}
