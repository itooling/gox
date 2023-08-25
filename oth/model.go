package oth

import (
	"time"

	"gorm.io/gorm"
)

const (
	codeOK   = 0
	codeFail = 1

	msgOk   = "ok"
	msgFail = "fail"

	defaultTenant = "one"
)

type Model struct {
	ID        uint           `json:"id" gorm:"column:id;primaryKey;comment:主键"`
	Tenant    string         `json:"-" gorm:"index;column:tenant;comment:租户"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime;default:current_timestamp;comment:创建日期"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime;default:current_timestamp;comment:更新日期"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:逻辑删除"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Tenant == "" {
		m.Tenant = defaultTenant
	}
	return nil
}

type Result struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

func Ok() *Result {
	r := Result{
		Message: msgOk,
		Code:    codeOK,
	}
	return &r
}

func OkMsg(msg string) *Result {
	r := Result{
		Message: msg,
		Code:    codeOK,
	}
	return &r
}

func OkCode(code int) *Result {
	r := Result{
		Message: msgOk,
		Code:    code,
	}
	return &r
}

func OkData(data any) *Result {
	r := Result{
		Message: msgOk,
		Code:    codeOK,
		Data:    data,
	}
	return &r
}

func OkDataCode(data any, code int) *Result {
	r := Result{
		Message: msgOk,
		Code:    code,
		Data:    data,
	}
	return &r
}

func OkDataMsg(data any, msg string) *Result {
	r := Result{
		Message: msg,
		Code:    codeOK,
		Data:    data,
	}
	return &r
}

func OkDataMsgCode(data any, msg string, code int) *Result {
	r := Result{
		Message: msg,
		Code:    code,
		Data:    data,
	}
	return &r
}

func Fail() *Result {
	r := Result{
		Message: msgFail,
		Code:    codeFail,
	}
	return &r
}

func FailMsg(msg string) *Result {
	r := Result{
		Message: msg,
		Code:    codeFail,
	}
	return &r
}

func FailMsgCode(msg string, code int) *Result {
	r := Result{
		Message: msg,
		Code:    code,
	}
	return &r
}

func FailMsgData(msg string, data any) *Result {
	r := Result{
		Message: msg,
		Code:    codeFail,
		Data:    data,
	}
	return &r
}

func FailMsgCodeData(msg string, code int, data any) *Result {
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
