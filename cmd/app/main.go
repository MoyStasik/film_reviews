package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Index")
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Login")
}

func signupHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Signup")
}

func run() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/signup", signupHandler)
	http.Handle("/", router)

	fmt.Println("Server is listening on port 8888")
	http.ListenAndServe(":8888", nil)
}

func main() {
	run()

}
