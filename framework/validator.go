/**
@Author: wei-g
@Date:   2020/12/3 2:50 下午
@Description: 让调用方更便捷的调用gin中封装的validator.Validate
*/

package framework

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	ValidateAssertErr = fmt.Errorf("gin Validator.Engine cannot be converted to *validator.Validate")
	Validate          *validator.Validate
	validateOnce      sync.Once
)

// AssertValidate 初始化 gin 使用的单例的 Validate,让调用方可以快速使用 binding.Validator 的进行注册校验规则
func AssertValidate() (*validator.Validate, error) {
	var err error
	validateOnce.Do(func() {
		var ok bool
		Validate, ok = binding.Validator.Engine().(*validator.Validate)
		if !ok {
			err = ValidateAssertErr
			return
		}
		return
	})
	return Validate, err
}

// RegisterValidation 注册定义校验规则函数
func RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) (err error) {
	err = Validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
	return
}
func RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{}) {
	Validate.RegisterStructValidation(fn, types...)
}
