---
title: 如何使用vbot做一个自动微信客服系统（一）环境、踩坑及运行
tags:
  - php7
  - vbot
  - 微信客服
id: '287'
categories:
  - - 后端开发
  - - 微信相关
date: 2018-10-12 09:52:49
cover: ../../static/uploads/2018/10/robot-customer-service-1024x661.jpg
---

## 目标

监听符合规则的群聊，发现某些关键字的时候，可以自动回复消息。

### 举个栗子

群聊名称规则为： WJ\_\*\*\*，监听到关键字 “博客” 时，自动回复我的博客地址，监听到关键字 “斗图” 时，随机回复一张图片，监听到关键字 “装X指南”，自动发送一篇pdf文档等。 如果配置了swoole，还可以让机器人每天早上给 客(妹)户(子) 发早安，每天晚上发晚安，天气热了降温，天气凉了添衣，特殊日子安抚安抚！o(\*￣▽￣\*)o

# Vbot是神魔

Vbot是一个基于PHP7的微信机器人，支持被动回复消息或者主动给用户发送消息。发消息的原理是使用微信网页端API，所以微信网页端能做什么，vbot就能做什么。 官网地址：[http://create.hanc.cc/](http://create.hanc.cc/) 官方示例：[https://github.com/HanSon/my-vbot](https://github.com/HanSon/my-vbot)

## 安装Vbot

要想运行vbot，首先得有一个运行环境，测试环境可以是虚拟机，因为要一直跑，所以运行环境需要是一直运行的服务器才行。 **【注意】**本教程基于Linux，因为Swoole扩展和 php-curl 扩展的特殊性，博主并没有在 Windows 上倒腾。 博主的运行环境：Centos 7

### PHP的安装

PHP的环境需要在 7.0 以上，为了方便安装，所以采用了yum + webstatic源

#### 使用webstatic源

执行指令：

```shell
rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
rpm -Uvh https://mirror.webtatic.com/yum/el7/webtatic-release.rpm
```

等待安装完毕之后，执行如下指令： **【注】**由于要编译安装某些 php 插件，所以需要安装 php-devel 以支持 phpize

```shell
yum install php71w-cli php71w-devel
```

执行完毕之后，检查下版本

```shell
php -v
```

执行结果如下：

```shell
[root@localhost ~]# php -v
PHP 7.1.20 (cli) (built: Jul 20 2018 08:31:34) ( NTS )
Copyright (c) 1997-2018 The PHP Group
Zend Engine v3.1.0, Copyright (c) 1998-2018 Zend Technologies
```

## 使CURL支持OPENSSL

**【重点大坑】**这个步骤千万要做，不然Vbot就无法主动发送消息了，会一直循环报一个特别奇葩的错误。 可以参考的博客：[https://www.cnblogs.com/showker/p/4706271.html](https://www.cnblogs.com/showker/p/4706271.html) 首先，查看一下现在curl的状态：

```shell
[root@localhost ~]# curl -V
curl 7.29.0 (x86_64-redhat-linux-gnu) libcurl/7.29.0 NSS/3.28.4 zlib/1.2.7 libidn/1.28 libssh2/1.4.3
Protocols: dict file ftp ftps gopher http https imap imaps ldap ldaps pop3 pop3s rtsp scp sftp smtp smtps telnet tftp 
Features: AsynchDNS GSS-Negotiate IDN IPv6 Largefile NTLM NTLM_WB SSL libz unix-sockets
```

可以清楚的看到 **NSS/3.28.4**，并不是期望的 **OPENSSL**

##### 安装openssl和openssl-devel

```shell
yum install openssl openssl-devel
```

##### 尝试更新curl

```shell
[root@localhost swoole-src-4.0.4]# yum update openssl curl
```

检查 `curl -V`，如果已经是 **OPENSSL** 模式，则跳过下一步，否则继续

##### 重新编译curl做替换

首先从github下载curl的源码包并解压

```shell
[root@localhost swoole-src-4.0.4]# wget https://github.com/curl/curl/archive/curl-7_61_0.tar.gz
--2018-08-29 19:26:56--  https://github.com/curl/curl/archive/curl-7_61_0.tar.gz
[root@localhost swoole-src-4.0.4]# tar -zxvf curl-7_61_0.tar.gz
[root@localhost swoole-src-4.0.4]# cd curl-curl-7_61_0/
[root@localhost curl-curl-7_61_0]#
```

开始 `./buildconf` `./configure` `make && make install` 【注意】./buildconf 是为了确认编译环境是否有缺失，如果差某些包，yum安装一下即可。 比如：

```shell
[root@localhost curl-curl-7_61_0]# ./buildconf 
buildconf: autoconf version 2.69 (ok)
buildconf: autom4te version 2.69 (ok)
buildconf: autoheader version 2.69 (ok)
buildconf: automake version 1.13.4 (ok)
buildconf: aclocal version 1.13.4 (ok)
buildconf: libtoolize not found.
  You need GNU libtoolize 1.4.2 or newer installed.
```

执行以下指令安装 libtool：

```shell
[root@localhost ~]# yum install libtool
```

再执行 `./buildconf` 就会出现如下所示的结果

```
[root@localhost curl-curl-7_61_0]# ./buildconf 
buildconf: autoconf version 2.69 (ok)
buildconf: autom4te version 2.69 (ok)
buildconf: autoheader version 2.69 (ok)
buildconf: automake version 1.13.4 (ok)
buildconf: aclocal version 1.13.4 (ok)
buildconf: libtoolize version 2.4.2 (ok)
buildconf: GNU m4 version 1.4.16 (ok)
buildconf: running libtoolize
libtoolize: putting auxiliary files in `.'.
libtoolize: copying file `./ltmain.sh'
libtoolize: putting macros in AC_CONFIG_MACRO_DIR, `m4'.
libtoolize: copying file `m4/libtool.m4'
libtoolize: copying file `m4/ltoptions.m4'
libtoolize: copying file `m4/ltsugar.m4'
libtoolize: copying file `m4/ltversion.m4'
libtoolize: copying file `m4/lt~obsolete.m4'
libtoolize: Remember to add `LT_INIT' to configure.ac.
buildconf: converting all mv to mv -f in local m4/libtool.m4
buildconf: running aclocal
buildconf: converting all mv to mv -f in local aclocal.m4
buildconf: running autoheader
buildconf: running autoconf
buildconf: running automake
configure.ac:125: installing './compile'
configure.ac:187: installing './config.guess'
configure.ac:187: installing './config.sub'
configure.ac:125: installing './install-sh'
configure.ac:133: installing './missing'
docs/examples/Makefile.am: installing './depcomp'
parallel-tests: installing './test-driver'
buildconf: OK
```

##### 编译安装

```shell
[root@localhost curl-curl-7_61_0]# ./configure --without-nss --with-ssl
[root@localhost curl-curl-7_61_0]# make && make install
```

**【特别注意】这里的 --without-nss --with-ssl 是指使用 SSL ，而不使用 NSS** 编译结果里边注意查看下 SSL support：

```
config.status: executing libtool commands
configure: Configured to build curl/libcurl:

  curl version:     7.61.0-DEV
  Host setup:       x86_64-unknown-linux-gnu
  Install prefix:   /usr/local
  Compiler:         gcc
  SSL support:      enabled (OpenSSL)
  SSH support:      no      (--with-libssh2)
  zlib support:     enabled
  brotli support:   no      (--with-brotli)
  GSS-API support:  no      (--with-gssapi)
  TLS-SRP support:  no      (--enable-tls-srp)
  resolver:         POSIX threaded
  IPv6 support:     enabled
  Unix sockets support: enabled
  IDN support:      no      (--with-{libidn2,winidn})
  Build libcurl:    Shared=yes, Static=yes
  Built-in manual:  enabled
  --libcurl option: enabled (--disable-libcurl-option)
  Verbose errors:   enabled (--disable-verbose)
  SSPI support:     no      (--enable-sspi)
  ca cert bundle:   /etc/pki/tls/certs/ca-bundle.crt
  ca cert path:     no
  ca fallback:      no
  LDAP support:     no      (--enable-ldap / --with-ldap-lib / --with-lber-lib)
  LDAPS support:    no      (--enable-ldaps)
  RTSP support:     enabled
  RTMP support:     no      (--with-librtmp)
  metalink support: no      (--with-libmetalink)
  PSL support:      no      (libpsl not found)
  HTTP2 support:    disabled (--with-nghttp2)
  Protocols:        DICT FILE FTP FTPS GOPHER HTTP HTTPS IMAP IMAPS POP3 POP3S RTSP SMB SMBS SMTP SMTPS TELNET TFTP
```

**【注意】**这一步请在确定没有问题之后再往下进行

##### 替换原有curl程序

替换之后，即可发现 curl 已经变成了 OPENSSL 模式，而非 NSS

```shell
[root@localhost curl-curl-7_61_0]# mv /usr/bin/curl /usr/bin/curl.bak
[root@localhost curl-curl-7_61_0]# ln /usr/local/bin/curl /usr/bin/curl
[root@localhost curl-curl-7_61_0]# curl -V
curl 7.61.0-DEV (x86_64-unknown-linux-gnu) libcurl/7.61.0-DEV OpenSSL/1.0.2k zlib/1.2.7
Release-Date: [unreleased]
Protocols: dict file ftp ftps gopher http https imap imaps pop3 pop3s rtsp smb smbs smtp smtps telnet tftp 
Features: AsynchDNS IPv6 Largefile NTLM NTLM_WB SSL libz UnixSockets HTTPS-proxy
```

## 重新编译php-curl插件

在 curl 变成 OPENSSL 模式之后，php-curl依旧是 NSS 模式，可以通过如下指令查证

```shell
[root@localhost curl-curl-7_61_0]# php -r 'phpinfo();'  grep 'SSL Version'
SSL Version => NSS/3.34
```

为了解决这个问题，我们需要手动重新编译 `php-curl`，这个步骤得在 curl 更新完成之后

##### 下载php对应版本的源码包

7.1的源码包:[http://cn2.php.net/get/php-7.1.21.tar.gz/from/this/mirror](http://cn2.php.net/get/php-7.1.21.tar.gz/from/this/mirror) 下载解压并进入 `php-7.1.21/ext/curl/` 目录

```shell
[root@localhost ~]# wget http://cn2.php.net/get/php-7.1.21.tar.gz/from/this/mirror
[root@localhost ~]# mv mirror php7.1.tar.gz
[root@localhost ~]# tar -zxvf php7.1
[root@localhost ~]# cd php-7.1.21/ext/curl/
```

##### 使用phpize准备编译

确定文件 php-config 的位置，php-config 用于显示php的安装信息

```shell
[root@localhost ~]# find / -name php-config
/usr/bin/php-config
```

在 `/ext/curl/` 源目录中执行 `phpize --with-php-config /usr/bin/php-config`

```shell
[root@localhost curl]# phpize --with-php-config /usr/bin/php-config
Configuring for:
PHP Api Version:         20160303
Zend Module Api No:      20160303
Zend Extension Api No:   320160303
```

##### 执行安装

`./configure` `make && make install` 执行编译安装过程 **【注意】**需要先行安装 `gcc` `gcc-c++` `glibc-headers`，否则报错，安装指令：`yum install gcc gcc-c++ glibc-headers`

```
[root@localhost curl]# ./configure 
......
configure: creating ./config.status
config.status: creating config.h
config.status: executing libtool commands
[root@localhost curl]# make && make install
.....
Installing shared extensions:     /usr/lib64/php/modules/
```

出现最后一句话说明安装成功 最后再检查下 php-curl 的模式

```shell
[root@localhost curl]# php -r 'phpinfo();'  grep 'SSL Version'
SSL Version => OpenSSL/1.0.2k
```

## 安装Swoole插件

Swoole是一个强大的异步框架，可以用于搭建各种服务，譬如 HTTP服务，WebSocket等，这里不需要我们去学习如何使用Swoole，只需要安装上就可以了。 Swoole官网：[https://www.swoole.com/](https://www.swoole.com/) 安装教程：[https://wiki.swoole.com/wiki/page/6.html](https://wiki.swoole.com/wiki/page/6.html)

#### 下载并使用源码包

如果没有安装 `wget` 工具，请先 `yum install wget`

```shell
wget https://github.com/swoole/swoole-src/archive/v4.0.4.tar.gz
tar -zxvf v4.0.4.tar.gz
cd swoole-src-4.0.4/
```

##### 使用phpize准备编译

在 Swoole 源目录中执行 `phpize --with-php-config /usr/bin/php-config`

```shell
[root@localhost swoole-src-4.0.4]# phpize --with-php-config /usr/bin/php-config
Configuring for:
PHP Api Version:         20160303
Zend Module Api No:      20160303
Zend Extension Api No:   320160303
```

##### 执行安装

`./configure` `make && make install` 执行编译安装过程

```shell
[root@localhost swoole-src-4.0.4]# ./configure 
checking for grep that handles long lines and -e... /usr/bin/grep
checking for egrep... /usr/bin/grep -E
checking for a sed that does not truncate output... /usr/bin/sed
checking for cc... cc
checking whether the C compiler works... yes
checking for C compiler default output file name... a.out
[root@localhost swoole-src-4.0.4]# make && make install
```

安装完毕之后大概是这样：

```
/bin/sh /root/swoole-src-4.0.4/libtool --mode=install cp ./swoole.la /root/swoole-src-4.0.4/modules
libtool: install: cp ./.libs/swoole.so /root/swoole-src-4.0.4/modules/swoole.so
libtool: install: cp ./.libs/swoole.lai /root/swoole-src-4.0.4/modules/swoole.la
libtool: finish: PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin:/sbin" ldconfig -n /root/swoole-src-4.0.4/modules
----------------------------------------------------------------------
Libraries have been installed in:
   /root/swoole-src-4.0.4/modules

If you ever happen to want to link against installed libraries
in a given directory, LIBDIR, you must either use libtool, and
specify the full pathname of the library, or use the `-LLIBDIR'
flag during linking and do at least one of the following:
   - add LIBDIR to the `LD_LIBRARY_PATH' environment variable
     during execution
   - add LIBDIR to the `LD_RUN_PATH' environment variable
     during linking
   - use the `-Wl,-rpath -Wl,LIBDIR' linker flag
   - have your system administrator add LIBDIR to `/etc/ld.so.conf'

See any operating system documentation about shared libraries for
more information, such as the ld(1) and ld.so(8) manual pages.
----------------------------------------------------------------------

Build complete.
Don't forget to run 'make test'.

Installing shared extensions:     /usr/lib64/php/modules/
Installing header files:          /usr/include/php/
```

##### 检查下Swoole是否存在

```shell
[root@localhost swoole-src-4.0.4]# ll /usr/lib64/php/modules/swoole.so 
-rwxr-xr-x. 1 root root 9228512 8月  29 19:10 /usr/lib64/php/modules/swoole.so
```

编辑 `/etc/php.ini` 文件，添加扩展

```shell
[root@localhost swoole-src-4.0.4]# vi /etc/php.ini
```

找到这一段，加上 `extension=swoole.so`

```
;;;;;;;;;;;;;;;;;;;;;;
; Dynamic Extensions ;
;;;;;;;;;;;;;;;;;;;;;;

; If you wish to have an extension loaded automatically, use the following
; syntax:
;
;   extension=modulename.extension
;
; For example, on Windows:
;
;   extension=msql.dll
;
; ... or under UNIX:
;
;
; ... or with a path:
;
;   extension=/path/to/extension/msql.so
;
; If you only provide the name of the extension, PHP will look for it in its
; default extension directory.
extension=swoole.so
```

使用 `php -m grep swoole` 检查扩展是否安装成功

```shell
[root@localhost swoole-src-4.0.4]# php -m  grep swoole
swoole
```

## 安装Redis

yum仓库里边就有，直接 yum 安装下就可以了

```
[root@localhost ~]# yum install redis
```

为了方便使用，编辑 `/etc/redis-config` 修改 `daemonize no` 为 `daemonize yes` 以守护进程的方式启动

```
################################# GENERAL #####################################

# By default Redis does not run as a daemon. Use 'yes' if you need it.
# Note that Redis will write a pid file in /var/run/redis.pid when daemonized.
daemonize yes
```

启动 redis-server 并连接测试

```
[root@localhost ~]# redis-server /etc/redis.conf 
[root@localhost ~]# redis-cli 
127.0.0.1:6379> set name hello
OK
127.0.0.1:6379> get name
"hello"
127.0.0.1:6379> exit
[root@localhost ~]# 
```

## 环境搭建完毕，开始构建代码

##### 使用 composer 安装vbot

composer 的安装不再赘述，参考官方文档：[https://pkg.phpcomposer.com/#how-to-install-composer](https://pkg.phpcomposer.com/#how-to-install-composer)

```
[root@localhost curl]# php -r "copy('https://install.phpcomposer.com/installer', 'composer-setup.php');"
[root@localhost curl]# php composer-setup.php
All settings correct for using Composer
Downloading...

Composer (version 1.6.5) successfully installed to: /root/php-7.1.21/ext/curl/composer.phar
Use it: php composer.phar

[root@localhost curl]# php -r "unlink('composer-setup.php');"
[root@localhost curl]# mv composer.phar /usr/bin/composer
[root@localhost curl]# composer -v
   ______
  / ____/___  ____ ___  ____  ____  ________  _____
 / /   / __ \/ __ `__ \/ __ \/ __ \/ ___/ _ \/ ___/
/ /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /
\____/\____/_/ /_/ /_/ .___/\____/____/\___/_/
                    /_/
Composer version 1.6.5 2018-05-04 11:44:59
```

推荐使用国内镜像，会快一些

```
[root@localhost my-bot]# composer config -g repo.packagist composer https://packagist.phpcomposer.com
[root@localhost my-bot]#
```

##### 安装vbot

vbot是一个脚本程序，不需要跑在 LNMP 环境下，只需要 PHP + redis 就可以了。 新建一个目录 `my-bot`，然后 `composer require hanson/vbot`

```
[root@localhost ~]# mkdir my-bot
[root@localhost ~]# cd my-bot/
[root@localhost my-bot]# composer require hanson/vbot
```

可能会出现如下报错

```
- illuminate/support v5.4.17 requires ext-mbstring * -> the requested PHP extension mbstring is missing from your system.
    - illuminate/support v5.4.13 requires ext-mbstring * -> the requested PHP extension mbstring is missing from your system.
    - illuminate/support v5.4.0 requires ext-mbstring * -> the requested PHP extension mbstring is missing from your system.
    - Installation request for hanson/vbot ^2.0 -> satisfiable by hanson/vbot[2.0.1, 2.0.10, 2.0.2, 2.0.3, 2.0.4, 2.0.5, 2.0.6, 2.0.7, 2.0.8, 2.0.9, v2.0].

  To enable extensions, verify that they are enabled in your .ini files:
    - /etc/php.ini
    - /etc/php.d/bz2.ini
    - /etc/php.d/calendar.ini
    - /etc/php.d/ctype.ini
    - /etc/php.d/curl.ini
    - /etc/php.d/exif.ini
    - /etc/php.d/fileinfo.ini
    - /etc/php.d/ftp.ini
    - /etc/php.d/gettext.ini
    - /etc/php.d/gmp.ini
    - /etc/php.d/iconv.ini
    - /etc/php.d/json.ini
    - /etc/php.d/phar.ini
    - /etc/php.d/shmop.ini
    - /etc/php.d/simplexml.ini
    - /etc/php.d/sockets.ini
    - /etc/php.d/tokenizer.ini
    - /etc/php.d/xml.ini
    - /etc/php.d/zip.ini
  You can also run `php --ini` inside terminal to see which files are used by PHP in CLI mode.

Installation failed, reverting ./composer.json to its original content.
```

根据错误提示得知，我们需要安装 `php-mbstring` 扩展，yum 安装下即可

```
[root@localhost my-bot]# yum install php71w-mbstring
```

安装完成之后再进行 `composer require hanson/vbot` 就可以了