package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// returns logo based on requested cup id
func GetLogo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	logo, err := ioutil.ReadFile(fmt.Sprintf("assets/Logo%s.png", vars["cupID"]))

	// if requested logo doesnt exist send fallback image
	if err != nil {
		fallback, err := ioutil.ReadFile("assets/fallbackLogo.png")
		if err != nil {
			ERROR(w, http.StatusUnprocessableEntity, nil)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(fallback)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(logo)
}
