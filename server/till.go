package server

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Till struct {
	Config *Config

	ok bool

	url   string
	phone string
}

func (t *Till) Load() error {
	if t.Config == nil {
		return nil
	}
	t.url = t.Config.TillURL
	t.phone = t.Config.TillTarget
	if t.url == "" || t.phone == "" {
		return nil
	}
	t.ok = true
	return nil
}

func (t *Till) Notify(err error) {
	if !t.ok {
		return
	}
	req := struct {
		Phone []string `json:"phone"`
		Text  string   `json:"text"`
	}{
		Phone: []string{t.phone},
		Text:  err.Error(),
	}
	reqBytes, _ := json.Marshal(req)
	reader := bytes.NewReader(reqBytes)
	http.Post(t.url, "application/json", reader)
}
