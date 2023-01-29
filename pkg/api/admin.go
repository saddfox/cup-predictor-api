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

// validate and add final result to database
func SubmitResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userID := uint(r.Context().Value("uid").(float64))

	switch vars["format"] {
	case "1":
		prediction := models.Format1{}

		err = json.Unmarshal(body, &prediction)
		if err != nil {
			ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		// check if user has already submitted prediction for this cup
		if db.DB.Where("user_ID = ? AND cup_ID = ?", userID, prediction.CupID).First(&prediction).Error == nil {
			ERROR(w, http.StatusConflict, errors.New("Prediction already exists"))
			return
		}

		err2 := models.ValidateFormat1(&prediction)
		if err2 != nil {
			fmt.Println("validation error")
			ERROR(w, http.StatusUnprocessableEntity, err2)
			return
		}

		prediction.UserID = userID
		prediction.EndResult = true

		result1 := db.DB.Create(&prediction)
		if result1.Error != nil {
			ERROR(w, http.StatusInternalServerError, result1.Error)
			return
		}

		// mark cup as finished
		result2 := db.DB.Debug().Model(&models.Cup{}).Where("ID = ?", prediction.CupID).Updates(map[string]interface{}{"Active": false, "Results": true})
		if result2.Error != nil {
			ERROR(w, http.StatusInternalServerError, result2.Error)
			return
		}

		JSON(w, http.StatusOK, prediction.ID)
		processPredictions1(prediction)
	case "2":
		prediction := models.Format2{}

		err = json.Unmarshal(body, &prediction)
		if err != nil {
			fmt.Println("unmarshal error")
			ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		// check if user has already submitted prediction for this cup
		if db.DB.Where("user_ID = ? AND cup_ID = ?", userID, prediction.CupID).First(&prediction).Error == nil {
			ERROR(w, http.StatusConflict, errors.New("Prediction already exists"))
			return
		}

		err2 := models.ValidateFormat2(&prediction)
		if err2 != nil {
			fmt.Println("validation error")
			ERROR(w, http.StatusUnprocessableEntity, err2)
			return
		}

		prediction.UserID = userID
		prediction.EndResult = true

		dbJson, err3 := json.Marshal(prediction.Results)
		if err3 != nil {
			ERROR(w, http.StatusUnprocessableEntity, err3)
			return
		}
		prediction.DbResults = dbJson

		result1 := db.DB.Create(&prediction)
		if result1.Error != nil {
			ERROR(w, http.StatusInternalServerError, result1.Error)
			return
		}

		// mark cup as finished
		result2 := db.DB.Debug().Model(&models.Cup{}).Where("ID = ?", prediction.CupID).Updates(map[string]interface{}{"Active": false, "Results": true})
		if result2.Error != nil {
			ERROR(w, http.StatusInternalServerError, result2.Error)
			return
		}
		JSON(w, http.StatusOK, prediction.ID)
		processPredictions2(prediction)
	}
}

// calculate scores for predictions
func processPredictions1(result models.Format1) {
	var predictions []models.Format1

	result1 := db.DB.Where("cup_id = ?", result.CupID).Find(&predictions)
	if result1.Error != nil {
		fmt.Println(result1.Error.Error())
		return
	}

	for _, prediction := range predictions {
		score := models.ScoreFormat1(prediction, result)
		db.DB.Model(&prediction).Update("points", score)
		fmt.Println("updated prediction: ", prediction.ID, "score", score)
	}

}

func processPredictions2(result models.Format2) {
	var predictions []models.Format2

	result1 := db.DB.Where("cup_id = ?", result.CupID).Find(&predictions)
	if result1.Error != nil {
		fmt.Println(result1.Error.Error())
		return
	}

	for _, prediction := range predictions {
		var r []int
		err := json.Unmarshal(prediction.DbResults, &r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		score := models.ScoreFormat2(r, result.Results)
		db.DB.Model(&prediction).Update("points", score)
		fmt.Println("updated prediction: ", prediction.ID, "score", score)
	}

}
