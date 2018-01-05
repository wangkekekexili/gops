package reporter

import (
	"errors"
	"os"

	"github.com/getsentry/raven-go"
	"github.com/wangkekekexili/gops/till"
)

type Module struct {
	Till *till.Module
}

func (m *Module) Load() error {
	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN == "" {
		return errors.New("SENTRY_DSN is empty")
	}
	raven.SetDSN(sentryDSN)
	return nil
}

func (m *Module) Err(err error) {
	raven.CaptureError(err, nil)
}

func (m *Module) ErrSync(err error) {
	raven.CaptureErrorAndWait(err, nil)
}

func (m *Module) ErrSMS(err error) {
	raven.CaptureError(err, nil)
	m.Till.Notify(err)
}
