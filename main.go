package main

import (
	"log"

	"webApps_order-management/client"
	"webApps_order-management/config"
	"webApps_order-management/handler"
	"webApps_order-management/routes"
)

func main() {
	// ---- Config ----
	cfg := config.LoadConfig()

	// ---- API client (this app talks to the backend, never a DB directly) ----
	api := client.NewAPIClient(cfg.APIBaseURL)

	// ---- Handlers ----
	h := &routes.Handlers{
		Auth:      handler.NewAuthHandler(api, cfg.CookieName, cfg.CookieSecure),
		Category:  handler.NewCategoryHandler(api),
		Product:   handler.NewProductHandler(api),
		Customer:  handler.NewCustomerHandler(api),
		Employee:  handler.NewEmployeeHandler(api),
		Supplier:  handler.NewSupplierHandler(api),
		Order:     handler.NewOrderHandler(api),
		OrderSale: handler.NewOrderSaleHandler(api),
	}

	// ---- Router ----
	router := routes.SetupRouter(cfg.JWTSecret, cfg.CookieName, h)

	log.Printf("web app starting on :%s (backend API: %s)", cfg.ServerPort, cfg.APIBaseURL)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
