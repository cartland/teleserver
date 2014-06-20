package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
)

var totallyFakeSecretKey = []byte("123")

func ServeHTTPWithHMAC(w http.ResponseWriter, r *http.Request) {
	json := r.FormValue("data")
	key := r.FormValue("key")

	mac := hmac.New(sha256.New, totallyFakeSecretKey)
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
	http.Error(w, "ok", 200)
}

func PostToURL(urlStr string, json string) (resp *http.Response, err error) {
	vals := make(url.Values)
	vals.Set("data", json)

	mac := hmac.New(sha256.New, totallyFakeSecretKey)
	mac.Write([]byte(json))
	vals.Set("key", (base64.StdEncoding.EncodeToString(mac.Sum(nil))))

	return http.PostForm(urlStr, vals)
}
