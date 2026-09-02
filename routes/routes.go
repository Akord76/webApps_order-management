package routes

import (
	"html/template"
	"net/http"

	"webApps_order-management/handler"
	"webApps_order-management/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Handlers bundles every page handler the router needs.
type Handlers struct {
	Auth                 *handler.AuthHandler
	User                 *handler.UserHandler
	Category             *handler.CategoryHandler
	Product              *handler.ProductHandler
	Customer             *handler.CustomerHandler
	Employee             *handler.EmployeeHandler
	Supplier             *handler.SupplierHandler
	Order                *handler.OrderHandler
	OrderDo              *handler.OrderDoHandler
	ConfigCompanyProfile *handler.ConfigCompanyProfileHandler
	SharingProfit        *handler.SharingProfitHandler
	CommitmentFee        *handler.CommitmentFeeHandler
}

// SetupRouter wires every page route.
//
// Role hierarchy mirrors the backend API:
//
//	ADMIN   : full access everywhere, including User Management
//	MANAGER : create/update/view Orders & Order Sales, full access to master data
//	USER    : view Orders & Order Sales only, full access to master data
func SetupRouter(jwtSecret, cookieName string, h *Handlers) *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static") // rute URL "/static" diarahkan ke direktori "./static"

	// 1. WAJIB: Daftarkan FuncMap SEBELUM LoadHTMLGlob!
	r.SetFuncMap(template.FuncMap{
		"formatOrderDoPrice": func(amount float64) string {
			p := message.NewPrinter(language.Indonesian)
			return p.Sprintf("%.0f", amount)
		},
		// Tambahkan helper perkalian jika diperlukan di template
		"mul": func(a int, b float64) float64 {
			return float64(a) * b
		},
	})

	r.HTMLRender = loadTemplates("template")

	r.Run(":8083")
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

	// User Management — ADMIN only.
	users := protected.Group("/users")
	users.Use(middleware.RequireRoles(middleware.RoleAdmin))
	{
		users.GET("", h.User.List)
		users.GET("/create", h.User.ShowCreate)
		users.POST("/create", h.User.Create)
		users.GET("/:id", h.User.Detail)
		users.GET("/:id/edit", h.User.ShowEdit)
		users.POST("/:id/edit", h.User.Update)
		users.POST("/:id/delete", h.User.Delete)
	}

	// Company Profile Settings — ADMIN only.
	configProfiles := protected.Group("/config-company-profiles")
	configProfiles.Use(middleware.RequireRoles(middleware.RoleAdmin))
	{
		configProfiles.GET("", h.ConfigCompanyProfile.List)
		configProfiles.GET("/create", h.ConfigCompanyProfile.ShowCreate)
		configProfiles.POST("/create", h.ConfigCompanyProfile.Create)
		configProfiles.GET("/:id", h.ConfigCompanyProfile.Detail)
		configProfiles.GET("/:id/edit", h.ConfigCompanyProfile.ShowEdit)
		configProfiles.POST("/:id/edit", h.ConfigCompanyProfile.Update)
		configProfiles.POST("/:id/delete", h.ConfigCompanyProfile.Delete)
	}

	commitmentFee := protected.Group("/commitmentfees")
	commitmentFee.Use(middleware.RequireRoles(middleware.RoleAdmin))
	{
		// Autocomplete endpoints
		commitmentFee.GET("/employees/search", h.CommitmentFee.EmployeeAutocomplete)
		commitmentFee.GET("/products/search", h.CommitmentFee.ProductAutocomplete)

		// CRUD endpoints
		commitmentFee.GET("", h.CommitmentFee.List)
		commitmentFee.GET("/create", h.CommitmentFee.ShowCreate)
		commitmentFee.POST("/create", h.CommitmentFee.Create)
		commitmentFee.GET("/:comID", h.CommitmentFee.Detail)
		commitmentFee.GET("/:comID/edit", h.CommitmentFee.ShowEdit)
		commitmentFee.POST("/:comID/edit", h.CommitmentFee.Update)
		commitmentFee.POST("/:comID/delete", h.CommitmentFee.Delete)
	}

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

	orderDos := protected.Group("/orderdo")
	{
		orderDos.GET("", h.OrderDo.List)
		orderDos.GET("/:orderDoID/:orderDoNo", h.OrderDo.Detail)

		write := orderDos.Group("")
		write.Use(middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleManager))
		{
			write.GET("/create", h.OrderDo.ShowCreate)
			write.POST("/create", h.OrderDo.Create)
			write.GET("/:orderDoID/:orderDoNo/edit", h.OrderDo.ShowEdit)
			write.POST("/:orderDoID/:orderDoNo/edit", h.OrderDo.Update)
			write.POST("/:orderDoID/:orderDoNo/delete", h.OrderDo.Delete)
			//Details
			//write.GET("/load_details/:orderDoID/:orderDoNo/details", h.OrderDo.)
			// write.POST("/act_addDetail/:orderDoNo/details", h.OrderDo.AddDetail)
			// write.POST("/details/:orderDoID/:orderDoNo/delete", h.OrderDo.DeleteDetail)
			write.POST("/act_delDetail/:orderDoID/:orderDoNo/details/:orderDoDetailID/:orderDoDetailNo/delete", h.OrderDo.DeleteDetail)

			//alamat Ambil detail hasil test postman = 200
			// http://localhost:8083/api/orderdo/load_details/:orderDoID/:orderDoNo/details

			// //alamat Hapus detail hasil test postman = 204
			// http://localhost:8083/api/orderdo/details/:orderDoDetailID/:orderDoDetailNo

			// //alamat add detail hasil test postman = 201
			// http://localhost:8083/api/orderdo/:orderDoNo/details

			// ini sudah test postman dan bekerja dengan baik, bisa di sesuaikan dengan ini ?

			//write.POST("/:orderDoNo/details", h.OrderDo.AddDetail)
			// Menangkap orderDoID dan orderDoNo dari URL halaman client
			write.POST("/:orderDoID/:orderDoNo/details", h.OrderDo.AddDetail)
			write.POST("/details/:orderDoDetailID/:orderDoDetailNo/delete", h.OrderDo.DeleteDetail)
			// Autocomplete proxies for the "Add Detail Line" form on the
			// detail page: browser calls these (cookie auth), the handler
			// forwards to the backend with the server-held Bearer token.
			write.GET("/lookup/customers/:custName", h.OrderDo.LookupCustomers)
			write.GET("/lookup/orders/:custID/:query", h.OrderDo.LookupOrders)
		}
	}

	// Sharing Profit — everyone can view, ADMIN/MANAGER can write (matches
	// the backend's /api/sharing-profits role split).
	sharingProfits := protected.Group("/sharing-profits")
	{
		sharingProfits.GET("", h.SharingProfit.List)
		sharingProfits.GET("/:num/:id", h.SharingProfit.Detail)

		spWrite := sharingProfits.Group("")
		spWrite.Use(middleware.RequireRoles(middleware.RoleAdmin, middleware.RoleManager))
		{
			spWrite.GET("/create", h.SharingProfit.ShowCreate)
			spWrite.POST("/create", h.SharingProfit.Create)
			spWrite.GET("/:num/:id/edit", h.SharingProfit.ShowEditForm)
			spWrite.POST("/:num/:id/edit", h.SharingProfit.Update)
			spWrite.POST("/:num/:id/delete", h.SharingProfit.Delete)
		}
	}

// spGroup := r.Group("/sharing-profits", middleware.RequireAuth())
// {
//     spGroup.GET("", h.SharingProfit.List)
//     spGroup.GET("/create", h.SharingProfit.ShowCreateForm)
//     spGroup.POST("/create", h.SharingProfit.Create)
    
//     // Route Edit & Update
//     spGroup.GET("/:num/:id/edit", h.SharingProfit.ShowEditForm)
//     spGroup.POST("/:num/:id/edit", h.SharingProfit.Update)
// }

	return r
}
