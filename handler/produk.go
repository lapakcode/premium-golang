package handler

import (
	"fmt"
	"net/http"
	"premium/helper"
	"premium/produk"
	"premium/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type produkHandler struct {
	service produk.Service
}

func NewProdukHandelr(service produk.Service) *produkHandler {
	return &produkHandler{service}
}

func (h *produkHandler) GetProduks(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	produks, err := h.service.GetProduks(userID)

	if err != nil {
		response := helper.APIResponse("Error get produk", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List produk", http.StatusOK, "error", produk.FormatProduks(produks))
	c.JSON(http.StatusOK, response)
}

func (h *produkHandler) GetProduk(c *gin.Context) {
	var input produk.GetProdukDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Gagal get detail produk", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	produkDetail, err := h.service.GetProdukByID(input)
	if err != nil {
		response := helper.APIResponse("Gagal get detail produk", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Detail produk", http.StatusOK, "sukses", produk.FormatProdukDetail(produkDetail))

	c.JSON(http.StatusOK, response)
	return
}

func (h *produkHandler) CreateProduk(c *gin.Context) {
	var input produk.CreateProdukInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("gagal create produk", http.StatusUnprocessableEntity, "Error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newProduk, err := h.service.CreateProduk(input)
	if err != nil {

		response := helper.APIResponse("gagal create produk", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Sukses create produk", http.StatusOK, "Sukses", produk.FormatProduk(newProduk))
	c.JSON(http.StatusOK, response)

}

func (h *produkHandler) UpdateProduk(c *gin.Context) {
	var inputID produk.GetProdukDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Gagal update produk", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData produk.CreateProdukInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("gagal update produk", http.StatusUnprocessableEntity, "Error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)

	inputData.User = currentUser

	updateProduk, err := h.service.UpdateProduk(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("gagal update produk", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Sukses update produk", http.StatusOK, "Sukses", produk.FormatProduk(updateProduk))
	c.JSON(http.StatusOK, response)
}

func (h *produkHandler) UploadGambar(c *gin.Context) {
	var input produk.CreateGambarProdukInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Gagal upload gambar produk", http.StatusUnprocessableEntity, "Error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Gagal upload gambar produk", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Gagal upload gambar produk", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.service.SaveGambarProduk(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Gagal upload gambar produk", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Sukses upload gambar produk", http.StatusOK, "sukses", data)

	c.JSON(http.StatusOK, response)
	return

}
