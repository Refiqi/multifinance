package model

import "time"

type User struct {
	NIK            string        `json:"nik" validation:"required" gorm:"primaryKey;unique;not null;index:idx_nik,idx_nik_loan_limit_id"`
	FullName       string        `json:"full_name" validation:"required" gorm:"not null"`
	LegalName      string        `json:"legal_name" validation:"required" gorm:"not null"`
	PlaceOfBirth   string        `json:"place_of_birth" validation:"required" gorm:"not null"`
	DateOfBirth    time.Time     `json:"date_of_birth" validation:"required" gorm:"not null"`
	Salary         float64       `json:"salary" validation:"required" gorm:"not null"`
	KTPPhotoURL    string        `json:"ktp_photo_url" validation:"required" gorm:"not null"`
	SelfiePhotoURL string        `json:"selfie_photo_url" validation:"required" gorm:"not null"`
	LoanLimitID    string 		 `json:"loan_limit_id" validation:"required" gorm:"index:idx_nik_loan_limit_id""`
	LoanLimit      LoanLimit 	 `json:"loan_limit" validation:"required" gorm:"foreignKey:LoanLimitID;constraint:OnDelete:CASCADE"`
	Transactions   []Transaction `json:"transactions" validation:"required" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}