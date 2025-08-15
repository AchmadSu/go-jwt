package controllers

import (
	"net/http"

	"example.com/m/bootstrap"
	"example.com/m/dto"
	"example.com/m/errs"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var input dto.CreateOrderInput
	resp := utils.NewResponse()
	message := "Failed to create order"
	if c.Bind(&input) != nil {
		err := errs.New("Body request invalid", http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	order, err := bootstrap.OrderService.CreateOrder(c.Request.Context(), &input)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}
	message = "Order has been registered successfully"
	resp.SetStatus(http.StatusOK).
		SetMessage(message).
		SetPayload(order).
		Send(c)
}

func GetOrders(c *gin.Context) {
	var pagination dto.PaginationRequest
	var orderPagination dto.PaginationOrderRequest
	resp := utils.NewResponse()
	message := "Failed to get order"

	if err := c.ShouldBindQuery(&pagination); err != nil {
		err = errs.New(utils.GetSafeErrorMessage(err, "Body request pagination invalid"), http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	if pagination.ID != nil {
		order, err := bootstrap.OrderService.GetOrder(&pagination)
		if err != nil {
			errResp := utils.PrintErrorResponse(resp, err, message)
			errResp.Send(c)
			return
		}
		message = "Get order by parameter query successfully"
		resp.SetMessage(message).
			SetPayload(order).
			Send(c)
		return
	}

	if err := c.ShouldBindQuery(&orderPagination); err != nil {
		err = errs.New(utils.GetSafeErrorMessage(err, "Body request order pagination invalid"), http.StatusBadRequest)
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	pg, err := bootstrap.OrderService.GetAllOrder(&pagination, &orderPagination)
	if err != nil {
		errResp := utils.PrintErrorResponse(resp, err, message)
		errResp.Send(c)
		return
	}

	message = "Get orders successfully"
	resp.SetMessage(message).
		SetPayload(pg.Data).
		SetMeta(gin.H{
			"page":      pg.Page,
			"totalPage": pg.TotalPages,
			"totalData": pg.Total,
		}).
		Send(c)
}
