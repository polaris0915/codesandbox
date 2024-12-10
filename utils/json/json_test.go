package json

import (
	"strconv"
	"testing"
)

func TestModelToJson(t *testing.T) {
	testStruct := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "john",
		Age:  18,
	}
	str, err := ModelToJson(testStruct)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ModelToJson(testStruct): %s\n", strconv.Quote(str))
}
