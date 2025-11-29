package pdf

import (
	"bytes"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"url_availability/internal/models"
)

// Генерация pdf файла

func Generate(tasks []*models.LinkTask) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	for _, task := range tasks {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, "Task #"+strconv.Itoa(task.TaskNum))
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 10)
		if len(task.Statuses) == 0 {
			pdf.Cell(0, 8, "No results yet")
			pdf.Ln(6)
		} else {
			for link, status := range task.Statuses {
				pdf.Cell(0, 8, link+" - "+string(status))
				pdf.Ln(6)
			}
		}
		pdf.Ln(10)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
