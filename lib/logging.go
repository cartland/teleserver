package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/stvnrhodes/broadcaster"
)

const maxFileSize = 1 << 20

func logName(path string, t time.Time) string {
	return fmt.Sprintf("%s_%s.txt", path, t.Format("2006-01-02_15:04:05"))
}

// write to file will attemt to write the contents of the reader to the file.
func writeToFile(path string, r io.Reader) error {
	filename := logName(path, time.Now())
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("Wrote %d bytes to %v", n, filename)
	return nil
}

// LogToFile logs the broadcasted date as JSON to a file.
func LogToFile(path string, b broadcaster.Caster) {
	buf := &bytes.Buffer{}
	e := json.NewEncoder(buf)

	// Pay attention to the kill signal so we can flush the buffer.
	flush := make(chan os.Signal, 1)
	signal.Notify(flush, os.Interrupt, os.Kill)

	msgs := b.Subscribe(nil)
	defer writeToFile(path, buf)
	for {
		select {
		case msg := <-msgs:
			// Fill up buffer with any data.
			if err := e.Encode(msg); err != nil {
				log.Print(err)
			}

			// Empty buffer to file if it's large enough.
			if buf.Len() >= maxFileSize {
				if err := writeToFile(path, buf); err != nil {
					log.Fatal(err)
				}
			}

		case <-flush:
			log.Println("Catching interrupt, flushing logs to disk")
			if err := writeToFile(path, buf); err != nil {
				log.Println(err)
			}
			os.Exit(0)
		}
	}
}
