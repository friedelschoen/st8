package format

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func UnmarshalConf(data map[string]string, out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("output must be a non-nil pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("output must point to a struct")
	}

	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("conf")
		if tag == "" {
			continue
		}

		valStr, ok := data[tag]
		if !ok {
			continue // Key not found, skip
		}

		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		// Custom method support: UnmarshalConf(string) error
		method := fv.Addr().MethodByName("UnmarshalConf")
		if method.IsValid() {
			res := method.Call([]reflect.Value{reflect.ValueOf(valStr)})
			if errVal := res[0]; !errVal.IsNil() {
				return fmt.Errorf("error unmarshalling %q: %w", tag, errVal.Interface().(error))
			}
			continue
		}

		// Parse built-in types
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(valStr)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid int for %q: %w", tag, err)
			}
			fv.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(valStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid uint for %q: %w", tag, err)
			}
			fv.SetUint(n)
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return fmt.Errorf("invalid float for %q: %w", tag, err)
			}
			fv.SetFloat(f)
		case reflect.Bool:
			switch strings.ToLower(valStr) {
			case "yes", "y", "true":
				fv.SetBool(true)
			case "no", "n", "false":
				fv.SetBool(false)
			default:
				return fmt.Errorf("invalid bool for %q: %s", tag, valStr)
			}
		default:
			return fmt.Errorf("unsupported field type %s for key %q", fv.Type(), tag)
		}
	}

	return nil
}
