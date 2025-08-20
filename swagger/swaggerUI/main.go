package main

import (
	"log"
	"net/http"
)

func main() {
	// Путь к папке с Swagger UI и вашим swagger.json/swagger.yaml
	fs := http.FileServer(http.Dir("./swagger-ui/dist"))

	http.Handle("/", fs)

	log.Println("Swagger UI доступен по адресу http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
