package main

import (
	"context"
	"encoding/json"
	// "net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Participant struct {
	UserID int     `json:"userId"`
	Amount float64 `json:"amount"` // Used for EXACT split
}

func main() {
	pool := ConnectDB()
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_URL")})

	InitPubSub(pool, rdb)

	r := gin.Default()

	r.POST("/expenses", func(c *gin.Context) {
		var req struct {
			Description  string        `json:"description"`
			Amount       float64       `json:"amount"`
			PayerID      int           `json:"payerId"`
			GroupID      int           `json:"groupId"`
			SplitType    string        `json:"splitType"`
			Participants []Participant `json:"participants"`
		}
		c.BindJSON(&req)

		// Strategy selection
		var strategy SplitStrategy
		if req.SplitType == "EQUAL" {
			strategy = EqualSplit{}
		} else {
			strategy = ExactSplit{}
		}

		shares := strategy.Calculate(req.Amount, req.Participants)
		
		// Save expense and publish
		var id int
		pool.QueryRow(context.Background(), "INSERT INTO expenses (description, amount, payer_id, group_id, split_type) VALUES ($1, $2, $3, $4, $5) RETURNING id", 
			req.Description, req.Amount, req.PayerID, req.GroupID, req.SplitType).Scan(&id)

		msg, _ := json.Marshal(ExpenseCreatedMsg{PayerID: req.PayerID, Shares: shares})
		rdb.Publish(context.Background(), "EXPENSE_CREATED", msg)

		c.JSON(201, gin.H{"message": "Expense added", "id": id})
	})

	r.Run(":3000")
}