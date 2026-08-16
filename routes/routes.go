package routes

import (
	"net/http"

	"webApps_order-management/handler"
	"webApps_order-management/middleware"

	"github.com/gin-gonic/gin"
)

// Handlers bundles every page handler the router needs.
type Handlers struct {
	Auth      *handler.AuthHandler
	Category  *handler.CategoryHandler
	Product   *handler.ProductHandler
	Customer  *handler.CustomerHandler
	Employee  *handler.EmployeeHandler
	Supplier  *handler.SupplierHandler
	Order     *handler.OrderHandler
	OrderSale *handler.OrderSaleHandler
}

// SetupRouter wires every page route.
//
// Role hierarchy mirrors the backend API:
//
//	ADMIN   : full access everywhere
//	MANAGER : create/update/view Orders & Order Sales, full access to master data
//	USER    : view Orders & Order Sales only, full access to master data
//
// (User account management itself is an ADMIN-only capability on the
// backend and isn't exposed as a page in this web app.)
func SetupRouter(jwtSecret, cookieName string, h *Handlers) *gin.Engine {
	r := gin.Default()

	// Custom loader: every template is addressed by its path relative to
	// template/ (e.g. "category_template/create_update.html"), which keeps
	// same-named files across modules (create_update.html appears in every
	// module folder) from colliding the way html/template's default
	// basename-only naming would.
	r.HTMLRender = loadTemplates("template")

	// ---- Public routes ----
	r.GET("/login", middleware.OptionalAuth(jwtSecret, cookieName), h.Auth.ShowLogin)
	r.POST("/login", h.Auth.Login)
	r.GET("/register", middleware.OptionalAuth(jwtSecret, cookieName), h.Auth.ShowRegister)
	r.POST("/register", h.Auth.Register)
	r.POST("/logout", h.Auth.Logout)

	// ---- Everything below requires a signed-in session ----
	protected := r.Group("")
	protected.Use(middleware.JWTAuth(jwtSecret, cookieName))

	protected.GET("/", func(c *gin.Context) {
		data := gin.H{"Title": "Home"}
		if v, ok := c.Get(middleware.ContextUsernameKey); ok {
			data["Username"] = v
		}
		if v, ok := c.Get(middleware.ContextRoleKey); ok {
			data["Role"] = v
		}
		c.HTML(http.StatusOK, "layout/home.html", data)
	})

	// Master data — any authenticated role.
	categories := protected.Group("/categories")
	{
		categories.GET("", h.Category.List)
		categories.GET("/create", h.Category.ShowCreate)
		categories.POST("/create", h.Category.Create)
		categories.GET("/:id", h.Category.Detail)
		categories.GET("/:id/edit", h.Category.ShowEdit)
		categories.POST("/:id/edit", h.Category.Update)
		categories.POST("/:id/delete", h.Category.Delete)
	}

	products := protected.Group("/products")
	{
		products.GET("", h.Product.List)
		products.GET("/create", h.Product.ShowCreate)
		products.POST("/create", h.Product.Create)
		products.GET("/:productNumber/:productID", h.Product.Detail)
		products.GET("/:productNumber/:productID/edit", h.Product.ShowEdit)
		products.POST("/:productNumber/:productID/edit", h.Product.Update)
		products.POST("/:productNumber/:productID/delete", h.Product.Delete)
	}

	customers := protected.Group("/customers")
	{
		customers.GET("", h.Customer.List)
		customers.GET("/create", h.Customer.ShowCreate)
		customers.POST("/create", h.Customer.Create)
		customers.GET("/:custNumber/:custID", h.Customer.Detail)
		customers.GET("/:custNumber/:custID/edit", h.Customer.ShowEdit)
		customers.POST("/:custNumber/:custID/edit", h.Customer.Update)
		customers.POST("/:custNumber/:custID/delete", h.Customer.Delete)
	}

	employees := protected.Group("/employees")
	{
		employees.GET("", h.Employee.List)
		employees.GET("/create", h.Employee.ShowCreate)
		employees.POST("/create", h.Employee.Create)
		employees.GET("/:employeeID/:cardNumber", h.Employee.Detail)
		employees.GET("/:employeeID/:cardNumber/edit", h.Employee.ShowEdit)
		employees.POST("/:employeeID/:cardNumber/edit", h.Employee.Update)
		employees.POST("/:employeeID/:cardNumber/delete", h.Employee.Delete)
	}

	suppliers := protected.Group("/suppliers")
	{
		suppliers.GET("", h.Supplier.List)
		suppliers.GET("/create", h.Supplier.ShowCreate)
		suppliers.POST("/create", h.Supplier.Create)
		suppliers.GET("/:supplierNumber/:supplierID", h.Supplier.Detail)
		suppliers.GET("/:supplierNumber/:supplierID/edit", h.Supplier.ShowEdit)
		suppliers.POST("/:supplierNumber/:supplierID/edit", h.Supplier.Update)
		suppliers.POST("/:supplierNumber/:supplierID/delete", h.Supplier.Delete)
	}

	// Orders — MANAGER/ADMIN write, everyone (incl. USER) reads.
	// The backend API enforces this too; RequireRoles here just keeps the
	// UI from offering actions the API would reject anyway.
	orders := protected.Group("/orders")
	{
		orders.GET("", h.Order.List)
		orders.GET("/:orderID/:orderNo", h.Order.Detail)

		write := orders.Group("")
		write.Use(middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleManager))
		{
			write.GET("/create", h.Order.ShowCreate)
			write.POST("/create", h.Order.Create)
			write.GET("/:orderID/:orderNo/edit", h.Order.ShowEdit)
			write.POST("/:orderID/:orderNo/edit", h.Order.Update)
			write.POST("/:orderID/:orderNo/delete", h.Order.Delete)
			write.POST("/:orderID/:orderNo/details", h.Order.AddDetail)
			write.POST("/:orderID/:orderNo/details/:orderDetailNo/delete", h.Order.DeleteDetail)
		}
	}

	orderSales := protected.Group("/order-sales")
	{
		orderSales.GET("", h.OrderSale.List)
		orderSales.GET("/:orderSaleNo/:orderSaleID", h.OrderSale.Detail)

		write := orderSales.Group("")
		write.Use(middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleManager))
		{
			write.GET("/create", h.OrderSale.ShowCreate)
			write.POST("/create", h.OrderSale.Create)
			write.GET("/:orderSaleNo/:orderSaleID/edit", h.OrderSale.ShowEdit)
			write.POST("/:orderSaleNo/:orderSaleID/edit", h.OrderSale.Update)
			write.POST("/:orderSaleNo/:orderSaleID/delete", h.OrderSale.Delete)
			write.POST("/:orderSaleNo/:orderSaleID/details", h.OrderSale.AddDetail)
			write.POST("/:orderSaleNo/:orderSaleID/details/:orderSaleDetailNo/:orderSaleDetailID/delete", h.OrderSale.DeleteDetail)
		}
	}

	return r
}
