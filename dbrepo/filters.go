package dbrepo

type Filters struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (f Filters) Limit() int {
	if f.Size < 1 {
		return 10
	}
	if f.Size > 100 {
		return 100
	}
	return f.Size
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.Size
}
