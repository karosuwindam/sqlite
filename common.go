package sqlite

import (
	"errors"
	"reflect"
)

// structToSlice(str) = interface{}
//
// 構造体から配列構造体に変換するツール
//
// str(interface{}) : ベースになる構造体
//
// 例文 data := structToSlice(Struct{}).([]Struct)
func structToSlice(str interface{}) interface{} {
	tSlice := reflect.SliceOf(reflect.TypeOf(str))
	vSlice := reflect.MakeSlice(tSlice, 0, 0)
	return vSlice.Interface()
}

// mapToStruct(s,i) = error
//
// map形式のデータから構造体のポインタ配列データに追加する
//
// s(map[string]interface{}) : 入力用のmap形式データ
// i(*[]interface{}) : 格納先のポインター配列、構造体
func mapToStruct(s map[string]interface{}, i interface{}) error {

	sv := reflect.ValueOf(i)
	if sv.Type().Kind() != reflect.Ptr {
		return errors.New("Don't struct pointer input i=" + sv.Type().Kind().String())
	}
	if len(s) == 0 {
		return nil
	}
	ii := sv.Elem().Interface()
	if reflect.TypeOf(ii).Kind() != reflect.Slice {
		return errors.New("Don't Slice input *i=" + reflect.TypeOf(ii).Kind().String())
	}
	tStruct := reflect.TypeOf(ii).Elem()
	vStruct := reflect.New(tStruct)
	ckStruct := reflect.TypeOf(vStruct.Elem().Interface())
	for i := 0; i < ckStruct.NumField(); i++ {
		f := ckStruct.Field(i)
		v := vStruct.Elem().FieldByName(f.Name)
		ss := s[f.Name]
		switch f.Type.Kind() {
		case reflect.Int & reflect.TypeOf(ss).Kind():
			v.SetInt(int64(ss.(int)))
		case reflect.String & reflect.TypeOf(ss).Kind():
			v.SetString(ss.(string))
		}
	}
	out := vStruct.Elem()
	v := sv.Elem()
	v.Set(reflect.Append(v, out))
	return nil
}
