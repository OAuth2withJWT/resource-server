package postgres

import (
	"database/sql"
	"log"

	"github.com/OAuth2withJWT/resource-server/app"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func (cr *CardRepository) GetCardsByUserId(userId int) ([]app.Card, error) {
	rows, err := cr.db.Query("SELECT * FROM cards WHERE user_id = $1", userId)
	if err != nil {
		return []app.Card{}, err
	}

	var cards []app.Card
	for rows.Next() {
		var card app.Card
		err := rows.Scan(&card.Id, &card.UserId, &card.CardNumber, &card.CurrentBalance)
		if err != nil {
			log.Fatal(err)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (cr *CardRepository) GetTotalBalanceByUserId(userId int) (app.BalanceResponse, error) {
	balance := app.BalanceResponse{
		UserId: userId,
	}
	err := cr.db.QueryRow("SELECT SUM(current_balance) AS total_balance FROM cards WHERE user_id = $1", userId).Scan(&balance.TotalBalance)
	if err != nil {
		return app.BalanceResponse{}, err
	}

	return balance, nil
}
