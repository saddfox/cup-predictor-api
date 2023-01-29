package models

import (
	"errors"
	"fmt"
	"sort"
)

// Model for FIFA World Cup style results
type Format1 struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	User      User
	UserID    uint `json:"-"`
	Cup       Cup
	CupID     uint
	Points    int  `json:"-"`
	EndResult bool `json:"-"`
	TeamA1    int
	TeamA2    int
	TeamA3    int
	TeamA4    int
	TeamB1    int
	TeamB2    int
	TeamB3    int
	TeamB4    int
	TeamC1    int
	TeamC2    int
	TeamC3    int
	TeamC4    int
	TeamD1    int
	TeamD2    int
	TeamD3    int
	TeamD4    int
	TeamE1    int
	TeamE2    int
	TeamE3    int
	TeamE4    int
	TeamF1    int
	TeamF2    int
	TeamF3    int
	TeamF4    int
	TeamG1    int
	TeamG2    int
	TeamG3    int
	TeamG4    int
	TeamH1    int
	TeamH2    int
	TeamH3    int
	TeamH4    int
	Round16_1 int
	Round16_2 int
	Round16_3 int
	Round16_4 int
	Round16_5 int
	Round16_6 int
	Round16_7 int
	Round16_8 int
	Quarter1  int
	Quarter2  int
	Quarter3  int
	Quarter4  int
	Semi1     int
	Semi2     int
	Third     int
	Final     int
}

// validates that a group is unique and values are for the right group
func validateGroup(team1 int, team2 int, team3 int, team4 int, sum int) bool {
	g := []int{team1, team2, team3, team4}
	sort.Ints(g)

	// once sorted we iterate through and make sure there are no duplicates
	for i := 0; i < 3; i++ {
		if g[i] == g[i+1] {
			return false
		}
	}

	// we check that the group has the right teams. group A should have teams 1-4, the sum is 10
	if g[0]+g[1]+g[2]+g[3] != sum {
		return false
	}

	return true
}

// validate correct teams are in corect groups and knockout brackets
func ValidateFormat1(c *Format1) error {
	if !validateGroup(c.TeamA1, c.TeamA2, c.TeamA3, c.TeamA4, 10) {
		return errors.New("Group A wrong values")
	}
	if !validateGroup(c.TeamB1, c.TeamB2, c.TeamB3, c.TeamB4, 26) {
		return errors.New("Group B wrong values")
	}
	if !validateGroup(c.TeamC1, c.TeamC2, c.TeamC3, c.TeamC4, 42) {
		return errors.New("Group C wrong values")
	}
	if !validateGroup(c.TeamD1, c.TeamD2, c.TeamD3, c.TeamD4, 58) {
		return errors.New("Group D wrong values")
	}
	if !validateGroup(c.TeamE1, c.TeamE2, c.TeamE3, c.TeamE4, 74) {
		return errors.New("Group E wrong values")
	}
	if !validateGroup(c.TeamF1, c.TeamF2, c.TeamF3, c.TeamF4, 90) {
		return errors.New("Group F wrong values")
	}
	if !validateGroup(c.TeamG1, c.TeamG2, c.TeamG3, c.TeamG4, 106) {
		return errors.New("Group G wrong values")
	}
	if !validateGroup(c.TeamH1, c.TeamH2, c.TeamH3, c.TeamH4, 122) {
		return errors.New("Group H wrong values")
	}

	if c.Round16_1 != c.TeamA1 && c.Round16_1 != c.TeamB2 {
		return errors.New("Round of 16 1 wrong values")
	}
	if c.Round16_2 != c.TeamC1 && c.Round16_2 != c.TeamD2 {
		return errors.New("Round of 16 2 wrong values")
	}
	if c.Round16_3 != c.TeamE1 && c.Round16_3 != c.TeamF2 {
		return errors.New("Round of 16 3 wrong values")
	}
	if c.Round16_4 != c.TeamG1 && c.Round16_4 != c.TeamH2 {
		return errors.New("Round of 16 4 wrong values")
	}
	if c.Round16_5 != c.TeamB1 && c.Round16_5 != c.TeamA2 {
		return errors.New("Round of 16 5 wrong values")
	}
	if c.Round16_6 != c.TeamD1 && c.Round16_6 != c.TeamC2 {
		return errors.New("Round of 16 6 wrong values")
	}
	if c.Round16_7 != c.TeamF1 && c.Round16_7 != c.TeamE2 {
		return errors.New("Round of 16 7 wrong values")
	}
	if c.Round16_8 != c.TeamH1 && c.Round16_8 != c.TeamG2 {
		return errors.New("Round of 16 8 wrong values")
	}

	if c.Quarter1 != c.Round16_1 && c.Quarter1 != c.Round16_2 {
		return errors.New("Quarter 1 wrong values")
	}
	if c.Quarter2 != c.Round16_3 && c.Quarter2 != c.Round16_4 {
		return errors.New("Quarter 2 wrong values")
	}
	if c.Quarter3 != c.Round16_5 && c.Quarter3 != c.Round16_6 {
		return errors.New("Quarter 3 wrong values")
	}
	if c.Quarter4 != c.Round16_7 && c.Quarter4 != c.Round16_8 {
		return errors.New("Quarter 4 wrong values")
	}

	if c.Semi1 != c.Quarter1 && c.Semi1 != c.Quarter2 {
		return errors.New("Semi 1 wrong values")
	}
	if c.Semi2 != c.Quarter3 && c.Semi2 != c.Quarter4 {
		return errors.New("Semi 2 wrong values")
	}

	if c.Third != c.Quarter1 && c.Third != c.Quarter2 && c.Third != c.Quarter3 && c.Third != c.Quarter4 && (c.Third == c.Semi1 || c.Third == c.Semi2) {
		return errors.New("Third wrong values")
	}

	if c.Final != c.Semi1 && c.Final != c.Semi2 {
		return errors.New("Final wrong values")
	}

	return nil
}

