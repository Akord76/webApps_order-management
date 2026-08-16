# Order Management — Web App (API consumer)

A server-rendered Go + Gin + Bootstrap 5 front-end that consumes the
`order-management` backend REST API. It holds no database connection of
its own — every page is populated by calling the API over HTTP, with the
signed-in user's JWT carried in an HttpOnly cookie.

## Architecture

```
Browser
  │
  ▼
routes/routes.go
  │
  ▼
middleware.JWTAuth (cookie)  -- redirects to /login if missing/expired
  │
  ▼
middleware.RequireRoles      -- 403 page for disallowed roles (Order/Order Sale writes)
  │
  ▼
handler/*.go                 -- parses form input, calls the API, renders a template
  │
  ▼
client/api_client.go         -- HTTP client, injects "Authorization: Bearer <jwt>"
  │
  ▼
order-management backend API (see the sibling `order-management/` project)
```

## Role behavior in the UI

Same hierarchy as the backend:

| Role    | Can see                                             |
|---------|------------------------------------------------------|
| ADMIN   | everything, including Order/Order Sale create/edit/delete |
| MANAGER | everything, including Order/Order Sale create/edit/delete |
| USER    | can view Orders/Order Sales; create/edit/delete controls are hidden and the routes are blocked by `middleware.RequireRoles` |

Master data (categories, products, customers, employees, suppliers) is
open to any signed-in role in this app — the backend still enforces its
own rules if that ever changes.

## Template naming

Every module folder has its own `create_update.html`, `page_view_*.html`,
and `details*.html`. Since Gin's default `LoadHTMLGlob` names templates by
base filename only, that would collide. `routes/template_loader.go` names
each template by its path relative to `template/` instead (e.g.
`"category_template/create_update.html"`), so handlers call
`c.HTML(200, "category_template/create_update.html", data)` explicitly per
module.

## Setup

```bash
cp .env.example .env   # point API_BASE_URL at your running backend
go mod tidy
go run main.go
```

Make sure the backend API (`order-management/`) is running first —
`API_BASE_URL` defaults to `http://localhost:8080/api`.

## Pages

| Path                              | Notes                                    |
|------------------------------------|-------------------------------------------|
| `/login`, `/register`, `/logout`   | public                                    |
| `/` | home dashboard, links to every module |
| `/categories`, `/products`, `/customers`, `/employees`, `/suppliers` | list/detail/create/edit/delete |
| `/orders`, `/order-sales` | list/detail for everyone; create/edit/delete + detail-line management restricted to ADMIN/MANAGER |

