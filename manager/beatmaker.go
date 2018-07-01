package manager

import (
	"github.com/jouir/pgbeat/base"
	"log"
	"time"
)

// Beatmaker manages beats sent to the instance
type Beatmaker struct {
	config   *base.Config
	db       *base.Db
	table    *base.Table
	beat     bool
	recovery bool
}

// NewBeatmaker creates a Beatmaker manager
func NewBeatmaker(config *base.Config) *Beatmaker {
	return &Beatmaker{
		config: config,
	}
}

// Run starts the Beatmaker
func (bm *Beatmaker) Run() {
	if bm.config.CreateDatabase {
		bm.createDatabase(bm.config.Database)
	}

	bm.table = base.NewTable(bm.config.Schema, bm.config.Table)
	bm.db = base.NewDb(bm.config.Dsn())

	log.Println("Connecting to database", bm.config.Database)
	bm.db.Connect()

	// Asynchronously checks for recovery mode
	bm.recovery = bm.db.InRecovery()
	go bm.checkRecovery()

	if bm.config.CreateTable {
		bm.createTable(bm.table)
	}

	bm.beat = bm.db.BeatExists(bm.table, bm.config.ID)
	bm.upsertBeat()
}

// checkRecovery looks for recovery mode and cache this informations
// Checks are in a different schedule to avoid hammering the instance
func (bm *Beatmaker) checkRecovery() {
	for {
		bm.recovery = bm.db.InRecovery()
		time.Sleep(time.Duration(bm.config.RecoveryInterval*1000) * time.Millisecond)
	}
}

// upsertBeat checks for beat existance and insert or update it
func (bm *Beatmaker) upsertBeat() {
	for {
		if bm.recovery {
			log.Println("Not inserting beat (recovery mode)")
		} else {
			if bm.beat {
				log.Println("Updating beat")
				bm.db.UpdateBeat(bm.table, bm.config.ID)
			} else {
				log.Println("Inserting beat")
				bm.db.InsertBeat(bm.table, bm.config.ID)
				bm.beat = true
			}
		}
		time.Sleep(time.Duration(bm.config.Interval*1000) * time.Millisecond)
	}
}

// createDatabase connects to instance (with or without database name) and
// create a database if it doesn't exit. It waits for instance to be writable.
func (bm *Beatmaker) createDatabase(name string) {
	var dsn string
	if bm.config.ConnectDatabase != "" {
		dsn = bm.config.DsnWithDatabase(bm.config.ConnectDatabase)
	} else {
		dsn = bm.config.DsnWithoutDatabase()
	}
	db := base.NewDb(dsn)

	log.Println("Connecting to instance to create database")
	db.Connect()

	// Asynchronously checks for recovery mode
	recovery := db.InRecovery()
	go func() {
		for {
			recovery = db.InRecovery()
			time.Sleep(time.Duration(bm.config.RecoveryInterval*1000) * time.Millisecond)
		}
	}()

	// Wait for instance to be writable
	for {
		if recovery {
			log.Println("Not creating database (recovery mode)")
		} else {
			if !db.DatabaseExists(name) {
				log.Println("Creating database", name)
				db.CreateDatabase(name)
			}
			return
		}
		time.Sleep(time.Duration(bm.config.Interval*1000) * time.Millisecond)
	}
}

// createTable create destination table if it doesn't exist
// it waits for instance to be writable
func (bm *Beatmaker) createTable(table *base.Table) {
	// Wait for instance to be writable
	for {
		if bm.recovery {
			log.Println("Not creating table (recovery mode)")
		} else {
			if !bm.db.TableExists(table) {
				log.Println("Creating table", table)
				bm.db.CreateTable(table)
			}
			return
		}
		time.Sleep(time.Duration(bm.config.Interval*1000) * time.Millisecond)
	}
}
