package helpers

import "github.com/jung-kurt/gofpdf"

func GenerateInvoicePDF() error {

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)

	// Add title
	pdf.Cell(40, 10, "Invoice")

	// Set font for the rest of the document
	pdf.SetFont("Arial", "", 12)

	// Add some content
	content := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque euismod quam ut ex semper gravida. Sed nec libero et ligula elementum bibendum in sit amet dui. Nullam ut sem ut erat iaculis sollicitudin."
	pdf.MultiCell(0.0, 10.0, content, "0", "L", true)

	// Save the PDF to a file
	err := pdf.OutputFileAndClose("invoice.pdf")
	if err != nil {
		return err
	}
	return nil
}
