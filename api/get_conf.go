package api

import (
	"awesomeProject1/consts"
	"awesomeProject1/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func getConf(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("%v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	m := new(model.ServiceConf)
	m.Name = consts.ServiceName
	m.Version = consts.ServiceVersion
	bytes, err := json.Marshal(m)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err.Error()), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err.Error()), http.StatusInternalServerError)
		return
	}
	return
}
