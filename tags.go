package main

import "strings"

func fieldTagValues(tag, tagName string) []string {
	tagPrefix := tagName + ":"
	if !strings.Contains(tag, tagPrefix) {
		return nil
	}
	tagValue := strings.Split(tag, tagPrefix)[1]
	leftQ := strings.Index(tagValue, `"`)
	if leftQ == -1 || leftQ == len(tagValue)-1 {
		return nil
	}
	rightQ := strings.Index(tagValue[leftQ+1:], `"`)
	if rightQ == -1 {
		return nil
	}
	tagValue = tagValue[leftQ+1 : leftQ+rightQ+1]
	return strings.Split(tagValue, ",")
}

func getAllTagValues(tag string) map[string][]string {
	tags := make(map[string][]string)
	for _, t := range strings.Fields(tag) {
		if !strings.Contains(t, ":") {
			continue
		}
		parts := strings.Split(t, ":")
		tags[parts[0]] = append(tags[parts[0]], parts[1])
	}
	return tags
}

type FieldTag map[string][]string

func ParseFieldTag(tag string) FieldTag {
	t := make(FieldTag)
	for _, fields := range strings.Fields(tag) {
		if !strings.Contains(fields, ":") {
			continue
		}
		parts := strings.Split(fields, ":")
		key := parts[0]
		vals := fieldTagValues(tag, key)
		t[key] = vals
	}
	return t
}

func (t FieldTag) GetAll(key string) []string {
	return t[key]
}

func (t FieldTag) GetFirst(key string) (string, bool) {
	if len(t[key]) == 0 {
		return "", false
	}
	return t[key][0], true
}
