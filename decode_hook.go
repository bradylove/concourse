package atc

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

var SanitizeDecodeHook = func(
	dataKind reflect.Kind,
	valKind reflect.Kind,
	data interface{},
) (interface{}, error) {
	if valKind == reflect.Map {
		if dataKind == reflect.Map {
			return sanitize(data)
		}
	}

	if valKind == reflect.String {
		if dataKind == reflect.String {
			return data, nil
		}

		if dataKind == reflect.Float64 {
			if f, ok := data.(float64); ok {
				return strconv.FormatFloat(f, 'f', -1, 64), nil
			}

			return nil, errors.New("impossible: float64 != float64")
		}

		// format it as JSON/YAML would
		return json.Marshal(data)
	}

	return data, nil
}

func sanitize(root interface{}) (interface{}, error) {
	switch rootVal := root.(type) {
	case map[interface{}]interface{}:
		sanitized := map[string]interface{}{}

		for key, val := range rootVal {
			str, ok := key.(string)
			if !ok {
				return nil, errors.New("non-string key")
			}

			sub, err := sanitize(val)
			if err != nil {
				return nil, err
			}

			sanitized[str] = sub
		}

		return sanitized, nil

	case []interface{}:
		sanitized := make([]interface{}, len(rootVal))
		for i, val := range rootVal {
			sub, err := sanitize(val)
			if err != nil {
				return nil, err
			}
			sanitized[i] = sub
		}
		return sanitized, nil

	default:
		return rootVal, nil
	}
}
