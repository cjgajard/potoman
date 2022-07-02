package potoman

import (
	"fmt"
	"regexp"
)

var (
	isHttpVer = regexp.MustCompile(`HTTP/(?P<major>\d)(?:\.(?P<minor>\d))?$`)
)

// getHttpVersion returns an Attr with the "major" number saved in `Key` and
// "minor" number in `Value`. If `content` is empty, it return version "1.0",
// otherwise it returns an error when `content` has not a HTTP/x.y pattern.
// ```
// getHttpVersion([]("HTTP/2.0") // => Attr{Key: "2", Value: "0"}
// ```
func getHttpVersion(content []byte) (Attr, error) {
	v := Attr{"1", "0"}
	if len(content) == 0 {
		return v, nil
	}
	matches := isHttpVer.FindSubmatch(content)
	if len(matches) >= 3 {
		v.Key = string(matches[isHttpVer.SubexpIndex("major")])
		if m := matches[isHttpVer.SubexpIndex("minor")]; m != nil {
			v.Value = string(m)
		}
		return v, nil
	}
	return v, fmt.Errorf("invalid prototype version %q", content)
}
