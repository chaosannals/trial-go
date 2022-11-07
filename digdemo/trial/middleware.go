package trial

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ApiAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		uuid := req.Header.Get("signature")
		if uuid != "123456789" {
			return c.JSON(
				http.StatusUnauthorized,
				map[string]any {
					"code": -1,
					"msg": "没有权限",
				},
			)
		}
		return next(c)
	}
}