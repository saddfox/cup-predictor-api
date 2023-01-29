package models

import (
	"errors"

	"gorm.io/datatypes"
)

// Model for Grand Slam style results
type Format2 struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	User      User
	UserID    uint `json:"-"`
	Cup       Cup
	CupID     uint
	Points    int `json:"-"`
	EndResult bool
	Results   []int          `gorm:"-"` // array of results as sent to api in json
	DbResults datatypes.JSON `json:"-"` // array of results as stored in db
}

// validate winners are in correct slots
func ValidateFormat2(c *Format2) error {
	// we expect 64 results for round 1, 32 for round 2, 16 for round 3,
	// 8 for round 4, 4 for quarters, 2 semis and 1 final. total should be 127 results
	if len(c.Results) != 127 {
		return errors.New("Wrong number of results")
	}
	// check if r1 winners are possible (first slot should be either player 1 or 2 (index 0 or 1))
	for i := 0; i < 64; i++ {
		if !(c.Results[i] == i*2 || c.Results[i] == i*2+1) {
			return errors.New("wrong r1 result")
		}
	}
	// check if r2 winners are possible (first r2 winner should be either first or second r1 winner)
	for i := 0; i < 32; i++ {
		if !(c.Results[i+64] == c.Results[i*2] || c.Results[i+64] == c.Results[i*2+1]) {
			return errors.New("wrong r2 result")

		}
	}
	for i := 0; i < 16; i++ {
		if !(c.Results[i+96] == c.Results[i*2+64] || c.Results[i+96] == c.Results[i*2+1+64]) {
			return errors.New("wrong r3 result")

		}
	}
	for i := 0; i < 8; i++ {
		if !(c.Results[i+112] == c.Results[i*2+96] || c.Results[i+112] == c.Results[i*2+1+96]) {
			return errors.New("wrong r4 result")
		}
	}
	for i := 0; i < 4; i++ {
		if !(c.Results[i+120] == c.Results[i*2+112] || c.Results[i+120] == c.Results[i*2+1+112]) {
			return errors.New("wrong quarter result")
		}
	}
	for i := 0; i < 2; i++ {
		if !(c.Results[i+124] == c.Results[i*2+120] || c.Results[i+124] == c.Results[i*2+1+120]) {
			return errors.New("wrong semi result")
		}
	}
	if !(c.Results[126] == c.Results[124] || c.Results[126] == c.Results[125]) {
		return errors.New("wrong final result")
	}

	return nil
}

// score user results vs final results
func ScoreFormat2(p []int, r []int) int {
	score := 0

	// 1 point for round 1
	for i := 0; i < 64; i++ {
		if p[i] == r[i] {
			score += 1
		}
	}
	// 1 point for round 2
	for i := 64; i < 96; i++ {
		if p[i] == r[i] {
			score += 1
		}
	}
	// 2 points for round 3
	for i := 96; i < 112; i++ {
		if p[i] == r[i] {
			score += 2
		}
	}
	// 3 points for round 4
	for i := 112; i < 120; i++ {
		if p[i] == r[i] {
			score += 3
		}
	}
	// 5 points for quarters
	for i := 120; i < 124; i++ {
		if p[i] == r[i] {
			score += 5
		}
	}
	// 8 points for semis
	for i := 124; i < 126; i++ {
		if p[i] == r[i] {
			score += 8
		}
	}
	// 12 points for winner
	if p[126] == r[126] {
		score += 12
	}

	return score
}
