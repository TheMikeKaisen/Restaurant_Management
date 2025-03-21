package routes

import ("github.com/gin-gonic/gin"
"github.com/TheMikeKaisen/Restaurant_Management/controllers")



func InvoiceRoutes(invoiceRoute *gin.Engine){
	invoiceRoute.GET("/invoice", controllers.GetInvoices())
	invoiceRoute.GET("/invoice/:invoice_id", controllers.GetInvoice())
	invoiceRoute.POST("/invoice", controllers.CreateInvoice())
	invoiceRoute.PATCH("invoice/:food_id", controllers.UpdateInvoice())
}


