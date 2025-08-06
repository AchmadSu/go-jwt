package helpers

type Status int

const (
	StatusInactive Status = 0
	StatusActive   Status = 1
	StatusDraft    Status = 2
)

func GetEntityStatusLabel(status Status) string {
	switch status {
	case StatusInactive:
		return "Inactive"
	case StatusActive:
		return "Active"
	case StatusDraft:
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
