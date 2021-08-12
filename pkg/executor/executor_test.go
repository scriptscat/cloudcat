package executor

import (
	"testing"
)

func TestIsolate(t *testing.T) {
	iso, _ := NewExecutor()

	ctx, _ := NewContext(iso, GmXmlHttpRequest(nil))

	iso.Run(ctx, "GM_xmlhttpRequest({url:'https://bbs.tampermonkey.net.cn/'})")

}
