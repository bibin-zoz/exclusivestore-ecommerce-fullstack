package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	ID int `json:"id"`
}

func SalesReporthandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "salesreport.html", nil)

}
