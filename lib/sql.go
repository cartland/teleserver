package lib

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/calsol/teleserver/messages"
	"github.com/stvnrhodes/broadcaster"
)

type DB struct {
	sql *sql.DB
}

// NewDB returns the sql database after creating any needed tables.
func NewDB(db *sql.DB) (*DB, error) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS messages (time INT, canid INT, data TEXT)")
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db DB) writeCAN(c *messages.CANPlus) error {
	data, err := json.Marshal(c.CAN)
	if err != nil {
		return err
	}
	_, err = db.sql.Exec(
		"INSERT INTO messages (time, canid, data) VALUES (?, ?, ?)",
		c.Time.UnixNano(), c.CANID, string(data),
	)
	return err
}

// WriteMessages logs the broadcasted date as JSON to the database.
func (db DB) WriteMessages(b broadcaster.Caster) {
	for msg := range b.Subscribe(nil) {
		var err error
		switch msg := msg.(type) {
		case messages.CANPlus:
			err = db.writeCAN(&msg)
		case *messages.CANPlus:
			err = db.writeCAN(msg)
		}
		if err != nil {
			log.Println(err)
		}
	}
}

func (db DB) GetLatest(canid uint16) (*messages.CANPlus, error) {
	row := db.sql.QueryRow("SELECT MAX(time), data FROM messages WHERE canid = ?", canid)
	var unixNanos int64
	var data []byte
	err := row.Scan(&unixNanos, &data)
	if err != nil {
		return nil, err
	}
	msg := messages.IDToMessage(canid)
	if err := json.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return &messages.CANPlus{CAN: msg, CANID: canid, Time: time.Unix(0, unixNanos)}, nil
}
