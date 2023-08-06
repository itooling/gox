package oth

type Page[T any] struct {
	Current int  `json:"current"`
	Size    int  `json:"size"`
	Total   int  `json:"total"`
	Count   int  `json:"count"`
	First   bool `json:"first"`
	Last    bool `json:"last"`
	List    []T  `json:"list"`
}

type PageParam struct {
	Current int `json:"current"`
	Size    int `json:"size"`
}

func IPage[T any](current, size, total int, list []T) *Page[T] {
	count := total/size + 1
	if total%size == 0 {
		count = total / size
	}
	return &Page[T]{
		Current: current,
		Size:    size,
		Total:   total,
		Count:   count,
		First:   current == 1,
		Last:    current == count,
		List:    list,
	}
}

func ToPage[T any](current, size int, data []T) *Page[T] {
	total := len(data)
	count := total/size + 1
	if total%size == 0 {
		count = total / size
	}
	beg := size * (current - 1)
	end := beg + size
	if current == count {
		end = total
	}
	list := data[beg:end]
	return IPage[T](current, size, total, list)
}
