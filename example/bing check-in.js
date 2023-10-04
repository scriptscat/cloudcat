// ==UserScript==
// @name         必应积分商城每日签到
// @namespace    wyz
// @description  每日自动完成任务获取必应积分奖励，可兑换实物
// @version      1.2.0
// @author       wyz
// @crontab      * 10-23 once * *
// @grant        GM_xmlhttpRequest
// @grant        GM_notification
// @grant        GM_getValue
// @grant        GM_setValue
// @connect      bing.com
// @connect      top.baidu.com
// @connect      sct.icodef.com
// @require      https://scriptcat.org/lib/946/%5E1.0.1/PushCat.js
// @exportValue  PushCat.AccessKey
// @cloudCat
// @exportCookie domain=.bing.com
// ==/UserScript==

/* ==UserConfig==
PushCat:
  AccessKey:
    title: 消息推送key
    description: 消息推送key https://sct.icodef.com/
    type: text
 ==/UserConfig== */

let accessKey = GM_getValue("PushCat.AccessKey");
if (accessKey == "") {
    // 由于脚本猫v0.13的UserConfig有bug，如果你需要消息推送服务的话，请在此手动设置
    GM_setValue("PushCat.AccessKey", "");
    accessKey = GM_getValue("PushCat.AccessKey");
}
// 消息推送: https://sct.icodef.com/
const push = new PushCat({
    accessKey,
});

function getMobileUA() {
    // 手机ua列表
    const ua = [
        // Android
        "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36 Edg/113.0.0.0",
        "Mozilla/5.0 (Linux; Android 12; HarmonyOS; CTR-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 11; AVA-PA00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36",
        "Mozilla/5.0 (Linux; Android 12; CPH2209) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.101 Mobile Safari/537.36",
        // iOS
        "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/113.0.0.0",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 16_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/99.0.4844.59 Mobile/15E148 Safari/604.1",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/97.0.4692.84 Mobile/15E148 Safari/604.1",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 15_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/97.0.4692.84 Mobile/15E148 Safari/604.1",
    ];
    let n = GM_getValue("mobileUA",);
    if (!n) {
        n = Math.floor(Math.random() * ua.length);
        GM_setValue("mobileUA", n);
    }
    // 随机数
    return ua[n];
}


function pushSend(title, content) {
    GM_log("推送消息", "info", {title, content});
    return new Promise(async resolve => {
        if (accessKey) {
            await push.send(title, content);
        }
        GM_notification({
            title: title,
            text: content,
        });
        resolve();
    })
}

function getSubstring(inputStr, startStr, endStr) {
    const startIndex = inputStr.indexOf(startStr);
    if (startIndex == -1) {
        return null;
    }
    const endIndex = inputStr.indexOf(endStr, startIndex + startStr.length);
    if (endIndex == -1) {
        return null;
    }
    return inputStr.substring(startIndex + startStr.length, endIndex);
}


function getRewardsInfo() {
    return new Promise((resolve, reject) => {
        // 获取今日签到信息
        GM_xmlhttpRequest({
            url: "https://rewards.bing.com",
            onload(resp) {
                if (resp.status == 200) {
                    resolve(resp);
                } else {
                    pushSend("必应每日签到失败", "请求返回错误: " + resp.status).then(() => reject());
                }
            }, onerror(e) {
                pushSend("必应每日签到失败", e || "未知错误").then(() => reject());
            }
        });
    })
}

function extractKeywords(inputStr) {
    const regex = /"indexUrl":"","query":"(.*?)"/g;
    const matches = [...inputStr.matchAll(regex)];
    return matches.map(match => match[1]);
}

let keywordList = [];
let keywordIndex = 0;

// 获取搜索关键字
function searchKeyword() {
    return new Promise((resolve, reject) => {
        if (keywordList.length == 0) {
            GM_xmlhttpRequest({
                url: "https://top.baidu.com/board?platform=pc&sa=pcindex_entry",
                onload(resp) {
                    if (resp.status == 200) {
                        keywordList = extractKeywords(resp.responseText);
                        resolve(keywordList[keywordIndex]);
                    } else {
                        pushSend('关键字获取失败', '热门词获取失败');
                        reject(new Error('关键字获取失败,' + resp.status));
                    }
                }
            });
        } else {
            keywordIndex++;
            if (keywordIndex > keywordList.length) {
                keywordIndex = 0;
            }
            resolve(keywordList[keywordIndex]);
        }
    }).then(k => k + new Date().getTime() % 1000);
}

let retryNum = 0;
let lastProcess = 0;
let domain = "www.bing.com";

function handler() {
    return getRewardsInfo().then(async resp => {
        // 获取今日已获取积分
        const data = resp.responseText;
        const dashboard = JSON.parse(getSubstring(data, "var dashboard = ", ";\r"));
        const pcAttributes = dashboard.userStatus.counters.pcSearch[0].attributes;
        if (dashboard.userStatus.counters.dailyPoint[0].pointProgress === lastProcess) {
            retryNum++;
            if (retryNum > 10) {
                await pushSend("必应每日签到错误", "请手动检查积分或者重新执行");
                return true;
            }
        } else {
            lastProcess = dashboard.userStatus.counters.dailyPoint[0].pointProgress;
        }
        if (parseInt(pcAttributes.progress) >= parseInt(pcAttributes.max)) {
            // 判断是否有手机
            if (dashboard.userStatus.counters.mobileSearch) {
                const mobileSearch = dashboard.userStatus.counters.mobileSearch[0].attributes;
                if (parseInt(mobileSearch.progress) < parseInt(mobileSearch.max)) {
                    // 进行一次手机搜索
                    GM_xmlhttpRequest({
                        url: "https://" + domain + "/search?q=" + await searchKeyword(),
                        onload(resp) {
                            const url = new URL(resp.finalUrl);
                            if (url.host != domain) {
                                domain = url.host;
                            }
                        },
                        headers: {
                            "User-Agent": getMobileUA()
                        }
                    });
                    return false;
                }
                GM_log("奖励信息", "info", {
                    pcProcess: pcAttributes.progress,
                    mobileProcess: mobileSearch.progress
                });
            } else {
                GM_log("奖励信息", "info", {pcProcess: pcAttributes.progress});
            }
            await pushSend("必应每日签到完成", "当前等级: " + dashboard.userStatus.levelInfo.activeLevel +
                "(" + dashboard.userStatus.levelInfo.progress + ")" +
                "\n可用积分: " + dashboard.userStatus.availablePoints + " 今日积分: " + dashboard.userStatus.counters.dailyPoint[0].pointProgress);
            return true;
        } else {
            // 进行一次搜索
            GM_xmlhttpRequest({
                url: "https://" + domain + "/search?q=" + await searchKeyword(),
                onload(resp) {
                    const url = new URL(resp.finalUrl);
                    if (url.host != domain) {
                        domain = url.host;
                    }
                }
            });
            return false;
        }
    });
}

return new Promise((resolve, reject) => {
    const h = async () => {
        try {
            const result = await handler();
            if (result) {
                resolve();
            } else {
                setTimeout(() => {
                    h();
                }, 1000 * (Math.floor(Math.random() * 4) + 10));
            }
        } catch (e) {
            pushSend('必应每日签到失败', '请查看错误日志手动重试');
            reject(e);
        }
    }
    h();
});

