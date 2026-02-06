package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(data any) {
	jsonPretty, _ := json.MarshalIndent(data, "", "\t")
	fmt.Printf("%s\n", jsonPretty)
}
