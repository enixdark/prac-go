package credentials

import (
	"errors"
	"fmt"
)

type ShirtCloner interface {
	GetClone(m int) (ItemInfoGetter, error)
}

type ShirtColor byte
type Shirt struct {
	Price float32
	SKU   string
	Color ShirtColor
}

const (
	White = 1
	Black = 2
	Blue  = 3
)

type ItemInfoGetter interface {
	GetInfo() string
}
type ShirtsCache struct {
	// Clone ShirtCloner
}

func (c *ShirtsCache) GetClone(m int) (ItemInfoGetter, error) {
	switch m {
	case White:
		newItem := *whitePrototype
		return &newItem, nil
	case Black:
		newItem := *blackPrototype
		return &newItem, nil
	case Blue:
		newItem := *bluePrototype
		return &newItem, nil
	default:
		return nil, errors.New("Shirt model not recognized")
	}
}

func GetShirtsCloner() ShirtCloner {
	return whitePrototype
}

var whitePrototype *Shirt = &Shirt{
	Price: 15.00,
	SKU:   "empty",
	Color: White,
}

var blackPrototype *Shirt = &Shirt{
	Price: 15.00,
	SKU:   "empty",
	Color: Black,
}

var bluePrototype *Shirt = &Shirt{
	Price: 15.00,
	SKU:   "empty",
	Color: Blue,
}

func (i *Shirt) GetPrice() float32 {
	return i.Price
}

func (s *Shirt) GetInfo() string {
	return fmt.Sprintf("Shirt with SKU '%s' and Color id %d that costs %f\n", s.SKU, s.Color, s.Price)
}
