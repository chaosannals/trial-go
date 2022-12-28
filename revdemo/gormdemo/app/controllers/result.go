package controllers

import (
	"net/http"

	"github.com/revel/revel"
)

type ResultController struct {
	*revel.Controller
	R *ResultHandler
}

func (c *ResultController) Init() {
	c.R = &ResultHandler{Controller: c}
}

type ResultHandler struct {
	Controller *ResultController
}

func (h *ResultHandler) Ok(data interface{}) revel.Result {
	h.Controller.Response.Status = http.StatusOK
	return h.Controller.RenderJSON(map[string]interface{}{
		"code": 0,
		"tip":  "ok",
		"info": data,
	})
}
