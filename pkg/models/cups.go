// team IDs to names

package models

import (
	"errors"
	"sort"

	"gorm.io/datatypes"
)

// model for various cups
type Cup struct {
	ID        uint           `gorm:"primary_key;auto_increment" json:"id"`
	Name      string         `json:"name"`
	Format    int            `json:"format"`
	Teams     datatypes.JSON `json:"teams,omitempty"`
	Active    bool           `json:"active"`             // submissions are closed, eg. cup has started
	Results   bool           `json:"results"`            // results are available, eg. cup has finished and admin has entered final results
	Submitted bool           `gorm:"-" json:"submitted"` // only used in api response to indicate user has submited result for this cup
	Points    int            `gorm:"-" json:"points"`    // only used in api response to show user his prediction final score
}

// model for incoming cup creation requests
type CreateCup struct {
	Name    string   `json:"name"`
	Format  int      `json:"format"`
	Teams   []string `json:"teams"`
	LogoUrl string   `json:"logoUrl"`
}

func ValidateNewCup(t []string, format int) error {
	if !(format == 1 || format == 2) {
		return errors.New("Worng format")
	}

	// different formats have different number of teams
	switch format {
	case 1:
		if len(t) != 32 {
			return errors.New("Worng number of teams")
		}
	case 2:
		if len(t) != 128 {
			return errors.New("Worng number of teams")
		}
	}

	// make a copy so we dont mess original order
	t2 := make([]string, len(t))
	copy(t2, t)
	sort.Strings(t2)

	// once sorted we iterate through and make sure there are no duplicates
	for i := 0; i < len(t2)-1; i++ {
		if t2[i] == t2[i+1] {
			return errors.New("Duplicate teams")
		}
	}
	return nil
}
