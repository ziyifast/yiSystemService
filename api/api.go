package api

import (
	"awesomeProject1/consts"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func interceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//resolve the cross origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func ServiceHandler() error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Errorf("%v\n", e)
		}
	}()
	mux := http.NewServeMux()
	mux.Handle("/yi-service/smoke", interceptor(http.HandlerFunc(smoke)))
	mux.Handle("/yi-service/conf", interceptor(http.HandlerFunc(getConf)))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", consts.ServicePort), mux); err != nil {
		log.Errorf("err=%v", err)
		return err
	}
	return nil
}
