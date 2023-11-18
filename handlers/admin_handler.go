package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	ID int `json:"id"`
}

// admin
func AdminHome(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "adminhome.html", nil)

}
func AdminLogin(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	c.HTML(http.StatusOK, "adminlogin.html", nil)

}

func AdminLoginPost(c *gin.Context) {

	Newmail := c.Request.FormValue("email")
	Newpassword := c.Request.FormValue("password")
	var compare models.Compare
	var data models.Invalid

	if Newmail == "" {
		data.EmailError = "Email should not be empty"
		c.HTML(http.StatusOK, "login.html", data)
		return
	}
	if Newpassword == "" {
		data.PasswordError = "password should not be empty"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}
	if err := db.DB.Raw("SELECT password, username,role,status FROM users WHERE email=$1", Newmail).Scan(&compare).Error; err != nil {
		fmt.Println(err)
		data.EmailError = "An error occurred while querying the database"
		c.HTML(http.StatusInternalServerError, "login.html", data)
		return
	}

	// Check if no user is found
	var count int64

	if result := db.DB.Model(&models.User{}).Where("email = ? ", Newmail).Count(&count); result.Error != nil || count == 0 {
		data.EmailError = "User not found! Re-check the Mailid"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	}
	if compare.Password != Newpassword {
		data.PasswordError = "check password again"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	}
	if compare.Role == "user" {
		data.RoleError = "click here for admin login -->"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	}
	if compare.Status != "active" {
		data.StatusError = "User is blocked"
		c.HTML(http.StatusBadRequest, "adminlogin.html", data)
		return
	} else {
		helpers.CreateToken(c, compare)
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Expires", "0")
		c.Redirect(http.StatusFound, "/admin/home")
		return

	}
}

// customer
func CustomerHandler(c *gin.Context) {
	// c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	// c.Header("Expires", "0")

	var users []models.User
	db.DB.Where("role=?", "user").Find(&users)

	// Pass data to the template
	c.HTML(http.StatusOK, "customers.html", gin.H{
		"Users": users,
	})

}

func UpdateStatusHandler(c *gin.Context) {

	userID := c.Query("user_id")

	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Status == "active" {
		user.Status = "blocked"
	} else {
		user.Status = "active"
	}

	if err := db.DB.Save(&user).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/customers")
}

func DeleteCustomerHandler(c *gin.Context) {

	var req DeleteRequest
	fmt.Println("sas")

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("sas")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customerID := req.ID

	var user models.User
	result := db.DB.Where("id = ?", customerID).Delete(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

//sellers section

func SellersHandler(c *gin.Context) {
	var sellers []models.User
	seller := "seller"
	db.DB.Where("role=?", seller).Find(&sellers)

	c.HTML(http.StatusOK, "sellers.html", gin.H{
		"Sellers": sellers,
	})
}

// category
func Categoryhandler(c *gin.Context) {
	var category []models.Categories
	db.DB.Find(&category)

	c.HTML(http.StatusOK, "categories.html", gin.H{
		"Category": category,
	})

}
func CategoryPost(c *gin.Context) {
	categoryName := c.PostForm("categoryName")
	status := c.PostForm("status")

	newCategory := models.Categories{
		CategoryName: categoryName,
		Status:       status,
	}

	result := db.DB.Create(&newCategory)

	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category added successfully",
		"category": newCategory,
		"redirect": "/admin/categories",
	})
}

func DeleteCategoryHandler(c *gin.Context) {

	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoryID := req.ID

	var category models.Categories
	result := db.DB.Where("id = ?", categoryID).Delete(&category)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func UpdateCategoryStatus(c *gin.Context) {

	ID := c.Query("id")

	var category models.Categories
	if err := db.DB.Where("id = ?", ID).First(&category).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if category.Status == "listed" {
		category.Status = "unlisted"
	} else {
		category.Status = "listed"
	}

	if err := db.DB.Save(&category).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category status"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/categories")
}

//Products

func ProductsHandler(c *gin.Context) {
	var products []models.Productview
	var category []models.Categories
	db.DB.Find(&category)

	db.DB.Table("products").Select("products.*, categories.category_name").
		Joins("JOIN categories ON products.category_id = categories.id").
		Find(&products)

	c.HTML(http.StatusOK, "productmanage.html", gin.H{
		"Products": products,
		"Category": category,
	})

}

func AddProduct(c *gin.Context) {

	productName := c.PostForm("productName")
	productDetails := c.PostForm("productDetails")
	storage := c.PostForm("storage")
	ram := c.PostForm("ram")
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	categoryID, _ := strconv.Atoi(c.PostForm("categoryID"))

	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	files := c.Request.MultipartForm.File["images"]

	newProduct := models.Products{
		ProductName:    productName,
		ProductDetails: productDetails,
		Storage:        storage,
		Ram:            ram,
		Stock:          stock,
		Price:          price,
		CategoryID:     categoryID,
	}

	result := db.DB.Create(&newProduct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Process each file
	for _, file := range files {
		// Open the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		newImage := models.Image{
			ProductID: newProduct.ID,
		}

		imageResult := db.DB.Create(&newImage)

		if imageResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": imageResult.Error.Error()})
			return
		}

		filename := fmt.Sprintf("%d_%s", newImage.ID, file.Filename)
		filepath := filepath.Join("static/uploads", filename)
		filepath = strings.ReplaceAll(filepath, "\\", "/")
		dst, err := os.Create(filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newImage.FilePath = filepath
		db.DB.Save(&newImage)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Product added successfully",
		"product":  newProduct,
		"redirect": "/admin/products",
	})

}

