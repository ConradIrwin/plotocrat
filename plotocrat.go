package main

import (
	"fmt"
	"github.com/ConradIrwin/plotocrat/db"
	"github.com/ConradIrwin/plotocrat/http"
	"os"
)

func main() {
	db.Setup(os.Getenv("MONGOHQ_URL"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("plotocrat listening on localhost:", port)
	http.Listen(port)
}
