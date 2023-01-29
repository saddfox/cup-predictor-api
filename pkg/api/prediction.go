package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saddfox/cup-predictor/pkg/db"
	"github.com/saddfox/cup-predictor/pkg/models"
)

// validate and add user prediction to database
func SubmitPrediction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userID := uint(r.Context().Value("uid").(float64))

	// admin cant submit predictions, only final results via admin endpoint
	if userID == 1 {
		ERROR(w, http.StatusForbidden, errors.New("Use admin endpoint to submit final result"))
		return
	}

	// use different models depending on format
	switch vars["format"] {
	case "1":
		prediction := models.Format1{}

		err = json.Unmarshal(body, &prediction)
		if err != nil {
			ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		// check if user has already submitted prediction for this cup
		if db.DB.Debug().Where("user_id = ? AND cup_id = ?", userID, prediction.CupID).First(&prediction).Error == nil {
			ERROR(w, http.StatusConflict, errors.New("Prediction already exists"))
			return
		}
		// check if submissions have been closed
		if db.DB.Debug().Where("id = ? AND active = ?", prediction.CupID, false).First(&models.Cup{}).Error == nil {
			ERROR(w, http.StatusConflict, errors.New("Submissions have been closed"))
			return
		}

		err2 := models.ValidateFormat1(&prediction)
		if err2 != nil {
			ERROR(w, http.StatusUnprocessableEntity, err2)
			return
		}

		prediction.UserID = userID

		result := db.DB.Create(&prediction)
		if result.Error != nil {
			ERROR(w, http.StatusInternalServerError, result.Error)
			return
		}
		response := response{}
		response.Status = fmt.Sprintf("Created: %d", prediction.ID)
		JSON(w, http.StatusCreated, response)
	case "2":
		prediction := models.Format2{}
		err = json.Unmarshal(body, &prediction)
		if err != nil {
			ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		// check if user has already submitted prediction for this cup
		if db.DB.Debug().Where("user_id = ? AND cup_id = ?", userID, prediction.CupID).First(&prediction).Error == nil {
			ERROR(w, http.StatusConflict, errors.New("Prediction already exists"))
			return
		}
		// check if submissions have been closed
		if db.DB.Debug().Where("id = ? AND active = ?", prediction.CupID, false).First(&models.Cup{}).Error == nil {
			ERROR(w, http.StatusConflict, errors.New("Submissions have been closed"))
			return
		}
		err2 := models.ValidateFormat2(&prediction)
		if err2 != nil {
			ERROR(w, http.StatusUnprocessableEntity, err2)
			return
		}
		prediction.UserID = userID
		dbJson, err3 := json.Marshal(prediction.Results)
		if err3 != nil {
			ERROR(w, http.StatusUnprocessableEntity, err3)
			return
		}
		prediction.DbResults = dbJson
		result := db.DB.Create(&prediction)
		if result.Error != nil {
			ERROR(w, http.StatusInternalServerError, result.Error)
			return
		}
		response := response{}
		response.Status = fmt.Sprintf("Created: %d", prediction.ID)
		JSON(w, http.StatusCreated, response)
	}
}
