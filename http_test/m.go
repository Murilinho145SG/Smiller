package main

import (
	"net/http"

	"github.com/Murilinho145SG/gouter"
	"github.com/Murilinho145SG/gouter/httpio"
	"github.com/Murilinho145SG/gouter/log"
)

func main() {
	r := gouter.NewRouter()

	r.Route("/post", func(w httpio.Writer, r *httpio.Request) {
		log.Info(r.RemoteAddr)

		b, err := r.Body.Read()
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info(string(b))
	})

	r.Route("/get", func(w httpio.Writer, r *httpio.Request) {
		log.Info(r.RemoteAddr)

		w.WriteWR([]byte("test"), 200)
	})

	gouter.Run("0.0.0.0:8080", r)
}
