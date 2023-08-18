package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth       = 450
	screenHeight      = 800
	backgroundPath    = "img/test-dot-map.png"
	imagePath1        = "img/hitomi-character-dot-1.png"
	imagePath2        = "img/hitomi-character-dot-2.png"
	imagePath3        = "img/hitomi-character-dot-3.png"
	imagePath4        = "img/hitomi-character-dot-2.png"
	nomalEggImagePath = "img/hitomi-tamago-nomal.png"
	rareEggImagePath1 = "img/hitomi-tamago-rare-1.png"
	rareEggImagePath2 = "img/hitomi-tamago-rare-2.png"
	kimiEggImagePath  = "img/hitomi-tamago-kimi.png"
	playerSpeed       = 7
	switchPeriod      = 15 // 画像を切り替える周期
	imageScale        = 4  // 画像のスケール倍率
	imageScale2       = 1
)

type Game struct {
	background      *ebiten.Image
	image1          *ebiten.Image
	image2          *ebiten.Image
	image3          *ebiten.Image
	image4          *ebiten.Image
	currentImage    *ebiten.Image
	nomalEggImage   *ebiten.Image
	rareEggImage1   *ebiten.Image
	rareEggImage2   *ebiten.Image
	kimiEggImage    *ebiten.Image
	currentEggImage *ebiten.Image
	imageWidth      int
	imageHeight     int
	playerX         float64
	playerY         float64
	frameCount      int
	keyPressed      bool
	addEggPosY      float64
	eggs            []Egg
	eggCounter      int
	eggCount        int
	isActive        bool
	canContorol     bool
	collectedStatus bool
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.IsKeyPressed(ebiten.KeyA) ||
		ebiten.IsKeyPressed(ebiten.KeyD) {
		g.keyPressed = true
	} else {
		g.keyPressed = false
	}

	if g.keyPressed {
		g.frameCount++
		if g.frameCount >= switchPeriod {
			g.frameCount = 0
			if g.currentImage == g.image1 {
				g.currentImage = g.image2
			} else if g.currentImage == g.image2 {
				g.currentImage = g.image3
			} else if g.currentImage == g.image3 {
				g.currentImage = g.image4
			} else if g.currentImage == g.image4 {
				g.currentImage = g.image1
			} else {
				g.currentImage = g.image1
			}
		}
	}

	if g.keyPressed {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			g.playerY -= playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			g.playerY += playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.playerX -= playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.playerX += playerSpeed
		}
	}

	// newPlayerX := g.playerX
	// newPlayerY := g.playerY

	// // キャラクターの位置が画面内に収まるように制限
	// newPlayerX = math.Max(newPlayerX, float64(g.imageWidth*imageScale2/2))
	// newPlayerX = math.Min(newPlayerX, float64(2048-g.imageWidth*imageScale2/2))
	// newPlayerY = math.Max(newPlayerY, float64(g.imageHeight*imageScale2/2))
	// newPlayerY = math.Min(newPlayerY, float64(screenHeight-g.imageHeight*imageScale2/2))

	// // キャラクターの位置を更新
	// g.playerX = newPlayerX
	// g.playerY = newPlayerY

	if shouldGenerateEgg() {
		g.eggCounter++
		if g.eggCounter%2 == 0 {
			g.currentEggImage = g.nomalEggImage
			g.isActive = true
			g.canContorol = false
		} else if g.eggCounter%7 == 0 {
			g.currentEggImage = g.rareEggImage1
			g.isActive = true
			g.canContorol = true
		} else {
			g.currentEggImage = g.kimiEggImage
			g.isActive = true
			g.canContorol = false
		}
		egg := generateEggPosition(g.currentEggImage, g.isActive, g.canContorol, g.collectedStatus)
		g.eggs = append(g.eggs, egg)
	}

	if g.currentEggImage != nil {
		g.eggCount++
		if g.eggCount >= switchPeriod {
			g.eggCount = 0
			if g.addEggPosY == 0 {
				g.addEggPosY = 4
			} else {
				g.addEggPosY = 0
			}

			for i := range g.eggs {
				if g.eggs[i].controlStatus {
					if g.eggCount%2 == 0 {
						if g.eggs[i].Eggtype == g.rareEggImage2 {
							g.eggs[i].Eggtype = g.rareEggImage1
						} else {
							g.eggs[i].Eggtype = g.rareEggImage2
						}
					}
				}
			}
		}
	}

	for i := range g.eggs {
		nomalEggSeFlag := false
		// fmt.Println(g.eggs[i].X, g.eggs[i].Y, (g.playerX)-(g.eggs[i].X-screenWidth/2), (g.playerY)-(g.eggs[i].Y-screenHeight/2))
		if !g.eggs[i].collectedStatus && math.Abs((g.playerX)-(g.eggs[i].X-screenWidth/2)) < 60 && math.Abs((g.playerY)-(g.eggs[i].Y-screenHeight/2)) < 60 {
			g.eggs[i].collectedStatus = true
			nomalEggSeFlag = true

			if nomalEggSeFlag {
				err := playSound("get-nomal-egg-se.wav")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	// fmt.Println(g.playerX, g.playerY)

	var remainingEggs []Egg
	for _, egg := range g.eggs {
		if !egg.collectedStatus {
			remainingEggs = append(remainingEggs, egg)
		}
	}
	g.eggs = remainingEggs

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-g.playerX/imageScale2, -g.playerY/imageScale2)
	screen.DrawImage(g.background, op)

	for _, egg := range g.eggs {
		op3 := &ebiten.DrawImageOptions{}
		op3.GeoM.Scale(imageScale*0.7, imageScale*0.7)
		op3.GeoM.Translate(-g.playerX+egg.X, -g.playerY+egg.Y+g.addEggPosY)
		screen.DrawImage(egg.Eggtype, op3)
	}

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(imageScale, imageScale)
	op2.GeoM.Translate(screenWidth/2-52, screenHeight/2-52)
	screen.DrawImage(g.currentImage, op2)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// ebiten.SetFullscreen(true)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Displaying Background with Player")

	// 背景画像のロード
	bgImg, _, err := ebitenutil.NewImageFromFile(backgroundPath)
	if err != nil {
		log.Fatal(err)
	}

	img1, _, err := ebitenutil.NewImageFromFile(imagePath1)
	if err != nil {
		log.Fatal(err)
	}

	img2, _, err := ebitenutil.NewImageFromFile(imagePath2)
	if err != nil {
		log.Fatal(err)
	}

	img3, _, err := ebitenutil.NewImageFromFile(imagePath3)
	if err != nil {
		log.Fatal(err)
	}

	img4, _, err := ebitenutil.NewImageFromFile(imagePath4)
	if err != nil {
		log.Fatal(err)
	}

	egg1, _, err := ebitenutil.NewImageFromFile(nomalEggImagePath)
	if err != nil {
		log.Fatal(err)
	}

	egg2, _, err := ebitenutil.NewImageFromFile(rareEggImagePath1)
	if err != nil {
		log.Fatal(err)
	}

	egg3, _, err := ebitenutil.NewImageFromFile(kimiEggImagePath)
	if err != nil {
		log.Fatal(err)
	}

	egg4, _, err := ebitenutil.NewImageFromFile(rareEggImagePath2)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		background:      bgImg,
		image1:          img1,
		image2:          img2,
		image3:          img3,
		image4:          img4,
		nomalEggImage:   egg1,
		rareEggImage1:   egg2,
		rareEggImage2:   egg4,
		kimiEggImage:    egg3,
		currentImage:    img1,
		currentEggImage: egg1,
		imageWidth:      img1.Bounds().Dx(),
		imageHeight:     img1.Bounds().Dy(),
	}

	// fmt.Printf("Background image size: %dx%d\n", bgImg.Bounds().Dx(), bgImg.Bounds().Dy())

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
