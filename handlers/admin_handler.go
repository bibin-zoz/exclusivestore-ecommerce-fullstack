package handlers

import (
	db "ecommercestore/database"
	"ecommercestore/helpers"
	"ecommercestore/models"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/nfnt/resize"
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
	var NameError models.Invalid

	if categoryName == "" {
		NameError.NameError = "Enter valid Category Name"
		c.JSON(http.StatusBadRequest, gin.H{
			"Errors": NameError,
		})
		return
	}

	newCategory := models.Categories{
		CategoryName: categoryName,
		Status:       status,
	}

	result := db.DB.Create(&newCategory)

	if result.Error != nil {
		NameError.NameError = "Category Already Exists"

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  result.Error.Error(),
			"Errors": NameError,
		})
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

	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID := req.ID

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

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category Status Updated Successfully",
		"redirect": "/admin/categories",
	})
}

//Products

func ProductsHandler(c *gin.Context) {
	var products []models.Products
	var category []models.Categories
	var brand []models.Brands
	db.DB.Find(&category)
	db.DB.Find(&brand)
	if err := db.DB.Preload("Variants").Preload("Category").Preload("Images").Preload("Brand").Find(&products).Error; err != nil {
		fmt.Println("failed to load products")
	}

	c.HTML(http.StatusOK, "productmanage.html", gin.H{
		"Products": products,
		"Category": category,
		"Brands":   brand,
	})

}

// add product
func AddProduct(c *gin.Context) {

	var Producterror models.Producterror

	productName := c.PostForm("productName")
	productDetails := c.PostForm("productDetails")
	storage := c.PostForm("storage")
	ram := c.PostForm("ram")
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	categoryID, _ := strconv.Atoi(c.PostForm("categoryID"))
	brandID, _ := strconv.Atoi(c.PostForm("brandID"))

	// validation
	if productName == "" || productDetails == "" || stock <= 0 || price <= 100 {
		if productName == "" {
			Producterror.ProductNameError = "invalid productname"
		}
		if productDetails == "" {
			Producterror.ProductDetailsError = "Details Cannot be empty"
		}
		if stock <= 0 {
			Producterror.StockError = "Enter Valid Stock Number"
		}
		if math.IsNaN(price) || price <= 100 {
			Producterror.PriceError = "Invalid price"
			fmt.Println(Producterror)

		}
		c.JSON(http.StatusBadRequest, gin.H{
			"Errors": Producterror,
		})
		fmt.Println(Producterror)
		return

	}
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	files := c.Request.MultipartForm.File["images"]

	newProduct := models.Products{
		ProductName:    productName,
		ProductDetails: productDetails,
		// Stock:          stock,
		// Price:          price,
		CategoryID: uint(categoryID),
		BrandID:    uint(brandID),
	}

	result := db.DB.Create(&newProduct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	newProductVariant := models.ProductVariants{
		ProductID: newProduct.ID,
		Processor: c.PostForm("processor"), // Assuming you get the processor value from the form
		Storage:   storage,
		Ram:       ram,
		Status:    "listed",
		Stock:     stock, // You can set the default value here or get it from the form
		MaxPrice:  price, // Assuming max price is the same as the regular price for now
	}
	newProductVariant.CreateSlug(productName)

	// Insert the new ProductVariant into the database
	resultVariant := db.DB.Create(&newProductVariant)
	if resultVariant.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resultVariant.Error.Error()})
		return
	}

	// Process each file
	for _, file := range files {
		isValid, detectedType := helpers.IsImageFile(file)
		if !isValid {
			fmt.Println("Unknown format:", detectedType)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid file format. Only image files are allowed."})
			return
		}

		fmt.Println("Valid format:", detectedType)

		// Continue processing the file since it's a valid image
		src, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		newImage := models.Image{
			ProductID: newProduct.ID,
		}

		imageResult := db.DB.Create(&newImage)
		if imageResult.Error != nil {
			fmt.Println("Error creating new image record:", imageResult.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": imageResult.Error.Error()})
			return
		}

		filename := fmt.Sprintf("%d_%s", newImage.ID, file.Filename)
		filepath := filepath.Join("static/uploads", filename)
		filepath = strings.ReplaceAll(filepath, "\\", "/")
		dst, err := os.Create(filepath)
		if err != nil {
			fmt.Println("Error creating destination file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer dst.Close()

		// Resize the image
		img, _, err := image.Decode(src)
		if err != nil {
			fmt.Println("Error decoding image:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Resize to your desired width and height
		resizedImg := resize.Resize(465, 576, img, resize.Lanczos3)

		// Save the resized image
		if err := jpeg.Encode(dst, resizedImg, nil); err != nil {
			fmt.Println("Error encoding image:", err)
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
	if err := db.DB.Preload("Brand").Preload("ProductVariants").Preload("Images").Find(&Product, ID).Error; err != nil {
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
	variantID, _ := strconv.Atoi(c.PostForm("variantID"))
	productName := c.PostForm("productName")
	productDetails := c.PostForm("productDetails")
	storage := c.PostForm("storage")
	ram := c.PostForm("ram")
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)
	categoryID, _ := strconv.ParseUint(c.PostForm("categoryID"), 10, 64)
	fmt.Println("variant", variantID)

	if variantID == 0 {
		fmt.Println("variant zero entered")
		// VariantID is 0, create a new ProductVariants record
		newProductVariant := models.ProductVariants{
			ProductID: uint(id),
			Processor: "tempprocessor", // Assuming you get the processor value from the form
			Storage:   storage,
			Ram:       ram,
			Status:    "listed",
			Stock:     stock, // You can set the default value here or get it from the form
			MaxPrice:  price, // Assuming max price is the same as the regular price for now
		}
		newProductVariant.CreateSlug(productName)

		// Insert the new ProductVariant into the database
		result := db.DB.Create(&newProductVariant)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result.Error != nil {
			fmt.Println("Error creating ProductVariants:", result.Error.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		fmt.Println("variant zero existed")
		c.JSON(http.StatusOK, gin.H{
			"message":  "variant added successfully",
			"redirect": "/admin/product?id=" + idStr,
		})
		return
	}

	result := db.DB.Model(&models.Products{}).Where("id=?", id).Updates(models.Products{
		ProductName:    productName,
		ProductDetails: productDetails,
		CategoryID:     uint(categoryID),
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		log.Println("prodcut model error")

		return
	}

	slugInput := fmt.Sprintf("%s-%s-%s", productName, storage, ram)
	newSlug := slug.MakeLang(slugInput, "en")
	db.DB.Where("slug=?", newSlug)

	result = db.DB.Model(&models.ProductVariants{}).Where("id=?", variantID).Updates(models.ProductVariants{
		Storage: storage,
		Ram:     ram,
		Stock:   stock,
		Price:   price,
		Slug:    newSlug,
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

		// Resize the image
		resizedImage, err := helpers.ResizeImage(src, 500, 500) // Adjust the dimensions as needed
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resizing image: " + err.Error()})
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

		// Save the resized image to the destination file
		if err := helpers.SaveResizedImage(dst, resizedImage, "jpeg"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newImage.FilePath = filepath
		db.DB.Save(&newImage)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Product updated successfully",
		"redirect": "/admin/product?id=" + idStr,
	})
}
