// ==UserScript==
// @name         value test
// @namespace    https://bbs.tampermonkey.net.cn/
// @version      0.1.0
// @description  try to take over the world!
// @author       You
// @crontab      */12 * * * * *
// @grant GM_setValue
// @grant GM_getValue
// ==/UserScript==

return new Promise((resolve) => {
    setTimeout(() => {
        // Your code here...
        GM_setValue("obj", {"test": 1});
        GM_setValue("arr", ["test", "test2"]);
        GM_setValue("bool", true);
        GM_setValue("num1", 12345);
        GM_setValue("num2", 123.45);
        GM_setValue("str", "string");

        GM_log(GM_getValue("obj", 1), "warn", {"test": 1});
        GM_log(GM_getValue("arr", 1), "warn", {"test": 1});
        GM_log(GM_getValue("bool", 1), "warn", {"test": 1});
        GM_log(GM_getValue("num1", 1), "warn", {"test": 1});
        GM_log(GM_getValue("num2", 1), "warn", {"test": 1});
        GM_log(GM_getValue("str", 1), "warn", {"test": 1});

        resolve();
    }, 2000);
});
