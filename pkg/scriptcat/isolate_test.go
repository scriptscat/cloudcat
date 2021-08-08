package scriptcat

import (
	"testing"
)

func TestIsolate(t *testing.T) {
	iso, _ := NewIsolate()

	ctx, _ := NewContext(iso, GmXmlHttpRequest())

	iso.Run(ctx, "GM_xmlhttpRequest({url:'https://bbs.tampermonkey.net.cn/'})")

}
