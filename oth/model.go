package oth

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

const (
	MsgOk   = "ok"
	MsgFail = "fail"
)

type Model struct {
	ID        uint           `json:"id" gorm:"column:id;primaryKey;comment:主键"`
	Tenant    string         `json:"tenant" gorm:"index;column:tenant;comment:租户"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime;default:current_timestamp;comment:创建日期"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime;default:current_timestamp;comment:更新日期"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index;comment:逻辑删除"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Tenant == "" {
		m.Tenant = "one"
	}
	return nil
}

type Result struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Success(data any, msg ...string) *Result {
	m := MsgOk
	if len(msg) > 0 {
		m = msg[0]
	}
	r := Result{
		Message: m,
		Code:    http.StatusOK,
		Data:    data,
	}
	return &r
}

func SuccessCode(data any, code int, msg ...string) *Result {
	m := MsgOk
	if len(msg) > 0 {
		m = msg[0]
	}
	r := Result{
		Message: m,
		Code:    code,
		Data:    data,
	}
	return &r
}

func Fail(msg string) *Result {
	r := Result{
		Message: msg,
	}
	return &r
}

func FailCode(msg string, code int) *Result {
	r := Result{
		Message: msg,
		Code:    code,
	}
	return &r
}

func FailData(msg string, data any) *Result {
	r := Result{
		Message: msg,
		Data:    data,
	}
	return &r
}

func FailCodeData(msg string, code int, data any) *Result {
	r := Result{
		Message: msg,
		Code:    code,
		Data:    data,
	}
	return &r
}

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
	if current < 1 {
		current = 1
	}
	if current > count {
		current = count
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
	if current < 1 {
		current = 1
	}
	if current > count {
		current = count
	}
	beg := size * (current - 1)
	end := beg + size
	if current == count {
		end = total
	}
	list := data[beg:end]
	return IPage[T](current, size, total, list)
}

type BaseError struct {
	Code    int
	Message string
}

func (u *BaseError) Error() string {
	return u.Message
}
