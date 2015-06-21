package core

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gocraft/web"
	"net/http"
	"regexp"
	"strings"
	"time"
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

// Thanks to https://gowalker.org/github.com/azer/url-router#PathToRegex
func PathToRegex(path string) (*regexp.Regexp, []string) {
	pattern, _ := regexp.Compile(":([A-Za-z0-9]+)")
	matches := pattern.FindAllStringSubmatch(path, -1)
	keys := []string{}

	for i := range matches {
		keys = append(keys, matches[i][1])
	}

	str := fmt.Sprintf("^%s\\/?$", strings.Replace(path, "/", "\\/", -1))

	str = pattern.ReplaceAllString(str, "([^\\/]+)")
	str = strings.Replace(str, ".", "\\.", -1)

	regex, _ := regexp.Compile(str)

	return regex, keys
}

func ProcessRequestTime(startTime time.Time) (duration int64, durationUnits string) {
	duration = time.Since(startTime).Nanoseconds()
	switch {
	case duration > 2000000:
		durationUnits = "ms"
		duration /= 1000000
	case duration > 1000:
		durationUnits = "μs"
		duration /= 1000
	default:
		durationUnits = "ns"
	}
	return
}

const (
	BACKEND_REQUEST_TYPE_SINGLE   = "single"
	BACKEND_REQUEST_TYPE_COMPOUND = "compound"
)

func LogRequest(rw web.ResponseWriter, req *web.Request, startTime time.Time) {

	duration, durationUnits := ProcessRequestTime(startTime)

	if rw.StatusCode() > 499 {
		log.WithFields(log.Fields{
			"type": "request",
			//"req":      req,
			"status":   rw.StatusCode(),
			"duration": fmt.Sprintf("%d%s", duration, durationUnits),
		}).Errorf("%d %s %s", rw.StatusCode(), req.Method, req.URL.Path)
	} else {
		log.WithFields(log.Fields{
			"type": "request",
			//"req":      req,
			"status":   rw.StatusCode(),
			"duration": fmt.Sprintf("%d%s", duration, durationUnits),
		}).Infof("%d %s %s", rw.StatusCode(), req.Method, req.URL.Path)
	}
}

func LogBackendRequest(err error, res *http.Response, method string, url string, startTime time.Time) {

	duration, durationUnits := ProcessRequestTime(startTime)

	var status int
	if res != nil {
		status = res.StatusCode
	}

	if err != nil || status > 499 {
		log.WithFields(log.Fields{
			"type": "backend.error",
			//"res":      res,
			"status":   status,
			"duration": fmt.Sprintf("%d%s", duration, durationUnits),
		}).Warnf("Backend: %d %s %s", status, method, url)
	} else {
		log.WithFields(log.Fields{
			"type": "backend.success",
			//"res":      res,
			"status":   res.StatusCode,
			"duration": fmt.Sprintf("%d%s", duration, durationUnits),
		}).Infof("Backend: %d %s %s", res.StatusCode, res.Request.Method, res.Request.URL.Path)
	}
}
