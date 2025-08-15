package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"example.com/m/config"
	"example.com/m/dto"
	"example.com/m/errs"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (dto.PublicUser, error) {
	userAny, exist := c.Get("user")
	if !exist {
		return dto.PublicUser{}, errs.New("user context not found", http.StatusNotFound)
	}

	user, ok := userAny.(dto.PublicUser)
	if !ok {
		return dto.PublicUser{}, errs.New("invalid user context type", http.StatusInternalServerError)
	}
	return user, nil
}

func ContainsString(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

func DifferenceUint(a, b []uint) []uint {
	m := make(map[uint]bool, len(b))
	for _, v := range b {
		m[v] = true
	}

	diff := []uint{}
	for _, v := range a {
		if !m[v] {
			diff = append(diff, v)
		}
	}
	return diff
}

func AssignedKeyModel(model any, data map[string]any) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errs.New("model must be pointer to struct", http.StatusBadRequest)
	}
	v = v.Elem()

	for key, value := range data {
		if value == nil {
			continue
		}

		val := reflect.ValueOf(value)

		// Dereference pointer
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
			value = val.Interface()
		}

		if val.Kind() == reflect.String && val.Len() == 0 {
			continue
		}

		if t, ok := value.(time.Time); ok && t.IsZero() {
			continue
		}

		if strings.EqualFold(key, "Price") {
			if num, ok := value.(float64); ok && num <= 0 {
				continue
			}
			if num, ok := value.(int); ok && num <= 0 {
				continue
			}
			if num, ok := value.(uint); ok && num == 0 {
				continue
			}
		}

		if strings.EqualFold(key, "Qty") {
			if num, ok := value.(float64); ok && num < 0 {
				continue
			}
			if num, ok := value.(int); ok && num < 0 {
				continue
			}
		}

		if (val.Kind() == reflect.Slice || val.Kind() == reflect.Map) && val.Len() == 0 {
			continue
		}

		field := v.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		if val.Type().AssignableTo(field.Type()) {
			field.Set(val)
		} else if val.Type().ConvertibleTo(field.Type()) {
			field.Set(val.Convert(field.Type()))
		}
	}

	return nil
}

func MergeDateTime(dateStr string, timeStr string) (time.Time, error) {
	combined := fmt.Sprintf("%s %s", dateStr, timeStr)
	result, err := time.ParseInLocation(string(config.LayoutDateTime), combined, time.Local)
	if err != nil {
		return time.Time{}, errs.New("invalid format date", http.StatusBadRequest)
	}
	return result, nil
}
