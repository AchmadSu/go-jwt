package helpers

import "example.com/m/config"

type Status int

func GetEntityStatusLabel(status Status) string {
	switch status {
	case config.Inactive:
		return "Inactive"
	case config.Active:
		return "Active"
	case config.Draft:
		return "Draft"
	default:
		return "Unknown"
	}
}

func SetEntityStatusLabel[T any](
	data []T,
	getStatus func(item *T) int,
	setLabel func(item *T, label string),
) []T {

	for i := range data {
		status := getStatus(&data[i])
		label := GetEntityStatusLabel(Status(status))
		setLabel(&data[i], label)
	}

	return data
}

// func GetEntityStatus(data string) int {
// 	var status int
// 	switch data {
// 	case "true":
// 		status = 1
// 	case "false":
// 		status = 0
// 	case "draft":
// 		status = 2
// 	default:
// 		status = 3
// 	}
// 	return status
// }
