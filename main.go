package main

import (
	"fmt"
	"net/http"

	"github.com/ConradIrwin/plotocrat/db"
	"github.com/ConradIrwin/plotocrat/env"
	"github.com/ConradIrwin/plotocrat/plotocrat"
)

func main() {
	ctx, cancel := env.Init()
	defer cancel()

	env.DB = db.Open(ctx)

	fmt.Println("plotocrat listening on " + env.Config.Listen)
	if err := http.ListenAndServe(env.Config.Listen, wrapHandler(plotocrat.Handler)); err != nil {
			panic(err)
	}
}

func wrapHandler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer func () {
			if r := recover(); r != nil {
				fmt.Println("panic:", r)
				w.WriteHeader(http.StatusInternalServerError)

				str := fmt.Sprintf("%#v", r)
				if e, ok := r.(error); ok {
					str = e.Error()
				}

				w.Write([]byte("internal-server-error:" + str))
			}
		}()

		f(w, r)
	})
}