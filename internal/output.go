package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

var writer map[string]func([]*Cookie) (string, error) = map[string]func([]*Cookie) (string, error){
	"json": writeJson,
	"curl": writeCurl,
}

func writeJson(c []*Cookie) (string, error) {
	result := map[string][]*Cookie{
		"cookie": c,
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func writeCurl(c []*Cookie) (string, error) {

	var result []string
	for _, v := range c {
		var r string
		if v.Flags == 5 || v.Flags == 4 {
			r += "#HttpOnly_"
		}
		r += fmt.Sprintf("%s\t", v.Domain)
		if v.Domain[0] == '.' {
			r += "TRUE\t"
		} else {
			r += "FALSE\t"
		}
		r += fmt.Sprintf("%s\t", v.Path)
		if v.Flags == 1 || v.Flags == 5 {
			r += "TRUE\t"
		} else {
			r += "FALSE\t"
		}
		result = append(result, r)
	}
	return strings.Join(result, "\n"), nil
}
