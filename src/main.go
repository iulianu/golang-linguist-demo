package main

import (
	"net/http"
	"goji.io"
	"goji.io/pat"
	"encoding/json"
	"strings"
	"bytes"
)

type MungingRequest struct {
	Text string
}

type MungingResponse struct {
	Text string `json:"text"`
}

func requestFormat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mr MungingRequest

	if jsonErr := decoder.Decode(&mr); jsonErr != nil {
		w.WriteHeader(400)
		return
	}

	mungingReader := strings.NewReader(mr.Text)
	mungingWriter := bytes.NewBuffer([]byte{})

	pairs, err := ParsePhrasePairStream([]byte(mr.Text))
	err = ParsePhrasePairReader(mungingWriter, mungingReader)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	resp := MungingResponse{Text: MungedPairs(pairs)}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(resp); err != nil {
		w.WriteHeader(500)
		return
	}
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/format"), requestFormat)

	http.ListenAndServe("localhost:9090", mux)
}