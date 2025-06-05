package api

import (
	"github.com/dongnguyen248/simple_bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	return util.IsSupportedCurrency(currency)
}
