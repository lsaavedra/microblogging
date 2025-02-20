package utils

type (
	PaginatedResponse[T any] struct {
		Total int `json:"total"`
		Items []T `json:"items"`
	}
)

func BuildResponse[T any](
	items []T,
) PaginatedResponse[T] {
	total := len(items)

	return PaginatedResponse[T]{
		Total: int(total),
		Items: items,
	}
}
