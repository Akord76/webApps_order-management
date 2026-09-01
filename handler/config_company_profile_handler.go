package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"webApps_order-management/client"
	"webApps_order-management/model"

	"github.com/gin-gonic/gin"
)

type ConfigCompanyProfileHandler struct {
	api *client.APIClient
}

func NewConfigCompanyProfileHandler(api *client.APIClient) *ConfigCompanyProfileHandler {
	return &ConfigCompanyProfileHandler{api: api}
}

func (h *ConfigCompanyProfileHandler) List(c *gin.Context) {
	var configs []model.ConfigCompanyProfile
	data := baseData(c, "Company Profile Settings")
	if err := h.api.Get("/config-company-profiles", token(c), &configs); err != nil {
		data["Error"] = "Failed to load company profile settings: " + err.Error()
	}
	data["Configs"] = configs
	c.HTML(http.StatusOK, "config_company_profile_template/page_view_config.html", data)
}

func (h *ConfigCompanyProfileHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	var config model.ConfigCompanyProfile
	if err := h.api.Get("/config-company-profiles/"+id, token(c), &config); err != nil {
		c.Redirect(http.StatusFound, "/config-company-profiles?err="+url.QueryEscape("Setting not found"))
		return
	}

	data := baseData(c, "Setting Detail")
	data["Config"] = config
	c.HTML(http.StatusOK, "config_company_profile_template/details_config.html", data)
}

func (h *ConfigCompanyProfileHandler) ShowCreate(c *gin.Context) {
	data := baseData(c, "New Setting")
	data["IsEdit"] = false
	data["Config"] = model.ConfigCompanyProfile{}
	c.HTML(http.StatusOK, "config_company_profile_template/create_update.html", data)
}

func (h *ConfigCompanyProfileHandler) Create(c *gin.Context) {
	configID, _ := strconv.Atoi(c.PostForm("config_id"))
	body := gin.H{
		"config_id":                configID,
		"config_company_profile":   c.PostForm("config_company_profile"),
		"config_company_profile_1": c.PostForm("config_company_profile_1"),
	}

	if err := h.api.Post("/config-company-profiles", token(c), body, nil); err != nil {
		data := baseData(c, "New Setting")
		data["IsEdit"] = false
		data["Config"] = model.ConfigCompanyProfile{
			ConfigID:              configID,
			ConfigCompanyProfile:  c.PostForm("config_company_profile"),
			ConfigCompanyProfile1: c.PostForm("config_company_profile_1"),
		}
		data["Error"] = "Failed to create setting: " + err.Error()
		c.HTML(http.StatusOK, "config_company_profile_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/config-company-profiles?ok="+url.QueryEscape("Setting created"))
}

func (h *ConfigCompanyProfileHandler) ShowEdit(c *gin.Context) {
	id := c.Param("id")

	var config model.ConfigCompanyProfile
	if err := h.api.Get("/config-company-profiles/"+id, token(c), &config); err != nil {
		c.Redirect(http.StatusFound, "/config-company-profiles?err="+url.QueryEscape("Setting not found"))
		return
	}

	data := baseData(c, "Edit Setting")
	data["IsEdit"] = true
	data["Config"] = config
	c.HTML(http.StatusOK, "config_company_profile_template/create_update.html", data)
}

func (h *ConfigCompanyProfileHandler) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	body := gin.H{
		"config_company_profile":   c.PostForm("config_company_profile"),
		"config_company_profile_1": c.PostForm("config_company_profile_1"),
	}

	if err := h.api.Put("/config-company-profiles/"+id, token(c), body, nil); err != nil {
		data := baseData(c, "Edit Setting")
		data["IsEdit"] = true
		data["Config"] = model.ConfigCompanyProfile{
			ConfigID:              idInt,
			ConfigCompanyProfile:  c.PostForm("config_company_profile"),
			ConfigCompanyProfile1: c.PostForm("config_company_profile_1"),
		}
		data["Error"] = "Failed to update setting: " + err.Error()
		c.HTML(http.StatusOK, "config_company_profile_template/create_update.html", data)
		return
	}

	c.Redirect(http.StatusFound, "/config-company-profiles?ok="+url.QueryEscape("Setting updated"))
}

func (h *ConfigCompanyProfileHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.api.Delete("/config-company-profiles/"+id, token(c)); err != nil {
		c.Redirect(http.StatusFound, "/config-company-profiles?err="+url.QueryEscape("Failed to delete setting: "+err.Error()))
		return
	}

	c.Redirect(http.StatusFound, "/config-company-profiles?ok="+url.QueryEscape("Setting deleted"))
}
