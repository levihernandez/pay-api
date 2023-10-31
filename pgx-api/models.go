package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type Customer struct {
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	LastName      string    `json:"last_name"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zip_code"`
	BankingStatus string    `json:"banking_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AccountType struct {
	ID           int64     `json:"id"`
	CustomerUUID string    `json:"customer_uuid"`
	AccountType  string    `json:"account_type"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Ledger struct {
	ID           int64     `json:"id"`
	SenderUUID   string    `json:"sender_uuid"`
	ReceiverUUID string    `json:"receiver_uuid"`
	Amount       float64   `json:"amount"`
	AccountType  string    `json:"account_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

func (u *Customer) Create(ctx context.Context, db *pgx.Conn) error {
	_, err := db.Exec(ctx, "INSERT INTO customers (uuid, name, last_name, state, country, zip_code, banking_status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		u.UUID, u.Name, u.LastName, u.State, u.Country, u.ZipCode, u.BankingStatus)
	return err
}
