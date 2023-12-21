package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InvoiceHandler(c *gin.Context) {
	// Set the file headers
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("invoice.pdf")
	if err := generateInvoicePDF(); err != nil {
		fmt.Println("Error generating PDF:", err)
		return
	}
}
