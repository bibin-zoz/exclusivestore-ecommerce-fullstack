package helpers

import "github.com/jung-kurt/gofpdf"

func GenerateInvoicePDF() error {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Invoice")

	pdf.SetFont("Arial", "", 12)

	content := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque euismod quam ut ex semper gravida. Sed nec libero et ligula elementum bibendum in sit amet dui. Nullam ut sem ut erat iaculis sollicitudin."
	pdf.MultiCell(0.0, 10.0, content, "0", "L", true)

	err := pdf.OutputFileAndClose("invoice.pdf")
	if err != nil {
		return err
	}
	return nil
}
