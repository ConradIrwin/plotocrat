package http

import (
	"../data"
	"fmt"
	"net/http"
)

func Listen(port string) {
	http.HandleFunc("/", index)
	fmt.Println("Listening on port %s", port)

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

		for key, file_headers := range req.MultipartForm.File {
			for _, file_header := range file_headers {
				file, err := file_header.Open()

				if err != nil {
					fmt.Println(err)
					http.Error(res, "Upload failed", http.StatusInternalServerError)
					return
				}

				raw, err := data.Upload(key, file)

				if err != nil {
					fmt.Println(err)
					http.Error(res, "Upload failed", http.StatusInternalServerError)
					return
				}

				fmt.Fprintln(res, "Got", key, raw.Uid, raw.Uploaded_at)
				for value := range raw.Series() {
					fmt.Fprintln(res, ">", value)
				}
			}
		}

	} else {
		http.Error(res, "Only GET and POST are supported", http.StatusMethodNotAllowed)
	}
}
