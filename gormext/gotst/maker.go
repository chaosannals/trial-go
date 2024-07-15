package main

import (
	"fmt"
	"regexp"
)

var NAME_PATTERN = regexp.MustCompile("json:\"(.+?)\"")

func convTsName(field *GoField) string {
	name := field.Name
	if field.Tag != "" {
		r := NAME_PATTERN.FindStringSubmatch(field.Tag)
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
	case "int32", "uint32", "int64", "uint64", "float32", "float64":
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

func convTsComment(field *GoField) string {
	if len(field.Comment2) > 0 {
		return fmt.Sprintf("\n/* %s */", field.Comment2)
	} else {
		return fmt.Sprintf("// %s", field.Comment)
	}
}

func makeTs(goType *GoStruct) {
	result := fmt.Sprintf("export type %s = {", goType.Name)
	for _, field := range goType.Fields {
		result = fmt.Sprintf(
			"%s\n    %s: %s; %s",
			result,
			convTsName(&field),
			convTsType(field.Type),
			convTsComment(&field),
		)
	}
	result = fmt.Sprintf("%s\n};", result)
	fmt.Println(result)
}
