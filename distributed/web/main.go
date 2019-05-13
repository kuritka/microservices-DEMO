package main

import (
	"microservices/distributed/web/controller"
	"net/http"
)

func main() {
	controller.Initialize()

	http.ListenAndServe(":3000", nil)
}
