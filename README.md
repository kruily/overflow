# Overflow 移除器
有一种场景：有a,b 两个请求,都用到了共同的结构体作为返回值。
但在b请求下，不需要对结构体中的某个参数做返回。可以使用overflow做移除。

## 下载
```shell
go get github.com/jingxiu1016/overflow
```

## 使用
```go
package main

import (
	"fmt"
	"github.com/jingxiu1016/overflow/core"
)

type Test struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	app, err := core.NewApp(&Test{
		Id:       1,
		Username: "jingxiu",
		Password: "testtestest",
	})
	if err != nil {
		fmt.Printf("error: %#v\n\n", err.Error())
		return
	}
	app.Overflow([]string{"password", "username"})

}

```
## 移除器
并未对不存在的键进行安全检测，因为不存在的键执行空操作
```go
app.Overflow([]string{...})
```
## 拿到结果
```go
app.Result()
或
app.Parse
```

## 支持结构体嵌套,切片组合
```go
package main
type SubElement struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Element struct {
	Type  string      `json:"type"`
	Value string      `json:"value"`
	Sub   *SubElement `json:"sub"`
}

type Test struct {
    Id       int64         `json:"id"`
    Username string        `json:"username"`
    Password string        `json:"password"`
    Element  *Element      `json:"element"`
    Sub      []*SubElement `json:"sub"`
    Idxs     []int64       `json:"idxs"`
    Tags     []string      `json:"tags"`
}
func main() {
	app, err := core.NewApp(&Test{
		Id:       1,
		Username: "jingxiu",
		Password: "jingxiu1016",
		Element: &Element{
			Type:  "test",
			Value: "this is test",
			Sub: &SubElement{
				Id:   1,
				Name: "element children",
			},
		},
		Sub: []*SubElement{
			{
				Id:   2,
				Name: "test children1",
			},
			{
				Id:   3,
				Name: "test children2",
			},
		},
		Idxs: []int64{2, 3, 4, 5, 6},
		Tags: []string{"test1", "test2"},
	})
	if err != nil {
		fmt.Printf("error: %#v\n\n", err.Error())
		return
	}
	

}

```

## 支持移除多层成员
| 以 **.** 作为成员父子项链接的方式
```go
app.Overflow([]string{"element.sub"})
```