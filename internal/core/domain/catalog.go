package domain

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantCard struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Address     string
	Phone       *string
	IsOpen      bool
	PhotoURL    *string
	CreatedAt   time.Time
}

type MenuDish struct {
	ID                 uuid.UUID
	CategoryID         uuid.UUID
	Name               string
	Description        *string
	PriceKopecks       int64
	CookingTimeMinutes *int
	PhotoURL           *string
}

type MenuCategory struct {
	ID        uuid.UUID
	Name      string
	SortOrder *int
	Dishes    []MenuDish
}

type Menu struct {
	RestID     uuid.UUID
	Categories []MenuCategory
}
