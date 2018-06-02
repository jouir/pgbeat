package manager

import (
	"github.com/jouir/pgbeat/base"
	"log"
	"time"
)

// Beatmaker manages beats sent to the instance
type Beatmaker struct {
	config *base.Config
	db     *base.Db
	beat   bool
}

// NewBeatmaker creates a Beatmaker manager
func NewBeatmaker(config *base.Config) *Beatmaker {
	return &Beatmaker{
		config: config,
	}
}

// Run fires up the Beatmaker
func (bm *Beatmaker) Run() {
	table := base.NewTable(bm.config.Schema, bm.config.Table)
	bm.db = base.NewDb(bm.config.Dsn())

	log.Println("Connecting to instance")
	bm.db.Connect()
	defer bm.db.Disconnect()

	if !bm.db.TableExists(table) {
		log.Println("Creating table", table)
		bm.db.CreateTable(table)
	}

	bm.beat = bm.db.BeatExists(table, bm.config.ID)

	for {
		if bm.beat {
			log.Println("Updating beat")
			bm.db.UpdateBeat(table, bm.config.ID)
		} else {
			log.Println("Inserting beat")
			bm.db.InsertBeat(table, bm.config.ID)
			bm.beat = true
		}
		time.Sleep(time.Duration(bm.config.Interval) * time.Millisecond)
	}
}

// Terminate cleans up connection
func (bm *Beatmaker) Terminate() {
	log.Println("Terminating")
	bm.db.Disconnect()
}