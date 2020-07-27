package bootstrap

import (
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/dongri/phonenumber"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func InitValidator() {

	// NOTE: ommitting allot of error checking for brevity

	en := en.New()
	uni = ut.New(en, en)

	trans, _ = uni.GetTranslator("en")

	validate = validator.New()
	validate.RegisterTagNameFunc(fieldNameFromJSON)
	en_translations.RegisterDefaultTranslations(validate, trans)
	_ = validate.RegisterValidation("nospace", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})
	_ = validate.RegisterValidation("alphanumandvalidchar", func(fl validator.FieldLevel) bool {
		s := "^[a-zA-Z0-9!@#$&()\\-`.+,/\\s]*$"
		r := regexp.MustCompile(s)
		return r.MatchString(fl.Field().String())
	})
	_ = validate.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
		return passwordOK(fl.Field().String())
	})
	_ = validate.RegisterValidation("emaildomain", func(fl validator.FieldLevel) bool {
		s := `^[a-zA-Z0-9_.+-]+@(?:(?:[a-zA-Z0-9-]+\.)?[a-zA-Z]+\.)?(sepulsa\.com|sepulsa\.id|alterra\.id)$`
		r := regexp.MustCompile(s)
		return r.MatchString(fl.Field().String())
	})
	_ = validate.RegisterValidation("mobileno", func(fl validator.FieldLevel) bool {
		if phonenumber.ParseWithLandLine(fl.Field().String(), "ID") == "" {
			return false
		} else {
			return true
		}
	})
	_ = validate.RegisterValidation("primitiveid", func(fl validator.FieldLevel) bool {
		_, err := primitive.ObjectIDFromHex(fl.Field().String())
		if err != nil {
			return false
		} else {
			return true
		}
	})

	_ = validate.RegisterTranslation("nospace", trans,
		func(ut ut.Translator) error {
			return ut.Add("nospace", "{0} can not be empty string", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			s, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return s
		},
	)

	_ = validate.RegisterTranslation("alphanumandvalidchar", trans,
		func(ut ut.Translator) error {
			return ut.Add("alphanumandvalidchar", "{0} must contain valid word, number, spaces, and \"!@#$&()\\-`.+,/\" char", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			s, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return s
		},
	)
	_ = validate.RegisterTranslation("strongpassword", trans,
		func(ut ut.Translator) error {
			return ut.Add("strongpassword", "{0} must contain at least one lowercase, uppercase, and number", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			s, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return s
		},
	)
	_ = validate.RegisterTranslation("emaildomain", trans,
		func(ut ut.Translator) error {
			return ut.Add("emaildomain", "{0} only accept [@sepulsa.com, @sepulsa.id, @alterra.id]", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			s, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return s
		},
	)
	_ = validate.RegisterTranslation("mobileno", trans,
		func(ut ut.Translator) error {
			return ut.Add("mobileno", "{0} invalid phone number format", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			s, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return s
		},
	)

	_ = validate.RegisterTranslation("primitiveid", trans,
		func(ut ut.Translator) error {
			return ut.Add("primitiveid", "{0} invalid mongo id format", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			s, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return s
		},
	)
}

// use the names which have been specified for JSON representations of structs, rather than normal Go field names
func fieldNameFromJSON(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func Validate(i interface{}) map[string]string {
	err := validate.Struct(i)
	if err != nil {
		return err.(validator.ValidationErrors).Translate(trans)
	}

	t := make(validator.ValidationErrorsTranslations)
	return t
}

func VarValidate(val, roles string) map[string]string {
	err := validate.Var(val, roles)
	if err != nil {
		return err.(validator.ValidationErrors).Translate(trans)
	}

	t := make(validator.ValidationErrorsTranslations)
	return t
}

func passwordOK(p string) bool {
	var mustHave = []func(rune) bool{
		unicode.IsUpper,
		unicode.IsLower,
		unicode.IsDigit,
	}

	for _, testRune := range mustHave {
		found := false
		for _, r := range p {
			if testRune(r) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
