package valid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vzh "github.com/go-playground/validator/v10/translations/zh"
)

// Validator struct
type Validator struct {
	Tag  string         // tag name
	Msg  string         // error msg
	Name string         // field name
	Func validator.Func // validator function
}

func NewValidator(tag, msg string, f validator.Func) *Validator {
	return &Validator{Tag: tag, Msg: msg, Func: f}
}

var (
	uni      *ut.UniversalTranslator
	tags     map[string]Validator
	trans    ut.Translator
	validate *validator.Validate
)

func init() {
	t := zh.New()
	uni = ut.New(t, t)
	tags = make(map[string]Validator)
	trans, _ = uni.GetTranslator("zh")
	reset()
}

// reset the validator
func reset() {
	var ok bool
	if validate, ok = binding.Validator.Engine().(*validator.Validate); ok {
		for _, t := range tags {
			validate.RegisterValidation(t.Tag, t.Func)
		}
	}
	vzh.RegisterDefaultTranslations(validate, trans)
}

// Translate translate error to error msg
func Translate(err error) error {
	if e, ok := err.(*validator.InvalidValidationError); ok {
		return e
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		var msg string
		for _, err := range errs {
			tag, field := err.Tag(), err.Field()
			if t, ex := tags[tag]; ex {
				msg = fmt.Sprintf("%s %s", field, t.Msg)
				if t.Name != "" {
					msg = strings.Replace(msg, field, t.Name, 1)
				}
			} else {
				msg = err.Translate(trans)
			}
		}
		if msg != "" {
			return errors.New(msg)
		}
	}
	return nil
}

// AddValidator add or reset the validator
func AddValidator(vs ...*Validator) {
	for _, v := range vs {
		if v.Tag != "" && v.Func != nil {
			tags[v.Tag] = *v
		}
	}
	reset()
}
