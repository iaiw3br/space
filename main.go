package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math"
	"math/rand"
	"time"
)

var (
	goldColor  = color.RGBA{255, 215, 0, 255}
	greyColor  = color.RGBA{128, 128, 128, 255}
	sizeOfStar = fyne.NewSize(40, 40)
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Planetary System simulation")
	myWindow.Resize(fyne.NewSize(800, 600))

	planetSlider := widget.NewSlider(1, 10)
	startButton := widget.NewButton("Start", func() {
		star := canvas.NewCircle(goldColor)
		star.StrokeWidth = 1
		star.StrokeColor = greyColor
		star.Resize(sizeOfStar)

		centerPosition := fyne.NewPos(400-star.Size().Width/2, 300-star.Size().Height/2)
		star.Move(centerPosition)

		container := container.NewWithoutLayout()
		container.Add(star)

		val := int(planetSlider.Value)
		planets := make([]*canvas.Circle, val)
		for i := 0; i < val; i++ {
			planet := canvas.NewCircle(
				// random color
				color.RGBA{
					uint8(rand.Intn(256)),
					uint8(rand.Intn(256)),
					uint8(rand.Intn(256)),
					255,
				})
			planet.Resize(fyne.NewSize(
				// random size between 10 and 50
				float32(rand.Intn(40)+10),
				float32(rand.Intn(40)+10),
			))
			planets[i] = planet
			container.Add(planet)
		}

		myWindow.SetContent(container)
		myWindow.SetFixedSize(true)

		ticker := time.NewTicker(time.Millisecond * 16) // 60 FPS

		go func() {
			// random radius
			radius := rand.Float64() * 100
			// random rotation
			angleIncrement := rand.Float64() * 2 * math.Pi

			for i := 1; i < len(planets); i++ {
				for {
					validPosition := true
					r := rand.Float64() * 250         // random radius
					a := rand.Float64() * 2 * math.Pi // random rotation
					for j := 0; j < i; j++ {
						if math.Abs(r-radius) < 40 && math.Abs(a-angleIncrement) < math.Pi/4 {
							validPosition = false
							break
						}
					}
					if validPosition {
						radius = r
						angleIncrement = a
						break
					}
				}
			}

			for {
				select {
				case t := <-ticker.C:
					for i, planet := range planets {
						angle := angleIncrement*float64(t.Second()) + float64(i)*math.Pi/2
						x := centerPosition.X + float32(math.Cos(angle)*radius)
						y := centerPosition.Y + float32(math.Sin(angle)*radius)
						planet.Move(fyne.NewPos(x, y))
					}
					canvas.Refresh(container)
				}
			}
		}()
	})

	stopButton := widget.NewButton("Stop", func() {})

	con := container.NewVBox(planetSlider, startButton, stopButton)
	myWindow.SetContent(con)

	myWindow.ShowAndRun()
}
