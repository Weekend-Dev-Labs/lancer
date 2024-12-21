package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type LancerValidator struct {
	validator *validator.Validate
}
func (lv *LancerValidator) Validate(i interface{}) error {
	if err := lv.validator.Struct(i); err != nil {
		return echo.
	}
}