package util

import (
	"os"
	"strings"
	"sync"

	"github.com/stvp/rollbar"
)

var (
	rollbarInit  sync.Once
	rollbarReady bool
)

func ReportError(err error) {
	rollbarInit.Do(func() {
		rollbarToken := strings.TrimSpace(os.Getenv("ROLLBAR_TOKEN"))
		if rollbarToken == "" {
			rollbarReady = false
			return
		}
		rollbar.Token = rollbarToken
		rollbar.Environment = "production"
	})
	if rollbarReady {
		rollbar.Error(rollbar.ERR, err)
	}
}
