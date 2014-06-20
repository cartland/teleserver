package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/stvnrhodes/broadcaster"
)

// Importer accepts messages with the proper key and puts them into the database
type Importer struct {
	db     *DB
	secret []byte
}

func NewImporter(db *DB, secret []byte) Importer {
	return Importer{db: db, secret: secret}
}

func (i Importer) writeToDB(data string) error {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(data), &parsed); err != nil {
		return err
	}

	id, ok := parsed["canID"].(float64)
	if !ok {
		return errors.New("Missing canID")
	}

	b, err := json.Marshal(parsed["CAN"])
	if err != nil {
		return err
	}

	timeStr, ok := parsed["time"].(string)
	if !ok {
		return errors.New("Missing time")
	}
	var t time.Time
	if err := t.UnmarshalText([]byte(timeStr)); err != nil {
		return err
	}

	return i.db.WriteMessage(t, uint16(id), b)
}

func (i Importer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json := r.FormValue("data")
	key := r.FormValue("key")

	mac := hmac.New(sha256.New, i.secret)
	mac.Write([]byte(json))

	hash, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if !hmac.Equal(mac.Sum(nil), hash) {
		http.Error(w, "Invalid key", 403)
		return
	}

	if err := i.writeToDB(json); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func PostOnBroadcast(b broadcaster.Caster, url string, secret []byte) {
	for msg := range b.Subscribe(nil) {
		if b, err := json.Marshal(msg); err != nil {
			log.Printf("Could not encode broadcast: %v", err)
		} else if resp, err := PostToURL(url, string(b), secret); err != nil {
			log.Printf("Could not send broadcast: %v", err)
		} else if resp.StatusCode != 200 {
			log.Printf("Bad status code: %v", resp)
			s, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(string(s))
				resp.Body.Close()
			}
		}
	}
}

func PostToURL(urlStr, json string, secret []byte) (resp *http.Response, err error) {
	vals := make(url.Values)
	vals.Set("data", json)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(json))
	vals.Set("key", (base64.StdEncoding.EncodeToString(mac.Sum(nil))))

	return http.PostForm(urlStr, vals)
}
