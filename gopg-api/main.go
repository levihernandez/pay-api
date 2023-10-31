package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"log"
	"net/http"
	"os"
	"time"
	"context"
)

var db *pg.DB

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		User:     config.User,
		Password: config.Password,
		Database: config.Database,
	})
	defer db.Close()

	r := gin.Default()
	r.POST("/users", createUser)
	r.GET("/users/:uuid/balances", checkBalances)
	r.POST("/payments", sendPayment)
	r.Run(":8088")
}

//... (Other functions and types remain the same)

func checkBalances(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	var balances []AccountType
	err := db.Model(&balances).Where("customer_uuid = ?", uuid).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balances"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balances": balances})
}

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

	_, err := db.Model(user).Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Update makeTransfer function as per your business logic using go-pg
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

	err := makeTransfer(payment.FromUUID, payment.ToUUID, payment.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make transfer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}




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


// ... (Other functions remain the same)
func makeTransfer(fromUUID, toUUID string, amount float64) error {
	return db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		// Retrieve sender and receiver account details.
		var fromAccount, toAccount AccountType
		if err := tx.Model(&fromAccount).Where("customer_uuid = ?", fromUUID).For("UPDATE").Select(); err != nil {
			return fmt.Errorf("failed to retrieve sender account: %v", err)
		}

		if err := tx.Model(&toAccount).Where("customer_uuid = ?", toUUID).For("UPDATE").Select(); err != nil {
			return fmt.Errorf("failed to retrieve receiver account: %v", err)
		}

		// Check if sender has enough balance.
		if fromAccount.Balance < amount {
			return fmt.Errorf("insufficient funds")
		}

		// Update balances.
		fromAccount.Balance -= amount
		toAccount.Balance += amount

		// Update sender account in database.
		if _, err := tx.Model(&fromAccount).WherePK().Update(); err != nil {
			return fmt.Errorf("failed to update sender account: %v", err)
		}

		// Update receiver account in database.
		if _, err := tx.Model(&toAccount).WherePK().Update(); err != nil {
			return fmt.Errorf("failed to update receiver account: %v", err)
		}

		// Record transaction in ledger.
		ledgerEntry := Ledger{
			SenderUUID:   fromUUID,
			ReceiverUUID: toUUID,
			Amount:       amount,
			AccountType:  fromAccount.AccountType, // Assuming the AccountType is same for both accounts, adjust if needed
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if _, err := tx.Model(&ledgerEntry).Insert(); err != nil {
			return fmt.Errorf("failed to record transaction in ledger: %v", err)
		}

		return nil
	})
}