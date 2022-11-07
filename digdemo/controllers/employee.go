package controllers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/chaosannals/trial-go-digdemo/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type EmployeeController struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewEmployeeController(
	db *gorm.DB,
	logger *zerolog.Logger,
) *EmployeeController {
	return &EmployeeController{
		db:     db,
		logger: logger,
	}
}

func (i *EmployeeController) List(c echo.Context) (err error) {
	var rows []models.EEmployee

	i.db.Model(&models.EEmployee{}).Select("*").Find(&rows)

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"rows":    rows,
			"message": "ok",
		},
	)
}

func (i *EmployeeController) Add(c echo.Context) (err error) {
	row := &models.EEmployee{}
	if err = c.Bind(row); err != nil {
		return
	}

	db := i.db.Create(row) // 返回的 db 不是 i.db
	if err = db.Error; err != nil {
		if ok, err := regexp.MatchString("Duplicate\\s+entry.+?ACCOUNT_UNIQUE", err.Error()); ok && err == nil {
			return c.JSON(
				http.StatusOK,
				map[string]any{
					"code":    -1,
					"account": row.Account,
					"message": "账号重复",
				},
			)
		}
		return
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"id":      row.ID,
			"message": "ok",
		},
	)
}

type EmployeeEditParam struct {
	models.EEmployee
	Mobilephones []string `json:"mobilephones"`
}

func (i *EmployeeController) Edit(c echo.Context) (err error) {
	i.logger.Trace().Msg("获取数据")
	param := &EmployeeEditParam{}
	if err = c.Bind(param); err != nil {
		return
	}

	row := &param.EEmployee
	i.logger.Trace().Msg("开始事务")
	tx := i.db.Begin()
	defer func() {
		if err != nil {
			i.logger.Trace().Msg("回滚")
			tx.Rollback()
		} else {
			i.logger.Trace().Msg("提交")
			tx.Commit()
		}
	}()

	if err = tx.Model(&models.EEmployee{}).
		Where("id = ?", row.ID).
		Updates(row).
		Error; err != nil {
		return
	}

	if err = tx.Where("employee_id = ?", row.ID).
		Delete(&models.EEmployeeMobilephone{}).
		Error; err != nil {
		return
	}

	var ems []models.EEmployeeMobilephone
	for _, m := range param.Mobilephones {
		ems = append(ems, models.EEmployeeMobilephone{
			EmployeeID:  row.ID,
			Mobilephone: m,
		})
	}
	
	if err = tx.Model(&models.EEmployeeMobilephone{}).
		Create(&ems).
		Error; err != nil {
		return
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"id":      row.ID,
			"message": "ok",
		},
	)
}

func (i *EmployeeController) Del(c echo.Context) (err error) {
	id := c.Param("id")
	if err = i.db.Delete(&models.EEmployee{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(
				http.StatusOK,
				map[string]any{
					"code":    0,
					"id":      id,
					"message": "不是有效的id",
				},
			)
		}
		return
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"message": "ok",
		},
	)
}
