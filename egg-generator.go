package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Egg struct {
	X               float64
	Y               float64
	Eggtype         *ebiten.Image
	activeStatus    bool
	controlStatus   bool
	collectedStatus bool
}

func generateEggPosition(currentEggImage *ebiten.Image, active bool, control bool, collect bool) Egg {
	egg := Egg{
		X:               rand.Float64() * float64(1024),
		Y:               rand.Float64() * float64(1024),
		Eggtype:         currentEggImage,
		activeStatus:    active,
		controlStatus:   control,
		collectedStatus: collect,
	}
	return egg
}

func shouldGenerateEgg() bool {
	eggGenerateFlag := false
	if !gameOverFlag && rand.Intn(100) < 1 {
		if rand.Intn(100) < 50 {
			eggGenerateFlag = true
		}
	}
	return eggGenerateFlag
}
