# CloudCat

> 一个用于 **[ScriptCat脚本猫](https://docs.scriptcat.org/)** 扩展云端执行脚本的服务

![](https://img.shields.io/github/stars/scriptscat/cloudcat.svg)![](https://img.shields.io/github/v/tag/scriptscat/cloudcat.svg?label=version&sort=semver)

## 安装

### linux

```bash
curl -sSL https://github.com/scriptscat/cloudcat/raw/main/deploy/install.sh | sudo bash
```

## 使用

```bash
# 查看帮助
ccatctl -h
# 安装脚本
ccatctl install -f example/bing\ check-in.js
# 查看脚本列表
ccatctl get script
# 导入cookie/value(storage name)
ccatctl import cookie 6a0bd33 example/cookie.json
# 运行脚本(脚本id)
ccatctl run 6a0bd33
```
