package profiler

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/stutiredboy/gobetween/logging"
)

func Start(bind string) {
	log := logging.For("profiler")

	log.Infof("Starting profiler: %v", bind)

	go func() {
		log.Error(http.ListenAndServe(bind, nil))
	}()
}
