package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	api *client.APIClient
}

func NewEmployeeHandler(api *client.APIClient) *EmployeeHandler {
	return &EmployeeHandler{api: api}
}

func (h *EmployeeHandler) List(c *gin.Context) {
	var employees []model.Employee
	data := baseData(c, "Employees")
	if err := h.api.Get("/employees", token(c), &employees); err != nil {
		data["Error"] = "Failed to load employees: " + err.Error()
	}
	data["Employees"] = employees
	c.HTML(http.StatusOK, "employee_template/page_view_employee.html", data)
}

func (h *EmployeeHandler) Detail(c *gin.Context) {
	path := "/employees/" + c.Param("employeeID") + "/" + c.Param("cardNumber")

	var employee model.Employee
	if err := h.api.Get(path, token(c), &employee); err != nil {
		c.Redirect(http.StatusFound, "/employees?err="+url.QueryEscape("Employee not found"))
		return
	}

	data := baseData(c, "Employee Detail")
	data["Employee"] = employee
	c.HTML(http.StatusOK, "employee_template/detailsemployee.html", data)
}

func (h *EmployeeHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Employee")
	data["IsEdit"] = false
	data["Employee"] = model.Employee{}
	c.HTML(http.StatusOK, "employee_template/create_update.html", data)
}

func (h *EmployeeHandler) formBody(c *gin.Context) gin.H {
	body := gin.H{
		"first_name":   c.PostForm("first_name"),
		"last_name":    c.PostForm("last_name"),
		"gender":       c.PostForm("gender"),
		"email":        c.PostForm("email"),
		"phone_number": c.PostForm("phone_number"),
		"bank_number":  c.PostForm("bank_number"),
		"photo_path":   c.PostForm("photo_path"),
	}
	if v, err := strconv.Atoi(c.PostForm("department_id")); err == nil {
		body["department_id"] = v
	}
	return body
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	cardNumber, _ := strconv.Atoi(c.PostForm("employee_card_number"))
	body := h.formBody(c)
	body["employee_card_number"] = cardNumber

	if err := h.api.Post("/employees", token(c), body, nil); err != nil {
		data := baseData(c, "New Employee")
		data["IsEdit"] = false
		data["Employee"] = model.Employee{EmployeeCardNumber: cardNumber, FirstName: c.PostForm("first_name")}
		data["Error"] = "Failed to create employee: " + err.Error()
		c.HTML(http.StatusOK, "employee_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/employees?ok="+url.QueryEscape("Employee created"))
}

func (h *EmployeeHandler) ShowEdit(c *gin.Context) {
	path := "/employees/" + c.Param("employeeID") + "/" + c.Param("cardNumber")

	var employee model.Employee
	if err := h.api.Get(path, token(c), &employee); err != nil {
		c.Redirect(http.StatusFound, "/employees?err="+url.QueryEscape("Employee not found"))
		return
	}

	data := baseData(c, "Edit Employee")
	data["IsEdit"] = true
	data["Employee"] = employee
	c.HTML(http.StatusOK, "employee_template/create_update.html", data)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	employeeID := c.Param("employeeID")
	cardNumber := c.Param("cardNumber")
	path := "/employees/" + employeeID + "/" + cardNumber

	body := h.formBody(c)

	if err := h.api.Put(path, token(c), body, nil); err != nil {
		eid, _ := strconv.Atoi(employeeID)
		cn, _ := strconv.Atoi(cardNumber)
		data := baseData(c, "Edit Employee")
		data["IsEdit"] = true
		data["Employee"] = model.Employee{EmployeeID: eid, EmployeeCardNumber: cn, FirstName: c.PostForm("first_name")}
		data["Error"] = "Failed to update employee: " + err.Error()
		c.HTML(http.StatusOK, "employee_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/employees?ok="+url.QueryEscape("Employee updated"))
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	path := "/employees/" + c.Param("employeeID") + "/" + c.Param("cardNumber")

	if err := h.api.Delete(path, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/employees?err="+url.QueryEscape("Failed to delete employee: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/employees?ok="+url.QueryEscape("Employee deleted"))
}
