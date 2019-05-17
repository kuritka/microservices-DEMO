package main

import (
	"fmt"
	"microservices/distributed/web/controller"
	"net/http"
)

func main() {
	controller.Initialize()

	fmt.Printf("http://localhost:3000/public Listening...")
	http.ListenAndServe(":3000", nil)

}
