

# CloudCat

> 一个用于 **[ScriptCat脚本猫](https://docs.scriptcat.org/)** 扩展云端执行脚本的服务

![](https://img.shields.io/github/stars/scriptscat/cloudcat.svg)![](https://img.shields.io/github/v/tag/scriptscat/cloudcat.svg?label=version&sort=semver)



## 编译



```shell
make build
```
\* Windows编译需要安装Mingw64



## 运行

### ScriptCat运行



```
scriptcat bilibili.zip
```



### Docker运行

> 请注意Docker中时区问题



```shell
docker pull codfrm/cloudcat:0.1.0
docker run -it -v $(PWD)/bilibili.zip:/cloudcat/bilibili.zip -v /etc/localtime:/etc/localtime -v /etc/timezone:/etc/timezone codfrm/cloudcat:0.1.0 exec bilibili.zip
```



