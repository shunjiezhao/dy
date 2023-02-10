package objects

import (
	"log"
	"net/http"
)

/**
ResponseWriter，Request我就不解释了
*/
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	log.Println(m)
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodPost {
		post(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
