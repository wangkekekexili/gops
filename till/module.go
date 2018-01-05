package till

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Module struct {
	ok bool

	url, phone string
}

func (m *Module) Load() error {
	m.url = os.Getenv("TILL_URL")
	m.phone = os.Getenv("TILL_TARGET")
	if m.url == "" || m.phone == "" {
		return nil
	}
	m.ok = true
	return nil
}

func (m *Module) Notify(err error) {
	if !m.ok {
		return
	}
	req := struct {
		Phone []string `json:"phone"`
		Text  string   `json:"text"`
	}{
		Phone: []string{m.phone},
		Text:  err.Error(),
	}
	reqBytes, _ := json.Marshal(req)
	reader := bytes.NewReader(reqBytes)
	http.Post(m.url, "application/json", reader)
}
