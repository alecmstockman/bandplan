package services

import (
	"bandplan/src/models"
	"fmt"

	"github.com/signintech/gopdf"
)

func (s Service) GenerateSetlistPDF(setlist models.Setlist, options models.SetlistPDFOptions) ([]byte, error) {
	fmt.Println("\n-------------------------------------------\nGenerateSetlistPDF ")

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

	y := 100.0

	if err := pdf.AddTTFFont("regular", "static/fonts/roboto-unhinted/Roboto-Regular.ttf"); err != nil {
		return nil, fmt.Errorf("GenerateSetlistPDF - add ttf font: %w", err)
	}

	if err := pdf.AddTTFFont("bold", "static/fonts/roboto-unhinted/Roboto-Bold.ttf"); err != nil {
		return nil, err
	}

	if err := pdf.SetFont("bold", "", 32); err != nil {
		return nil, err
	}

	width, err := pdf.MeasureTextWidth(setlist.Name)
	if err != nil {
		return nil, fmt.Errorf("GenerateSetlistPDF - measure text width: %w", err)
	}

	x := (pageWidth - width) / 2
	pdf.SetXY(x, y)

	if err := pdf.Text(setlist.Name); err != nil {
		return nil, err
	}

	pdf.SetLineWidth(1)
	pdf.Line(x, y+4, x+width, y+4)

	if err := pdf.SetFont("bold", "", 22); err != nil {
		return nil, fmt.Errorf("GenerateSetlistPDF - set font: %w", err)
	}

	if setlist.Songs == nil {
		fmt.Println("error, no songs")
		return []byte{}, fmt.Errorf("GenerateSetlistPDF - no songs in setlist")
	}

	y += 36

	count := 1

	for _, item := range setlist.Songs {

		if y > pageHeight-60 {
			pdf.AddPage()
			y = 100
		}

		if item.ItemType == models.SetlistItemSong {
			song := *item.Song

			if err := pdf.SetFont("bold", "", 24); err != nil {
				return nil, fmt.Errorf("GenerateSetlistPDF - set font: %w", err)
			}

			title := song.Title

			if options.Numbered {
				title = fmt.Sprintf("%d. %s", count, title)
				count++
			}

			width, err := pdf.MeasureTextWidth(title)
			if err != nil {
				return nil, fmt.Errorf("GenerateSetlistPDF - measure text width: %w", err)
			}

			x := (pageWidth - width) / 2
			pdf.SetXY(x, y)

			if err := pdf.Text(title); err != nil {
				return nil, fmt.Errorf("GenerateSetlistPDF - pdf.Text: %w", err)
			}

			y += 36

		} else if item.ItemType == models.SetlistItemTransition && options.Transitions == true {
			transition := *item.Transition

			x = pageWidth

			if err := pdf.SetFont("bold", "", 22); err != nil {
				return nil, err
			}

			boldLabel := fmt.Sprintf("%s ", transition.Title)
			widthBold, err := pdf.MeasureTextWidth(boldLabel)
			if err != nil {
				return nil, err
			}

			regularLabel := fmt.Sprintf("(transition)")
			widthRegular, err := pdf.MeasureTextWidth(regularLabel)

			x = (pageWidth - widthBold - widthRegular) / 2
			pdf.SetXY(x, y)
			if err := pdf.Text(boldLabel); err != nil {
				return nil, err
			}

			if err := pdf.SetFont("regular", "", 20); err != nil {
				return nil, err
			}

			pdf.SetXY(x+widthBold, y)
			if err := pdf.Text(regularLabel); err != nil {
				return nil, err
			}

			y += 36

		} else if item.ItemType == models.SetlistItemBreak && options.Breaks == true {
			breakItem := *item.Break

			if err := pdf.SetFont("regular", "", 20); err != nil {
				return nil, fmt.Errorf("GenerateSetlistPDF - set font: %w", err)
			}

			breakName := fmt.Sprintf("* %s *", breakItem.Title)

			width, err := pdf.MeasureTextWidth(breakName)
			if err != nil {
				return nil, fmt.Errorf("GenerateSetlistPDF - measure text width: %w", err)
			}

			x := (pageWidth - width) / 2
			pdf.SetXY(x, y)

			if err := pdf.Text(breakName); err != nil {
				return nil, fmt.Errorf("GenerateSetlistPDF - pdf.Text: %w", err)
			}

			y += 36
		}
	}

	return pdf.GetBytesPdfReturnErr()
}
