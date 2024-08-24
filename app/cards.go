package app

type CardService struct {
	repository CardRepository
}

func NewCardService(cr CardRepository) *CardService {
	return &CardService{
		repository: cr,
	}
}

type CardRepository interface {
	GetCardsByUserId(user_id int) ([]Card, error)
	GetTotalBalanceByUserId(user_id int) (BalanceResponse, error)
}

type Card struct {
	Id             int     `json:"card_id"`
	UserId         int     `json:"user_id"`
	CardNumber     string  `json:"card_number"`
	CurrentBalance float64 `json:"current_balance"`
}

type BalanceResponse struct {
	UserId     int     `json:"user_id"`
	TotalValue float64 `json:"total_balance"`
}

func (s *CardService) GetCardsByUserId(userId int) ([]Card, error) {
	cards, err := s.repository.GetCardsByUserId(userId)
	if err != nil {
		return []Card{}, err
	}

	return cards, nil
}

func (s *CardService) GetTotalBalanceByUserId(userId int) (BalanceResponse, error) {
	balance, err := s.repository.GetTotalBalanceByUserId(userId)
	if err != nil {
		return BalanceResponse{}, err
	}

	return balance, nil
}
