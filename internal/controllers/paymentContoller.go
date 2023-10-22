package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pattern/internal/models"
	"pattern/internal/repositories"
	"pattern/internal/services"
)

type PaymentController struct {
	paymentService services.Purchase
	userService    repositories.UserRepository
}

func NewPaymentController(paymentService services.Purchase, userService repositories.UserRepository) *PaymentController {
	return &PaymentController{paymentService: paymentService, userService: userService}
}

func (c *PaymentController) MakePayment(ctx *gin.Context) {
	var payment models.Payment
	if err := ctx.BindJSON(&payment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	paymentMethod := ctx.Param("method")
	userEmail, _ := ctx.Get("email")
	email, _ := userEmail.(string)

	user, err := c.userService.GetUserByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	payment.User = *user
	if _, err := c.paymentService.ProcessPayment(&payment, paymentMethod); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newPayment := models.Payment{
		UserID:   user.ID,
		Amount:   payment.Amount,
		Currency: payment.Currency,
		User:     payment.User,
	}

	if err := c.paymentService.CreatePayment(&newPayment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.AddPaymentToUser(user.ID, &newPayment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}
