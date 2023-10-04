package scriptcat

import (
	"errors"
	"regexp"
	"strings"
)

func ParseMeta(script string) string {
	return regexp.MustCompile(`// ==UserScript==([\\s\\S]+?)// ==/UserScript==`).FindString(script)
}

func ParseMetaToJson(meta string) map[string][]string {
	reg := regexp.MustCompile("(?im)^//\\s*@(.+?)([\r\n]+|$|\\s+(.+?)$)")
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

// ConvCron 转换cron表达式
func ConvCron(cron string) (string, error) {
	// 对once进行处理
	unit := strings.Split(cron, " ")
	if len(unit) == 5 {
		unit = append([]string{"0"}, unit...)
	}
	if len(unit) != 6 {
		return "", errors.New("cron format error: " + cron)
	}
	for i, v := range unit {
		if v == "once" {
			unit[i] = "*"
			i -= 1
			for ; i >= 0; i-- {
				if unit[i] == "*" {
					unit[i] = "0"
				}
			}
			break
		} else if strings.Contains(v, "-") {
			// 取最小的时间
			unit[i] = strings.Split(v, "-")[0]
		}
	}
	return strings.Join(unit, " "), nil
}
