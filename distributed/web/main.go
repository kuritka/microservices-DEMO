package main

import (
	"fmt"
	"microservices/distributed/web/controller"
	"net/http"
)

func main() {
	controller.Initialize()

	fmt.Printf("Listening...")
	http.ListenAndServe(":3000", nil)

}
