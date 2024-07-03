package main

import (
	"fmt"
	"regexp"
)

// var typeMap = map[string]string {
// 	"uint32": "number",
// 	"string": "string",
// 	""
// }

var NAME_PATTERN = regexp.MustCompile("json:\"(.+?)\"")

func convTsName(field FieldInfo) string {
	name := field.Name
	if field.Comment != "" {
		r := NAME_PATTERN.FindStringSubmatch(field.Comment)
		if len(r) == 2 {
			name = r[1]
		}
	}
	if field.Type[0] == "*" {
		return fmt.Sprintf("%s?", name)
	}
	return name
}

func convTsType(fieldTypes []string) string {
	switch fieldTypes[0] {
	case "uint32":
		return "number"
	case "time":
		return "string"
	case "*":
		return convTsType(fieldTypes[1:])
	case "[]":
		t := convTsType(fieldTypes[1:])
		return fmt.Sprintf("[]%s", t)
	}
	return fieldTypes[0]
}

func makeTs(goType *TypeInfo) {
	result := fmt.Sprintf("export type %s = {", goType.Name)
	for _, field := range goType.Fields {
		result = fmt.Sprintf("%s\n    %s: %s;", result, convTsName(field), convTsType(field.Type))
	}
	result = fmt.Sprintf("%s\n};", result)
	fmt.Println(result)
}
