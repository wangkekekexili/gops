package server

import (
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type Reporter struct {
	Config *Config
	Till   *Till
}

func (r *Reporter) Load() error {
	if r.Config == nil || r.Config.SentryDSN == "" {
		return errors.New("cannot load sentry")
	}
	raven.SetDSN(r.Config.SentryDSN)
	return nil
}

func (r *Reporter) Err(err error) {
	raven.CaptureError(err, nil)
}

func (r *Reporter) ErrSync(err error) {
	raven.CaptureErrorAndWait(err, nil)
}

func (r *Reporter) ErrSMS(err error) {
	raven.CaptureError(err, nil)
	r.Till.Notify(err)
}
