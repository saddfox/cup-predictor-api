package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/saddfox/cup-predictor/pkg/db"
	"github.com/saddfox/cup-predictor/pkg/models"
	"gorm.io/datatypes"
)

// returns teams/players of requested cup
func GetCup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var cup models.Cup

	result := db.DB.First(&cup, vars["cupID"])
	if result.Error != nil {
		ERROR(w, http.StatusUnprocessableEntity, result.Error)
		return
	}

	cupTeams := struct {
		Teams datatypes.JSON `json:"teams"`
	}{
		Teams: cup.Teams,
	}

	JSON(w, http.StatusOK, cupTeams)
}

// returns a user specific response with all cup info, results status and points
func GetAllCups(w http.ResponseWriter, r *http.Request) {
	var cups []models.Cup
	var userPredictions []models.Format1

	// omit Teams field
	result1 := db.DB.Omit("Teams").Find(&cups)
	if result1.Error != nil {
		ERROR(w, http.StatusInternalServerError, result1.Error)
		return
	}
	result2 := db.DB.Debug().Select("UserID", "CupID", "Points").Where("user_id = ?", r.Context().Value("uid")).Find(&userPredictions)
	if result2.Error != nil {
		ERROR(w, http.StatusInternalServerError, result2.Error)
		return
	}

	// if user has submitted prediction for a given cup set submitted = true
	for _, prediction := range userPredictions {
		for i, cup := range cups {
			if prediction.CupID == cup.ID {
				cups[i].Submitted = true
				cups[i].Points = prediction.Points
			}
		}
	}

	JSON(w, http.StatusOK, cups)
}

// mark cup as locked
func LockCup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result2 := db.DB.Debug().Model(&models.Cup{}).Where("ID = ?", vars["cupID"]).Updates(map[string]interface{}{"Active": false})
	if result2.Error != nil {
		ERROR(w, http.StatusInternalServerError, result2.Error)
		return
	}
	JSON(w, http.StatusOK, response{Status: "OK"})
}

func AddCup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	newCup := models.CreateCup{}
	err = json.Unmarshal(body, &newCup)
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err2 := models.ValidateNewCup(newCup.Teams, newCup.Format)
	if err2 != nil {
		ERROR(w, http.StatusUnprocessableEntity, err2)
		return
	}
	cup := models.Cup{}
	cup.Name = newCup.Name
	cup.Active = true
	cup.Format = newCup.Format
	j, err3 := json.Marshal(newCup.Teams)
	if err3 != nil {
		ERROR(w, http.StatusUnprocessableEntity, err3)
		return
	}
	cup.Teams = j
	result := db.DB.Create(&cup)
	if result.Error != nil {
		ERROR(w, http.StatusInternalServerError, result.Error)
		return
	}
	JSON(w, http.StatusCreated, cup)

	downloadLogo(newCup.LogoUrl, cup.ID)
}

// download logo from url, resize it to 500px width and save it as png
func downloadLogo(url string, cupID uint) {
	fmt.Println("downloading logo from: ", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer response.Body.Close()

	file, err := os.Create(fmt.Sprintf("assets/temp-logo%d.png", cupID))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	file.Close()

	original, err := imaging.Open(fmt.Sprintf("assets/temp-logo%d.png", cupID))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	resized := imaging.Resize(original, 500, 0, imaging.Linear)
	out, err := os.Create(fmt.Sprintf("assets/Logo%d.png", cupID))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = imaging.Encode(out, resized, imaging.PNG)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	out.Close()
	e := os.Remove(fmt.Sprintf("assets/temp-logo%d.png", cupID))
	if e != nil {
		fmt.Println(e)
	}

}

type scoreboard struct {
	Username string
	Score    int
}

func GetResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var cup models.Cup

	result1 := db.DB.Select("format").First(&cup, vars["cupID"])
	if result1.Error != nil {
		ERROR(w, http.StatusInternalServerError, result1.Error)
		return
	}

	sb := []scoreboard{}
	switch cup.Format {
	case 1:
		result := db.DB.Debug().Table("format1").Order("format1.points desc").Not("users.id = ?", "1").Joins("JOIN users ON users.id=format1.user_id").Select("users.name as username, format1.points as score").Where("cup_id = ?", vars["cupID"]).Find(&sb)
		if result.Error != nil {
			ERROR(w, http.StatusInternalServerError, result.Error)
			return
		}
	case 2:
		result := db.DB.Debug().Table("format2").Order("format2.points desc").Not("users.id = ?", "1").Joins("JOIN users ON users.id=format2.user_id").Select("users.name as username, format2.points as score").Where("cup_id = ?", vars["cupID"]).Find(&sb)
		if result.Error != nil {
			ERROR(w, http.StatusInternalServerError, result.Error)
			return
		}

	}
	fmt.Println(sb)
	JSON(w, http.StatusOK, sb)

}
