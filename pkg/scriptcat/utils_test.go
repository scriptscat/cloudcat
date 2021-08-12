package scriptcat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseMeta(t *testing.T) {
	ret := ParseMeta(`// ==UserScript==
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
// ==/UserScript==
script code
`)

	assert.Equal(t, `// ==UserScript==
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
// ==/UserScript==`, ret)

}
