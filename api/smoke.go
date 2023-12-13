package api

import (
	"fmt"
	"net/http"
	"time"
)

func smoke(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	s := time.Now().String()
	_, err := w.Write([]byte(s))
	if err != nil {
		fmt.Errorf("%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
