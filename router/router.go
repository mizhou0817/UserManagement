package router

import (
	"github.com/didip/tollbooth"
	"github.com/gorilla/mux"
	"go-postgres/middleware"
	"net/http"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	limiter := tollbooth.NewLimiter(10, nil)
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// get the client IP addr
		clientIP := r.RemoteAddr

		// use the client IP addr as the key of the limiter
		httpError := tollbooth.LimitByKeys(limiter, []string{clientIP})
		if httpError != nil {
			http.Error(w, httpError.Message, httpError.StatusCode)
			return
		}

		middleware.CreateUser(w, r)

	}).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users", middleware.GetAllUser).Methods("GET", "OPTIONS")

	return router
}
