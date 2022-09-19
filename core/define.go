package core

import "errors"

type (
	Application struct {
		Origin     interface{}
		Parse      map[string]interface{}
		StructName string
	}
)

var (
	ErrorStructType = errors.New("非结构体类型")
	//ErrorStructNilPointer = errors.New("结构体·空指针异常")
)

const TAGNAME = "json"
