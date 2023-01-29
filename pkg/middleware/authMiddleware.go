package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/saddfox/cup-predictor/pkg/api"
	"github.com/saddfox/cup-predictor/pkg/auth"
	"github.com/saddfox/cup-predictor/pkg/db"
	"github.com/saddfox/cup-predictor/pkg/models"
)

// authorization middleware. handles cors and validates jwt token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			return
		}

		fmt.Println(r.Header)

		ctx := r.Context()
		token := r.Header.Get("Authorization")

		uid, err := auth.ValidateToken(token)
		if err != nil {
			api.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		ctx = context.WithValue(ctx, "uid", uid)

		next(w, r.WithContext(ctx))
	}
}

// admin authorization middleware. handles cors, validates jwt token and checks if user is admin
func AuthMiddlewareAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			return
		}
		fmt.Println(r.Header)

		ctx := r.Context()
		token := r.Header.Get("Authorization")

		uid, err := auth.ValidateToken(token)
		if err != nil {
			api.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		var user models.User
		result := db.DB.Debug().First(&user, uid)

		if result.Error != nil {
			api.ERROR(w, http.StatusUnauthorized, result.Error)
			return
		}

		ctx = context.WithValue(ctx, "uid", uid)
		if user.Admin != true {
			api.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		next(w, r.WithContext(ctx))
	}
}
