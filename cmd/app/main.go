package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ssov1 "github.com/Lesha222/protos/gen/go/sso"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	pool pgxpool.Pool

	client ssov1.AuthClient
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

type userReviews struct {
	UserId        int //заменить на name
	FilmId        int
	ReviewContent string
	ReviewValue   string
}

type usersReviewsDTO struct {
	FilmsReview []userReviews
}

type filmsReviewsDTO struct {
	FilmsReview []filmReview
}

type filmSendReview struct {
	UserId        int
	FilmId        int
	ReviewContent string
	ReviewValue   string
}

type filmReviewDTO struct {
	Film filmSendReview
}

type loginDTO struct {
	appId    int
	Email    string
	Password string
}

type UserClaims struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
	AppID string `json:"app_id"`
	jwt.RegisteredClaims
}

type registerDTO struct {
	Email    string
	Name     string
	Password string
}

func (h *Handler) indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Index")
}

func VerifyToken(tokenString string, appSecret string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.Email, nil
	}

	return "", fmt.Errorf("invalid token")
}

func (h *Handler) loginHandler(w http.ResponseWriter, req *http.Request) {
	var creds loginDTO
	if err := json.NewDecoder(req.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	conn, err := h.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()
	loginResponse := &ssov1.LoginResponse{}
	loginResponse, err = h.client.Login(req.Context(), &ssov1.LoginRequest{
		AppId:    1,
		Email:    creds.Email,
		Password: creds.Password,
	})

	if err != nil {
		log.Printf("Ошибка логина %v", err)
	} else {
		log.Printf("%v", loginResponse)
	}
	jsonResp, err := json.Marshal(loginResponse)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	var secretKey string = "GOIDA"
	fmt.Println(VerifyToken(string(jsonResp), secretKey))

}

func (h *Handler) signupHandler(w http.ResponseWriter, req *http.Request) {
	var creds registerDTO
	if err := json.NewDecoder(req.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	conn, err := h.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}
	registerResponse := &ssov1.RegisterResponse{}
	defer conn.Release()
	registerResponse, err = h.client.Register(req.Context(), &ssov1.RegisterRequest{
		Email:    creds.Email,
		Name:     creds.Name,
		Password: creds.Password,
	})

	if err != nil {
		log.Printf("Ошибка регистрации %v", err)
	} else {
		log.Printf("%v", registerResponse)
	}
	jsonResp, err := json.Marshal(registerResponse)

	if err != nil {
		log.Printf("error happened in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)

}

// GET all reviews for film by query param id
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

	filmReviewsDto := filmsReviewsDTO{
		FilmsReview: filmsReview,
	}

	jsonResp, err := json.Marshal(filmReviewsDto)

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

// GET all films
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
	log.Printf("jsonesp: %v", jsonResp)
	log.Printf("Scan row: %v", films)

}

// POST review on film
func (h *Handler) filmReviewHandler(w http.ResponseWriter, req *http.Request) {
	var creds filmReviewDTO
	if err := json.NewDecoder(req.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	conn, err := h.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()
	rows := h.pool.QueryRow(context.Background(), `INSERT INTO review (
																user_id,
																film_id,
																review_content,
																review_value
																) VALUES ($1, $2, $3, $4);`,
		creds.Film.UserId,
		creds.Film.FilmId,
		creds.Film.ReviewContent,
		creds.Film.ReviewValue)
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}

	log.Printf("film review post: %v", creds.Film.ReviewValue)
	log.Printf("film review request: %v", rows.Scan())

}

// GET all reviews by userId
func (h *Handler) userReviewsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("Пришел запрос: %v\n", req.URL.Query())
	var userId = req.URL.Query().Get("user_id")
	conn, err := h.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), `SELECT review.user_id, review.film_id, review_content, review_value FROM review
														WHERE review.user_id = $1;`, userId)
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
	}

	filmsReview := []userReviews{}
	for rows.Next() {
		var userId int
		var filmId int
		var reviewContent string
		var reviewValue string

		err = rows.Scan(&userId, &filmId, &reviewContent, &reviewValue)

		if err != nil {
			log.Printf("Failed to scan row: %v", err)
		}

		filmsReview = append(filmsReview, userReviews{
			UserId:        userId,
			FilmId:        filmId,
			ReviewContent: reviewContent,
			ReviewValue:   reviewValue,
		})
	}

	filmReviewsDto := usersReviewsDTO{
		FilmsReview: filmsReview,
	}

	jsonResp, err := json.Marshal(filmReviewsDto)

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

func run(addr string) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, "postgres://postgres:1906@postgres_reviews:5432/reviews")

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	grpcConnAuth, err := grpc.Dial(
		"films_auth:44444",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	defer grpcConnAuth.Close()
	authClient := ssov1.NewAuthClient(grpcConnAuth)

	handler := Handler{
		pool:   *pool,
		client: authClient,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handler.indexHandler)
	router.HandleFunc("/login", handler.loginHandler).Methods("POST")
	router.HandleFunc("/signup", handler.signupHandler).Methods("POST")
	router.HandleFunc("/films", handler.filmsHandler).Methods("GET")
	router.HandleFunc("/filmReviews", handler.filmReviewsHandler).Methods("GET")
	router.HandleFunc("/filmReview", handler.filmReviewHandler).Methods("POST")
	router.HandleFunc("/userReviews", handler.userReviewsHandler).Methods("GET")
	http.Handle("/", router)

	fmt.Println("Server is listening on port 8888")
	http.ListenAndServe(addr, nil)
}

func main() {

	run(":8888")

}
