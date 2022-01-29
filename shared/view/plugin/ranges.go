package plugin

import (
	"fmt"
	"html/template"
)

func Ranges() template.FuncMap {
	f := make(template.FuncMap)

	f["RN"] = func(t int) []int {
		var values []int
		for i := 1; i <= t; i++ {
			values = append(values, i)
		}
		return values
	}

	return f
}
func Add() template.FuncMap {
	f := make(template.FuncMap)

	f["add"] = func(t int) int {

		return t + 1
	}

	return f
}

func GB() template.FuncMap {
	f := make(template.FuncMap)

	f["gb"] = func(t uint64) int {

		return int(t / 1024 / 1024 / 1024)
	}

	return f
}
func Uptime() template.FuncMap {
	f := make(template.FuncMap)

	f["uptime"] = func(t uint64) string {
		timeInMilliSeconds := t
		seconds := timeInMilliSeconds / 1000
		minutes := seconds / 60
		hours := minutes / 60
		//days := hours / 24

		//days, hours%24, minutes%60, seconds%60
		return fmt.Sprintf("hours%d", hours)
	}

	return f
}

func Sub() template.FuncMap {
	f := make(template.FuncMap)

	f["sub"] = func(t int) int {

		return t - 1
	}

	return f
}
