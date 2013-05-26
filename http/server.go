package http

import (
	"encoding/json"
	"fmt"
	"github.com/ConradIrwin/plotocrat/db"
	"github.com/ConradIrwin/plotocrat/models"
	"github.com/gorilla/mux"
	"net/http"
)

func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", upload).Methods("POST")
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/{uid:[a-f0-9]{20}}.txt", download(asTxt)).Methods("GET").Name("txt")
	r.HandleFunc("/{uid:[a-f0-9]{20}}.tsv", download(asTsv)).Methods("GET").Name("tsv")
	r.HandleFunc("/{uid:[a-f0-9]{20}}.json", download(asJson)).Methods("GET").Name("json")
	r.HandleFunc("/{uid:[a-f0-9]{20}}.svg", download(asSvg)).Methods("GET").Name("svg")

	return r
}

func Listen(port string) {
	http.Handle("/", router())

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world", req.URL.Path)
}

func upload(res http.ResponseWriter, req *http.Request) {
	// Limit the size of uploaded files: https://code.google.com/p/go/issues/detail?id=2093
	req.Body = http.MaxBytesReader(res, req.Body, 10*1024*1024)

	err := req.ParseMultipartForm(10 * 1024 * 1024)

	if err != nil {
		fmt.Println(err)
		http.Error(res, "Invalid POST", http.StatusBadRequest)
		return
	}

	for key, fileHeaders := range req.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()

			if err != nil {
				fmt.Println(err)
				http.Error(res, "Upload failed", http.StatusInternalServerError)
				return
			}

			plot, err := models.NewPlot(key, file)

			if err != nil {
				fmt.Println(err)
				http.Error(res, "Parsing data failed", http.StatusInternalServerError)
				return
			}

			err = db.SavePlot(plot)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Writing to database failed", http.StatusInternalServerError)
				return
			}

			url, err := router().Get("tsv").URL("uid", plot.Uid)
			if err != nil {
				panic(err)
			}

			http.Redirect(res, req, url.String(), http.StatusMovedPermanently)
			fmt.Fprintln(res, "http://"+req.Host+url.String())
		}
	}
}

func download(handler func(*models.Plot, http.ResponseWriter)) func(http.ResponseWriter, *http.Request) {
	return (func(res http.ResponseWriter, req *http.Request) {
		plot, err := db.LoadPlot(mux.Vars(req)["uid"])

		if err != nil {
			fmt.Println(err)
			http.Error(res, "Reading from database failed", http.StatusNotFound)
			return
		}

		handler(plot, res)
	})
}

func asTxt(plot *models.Plot, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(res, plot.Data)
}

func asJson(plot *models.Plot, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(res)

	err := enc.Encode(plot.Values())
	if err != nil {
		panic(err)
	}
}

func asSvg(plot *models.Plot, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "image/svg+xml")
	viz := models.NewVisualization(plot)
	viz.WriteTo(res)
}

func asTsv(plot *models.Plot, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, value := range plot.Values() {
		fmt.Fprintln(res, value)
	}
}
