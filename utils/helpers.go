package utils

import (
	"net/http"

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

// func SetParsedUintUsersRelation(data map[string]string) map[string]uint {
// 	uIntMap := make(map[string]uint)
// 	for key, strVal := range data {
// 		if key == "creator_id" || key == "modifier_id" {
// 			parsedUint64, err := strconv.ParseUint(strVal, 10, 64)
// 			if err != nil {
// 				parsedUint64 = 0
// 			}
// 			uIntMap[key] = uint(parsedUint64)
// 		}
// 	}
// 	return uIntMap
// }

// func parseInt(val any, defaultVal int) int {
// 	switch v := val.(type) {
// 	case int:
// 		return v
// 	case int64:
// 		return int(v)
// 	case float64:
// 		return int(v)
// 	case string:
// 		i, err := strconv.Atoi(v)
// 		if err == nil {
// 			return i
// 		}
// 	}
// 	return defaultVal
// }

// func parseUint(val any) uint {
// 	switch v := val.(type) {
// 	case uint:
// 		return v
// 	case int:
// 		return uint(v)
// 	case int64:
// 		return uint(v)
// 	case float64:
// 		return uint(v)
// 	case string:
// 		i, err := strconv.ParseUint(v, 10, 64)
// 		if err == nil {
// 			return uint(i)
// 		}
// 	}
// 	return 0
// }