// score user results vs final results
func ScoreFormat1(prediction Format1, result Format1) int {
	score := 0
	groupNames := []string{"A", "B", "C", "D", "E", "F", "G", "H"}

	pred := gamesToMap1(prediction)
	res := gamesToMap1(result)

	// handle group stage
	for _, letter := range groupNames {
		for i := 1; i <= 4; i++ {
			r := res[fmt.Sprintf("Team%s%d", letter, i)]
			r1 := res[fmt.Sprintf("Team%s1", letter)]
			r2 := res[fmt.Sprintf("Team%s2", letter)]
			p := pred[fmt.Sprintf("Team%s%d", letter, i)]

			// if prediction exactly matches result
			if p == r {
				score += 1
			}
			// if correctly predicted who will make it out of group
			if p == r1 || p == r2 {
				score += 4
			}
		}
	}

	for i := 1; i <= 8; i++ {
		r := res[fmt.Sprintf("Round16_%d", i)]
		p := pred[fmt.Sprintf("Round16_%d", i)]
		// bonus points if user correctly predicted exact match
		if p == r {
			score += 1
		}
		// points if team pregressed from round of 16
		for j := 1; j <= 8; j++ {
			p2 := pred[fmt.Sprintf("Round16_%d", j)]
			if r == p2 {
				score += 3
			}
		}
	}
	for i := 1; i <= 4; i++ {
		r := res[fmt.Sprintf("Quarter%d", i)]
		p := pred[fmt.Sprintf("Quarter%d", i)]
		// bonus points if user correctly predicted exact match
		if p == r {
			score += 1
		}
		// points if team pregressed from quarter finals
		for j := 1; j <= 4; j++ {
			p2 := pred[fmt.Sprintf("Quarter%d", j)]
			if r == p2 {
				score += 3
			}
		}
	}
	for i := 1; i <= 2; i++ {
		r := res[fmt.Sprintf("Semi%d", i)]
		p := pred[fmt.Sprintf("Semi%d", i)]
		// bonus points if user correctly predicted exact match
		if p == r {
			score += 1
		}
		// points if team pregressed from semi finals
		for j := 1; j <= 2; j++ {
			p2 := pred[fmt.Sprintf("Semi%d", j)]
			if r == p2 {
				score += 3
			}
		}
	}
	// points for predicting third place
	if pred["Third"] == res["Third"] {
		score += 2
	}
	// points for predicting winner
	if pred["Final"] == res["Final"] {
		score += 6
	}
	return score
}

// helper function to map values form format1 object to a map for easier processing
func gamesToMap1(in Format1) map[string]int {
	matchMap := map[string]int{}
	matchMap["TeamA1"] = in.TeamA1
	matchMap["TeamA2"] = in.TeamA2
	matchMap["TeamA3"] = in.TeamA3
	matchMap["TeamA4"] = in.TeamA4
	matchMap["TeamB1"] = in.TeamB1
	matchMap["TeamB2"] = in.TeamB2
	matchMap["TeamB3"] = in.TeamB3
	matchMap["TeamB4"] = in.TeamB4
	matchMap["TeamC1"] = in.TeamC1
	matchMap["TeamC2"] = in.TeamC2
	matchMap["TeamC3"] = in.TeamC3
	matchMap["TeamC4"] = in.TeamC4
	matchMap["TeamD1"] = in.TeamD1
	matchMap["TeamD2"] = in.TeamD2
	matchMap["TeamD3"] = in.TeamD3
	matchMap["TeamD4"] = in.TeamD4
	matchMap["TeamE1"] = in.TeamE1
	matchMap["TeamE2"] = in.TeamE2
	matchMap["TeamE3"] = in.TeamE3
	matchMap["TeamE4"] = in.TeamE4
	matchMap["TeamF1"] = in.TeamF1
	matchMap["TeamF2"] = in.TeamF2
	matchMap["TeamF3"] = in.TeamF3
	matchMap["TeamF4"] = in.TeamF4
	matchMap["TeamG1"] = in.TeamG1
	matchMap["TeamG2"] = in.TeamG2
	matchMap["TeamG3"] = in.TeamG3
	matchMap["TeamG4"] = in.TeamG4
	matchMap["TeamH1"] = in.TeamH1
	matchMap["TeamH2"] = in.TeamH2
	matchMap["TeamH3"] = in.TeamH3
	matchMap["TeamH4"] = in.TeamH4
	matchMap["Round16_1"] = in.Round16_1
	matchMap["Round16_2"] = in.Round16_2
	matchMap["Round16_3"] = in.Round16_3
	matchMap["Round16_4"] = in.Round16_4
	matchMap["Round16_5"] = in.Round16_5
	matchMap["Round16_6"] = in.Round16_6
	matchMap["Round16_7"] = in.Round16_7
	matchMap["Round16_8"] = in.Round16_8
	matchMap["Quarter1"] = in.Quarter1
	matchMap["Quarter2"] = in.Quarter2
	matchMap["Quarter3"] = in.Quarter3
	matchMap["Quarter4"] = in.Quarter4
	matchMap["Semi1"] = in.Semi1
	matchMap["Semi2"] = in.Semi2
	matchMap["Third"] = in.Third
	matchMap["Final"] = in.Final
	return matchMap
}
