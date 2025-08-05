package controllers

import (
	"net/http"

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
	product, err := bootstrap.ProductService.Create(c.Request.Context(), &input)
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

func GetProducts(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")
	code := c.Query("code")
	creatorId := c.Query("creator_id")
	modifierId := c.Query("modifier_id")
	resp := utils.NewResponse()
	message := "Failed to fetch product data"
	objectData := map[string]string{
		"id":          id,
		"name":        name,
		"code":        code,
		"creator_id":  creatorId,
		"modifier_id": modifierId,
	}

	if name != "" || id != "" || code != "" {
		product, err := bootstrap.ProductService.GetProduct(objectData)
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

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		err = errs.New(utils.GetSafeErrorMessage(err, "Body request pagination invalid"), http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	pg, err := bootstrap.ProductService.GetAllProducts(&pagination, objectData)
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
