package lib

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/stvnrhodes/broadcaster"
)

func LogToFile(path string, b broadcaster.Caster) {
	f, err := os.OpenFile(path+"_"+time.Now().Format("2006-01-02_15:04:05"), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		e := json.NewEncoder(f)
		for msg := range b.Subscribe(nil) {
			if err := e.Encode(msg); err != nil {
				log.Print(err)
			}
		}
	}()
}
