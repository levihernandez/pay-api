package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
)

var db *pgxpool.Pool

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err = initDB(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	r.POST("/users", createUser)
	r.GET("/users/:uuid/balances", checkBalances)
	r.POST("/payments", sendPayment)
	r.Run(":8087")
}

/*
type AccountType struct {
	ID          int64     `json:"id"`
	CustomerUUID string    `json:"customer_uuid"`
	AccountType string    `json:"account_type"`
	Balance     float64   `json:"balance"`
}*/

func checkBalances(c *gin.Context) {
	uuid := c.Param("uuid")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	var balances []AccountType
	rows, err := db.Query(context.Background(), "SELECT id, customer_uuid, account_type, balance FROM account_types WHERE customer_uuid = $1", uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balances"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var accountType AccountType
		if err := rows.Scan(&accountType.ID, &accountType.CustomerUUID, &accountType.AccountType, &accountType.Balance); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balances"})
			return
		}
		balances = append(balances, accountType)
	}

	c.JSON(http.StatusOK, gin.H{"balances": balances})
}

/*
type PaymentRequest struct {
	FromUUID string  `json:"from_uuid"`
	ToUUID   string  `json:"to_uuid"`
	Amount   float64 `json:"amount"`
}*/

func sendPayment(c *gin.Context) {
	var payment PaymentRequest
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if payment.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	if payment.FromUUID == payment.ToUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender and receiver cannot be the same"})
		return
	}

	err := makeTransfer(db, payment.FromUUID, payment.ToUUID, payment.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make transfer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}


/*
type Customer struct {
	UUID          string `json:"uuid" db:"uuid"`
	Name          string `json:"name" db:"name"`
	LastName      string `json:"last_name" db:"last_name"`
	State         string `json:"state" db:"state"`
	Country       string `json:"country" db:"country"`
	ZipCode       string `json:"zip_code" db:"zip_code"`
	BankingStatus string `json:"banking_status" db:"banking_status"`
}*/

/*
type UserCreationRequest struct {
	Name          string `json:"name"`
	LastName      string `json:"last_name"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
	BankingStatus string `json:"banking_status"`
}*/


func createUser(c *gin.Context) {
	var request UserCreationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user := &Customer{
		Name:          request.Name,
		LastName:      request.LastName,
		State:         request.State,
		Country:       request.Country,
		ZipCode:       request.ZipCode,
		BankingStatus: request.BankingStatus,
	}

	_, err := db.Exec(context.Background(), "INSERT INTO customers (name, last_name, state, country, zip_code, banking_status) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Name, user.LastName, user.State, user.Country, user.ZipCode, user.BankingStatus)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}



//... (Other functions remain the same)

func loadConfig(filename string) (*DBConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config DBConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	return &config, nil
}


// User represents a user in the system.
type User struct {
	UUID    string
	Name    string
	Balance float64
}

// makeTransfer handles transferring a specified amount of money from one user to another.
func makeTransfer(db *pgxpool.Pool, fromUUID, toUUID string, amount float64) error {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Retrieve sender details.
	var fromUser User
	err = tx.QueryRow(ctx, "SELECT uuid, name, balance FROM customers WHERE uuid = $1 FOR UPDATE", fromUUID).Scan(&fromUser.UUID, &fromUser.Name, &fromUser.Balance)
	if err != nil {
		return fmt.Errorf("failed to retrieve sender details: %v", err)
	}

	// Check if sender has enough balance.
	if fromUser.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// Retrieve receiver details.
	var toUser User
	err = tx.QueryRow(ctx, "SELECT uuid, name, balance FROM customers WHERE uuid = $1 FOR UPDATE", toUUID).Scan(&toUser.UUID, &toUser.Name, &toUser.Balance)
	if err != nil {
		return fmt.Errorf("failed to retrieve receiver details: %v", err)
	}

	// Update balances.
	fromUser.Balance -= amount
	toUser.Balance += amount

	// Update sender balance in database.
	_, err = tx.Exec(ctx, "UPDATE customers SET balance = $1 WHERE uuid = $2", fromUser.Balance, fromUser.UUID)
	if err != nil {
		return fmt.Errorf("failed to update sender balance: %v", err)
	}

	// Update receiver balance in database.
	_, err = tx.Exec(ctx, "UPDATE customers SET balance = $1 WHERE uuid = $2", toUser.Balance, toUser.UUID)
	if err != nil {
		return fmt.Errorf("failed to update receiver balance: %v", err)
	}

	// Record transaction in ledger.
	_, err = tx.Exec(ctx, "INSERT INTO ledger (sender_uuid, receiver_uuid, amount) VALUES ($1, $2, $3)", fromUser.UUID, toUser.UUID, amount)
	if err != nil {
		return fmt.Errorf("failed to record transaction in ledger: %v", err)
	}

	// Commit the transaction.
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
