package products

import (
	"time"
)

type Category int

const (
	CategoryArabica Category = iota + 1
	CategoryRobusta
	CategoryBlend
	CategoryDecaf
	CategoryEspresso
	CategoryFilter
)

func (c Category) String() string {
	switch c {
	case CategoryArabica:
		return "arabica"
	case CategoryRobusta:
		return "robusta"
	case CategoryBlend:
		return "blend"
	case CategoryDecaf:
		return "decaf"
	case CategoryEspresso:
		return "espresso"
	case CategoryFilter:
		return "filter"
	default:
		return "unknown"
	}
}

type RoastLevel string

const (
	RoastLevelLight    RoastLevel = "light"
	RoastLevelMedium   RoastLevel = "medium"
	RoastLevelDark     RoastLevel = "dark"
	RoastLevelEspresso RoastLevel = "espresso"
)

type Product struct {
	ID          string      `json:"id" gorm:"primaryKey;type:varchar(100)"`
	Name        string      `json:"name" gorm:"not null"`
	NameEn      *string     `json:"name_en,omitempty"`
	Description *string     `json:"description,omitempty"`
	Price       int         `json:"price" gorm:"not null;check:price > 0"`
	ImageURL    *string     `json:"image_url,omitempty"`
	Category    Category    `json:"category" gorm:"type:smallint;not null"`
	InStock     int         `json:"in_stock" gorm:"default:0;check:in_stock >= 0"`
	Weight      *float64    `json:"weight,omitempty"`
	Origin      *string     `json:"origin,omitempty"`
	RoastLevel  *RoastLevel `json:"roast_level,omitempty"`
	Bitterness  *int        `json:"bitterness,omitempty" gorm:"check:bitterness >= 1 AND bitterness <= 5"`
	Body        *int        `json:"body,omitempty" gorm:"check:body >= 1 AND body <= 5"`
	Acidity     *int        `json:"acidity,omitempty" gorm:"check:acidity >= 1 AND acidity <= 5"`
	Sweetness   *int        `json:"sweetness,omitempty" gorm:"check:sweetness >= 1 AND sweetness <= 5"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Product) TableName() string {
	return "products"
}
