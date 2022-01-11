package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/teandresmith/restaurant-project/controllers"
)


func InvoiceRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices", controllers.GetInvoices())
	incomingRoutes.GET("/invoice/:invoice_id", controllers.GetInvoice())
	incomingRoutes.POST("/invoices", controllers.CreateInvoice())
	incomingRoutes.PATCH("invoices/:invoice_id", controllers.UpdateInvoice())
	incomingRoutes.DELETE("invoices/:invoice_id", controllers.DeleteInvoice())
}