package models

type Group struct {
	ID            string `gorm:"primaryKey"`
	Name          string
	Description   string
	Wallet        string
	Members       []Member       `gorm:"foreignKey:GroupID"`
	Contributions []Contribution `gorm:"foreignKey:GroupID"`
}

type Member struct {
	ID      string `gorm:"primaryKey"`
	GroupID string
	Wallet  string
}

type Contribution struct {
	ID       string `gorm:"primaryKey"`
	GroupID  string
	MemberID string
	Amount   float64
}
