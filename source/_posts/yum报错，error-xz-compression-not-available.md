---
title: 'YUM报错，Error: xz compression not available'
tags:
  - Linux
  - yum
id: '305'
categories:
  - - Linux
date: 2018-11-09 12:05:09
---



## 遇到的问题

博主的服务器版本是 Centos6.8，今天想更新以下PHP版本，然后去安装最新的 `webstatic`，于是找到如下两条命令

```shell
rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm  
rpm -Uvh https://mirror.webtatic.com/yum/el7/webtatic-release.rpm
```

没有仔细看安装的RPM地址，然后就执行了，安装过程倒是没有报错，不过当我在执行 `yum search`、 `yum install` 的时候问题就出来了

```shell
[root@JeffreyWang nextcloud]# yum search php72w
Loaded plugins: fastestmirror, security
Repository epel is listed more than once in the configuration
Loading mirror speeds from cached hostfile
 * webtatic: sp.repo.webtatic.com
**Error: xz compression not available**
```

报错，**Error: xz compression not available**

## 修复问题

### 问题出现原因

YUM源版本不符

### 解决思路

*   将错误的源删掉
*   删除yum缓存和webstatic缓存，否则问题还会出现
*   安装正确的yum源

所以，执行如下命令：

```shell
# 删除错误源
[root@JeffreyWang nextcloud]# yum remove epel-release
Loaded plugins: fastestmirror, security
Setting up Remove Process
Resolving Dependencies
There are unfinished transactions remaining. You might consider running yum-complete-transaction first to finish them.
--> Running transaction check
---> Package epel-release.noarch 0:7-11 will be erased
--> Processing Dependency: epel-release >= 7 for package: webtatic-release-7-3.noarch
--> Running transaction check
---> Package webtatic-release.noarch 0:7-3 will be erased

--> Finished Dependency Resolution
# 删除缓存，注意：与 epel7 相关的缓存均需要删除，比如我这里安装的 webstatic，否则报错依旧
[root@JeffreyWang nextcloud]# rm -rf /var/cache/yum/x86_64/6/epel/
[root@JeffreyWang nextcloud]# rm -rf /var/cache/yum/x86_64/6/webtatic/*
# 安装对应系统的最新 epel
[root@JeffreyWang nextcloud]# rpm -Uvh http://mirror.webtatic.com/yum/el6/latest.rpm
```

## 总结

网上博客里边的指令不要随意跑！！！