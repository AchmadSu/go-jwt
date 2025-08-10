package controllers

import (
	"net/http"
	"strconv"

	"example.com/m/bootstrap"
	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var input dto.CreateProductInput
	resp := utils.NewResponse()
	message := "Failed to create product"
	if c.Bind(&input) != nil {
		err := errs.New("Body request invalid", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	product, err := bootstrap.ProductService.CreateProduct(c.Request.Context(), &input)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	message = "Product has registered successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(product).
		Send(c)
}

func UpdateProduct(c *gin.Context) {
	var input dto.UpdateProductInput
	resp := utils.NewResponse()
	id := c.Query("id")
	parsedID, err := strconv.Atoi(id)
	message := "Failed to create product"
	if c.Bind(&input) != nil {
		err := errs.New("Body request invalid", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	if err != nil {
		err := errs.New("product ID is not a number", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	product, err := bootstrap.ProductService.UpdateProduct(parsedID, c.Request.Context(), &input)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	message = "Product has updated successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(product).
		Send(c)
}

func GetProducts(c *gin.Context) {
	var pagination dto.PaginationRequest
	resp := utils.NewResponse()
	message := "Failed to get product"

	if err := c.ShouldBindQuery(&pagination); err != nil {
		err = errs.New(utils.GetSafeErrorMessage(err, "Body request pagination invalid"), http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	if pagination.Name != "" || pagination.ID != nil || pagination.Code != "" {
		product, err := bootstrap.ProductService.GetProduct(&pagination)
		if err != nil {
			errResp := utils.PrintErrorResponse(resp, err, message)
			errResp.Send(c)
			return
		}
		message = "Get product by parameter query successfully"
		resp.SetMessage(message).
			SetPayload(product).
			Send(c)
		return
	}

	pg, err := bootstrap.ProductService.GetAllProducts(&pagination)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	message = "Get products successfully"
	resp.SetMessage(message).
		SetPayload(pg.Data).
		SetMeta(gin.H{
			"page":      pg.Page,
			"totalPage": pg.TotalPages,
			"totalData": pg.Total,
		}).
		Send(c)
}
