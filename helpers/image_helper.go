package helpers

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/nfnt/resize"
)

func IsImageFile(fileHeader *multipart.FileHeader) (bool, string) {

	file, err := fileHeader.Open()
	if err != nil {
		return false, ""
	}
	defer file.Close()


	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, ""
	}


	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return false, ""
	}


	fileType := http.DetectContentType(buffer)
	allowedImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/jpg":  true,

		
	}

	if allowedImageTypes[fileType] {
		return true, fileType
	}

	
	return false, fileType
}
func ResizeImage(src io.Reader, width, height uint) (image.Image, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
	return resizedImg, nil
}

func SaveResizedImage(dst io.Writer, resizedImg image.Image, format string) error {
	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(dst, resizedImg, nil)
		if err != nil {
			return fmt.Errorf("error encoding JPEG: %s", err.Error())
		}
		return nil
	case "png":
		return png.Encode(dst, resizedImg)
	case "gif":
		return gif.Encode(dst, resizedImg, nil)
	// Add more formats as needed
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
}
