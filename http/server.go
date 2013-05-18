package http

import (
	"github.com/ConradIrwin/plotocrat/data"
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
					http.Error(res, "Upload failed", http.StatusInternalServerError)
					return
				}

				fmt.Fprintln(res, "Got", key, plot.Uid, plot.UploadedAt)
				for _, value := range plot.Data {
					fmt.Fprintln(res, ">", value)
				}
			}
		}

	} else {
		http.Error(res, "Only GET and POST are supported", http.StatusMethodNotAllowed)
	}
}
