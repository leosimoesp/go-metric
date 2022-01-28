package log

import log "github.com/sirupsen/logrus"

func Logger() *log.Logger {
	res := log.New()
	res.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	return res
}
