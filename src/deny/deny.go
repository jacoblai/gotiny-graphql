package deny

import (
	"bytes"
	"encoding/json"
)

var KindWord = []byte("$")

func InjectionPass(word []byte) bool {
	if bytes.Contains(word, KindWord) {
		return false
	}
	return true
}

func Injection(input interface{}) bool {
	bts, _ := json.Marshal(input)
	return InjectionPass(bts)
}

//var KindWord = [][]byte{[]byte("$"), []byte("{"), []byte("}")}
//
//func InjectionPass(word []byte) bool {
//	for _, v := range KindWord {
//		if bytes.Contains(word, v) {
//			return false
//		}
//	}
//	return true
//}
//
//func Injection(src interface{}) *bool {
//	res := false
//	if src == nil {
//		return &res
//	}
//	original := reflect.ValueOf(src)
//	res1 := &res
//	loopRecursive(original, &res1)
//
//	return res1
//}
//
//func loopRecursive(src reflect.Value, res **bool) {
//	switch src.Kind() {
//	case reflect.Ptr:
//		originalValue := src.Elem()
//		if !originalValue.IsValid() {
//			return
//		}
//		loopRecursive(originalValue, res)
//	case reflect.Interface:
//		if src.IsNil() {
//			return
//		}
//		originalValue := src.Elem()
//		loopRecursive(originalValue, res)
//	case reflect.Struct:
//		_, ok := src.Interface().(time.Time)
//		if ok {
//			return
//		}
//		for i := 0; i < src.NumField(); i++ {
//			if src.Type().Field(i).PkgPath != "" {
//				continue
//			}
//			loopRecursive(src.Field(i), res)
//		}
//	case reflect.Slice:
//		if src.IsNil() {
//			return
//		}
//		w, ok := src.Interface().([]byte)
//		if ok {
//			re := InjectionPass(w)
//			*res = &re
//			log.Println(string(w), *(*res))
//			return
//		}
//		for i := 0; i < src.Len(); i++ {
//			loopRecursive(src.Index(i), res)
//		}
//	case reflect.Map:
//		if src.IsNil() {
//			return
//		}
//		for _, key := range src.MapKeys() {
//			originalValue := src.MapIndex(key)
//			loopRecursive(originalValue, res)
//		}
//	case reflect.String:
//		re := InjectionPass([]byte(src.Interface().(string)))
//		*res = &re
//		log.Println(src.Interface().(string), *(*res))
//		return
//	}
//}
