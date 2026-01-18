package tags

import "strings"

type FieldTag map[string]string

func ParseFieldTag(tag string) FieldTag {
	t := make(FieldTag)
	for _, fields := range strings.Fields(tag) {
		if !strings.Contains(fields, ":") {
			continue
		}
		parts := strings.Split(fields, ":")
		key := parts[0]
		if vals := fieldTagValues(tag, key); vals != "" {
			t[key] = vals
		}
	}
	return t
}

func (t FieldTag) GetAll(key string) []string {
	if val, ok := t[key]; ok {
		return strings.Split(val, ",")
	}
	return []string{}
}

func (t FieldTag) GetFirst(key string) (string, bool) {
	if val, ok := t[key]; ok {
		parts := strings.SplitN(val, ",", 2)
		return parts[0], true
	}
	return "", false
}

func (t FieldTag) GetString(key string) (string, bool) {
	if val, ok := t[key]; ok {
		return val, true
	}
	return "", false
}

func fieldTagValues(tag, tagName string) string {
	tagPrefix := tagName + ":"
	if !strings.Contains(tag, tagPrefix) {
		return ""
	}
	tagValue := strings.Split(tag, tagPrefix)[1]
	leftQ := strings.Index(tagValue, `"`)
	if leftQ == -1 || leftQ == len(tagValue)-1 {
		return ""
	}
	rightQ := strings.Index(tagValue[leftQ+1:], `"`)
	if rightQ == -1 {
		return ""
	}
	return tagValue[leftQ+1 : leftQ+rightQ+1]
}
