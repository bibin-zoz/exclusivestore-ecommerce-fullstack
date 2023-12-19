package helpers

import (
	"ecommercestore/models"
	"fmt"
	"time"
)

func UpdateExpiredOffers() {

	var offers []models.CategoryOffer
	for _, offer := range offers {
		if offer.ExpiryAt.Before(time.Now()) {

			fmt.Printf("Offer with ID %d is expired\n", offer.ID)

			offer.Status = "expired"
		}
	}
}
