package models

type Product struct {
	ID          *string `json:"id,omitempty" bson:"_id,omitempty"`
	Name        *string `json:"title"`
	Category    *string `json:"category"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	Price       *int    `json:"price"`
}
