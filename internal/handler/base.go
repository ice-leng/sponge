package handler

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/validator"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validatorV10 "github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Trans ut.Translator
)

// InitTrans initialize the translator
func InitTrans() {
	v := validator.Init()
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	Trans = trans
	_ = zhTranslations.RegisterDefaultTranslations(v.Validate, trans)

	v.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		label := fld.Tag.Get("label")
		if label == "" {
			return fld.Name
		}
		return label
	})

	binding.Validator = v
}

type baseHandler struct{}

// getCurrentUid get the current user id from context
func (b *baseHandler) getCurrentUid(c *gin.Context) uint64 {
	uid, ok := c.Get("id")
	if !ok {
		return 0
	}
	return uid.(uint64)
}

// isErrcode check if the error is errcode.Error
func (b *baseHandler) isErrcode(err error) (*errcode.Error, bool) {
	ec := errcode.ParseError(err)
	if ec.Code() > 0 {
		return ec, true
	}
	return nil, false
}

// getIdFromPath get the id from path parameter
func (b *baseHandler) getIdFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

// getValidatorErrorMsg get the validator error message
func (b *baseHandler) getValidatorErrorMsg(err error) string {
	var verrs validatorV10.ValidationErrors
	if errors.As(err, &verrs) {
		return verrs[0].Translate(Trans)
	}
	return err.Error()
}
