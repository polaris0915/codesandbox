package json

import "github.com/goccy/go-json"

func JsonToModel(jsonStr string, obj any) error {
	if err := json.Unmarshal([]byte(jsonStr), obj); err != nil {
		return err
	}
	// fmt.Printf("obj: %+v", obj)
	return nil
}

func RawModelToJson(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func ModelToJson(obj any) (string, error) {
	data, err := RawModelToJson(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
