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

func CreateStock(c *gin.Context) {
	var input dto.CreateStockInput
	resp := utils.NewResponse()
	message := "Failed to create stock"
	if c.Bind(&input) != nil {
		err := errs.New("Body request invalid", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	stock, err := bootstrap.StockService.CreateStock(c.Request.Context(), &input)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	message = "Stock has been registered successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(stock).
		Send(c)
}

func UpdateStock(c *gin.Context) {
	var input dto.UpdateStockInput
	resp := utils.NewResponse()
	id := c.Query("id")
	parsedID, err := strconv.Atoi(id)
	message := "Failed to update stock"
	if c.Bind(&input) != nil {
		err := errs.New("Body request invalid", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	if err != nil {
		err := errs.New("stock ID is not a number", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	stock, err := bootstrap.StockService.UpdateStock(parsedID, c.Request.Context(), &input)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	message = "Stock has been updated successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(stock).
		Send(c)
}

func GetStocks(c *gin.Context) {
	var pagination dto.PaginationRequest
	resp := utils.NewResponse()
	message := "Failed to get stock"

	if err := c.ShouldBindQuery(&pagination); err != nil {
		err = errs.New(utils.GetSafeErrorMessage(err, "Body request pagination invalid"), http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	if pagination.Name != "" || pagination.ID != nil || pagination.Code != "" {
		stock, err := bootstrap.ProductService.GetProduct(&pagination)
		if err != nil {
			errResp := utils.PrintErrorResponse(resp, err, message)
			errResp.Send(c)
			return
		}
		message = "Get stock by parameter query successfully"
		resp.SetMessage(message).
			SetPayload(stock).
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
