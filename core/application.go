package core

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// NewApp 新建应用
func NewApp(obj interface{}) (*Application, error) {
	value, err := correction(reflect.ValueOf(obj))
	if err != nil {
		return nil, err
	}
	if value.Kind() != reflect.Struct {
		return nil, ErrorStructType
	}
	app := &Application{Origin: obj, Parse: make(map[string]interface{}), StructName: value.Type().Name()}
	res := parseStruct(obj)
	app.Parse = res
	return app, nil
}

// Overflow 移除器
func (a *Application) Overflow(arr []string) {
	for _, item := range arr {
		if strings.Contains(item, ".") {
			arr := strings.Split(item, ".")
			idx := int64(0)
			roll(arr, a.Parse, &idx)
		} else {
			delete(a.Parse, item)
		}
	}
}

// Result 拿到结果
func (a *Application) Result() map[string]interface{} {
	return a.Parse
}

// 解析结构体
func parseStruct(obj interface{}) map[string]interface{} {
	value, err := correction(reflect.ValueOf(obj))
	if err != nil {
		return map[string]interface{}{}
	}
	types := value.Type()
	// 解析到Parse 字段
	result := make(map[string]interface{})
	ch := make(chan map[string]interface{})
	defer close(ch)
	go worker(ch, result)
	wg := sync.WaitGroup{}
	wg.Add(value.NumField())
	for i := 0; i < value.NumField(); i++ {

		go func(v reflect.Value, i int) {
			//v := value.Field(i)
			tag := types.Field(i).Tag.Get(TAGNAME)
			if strings.Contains(tag, ",") {
				tag = strings.TrimSpace(strings.Split(tag, ",")[0])
			}
			switch v.Kind() {
			case reflect.Ptr:
				res := parseStruct(v.Interface())
				//result[tag] = res
				ch <- map[string]interface{}{tag: res}
			case reflect.Slice:
				if !v.IsNil() && v.Index(0).Kind() != reflect.Struct && v.Index(0).Kind() != reflect.Ptr {
					// result[tag] = v.Interface()
					ch <- map[string]interface{}{tag: v.Interface()}
				} else {
					wg := sync.WaitGroup{}
					wg.Add(v.Len())
					slice := make([]map[string]interface{}, 0)
					for i := 0; i < v.Len(); i++ {
						go func(item reflect.Value) {
							res := parseStruct(item.Interface())
							slice = append(slice, res)
							wg.Done()
						}(v.Index(i))
					}
					wg.Wait()
					//result[tag] = slice
					ch <- map[string]interface{}{tag: slice}
				}
			//case reflect.
			default:
				//result[tag] = v.Interface()
				if v.IsValid() {
					ch <- map[string]interface{}{tag: v.Interface()}
				} else {
					ch <- map[string]interface{}{}
				}
			}

			wg.Done()
		}(value.Field(i), i)
	}
	wg.Wait()
	ch <- nil
	return result
}

// 移除器中-递归删除
func roll(arr []string, sc map[string]interface{}, idx *int64) {
	//_, ok := sc[arr[*idx]]
	//if !ok {
	//	panic(ErrorDeleteKey)
	//}
	if *idx < int64(len(arr)-1) {
		d := sc
		m, ok := d[arr[*idx]].(map[string]interface{})
		if !ok {
			fmt.Printf("未找到的key：%#v\n", arr[*idx])
			return
		}
		next := *idx + 1
		roll(arr[*idx:], m, &next)
	} else {
		delete(sc, arr[*idx])
	}
}

// 反射修正
func correction(value reflect.Value) (reflect.Value, error) {
	if value.IsZero() {
		return value, ErrorStructNilPointer
	}
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value, nil
}

// 管道接收
func worker(ch chan map[string]interface{}, result map[string]interface{}) {
	for {
		select {
		case res, ok := <-ch:
			if !ok {
				return
			}
			for k, v := range res {
				result[k] = v
			}
		}
	}
}
