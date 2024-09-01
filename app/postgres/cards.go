package postgres

import (
	"database/sql"

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
	rows, err := cr.db.Query("SELECT id, user_id, card_number, current_balance, TO_CHAR(expiration_date, 'DD/MM') AS expiration_date, card_type FROM cards WHERE user_id = $1", userId)
	if err != nil {
		return []app.Card{}, err
	}
	defer rows.Close()

	var cards []app.Card
	for rows.Next() {
		var card app.Card
		err := rows.Scan(&card.Id, &card.UserId, &card.CardNumber, &card.CurrentBalance, &card.ExpirationDate, &card.CardType)
		if err != nil {
			return []app.Card{}, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (cr *CardRepository) GetTotalBalanceByUserId(userId int) (app.BalanceResponse, error) {
	balance := app.BalanceResponse{
		UserId: userId,
	}
	var totalBalance sql.NullFloat64
	err := cr.db.QueryRow("SELECT SUM(current_balance) AS total_balance FROM cards WHERE user_id = $1", userId).Scan(&totalBalance)
	if err != nil {
		return app.BalanceResponse{}, err
	}

	if totalBalance.Valid {
		balance.TotalValue = totalBalance.Float64
	} else {
		return app.BalanceResponse{}, &app.InvalidUserIdError{UserId: userId}
	}
	return balance, nil
}
