package model

import "time"

type Transaction struct {
	ID 			  string  `json:"id" validation:"required" gorm:"primaryKey;unique;not null;index:idx_transaction_id_user_id"`
	UserID    	  string  `json:"user_id" validation:"required" gorm:"not null;index:idx_transaction_id_user_id"`
	OTR           float64 `json:"otr" validation:"required" gorm:"not null"`
	AdminFee      float64 `json:"admin_fee" validation:"required" gorm:"not null"`
	Installments  int     `json:"installments" validation:"required" gorm:"not null"`
	Interest      float64 `json:"interest" validation:"required" gorm:"not null"`
	AssetName     string  `json:"asset_name" validation:"required" gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}