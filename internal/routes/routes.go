package routes

import "github.com/gin-gonic/gin"

type Routes struct {
	pay  controllers.PaymentController
	auth controllers.AuthController
}

func NewRoutes(pay controllers.PaymentController, auth controllers.AuthController) *Routes {
	return &Routes{pay: pay, auth: auth}
}

func (r *Routes) SetupRouter(router *gin.Engine) *gin.Engine {
	app := router.Group("/app/v1")
	{
		paymentRouter := app.Group("/payments")
		{
			paymentRouter.POST("/make/:method", middlewares.RequireAuthMiddleware, r.pay.MakePayment)
		}

		authRouter := app.Group("/auth")
		{
			authRouter.POST("/register", r.auth.RegisterUserHandler)
			authRouter.POST("/login", r.auth.LoginHandler)
		}

	}
	return router
}
