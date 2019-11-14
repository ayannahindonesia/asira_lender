package admin_handlers

import (
	"asira_lender/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetImageB64String(c echo.Context) error {
	defer c.Request().Body.Close()
	// err := validatePermission(c, "core_view_image")
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusForbidden, err, fmt.Sprintf("%s", err))
	// }
	image := models.Image{}

	imageId, _ := strconv.Atoi(c.Param("image_id"))
	err := image.FindbyID(imageId)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Gambar Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, image)
}
