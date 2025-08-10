package utils

import (
	"net/http"
	"reflect"
	"time"

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

func AssignedKeyModel(model interface{}, data map[string]any) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errs.New("model must be pointer to struct", http.StatusBadRequest)
	}
	v = v.Elem()

	for key, value := range data {
		// Skip value nil
		if value == nil {
			continue
		}

		val := reflect.ValueOf(value)

		// Skip empty string
		if val.Kind() == reflect.String && val.Len() == 0 {
			continue
		}

		//skip empty date time
		if t, ok := value.(time.Time); ok && t.IsZero() {
			continue
		}

		if key == "Price" {
			if qtyVal, ok := value.(int); ok && qtyVal <= 0 {
				continue
			}
			if qtyVal, ok := value.(uint); ok && qtyVal == 0 {
				continue
			}
			if qtyVal, ok := value.(float64); ok && qtyVal <= 0 {
				continue
			}
		}

		if key == "Qty" {
			if qtyVal, ok := value.(int); ok && qtyVal < 0 {
				continue
			}
			if qtyVal, ok := value.(float64); ok && qtyVal < 0 {
				continue
			}
		}

		// Skip empty slice/map
		if (val.Kind() == reflect.Slice || val.Kind() == reflect.Map) && val.Len() == 0 {
			continue
		}

		field := v.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			continue // skip if field not found
		}

		// Assign value if data type is correct
		if val.Type().AssignableTo(field.Type()) {
			field.Set(val)
		} else if val.Type().ConvertibleTo(field.Type()) {
			field.Set(val.Convert(field.Type()))
		}
	}
	return nil
}
