package main

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseCreatedMsg struct {
	PayerID int     `json:"payerId"`
	Shares  []Share `json:"shares"`
}

func InitPubSub(db *pgxpool.Pool, rdb *redis.Client) {
	pubsub := rdb.Subscribe(context.Background(), "EXPENSE_CREATED")
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var data ExpenseCreatedMsg
			json.Unmarshal([]byte(msg.Payload), &data)
			processNetting(db, data)
		}
	}()
}

func processNetting(db *pgxpool.Pool, data ExpenseCreatedMsg) {
	ctx := context.Background()
	for _, share := range data.Shares {
		if share.UserID == data.PayerID { continue }

		var oppAmt float64
		err := db.QueryRow(ctx, "SELECT amount FROM balances WHERE user_id = $1 AND owes_to = $2", data.PayerID, share.UserID).Scan(&oppAmt)

		if err == nil { // Opposite debt exists
			if oppAmt > share.Amount {
				db.Exec(ctx, "UPDATE balances SET amount = amount - $1 WHERE user_id = $2 AND owes_to = $3", share.Amount, data.PayerID, share.UserID)
			} else if oppAmt == share.Amount {
				db.Exec(ctx, "DELETE FROM balances WHERE user_id = $1 AND owes_to = $2", data.PayerID, share.UserID)
			} else {
				db.Exec(ctx, "DELETE FROM balances WHERE user_id = $1 AND owes_to = $2", data.PayerID, share.UserID)
				db.Exec(ctx, "INSERT INTO balances (user_id, owes_to, amount) VALUES ($1, $2, $3) ON CONFLICT (user_id, owes_to) DO UPDATE SET amount = balances.amount + EXCLUDED.amount", share.UserID, data.PayerID, share.Amount-oppAmt)
			}
		} else { // No opposite debt
			db.Exec(ctx, "INSERT INTO balances (user_id, owes_to, amount) VALUES ($1, $2, $3) ON CONFLICT (user_id, owes_to) DO UPDATE SET amount = balances.amount + EXCLUDED.amount", share.UserID, data.PayerID, share.Amount)
		}
	}
}