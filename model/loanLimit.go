package model

import "time"

type LoanLimit struct {
	ID 			 string   `json:"id" validation:"required" gorm:"primaryKey;unique;not null;index"`
	Tenor1 		 float64  `json:"tenor_1" validation:"required" gorm:"not null"`
	Tenor2 		 float64  `json:"tenor_2" validation:"required" gorm:"not null"`
	Tenor3 		 float64  `json:"tenor_3" validation:"required" gorm:"not null"`
	Tenor6 		 float64  `json:"tenor_6" validation:"required" gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func GetDefaultLoanLimit() LoanLimit {
	return LoanLimit{
		ID:        "default",
		Tenor1:    10,
		Tenor2:    20,
		Tenor3:    30,
		Tenor6:    40,
	}
}