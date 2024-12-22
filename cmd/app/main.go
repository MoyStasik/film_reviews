package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Handler struct {
	pool pgxpool.Pool

	//client ssov1.AuthClient
}

type filmStruct struct {
	FilmId       int
	FilmName     string
	CategoryName string
	ImgPath      string
	ReviewsValue int
}

type filmsDTO struct {
	Films []filmStruct
}

type filmReview struct {
	UserId        int //заменить на name
	ReviewContent string
	ReviewValue   string
}

type filmsReviewDTO struct {
	FilmsReview []filmReview
}

func (h *Handler) indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Index")
}

func (h *Handler) loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Login")
}

func (h *Handler) signupHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Signup")

}

func (h *Handler) filmReviewsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Пришел запрос: %v\n", req.URL.Query())
	var filmId = req.URL.Query().Get("film_id")
	conn, err := h.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), `SELECT review.user_id, review_content, review_value FROM review
														WHERE review.film_id = $1;`, filmId)
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}

	filmsReview := []filmReview{}
	for rows.Next() {
		var userId int
		var reviewContent string
		var reviewValue string

		err = rows.Scan(&userId, &reviewContent, &reviewValue)

		if err != nil {
			log.Printf("Failed to scan row: %v", err)
		}

		filmsReview = append(filmsReview, filmReview{
			UserId:        userId,
			ReviewContent: reviewContent,
			ReviewValue:   reviewValue,
		})
	}

	filmReviewDto := filmsReviewDTO{
		FilmsReview: filmsReview,
	}

	jsonResp, err := json.Marshal(filmReviewDto)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	log.Printf("jsonResp: %v", jsonResp)
	log.Printf("Scan row: %v", filmsReview)
}

func (h *Handler) filmsHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := h.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), `SELECT film.id, film_name, category_name, img_path, reviews_value FROM film
													JOIN film_category ON film.film_category_id = film_category.id;`)
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}

	films := []filmStruct{}
	for rows.Next() {

		var filmId int
		var filmName string
		var categoryName string
		var imgPath string
		var reviewsValue int
		err = rows.Scan(&filmId, &filmName, &categoryName, &imgPath, &reviewsValue)

		if err != nil {
			log.Printf("Failed to scan row: %v", err)
		}

		films = append(films, filmStruct{
			FilmId:       filmId,
			FilmName:     filmName,
			CategoryName: categoryName,
			ImgPath:      imgPath,
			ReviewsValue: reviewsValue,
		})

	}
	filmsDto := filmsDTO{
		Films: films,
	}
	jsonResp, err := json.Marshal(filmsDto)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	log.Printf("jsonResp: %v", jsonResp)
	log.Printf("Scan row: %v", films)

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
	router.HandleFunc("/films", handler.filmsHandler).Methods("GET")
	router.HandleFunc("/filmReviews", handler.filmReviewsHandler).Methods("GET")
	http.Handle("/", router)

	fmt.Println("Server is listening on port 8888")
	http.ListenAndServe(addr, nil)
}

func main() {

	run(":8888")

}
