package json

import "github.com/goccy/go-json"

func JsonToModel(jsonStr string, obj any) error {
	if err := json.Unmarshal([]byte(jsonStr), obj); err != nil {
		return err
	}
	// fmt.Printf("obj: %+v", obj)
	return nil
}
