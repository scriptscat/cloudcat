package scriptcat

import (
	"regexp"
	"strings"
)

func ParseMeta(script string) string {
	return regexp.MustCompile("\\/\\/ ==UserScript==([\\s\\S]+?)\\/\\/ ==\\/UserScript==").FindString(script)
}

func ParseMetaToJson(meta string) map[string][]string {
	reg := regexp.MustCompile("(?im)^//\\s*@(.+?)($|\\s+(.+?)$)")
	list := reg.FindAllStringSubmatch(meta, -1)
	ret := make(map[string][]string)
	for _, v := range list {
		v[1] = strings.ToLower(v[1])
		if _, ok := ret[v[1]]; !ok {
			ret[v[1]] = make([]string, 0)
		}
		ret[v[1]] = append(ret[v[1]], strings.TrimSpace(v[3]))
	}
	return ret
}
