package restaurants_transport_http

import (
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type RestaurantResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Address     string    `json:"address"`
	Phone       *string   `json:"phone"`
	IsOpen      bool      `json:"is_open"`
	PhotoURL    *string   `json:"photo_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type MenuDishDTO struct {
	ID                 uuid.UUID `json:"id"`
	CategoryID         uuid.UUID `json:"category_id"`
	Name               string    `json:"name"`
	Description        *string   `json:"description"`
	PriceKopecks       int64     `json:"price_kopecks"`
	CookingTimeMinutes *int      `json:"cooking_time_minutes"`
	PhotoURL           *string   `json:"photo_url"`
}

type MenuCategoryDTO struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	SortOrder *int          `json:"sort_order"`
	Dishes    []MenuDishDTO `json:"dishes"`
}

type MenuResponseDTO struct {
	RestID     uuid.UUID         `json:"restaurant_id"`
	Categories []MenuCategoryDTO `json:"categories"`
}

func restaurantToDTO(rest domain.RestaurantCard) RestaurantResponseDTO {
	return RestaurantResponseDTO{
		ID:          rest.ID,
		Name:        rest.Name,
		Description: rest.Description,
		Address:     rest.Address,
		Phone:       rest.Phone,
		IsOpen:      rest.IsOpen,
		PhotoURL:    rest.PhotoURL,
		CreatedAt:   rest.CreatedAt,
	}
}

func restaurantsToDTO(rests []domain.RestaurantCard) []RestaurantResponseDTO {
	res := make([]RestaurantResponseDTO, len(rests))
	for k, v := range rests {
		res[k] = restaurantToDTO(v)
	}

	return res
}

func menuToDTO(menu domain.Menu) MenuResponseDTO {
	cats := make([]MenuCategoryDTO, len(menu.Categories))
	for k, cat := range menu.Categories {
		dishes := make([]MenuDishDTO, len(cat.Dishes))
		for i, dish := range cat.Dishes {
			dishes[i] = MenuDishDTO{
				ID:                 dish.ID,
				CategoryID:         dish.CategoryID,
				Name:               dish.Name,
				Description:        dish.Description,
				PriceKopecks:       dish.PriceKopecks,
				CookingTimeMinutes: dish.CookingTimeMinutes,
				PhotoURL:           dish.PhotoURL,
			}
		}

		cats[k] = MenuCategoryDTO{
			ID:        cat.ID,
			Name:      cat.Name,
			SortOrder: cat.SortOrder,
			Dishes:    dishes,
		}
	}

	return MenuResponseDTO{
		RestID:     menu.RestID,
		Categories: cats,
	}
}
