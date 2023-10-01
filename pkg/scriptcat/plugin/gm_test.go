package plugin

import (
	"fmt"
	"net/url"
	"testing"
)

func TestName(t *testing.T) {
	t.Logf("A")
	a := "https://www.baidu.com/?q=你好&b=%E4%BD%A0%E5%A5%BD"
	u, _ := url.Parse(a)
	fmt.Println(u, u.String(), u.RequestURI(), u.Query().Encode())
	t.Log(a)
}
