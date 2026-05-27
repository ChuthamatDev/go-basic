package controllers

import (
	"math/rand"

	m "go-fiber-test/models"
)

func buildDogsResult(dogs []m.Dogs) m.ResultDataV3 {
	result := m.ResultDataV3{
		Data:       make([]m.DogsRes, 0, len(dogs)),
		Name:       "golang-test",
		Count:      len(dogs),
		SumRed:     0,
		SumGreen:   0,
		SumPink:    0,
		SumNoColor: 0,
	}

	for _, dog := range dogs {
		color := resolveDogColor(dog.DogID)

		switch color {
		case "red":
			result.SumRed++
		case "green":
			result.SumGreen++
		case "pink":
			result.SumPink++
		default:
			result.SumNoColor++
		}

		result.Data = append(result.Data, m.DogsRes{
			Name:  dog.Name,
			DogID: dog.DogID,
			Type:  color,
			Color: color,
		})
	}

	return result
}

func resolveDogColor(dogID int) string {
	switch {
	case dogID >= 10 && dogID <= 50:
		return "red"
	case dogID >= 100 && dogID <= 150:
		return "green"
	case dogID >= 200 && dogID <= 250:
		return "pink"
	default:
		return "no color"
	}
}

func generateDogID(color string) int {
	switch color {
	case "red":
		return rand.Intn(41) + 10
	case "green":
		return rand.Intn(51) + 100
	case "pink":
		return rand.Intn(51) + 200
	default:
		return rand.Intn(9) + 1
	}
}
