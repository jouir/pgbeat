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
	done     chan bool
	beat     bool
	recovery bool
	table    *base.Table
}

// NewBeatmaker creates a Beatmaker manager
func NewBeatmaker(config *base.Config, done chan bool) *Beatmaker {
	return &Beatmaker{
		config: config,
		done:   done,
	}
}

// Run starts the Beatmaker
func (bm *Beatmaker) Run() {
	bm.table = base.NewTable(bm.config.Schema, bm.config.Table)
	bm.db = base.NewDb(bm.config.Dsn())

	log.Println("Connecting to instance")
	bm.db.Connect()
	defer bm.terminate()

	// Initial recovery check
	bm.recovery = bm.db.InRecovery()

	// Further recovery checks in asynchronous mode
	go bm.checkRecovery()

	if !bm.db.TableExists(bm.table) {
		log.Println("Creating table", bm.table)
		bm.db.CreateTable(bm.table)
	}

	bm.beat = bm.db.BeatExists(bm.table, bm.config.ID)
	go bm.upsertBeat()

	<-bm.done
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

// terminate cleans up connection
func (bm *Beatmaker) terminate() {
	log.Println("Terminating")
	bm.db.Disconnect()
}
