package errors

import (
	goErrors "errors"
	"fmt"
)

// 系统错误接口
type SysErrorInterface interface {
	Error() string
	String() string
	Status() int
}

// 普通错误接口
type NormalErrorInterface interface {
	Error() string
	Status() int
}

// 兼容golang errors对象
func New(text string) error {
	return goErrors.New(text)
}

// 创建系统错误
func NewSys(options ...interface{}) error {
	var (
		status    int
		simpleMsg string
		fullMsg   string
	)

	for _, v := range options {
		switch opt := v.(type) {
		case int:
			status = opt
		default:
			simpleMsg = "系统错误"
			fullMsg = fmt.Sprintf("系统错误: %v", opt)
		}
	}

	return &SysError{
		status:    status,
		simpleMsg: simpleMsg,
		fullMsg:   fullMsg,
	}
}

// 系统错误实现
type SysError struct {
	status    int
	simpleMsg string
	fullMsg   string
}

func (e *SysError) Error() string {
	return e.simpleMsg
}

func (e *SysError) String() string {
	return e.fullMsg
}

func (e *SysError) Status() int {
	return e.status
}

// 创建普通错误
func NewNormal(options ...interface{}) error {
	var (
		status int
		msg    string
	)

	for _, v := range options {
		switch opt := v.(type) {
		case int:
			status = opt
		default:
			msg = fmt.Sprintf("%v", opt)
		}
	}

	return &NormalError{
		status: status,
		msg:    msg,
	}
}

// 普通错误实现
type NormalError struct {
	status int
	msg    string
}

func (e *NormalError) Error() string {
	return e.msg
}

func (e *NormalError) Status() int {
	return e.status
}
