package main

import (
	"github.com/go-pg/pg/v10"
	"time"
)

type Customer struct {
	UUID          string    `json:"uuid" pg:",pk"`
	Name          string    `json:"name"`
	LastName      string    `json:"last_name"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zip_code"`
	BankingStatus string    `json:"banking_status"`
	CreatedAt     time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt     time.Time `json:"updated_at" pg:"default:now()"`
}

type AccountType struct {
	ID           int64     `json:"id" pg:",pk"`
	CustomerUUID string    `json:"customer_uuid"`
	AccountType  string    `json:"account_type"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt    time.Time `json:"updated_at" pg:"default:now()"`
}

type Ledger struct {
	ID           int64     `json:"id" pg:",pk"`
	SenderUUID   string    `json:"sender_uuid"`
	ReceiverUUID string    `json:"receiver_uuid"`
	Amount       float64   `json:"amount"`
	AccountType  string    `json:"account_type"`
	CreatedAt    time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt    time.Time `json:"updated_at" pg:"default:now()"`
}

type UserCreationRequest struct {
	Name          string `json:"name"`
	LastName      string `json:"last_name"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
	BankingStatus string `json:"banking_status"`
}

type PaymentRequest struct {
	FromUUID string  `json:"from_uuid"`
	ToUUID   string  `json:"to_uuid"`
	Amount   float64 `json:"amount"`
}

func (u *Customer) Create(db *pg.DB) error {
	_, err := db.Model(u).Insert()
	return err
}
