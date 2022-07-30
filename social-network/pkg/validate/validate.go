package validate

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"social-network/pkg/logger"
	"strconv"
	"strings"
)

func Struct(domain interface{}, structWData interface{}) interface{} {
	val := reflect.ValueOf(domain) // could be any underlying type

	// if its a pointer, resolve its value
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	// should double-check we now have a struct (could still be anything)
	if val.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}

	// now we grab our values as before (note: I assume table name should come from the struct type)
	structType := val.Type()

	m := map[string]interface{}{}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag

		fieldName := field.Name
		fieldType := tag.Get("validate")

		m[fieldName] = fieldType
	}

	strWData := fmt.Sprintf("%+v\n", structWData)

	return toJson(match(m, strWData))
}

func match(keysWithRules map[string]interface{}, strWithData string) map[string][]map[string]string {
	errors := make(map[string][]map[string]string)

	for key, value := range keysWithRules {
		v := strings.Split(fmt.Sprintf("%v", value), ",")
		for _, s := range v {

			re := regexp.MustCompile(fmt.Sprintf(`(%s):([^\s|}]+)?`, key))

			findInString := re.FindString(strWithData)

			separated := strings.Split(findInString, ":")

			if strings.Contains(s, "required") && key == separated[0] {
				if !required(key, separated[1], s) {
					errorM := make(map[string]string)
					errorM["msg"] = "Field is required"
					errorM["param"] = key
					errors["errors"] = append(errors["errors"], errorM)
				}
			}
			if strings.Contains(s, "max") && key == separated[0] {
				b, maxVal := max(key, separated[1], s)
				if !b {
					errorM := make(map[string]string)
					errorM["msg"] = "Max field is " + maxVal
					errorM["param"] = key
					errors["errors"] = append(errors["errors"], errorM)
				}

			}
			if strings.Contains(s, "min") && key == separated[0] {
				b, minVal := min(key, separated[1], s)
				if !b {
					errorM := make(map[string]string)
					errorM["msg"] = "Min field is " + minVal
					errorM["param"] = key
					errors["errors"] = append(errors["errors"], errorM)
				}
			}
		}
	}

	return errors
}

func toJson(errStruct map[string][]map[string]string) interface{} {
	jsonStr, err := json.Marshal(errStruct)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	rawIn := json.RawMessage(jsonStr)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var s interface{}
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		panic(err)
	}

	return s
}

func required(key, val, rule string) bool {
	return len(val) > 0
}
func max(key, val, rule string) (bool, string) {
	separate := strings.Split(rule, "=")

	m, err := strconv.Atoi(separate[1])
	if err != nil {
		logger.ErrorLogger.Println(err)
		return false, ""
	}
	return len(val) <= m, separate[1]
}
func min(key, val, rule string) (bool, string) {
	separate := strings.Split(rule, "=")

	m, err := strconv.Atoi(separate[1])
	if err != nil {
		logger.ErrorLogger.Println(err)
		return false, ""
	}
	return len(val) >= m, separate[1]
}
