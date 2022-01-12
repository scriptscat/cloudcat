package executor

import (
	"fmt"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"rogchap.com/v8go"
)

func TestGmXmlHttpRequest(t *testing.T) {
	iso, _ := NewExecutor(WithLogger(logrus.StandardLogger().Logf), GmNotification(), GmXmlHttpRequest(nil), Console())

	ret, err := iso.Run(`
function app() {
    return new Promise(resolve => {
        GM_xmlhttpRequest({
            method: 'GET',
            url: 'https://api.bilibili.com/x/web-interface/nav',
            headers: {
                "Referer": "https://www.bilibili.com",
                "Origin": "https://www.bilibili.com"
            },
            onload(xhr) {
				console.log("qwe123",xhr,xhr.response,xhr.response.code);
                switch (xhr.response.code) {
                    case -101:
                        GM_notification({
                            title: 'bilibili自动签到 - ScriptCat',
                            text: '哔哩哔哩签到失败,账号未登录,请先登录',
                        });
                        resolve();
                        break;
                    default:
                        resolve();
                }
            },
            onerror() {
                console.log("error");
                resolve();
            }
        })
    })
}
app();
`)
	assert.Nil(t, err)
	p, err := ret.AsPromise()
	assert.Nil(t, err)
	l := sync.WaitGroup{}
	l.Add(1)
	p.Then(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		l.Done()
		return nil
	})
	p.Catch(func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		l.Done()
		return nil
	})
	l.Wait()
	assert.Equal(t, v8go.Fulfilled, p.State())
	fmt.Println(p.Result().String())
}
