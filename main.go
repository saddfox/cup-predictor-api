package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/saddfox/cup-predictor/pkg/api"
	"github.com/saddfox/cup-predictor/pkg/db"
	"github.com/saddfox/cup-predictor/pkg/middleware"
)

// setup router routes and middleware
func setupRoutes(r *mux.Router) {
	r.HandleFunc("/api/register", api.CreateUser).Methods("POST")                                                                 // register user
	r.HandleFunc("/api/login", api.LoginUser).Methods("POST")                                                                     // login user
	r.HandleFunc("/api/getCup/{cupID}", middleware.AuthMiddleware(api.GetCup)).Methods("GET", "OPTIONS")                          // get cup teams
	r.HandleFunc("/api/getAllCups", middleware.AuthMiddleware(api.GetAllCups)).Methods("GET", "OPTIONS")                          // get all cups & status
	r.HandleFunc("/api/getLogo/{cupID}", api.GetLogo).Methods("GET")                                                              // get logo
	r.HandleFunc("/api/getResults/{cupID}", middleware.AuthMiddleware(api.GetResults)).Methods("GET", "OPTIONS")                  // get all results for cup
	r.HandleFunc("/api/submitPrediction/{format}", middleware.AuthMiddleware(api.SubmitPrediction)).Methods("POST")               // submit user prediction
	r.HandleFunc("/api/admin/login", api.LoginAdmin).Methods("POST", "OPTIONS")                                                   // admin login
	r.HandleFunc("/api/admin/submitResult/{format}", middleware.AuthMiddlewareAdmin(api.SubmitResult)).Methods("POST", "OPTIONS") // submit admin results
	r.HandleFunc("/api/admin/lock/{cupID}", middleware.AuthMiddlewareAdmin(api.LockCup)).Methods("POST", "OPTIONS")               // lock submissions for cupID
	r.HandleFunc("/api/admin/addCup", middleware.AuthMiddlewareAdmin(api.AddCup)).Methods("POST", "OPTIONS")                      // add new cup
}

// initialize db and seed demo data
func init() {
	// try to read .env file
	// in docker we just use ENV variables and this WILL throw error
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	db.ConnectDB()
	if err := db.SeedUsers(); err != nil {
		fmt.Println(err.Error())
	}
	if err := db.SeedTeams(); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	router := mux.NewRouter()
	setupRoutes(router)

	fmt.Println("Server is ready")
	http.ListenAndServe("0.0.0.0:7000", router)
}