func UpdateProductStatus(c *gin.Context) {

	ID := c.Query("id")

	var product models.Products
	if err := db.DB.Where("id = ?", ID).First(&product).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if product.Status == "listed" {
		product.Status = "unlisted"
	} else {
		product.Status = "listed"
	}

	if err := db.DB.Save(&product).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product status"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/products")
}

func DeleteProductHandler(c *gin.Context) {

	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productID := req.ID
	fmt.Println(productID)

	var product models.Products

	if err := db.DB.Preload("Images").First(&product, productID).Error; err != nil {
		return
	}

	// Delete the product and its associated images
	if err := db.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func ProductDetailsHandler(c *gin.Context) {
	ID, _ := strconv.Atoi(c.Query("id"))
	var Product models.Products
	var category []models.Categories
	if err := db.DB.Preload("Images").Find(&Product, ID).Error; err != nil {

		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Product not found"})
		return
	}
	if err := db.DB.Find(&category).Error; err != nil {

		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Product not found"})
		return
	}

	log.Println(category)

	c.HTML(http.StatusOK, "productedit.html", gin.H{
		"Product":  Product,
		"Category": category,
	})

	// // Return the result as JSON
	// if c.Request.Header.Get("Accept") == "application/json" {
	// 	// Return JSON if the client accepts JSON
	// 	c.JSON(http.StatusOK, gin.H{"Products": product})
	// } else {
	// 	// Return HTML if the client accepts HTML or doesn't specify a preference
	// 	c.HTML(http.StatusOK, "productedit.html", gin.H{
	// 		"Products": product,
	// 	})
	// }

}

func ProductUpdateHandler(c *gin.Context) {
	// Get product ID from the form data
	idStr := c.PostForm("id")
	id, _ := strconv.Atoi(idStr)

	// Parse other form data
	productName := c.PostForm("productName")
	productDetails := c.PostForm("productDetails")
	storage := c.PostForm("storage")
	ram := c.PostForm("ram")
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	categoryID, _ := strconv.Atoi(c.PostForm("categoryID"))

	// Update product details
	result := db.DB.Model(&models.Products{}).Where("id=?", id).Updates(models.Products{
		ProductName:    productName,
		ProductDetails: productDetails,
		Storage:        storage,
		Ram:            ram,
		Stock:          stock,
		Price:          price,
		CategoryID:     categoryID,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	deleteImageIDs := c.PostFormArray("deleteImages")
	fmt.Println("Delete Image IDs:", deleteImageIDs)

	for _, deleteImageIDStr := range deleteImageIDs {
		fmt.Println("Deleting Image ID:", deleteImageIDStr)
		deleteImageID, err := strconv.Atoi(deleteImageIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
			return
		}

		// Delete the image record from the database
		deleteResult := db.DB.Delete(&models.Image{}, deleteImageID)
		if deleteResult.Error != nil {
			fmt.Println("Error deleting image:", deleteResult.Error.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": deleteResult.Error.Error()})
			return
		}

	}

	// Handle image updates
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	files := c.Request.MultipartForm.File["images"]

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		newImage := models.Image{
			ProductID: uint(id),
		}

		imageResult := db.DB.Create(&newImage)

		if imageResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": imageResult.Error.Error()})
			return
		}

		filename := fmt.Sprintf("%d_%s", newImage.ID, file.Filename)
		filepath := filepath.Join("static/uploads", filename)
		filepath = strings.ReplaceAll(filepath, "\\", "/")
		dst, err := os.Create(filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newImage.FilePath = filepath
		db.DB.Save(&newImage)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Product updated successfully",
		"redirect": "/admin/products",
	})
}

func UploadHandler(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form data"})
		return
	}

	// Extract form values
	productName := c.Request.FormValue("productName")
	productDetails := c.Request.FormValue("productDetails")
	storage := c.Request.FormValue("storage")
	ram := c.Request.FormValue("ram")
	stock := c.Request.FormValue("stock")
	price := c.Request.FormValue("price")

	// Convert stock and price to appropriate types
	stockInt, err := strconv.Atoi(stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock value"})
		return
	}

	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
		return
	}

	// Create a new product
	newProduct := models.Products{
		ProductName:    productName,
		ProductDetails: productDetails,
		Storage:        storage,
		Ram:            ram,
		Stock:          stockInt,
		Price:          priceFloat,
	}

	// Save the product to the database
	if err := db.DB.Create(&newProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create product"})
		return
	}

	// Get the ID of the newly created product
	productID := newProduct.ID

	// Handle file uploads
	files := c.Request.MultipartForm.File["images"]
	for _, file := range files {
		// Save the file to the uploads directory
		filePath := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
			return
		}

		// Replace backslashes with forward slashes in the file path
		filePath = strings.Replace(filePath, `\`, "/", -1)

		// Create a new image record in the database
		newImage := models.Image{
			ProductID: productID,
			FilePath:  filePath,
		}

		// Save the image record to the database
		if err := db.DB.Create(&newImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create image record"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
}

// GetImagesHandler retrieves and displays images from the database
func GetImagesHandler(c *gin.Context) {
	var images []models.Image
	if err := db.DB.Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to query the database"})
		return
	}

	// Replace forward slashes with backslashes in file paths
	fmt.Println("img", images)

	c.HTML(http.StatusOK, "images.html", gin.H{"images": images})
}
