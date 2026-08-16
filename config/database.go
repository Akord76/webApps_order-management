package config

// database.go is intentionally a stub.
//
// This web app is a pure API client: every screen renders data fetched
// from the order-management REST API (see client/api_client.go), and it
// never opens its own SQL Server connection. The file is kept here so the
// project layout mirrors the backend service.
//
// If a future requirement needs local persistence in this app (e.g. a
// server-side cache, audit log, or session store backed by a real
// database instead of the JWT cookie), wire it up here following the same
// pattern as the backend's config/database.go (gorm.Open + connection
// pool settings) and call it from main.go.
