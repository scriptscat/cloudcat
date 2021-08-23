package scriptcat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMetaToJson(t *testing.T) {
	ret := ParseMetaToJson(`// ==UserScript==
// @name         bilibili自动签到
// @namespace    wyz
// @version      1.1.2
// @author       wyz
// @crontab * * once * *
// @debug
// @grant GM_xmlhttpRequest
// @grant GM_notification
// @connect api.bilibili.com
// @connect api.live.bilibili.com
// @cloudCat
// @exportCookie domain=api.bilibili.com
// @exportCookie domain=api.live.bilibili.com
// ==/UserScript==`)

	assert.Equal(t, 2, len(ret["grant"]))

}
