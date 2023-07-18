// ==UserScript==
// @name         xhr test
// @namespace    https://bbs.tampermonkey.net.cn/
// @version      0.1.0
// @description  try to take over the world!
// @author       You
// @background
// @grant GM_xmlhttpRequest
// @connect  bbs.tampermonkey.net.cn
// ==/UserScript==

return new Promise((resolve, reject) => {
    GM_xmlhttpRequest({
        url: "https://bbs.tampermonkey.net.cn/",
        method: "POST",
        data: "test",
        headers: {
            "referer": "http://www.example.com/",
            "origin": "www.example.com",
            // 为空将不会发送此header
            "sec-ch-ua-mobile": "",
        },
        onload(resp) {
            GM_log("onload", "info", {resp: resp});
        },
        onreadystatechange(resp) {
            GM_log("onreadystatechange", "info", {resp: resp});
        },
        onloadend(resp) {
            GM_log("onloadend", "info", {resp: resp});
        },
    });
    resolve();
});
