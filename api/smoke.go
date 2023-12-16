package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func smoke(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	s := time.Now().String()
	_, err := w.Write([]byte(s))
	if err != nil {
		log.Errorf("%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Infof("handle msg success...")
	return
}
