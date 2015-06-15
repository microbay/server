package core

import (
	"regexp"
)

type URLParams map[string]string

func Params(url string, regex *regexp.Regexp, keys []string) URLParams {
	match := regex.FindAllStringSubmatch(url, -1)[0][1:]
	result := make(URLParams)
	for i := range match {
		if len(keys) <= i {
			break
		}
		result[keys[i]] = match[i]
	}
	return result
}
