package lib

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/calsol/teleserver/msgs"
	"github.com/stvnrhodes/broadcaster"
)

// DB represents an SQL database for the teleserver.
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

// WriteCAN writes the CAN message to the database.
func (db DB) WriteCAN(c *msgs.CANPlus) error {
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
		case msgs.CANPlus:
			err = db.WriteCAN(&msg)
		case *msgs.CANPlus:
			err = db.WriteCAN(msg)
		}
		if err != nil {
			log.Println(err)
		}
	}
}

// GetLatest returns the most recent message with the given id
func (db DB) GetLatest(canid uint16) (*msgs.CANPlus, error) {
	row := db.sql.QueryRow("SELECT MAX(time), data FROM messages WHERE canid = ?", canid)
	var unixNanos int64
	var data []byte
	if err := row.Scan(&unixNanos, &data); err != nil {
		return nil, err
	}
	msg := msgs.IDToMessage(canid)
	if err := json.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return &msgs.CANPlus{CAN: msg, CANID: canid, Time: time.Unix(0, unixNanos).UTC()}, nil
}

// GetSince returns all messages with the given id that have happened since the
// given duration.
func (db DB) GetSince(d time.Duration, canid uint16) ([]*msgs.CANPlus, error) {
	t := time.Now().Add(-d).UnixNano()
	rows, err := db.sql.Query("SELECT time, data FROM messages WHERE canid = ? AND time > ? ORDER BY time", canid, t)
	msg := msgs.IDToMessage(canid)
	if err != nil {
		return nil, err
	}
	var results []*msgs.CANPlus
	for rows.Next() {
		msg := msg.New()
		var unixNanos int64
		var data []byte
		if err := rows.Scan(&unixNanos, &data); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, msg); err != nil {
			return nil, err
		}
		results = append(results, &msgs.CANPlus{CAN: msg, CANID: canid, Time: time.Unix(0, unixNanos).UTC()})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
