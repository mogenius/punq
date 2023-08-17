package structs

import (
	"sort"
	"strings"
)

type UniqueStringArray struct {
	helperDict map[string]bool
	result     []string
}

func NewUniqueStringArray() UniqueStringArray {
	d := UniqueStringArray{}
	d.helperDict = map[string]bool{}
	d.result = []string{}
	return d
}

func (d *UniqueStringArray) Add(s string) {
	if d.helperDict[s] {
		return // Already in the map
	}
	d.result = append(d.result, s)
	d.helperDict[s] = true
}

func (d *UniqueStringArray) Display() string {
	sort.Strings(d.result)
	return strings.Join(d.result, ", ")
}
