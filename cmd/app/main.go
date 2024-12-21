package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Handler struct {
	pool pgxpool.Pool
}

func (h *Handler) indexHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := h.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), "SELECT id, film_name from film")
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}

	for rows.Next() {
		var id int
		var film_name string
		err = rows.Scan(&id, &film_name)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
		}
		log.Printf("ID: %d, Film_name: %s\n", id, film_name)
	}

	fmt.Fprint(w, "Index")
}

func (h *Handler) loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Login")
}

func (h *Handler) signupHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Signup")
}

func run(addr string) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, "postgres://postgres:1906@postgres_reviews:5432/reviews")

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	handler := Handler{
		pool: *pool,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handler.indexHandler)
	router.HandleFunc("/login", handler.loginHandler)
	router.HandleFunc("/signup", handler.signupHandler)
	http.Handle("/", router)

	fmt.Println("Server is listening on port 8888")
	http.ListenAndServe(addr, nil)
}

func main() {

	run(":8888")

}
