package services

import (
	"bandplan/src/models"
	"fmt"

	"github.com/signintech/gopdf"
)

func (s Service) GenerateSetlistPDF(setlist models.Setlist, options models.SetlistPDFOptions) ([]byte, error) {
	fmt.Println("GenerateSetlistPDF ")

	pdf := gopdf.GoPdf{}

	const pageWidth = 612.0
	const pageHeight = 792.0

	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{
			W: pageWidth,
			H: pageHeight,
		},
	})
	pdf.AddPage()

	if err := pdf.AddTTFFont("regular", "static/fonts/roboto-unhinted/Roboto-Regular.ttf"); err != nil {
		return nil, fmt.Errorf("GenerateSetlistPDF - add ttf font: %w", err)
	}

	y := 100.0

	fmt.Println("setlist.songs: ", setlist.Songs)

	if setlist.Songs == nil {
		fmt.Println("error, no songs")
		return []byte{}, fmt.Errorf("GenerateSetlistPDF - no songs in setlist")
	}

	for _, song := range setlist.Songs {
		fmt.Println("\nsong: ", song)
		// Start another page when necessary.
		if y > pageHeight-60 {
			pdf.AddPage()
			y = 100
		}

		width, err := pdf.MeasureTextWidth(song.Song.Title)
		if err != nil {
			return nil, fmt.Errorf("GenerateSetlistPDF - measure text width: %w", err)
		}

		// Half of the remaining width puts the text in the center.
		x := (pageWidth - width) / 2
		pdf.SetXY(x, y)

		if err := pdf.Text(song.Song.Title); err != nil {
			return nil, fmt.Errorf("GenerateSetlistPDF - pdf.Text: %w", err)
		}

		y += 36
	}

	return pdf.GetBytesPdfReturnErr()
}
