package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SelectionApp struct {
	screenshot     *ebiten.Image
	startX, startY int
	endX, endY     int
	isDragging     bool
	selectionDone  bool
	resultImage    image.Image
}

func NewSelectionApp(screenshotImg image.Image) *SelectionApp {
	return &SelectionApp{
		screenshot: ebiten.NewImageFromImage(screenshotImg),
	}
}

func (a *SelectionApp) Update() error {
	// Exit on ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if a.isDragging {
			// Cancel current selection
			a.isDragging = false
			a.startX, a.startY = 0, 0
			a.endX, a.endY = 0, 0
			return nil
		}
		return ebiten.Termination
	}

	// Mouse Logic
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !a.isDragging {
			a.isDragging = true
			a.startX, a.startY = x, y
		}
		a.endX, a.endY = x, y
	} else if a.isDragging {
		// Mouse released
		a.isDragging = false
		a.selectionDone = true
		a.cropImage()
		return ebiten.Termination
	}

	return nil
}

func (a *SelectionApp) cropImage() {
	// Normalize coordinates
	x0, y0, x1, y1 := a.startX, a.startY, a.endX, a.endY
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	// Ensure bounds are within the screenshot dimensions
	sW, sH := a.screenshot.Bounds().Dx(), a.screenshot.Bounds().Dy()
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x1 > sW {
		x1 = sW
	}
	if y1 > sH {
		y1 = sH
	}

	width := x1 - x0
	height := y1 - y0

	if width <= 0 || height <= 0 {
		return // Invalid selection
	}

	rect := image.Rect(x0, y0, x1, y1)
	sub := a.screenshot.SubImage(rect).(*ebiten.Image)

	// Create a standard Go image to hold the pixels
	// This must be done BEFORE the Ebiten context is destroyed (i.e. before RunGame returns)
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	// ReadPixels reads the pixels from the image to the byte slice.
	// It expects the slice to be big enough.
	sub.ReadPixels(rgba.Pix)

	a.resultImage = rgba
}

func (a *SelectionApp) Draw(screen *ebiten.Image) {
	// Draw the screenshot
	screen.DrawImage(a.screenshot, nil)

	// Draw the semi-transparent overlay
	// We want to darken everything EXCEPT the selection.
	// Ebiten doesn't have a simple "inverse clip", so we draw 4 rectangles around the selection
	// OR we draw a full overlay and "cut out" the hole.

	overlayColor := color.RGBA{0, 0, 0, 150} // Dark overlay

	sW, sH := screen.Bounds().Dx(), screen.Bounds().Dy()

	if !a.isDragging {
		// Draw full overlay if not dragging (or maybe waiting for start)
		vector.DrawFilledRect(screen, 0, 0, float32(sW), float32(sH), overlayColor, false)
		return
	}

	// Normalize coords
	x0, y0, x1, y1 := float32(a.startX), float32(a.startY), float32(a.endX), float32(a.endY)
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	// Draw 4 rectangles to exclude the selection area
	// Top
	vector.DrawFilledRect(screen, 0, 0, float32(sW), y0, overlayColor, false)
	// Bottom
	vector.DrawFilledRect(screen, 0, y1, float32(sW), float32(sH)-y1, overlayColor, false)
	// Left (between top/bottom)
	vector.DrawFilledRect(screen, 0, y0, x0, y1-y0, overlayColor, false)
	// Right (between top/bottom)
	vector.DrawFilledRect(screen, x1, y0, float32(sW)-x1, y1-y0, overlayColor, false)

	// Draw selection border
	vector.StrokeRect(screen, x0, y0, x1-x0, y1-y0, 2, color.White, false)
}

func (a *SelectionApp) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (a *SelectionApp) Result() image.Image {
	return a.resultImage
}
