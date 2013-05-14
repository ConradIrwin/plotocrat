package main

import (
	"./db"
	"./http"
	"os"
)

func main() {
	print("Heloo !")
	db.Setup(os.Getenv("MONGOHQ_URL"))
	http.Listen(os.Getenv("PORT"))
}
