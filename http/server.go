package http

import (
	"github.com/ConradIrwin/plotocrat/data"
	"github.com/ConradIrwin/plotocrat/db"
	"fmt"
	"net/http"
)

func Listen(port string) {
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Fprintln(res, "hello, world")
	} else if req.Method == "POST" {

		// Limit the size of uploaded files: https://code.google.com/p/go/issues/detail?id=2093
		req.Body = http.MaxBytesReader(res, req.Body, 10 * 1024 * 1024)

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

				plot, err := data.Parse(key, file)

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

				plot, err = db.LoadPlot(plot.Uid)
				if err != nil {
					fmt.Println(err)
					http.Error(res, "Reading from database failed", http.StatusInternalServerError)
					return
				}

				fmt.Fprintln(res, "Gotterred", key, plot.Uid, plot.UploadedAt)
				for _, value := range plot.Values() {
					fmt.Fprintln(res, ">", value)
				}
			}
		}

	} else {
		http.Error(res, "Only GET and POST are supported", http.StatusMethodNotAllowed)
	}
}
