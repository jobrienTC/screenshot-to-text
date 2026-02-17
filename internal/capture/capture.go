package capture

import (
	"errors"
	"image"
	"image/draw"

	"github.com/kbinani/screenshot"
)

// GetScreenShot captures all active displays and returns a single stitched image.
func GetScreenShot() (*image.RGBA, error) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return nil, errors.New("no active display found")
	}

	var allBounds image.Rectangle = image.Rect(0, 0, 0, 0)

	// Calculate the bounding rectangle of all screens
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		allBounds = allBounds.Union(bounds)
	}

	// Create a canvas to hold the stitched image
	// We use the Union bounds to handle multi-monitor setups correctly (negative coordinates etc)
	img := image.NewRGBA(allBounds)

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		capturedImg, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, err
		}
		// Draw the captured screen onto the canvas at the correct offset
		draw.Draw(img, bounds, capturedImg, bounds.Min, draw.Src)
	}

	return img, nil
}
