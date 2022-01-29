package plugin

import (
	"html/template"
)

func Ranges() template.FuncMap {
	f := make(template.FuncMap)

	f["RN"] = func(t int) [] int {
		var values []int
		for i:=1;i <= t ;i++  {
			values = append(values,i)
		}
		return values
	}

	return f
}
func Add() template.FuncMap {
	f := make(template.FuncMap)

	f["add"] = func(t int)  int {

		return t+1
	}

	return f
}



func Sub() template.FuncMap {
	f := make(template.FuncMap)

	f["sub"] = func(t int)  int {

		return t-1
	}

	return f
}
