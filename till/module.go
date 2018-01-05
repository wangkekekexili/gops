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

func (t *Module) Load() error {
	t.url = os.Getenv("TILL_URL")
	t.phone = os.Getenv("TILL_TARGET")
	if t.url == "" || t.phone == "" {
		return nil
	}
	t.ok = true
	return nil
}

func (t *Module) Notify(err error) {
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

