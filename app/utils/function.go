package utils

import (
	"encoding/json"
	"fmt"
	"goframe/app/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

/**
 * [获取集合所有key]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func GetMpaKeys(m map[string]string) []string {
	j := 0
	keys := make([]string, len(m))

	for k := range m {
		keys[j] = k
		j++
	}

	return keys
}

/**
 * [获取集合所有value]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func GetMpaValues(m map[string]string) []string {
	j := 0
	values := make([]string, len(m))

	for _, v := range m {
		values[j] = v
		j++
	}

	return values
}

/**
 * [获取真实客户ip]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func GetClientIp(r *http.Request) string {
	IPAddress := r.Header.Get("X-Forwarded-Fo")

	if len(IPAddress) > 0 {
		ipArr := strings.Split(IPAddress, ",")
		return ipArr[0]
	}

	IPAddress = r.Header.Get("X-Real-Ip")
	if len(IPAddress) > 0 {
		return IPAddress
	}

	IPAddress = r.RemoteAddr
	if len(IPAddress) > 0 {
		return IPAddress
	}

	return "0.0.0.0"
}

/**
 * [http参数绑定]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func BindParam(r *http.Request, ptr interface{}) error {

	if err := r.ParseForm(); err != nil {
		return err
	}

	// 创建字段映射表，键为有效名称
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("form")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	//解析get参数
	queryForm, _ := url.ParseQuery(r.URL.RawQuery)
	if len(queryForm) > 0 {
		for name, values := range queryForm {
			f := fields[name]

			if !f.IsValid() {
				continue // 忽略不能识别的 HTTP 参数
			}

			for _, value := range values {
				if f.Kind() == reflect.Slice {
					elem := reflect.New(f.Type().Elem()).Elem()
					if err := populate(elem, value); err != nil {
						return fmt.Errorf("%s: %v", name, err)
					}
					f.Set(reflect.Append(f, elem))
				} else {
					if err := populate(f, value); err != nil {
						return fmt.Errorf("%s: %v", name, err)
					}
				}
			}
		}
	}

	//解析post:form-data参数
	r.ParseMultipartForm(32 << 20)
	if r.MultipartForm != nil {
		for name, values := range r.MultipartForm.Value {
			f := fields[name]

			if !f.IsValid() {
				continue // 忽略不能识别的 HTTP 参数
			}

			for _, value := range values {
				if f.Kind() == reflect.Slice {
					elem := reflect.New(f.Type().Elem()).Elem()
					if err := populate(elem, value); err != nil {
						return fmt.Errorf("%s: %v", name, err)
					}
					f.Set(reflect.Append(f, elem))
				} else {
					if err := populate(f, value); err != nil {
						return fmt.Errorf("%s: %v", name, err)
					}
				}
			}
		}
	}

	//解析post:body-json数据
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	if len(body) > 0 {
		err = json.Unmarshal([]byte(body), &ptr)
		if err != nil {
			logger.Logger.Println(err.Error())
		}
	}

	return nil
}

func populate(v reflect.Value, value string) error {

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

/**
 * [接口返回值]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func ReturnJson(w http.ResponseWriter, code int64, msg string, data ...interface{}) {
	type JsonRes struct {
		Code int64       `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	var obj struct {
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var jsonRes JsonRes
	if len(data) == 0 {
		jsonRes = JsonRes{
			Code: code,
			Msg:  msg,
			Data: &obj,
		}
	} else {
		jsonRes = JsonRes{
			Code: code,
			Msg:  msg,
			Data: data,
		}
	}

	res, _ := json.Marshal(&jsonRes)
	fmt.Fprintf(w, "%s", string(res))
}

/**
 * [打印]
 * ------------------------------------------------------------------------------
 * @author  github
 * ------------------------------------------------------------------------------
 * @version date:2999-01-01
 * ------------------------------------------------------------------------------
 */
func Dump(obj interface{}) error {
	if obj == nil {
		fmt.Println("nil")
		return nil
	}
	switch obj.(type) {
	case bool:
		fmt.Println(obj.(bool))
	case int:
		fmt.Println(obj.(int))
	case float64:
		fmt.Println(obj.(float64))
	case string:
		fmt.Println(obj.(string))
	case map[string]interface{}:
		for k, v := range obj.(map[string]interface{}) {
			fmt.Printf("%s: ", k)
			err := Dump(v)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Unsupported type: %v", obj)
	}

	return nil
}
