package creational

import (
	"errors"
	"fmt"
)

type Car interface {
	GetDoors() int
}

type LuxuryCar struct{}

type FamilyCar struct{}

type CarFactory struct{}

const (
	LuxuryCarType = 1
	FamilyCarType = 2
)

func (c *CarFactory) GetVehicle(v int) (Vehicle, error) {
	switch v {
	case LuxuryCarType:
		return new(LuxuryCar), nil
	case FamilyCarType:
		return new(FamilyCar), nil
	default:
		return nil, errors.New(fmt.Sprintf("Vehicle of type %d not recognized\n", v))
	}
}

func (l *LuxuryCar) GetDoors() int {
	return 4
}

func (l *LuxuryCar) GetWheels() int {
	return 4
}

func (l *LuxuryCar) GetSeats() int {
	return 5
}

func (l *FamilyCar) GetDoors() int {
	return 5
}

func (l *FamilyCar) GetWheels() int {
	return 4
}
func (l *FamilyCar) GetSeats() int {
	return 5
}
