package utils

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// 注册自定义验证规则
	_ = validate.RegisterValidation("password", validatePassword)
	_ = validate.RegisterValidation("username", validateUsername)
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// validatePassword 自定义密码验证
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// 密码长度至少6位
	if len(password) < 6 {
		return false
	}

	// 检查是否包含数字和字母
	var hasLetter, hasNumber bool
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}

// validateUsername 自定义用户名验证
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// 用户名长度3-20位，只能包含字母、数字、下划线
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}
