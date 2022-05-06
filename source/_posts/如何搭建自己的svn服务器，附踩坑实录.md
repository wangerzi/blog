---
title: 如何搭建自己的SVN服务器，附踩坑实录
tags:
  - SVN
id: '255'
categories:
  - - Linux
date: 2018-07-20 10:23:53
cover: /static/uploads/git.jpg
---

# 目标

1.  在Linux服务器上搭建SVN，实现服务器本地使用SVN自动同步最新代码
2.  新建N个账号，给产品、程序员在公网环境下使用



# 准备工作

首先，，你得有一台服务器，如果这个SVN仓库需要能在公网访问，服务器得需要有一个公网IP。 如果没有公网环境，虚拟机来做测试也是可以的，甚至家里边用不上的电脑加花生壳软件，再做一个映射都行。 ![](/static/uploads/2018/07/timg-5-300x231.jpg)

# SVN简介

SVN是subversion的简称，是一个操作简单的中心化的版本控制器，Windows上使用TortoiseSVN软件或者Linux上使用subversion工具，就可以方便的在各个客户端同步代码、资料等。[TortoiseSVN下载](https://tortoisesvn.net/downloads.html "TortoiseSVN下载") ![](/static/uploads/2018/07/9e18b6f2f0fbb72f8c8c374218f3436c.gif) ![](/static/uploads/2018/07/790909ad2936bd191a9ead8940c0d568.png)

# 操作步骤

## 安装subversion

先确定服务器是否有安装subversion

```
[root@localhost ~]# svnversion --version
```

如下图所示，是未安装的情况 ![](/static/uploads/2018/07/278260c66dc742a58eefce599cf4d745.png) 如果未安装，使用yum install安装一下即可

```
[root@localhost ~]# yum install subversion
```

![](/static/uploads/2018/07/0bab54e406d396fb4bba6e41f8e6eb8a.png) 安装成功后，subversion -v确定下，如下图所示就是安装成功

```
[root@localhost ~]# svnversion --version
```

![](/static/uploads/2018/07/cc9a499649ac9674e9fe1d5572df2e6b.png)

## 新建SVN仓库

创建SVN仓库的命令是svnadmin create 仓库路径 比如，我想让我的`/var/svn/`作为我SVN仓库的根目录，那么可以分为三个步骤进行仓库创建操作

```
创建文件夹 => 进入文件夹 => 创建仓库
```

示例如下：

```
[root@localhost ~]# mkdir /var/svn
[root@localhost ~]# cd /var/svn/
[root@localhost svn]# svnadmin create test
```

执行结果如下，在`/var/svn/test`下创建了一个仓库 ![](/static/uploads/2018/07/c848cc30a22f99fc77c591341cfca9ad.png)

## 配置SVN

创建完毕后，`/var/svn/test`目录下出现conf、db、hooks等目录 其中，conf文件夹下包含三个重要的文件`conf/passwd`、`conf/authz`、`conf/svnserve`

文件名

描述

conf/passwd

用户文件，配置用户名和密码

conf/authz

用户权限文件，用于配置用户权限

conf/svnserve.conf

配置启动参数

### 用户文件passwd

首先，编辑passwd文件，里边只有一个\[users\]模块，含义是 用户名 = 密码 比如，我们有Jeffrey、Sally、Jack、Pony四个人，密码是名字+123，就可以这样写

```
[root@localhost conf]# vim passwd
```

![](/static/uploads/2018/07/e47e510b09981bc35dbe67eb8f43de4d.png)

### 权限文件authz

权限文件用于指定访问权限，可指定分组权限或个人权限 比如： Jeffrey和Sally是程序员，属于programers分组，只读`/doc/`目录，读写`/src/`目录。 Jack是产品经理，属于PM分组，能读写`/doc/`，不能访问`/src/`。 Pony是老板，属于Mangager分组，能读写`/doc`、`/src/`。 **样例如下**： \[groups\]模块里边配置分组信息 \[/dir\]模块里边配置分组权限/个人权限信息，r代表读，w代表写

```
配置说明：
分组名 = rw   # 分组权限配置
&用户名 = rw  # 单个用户权限配置
* = rw        # 所有用户权限配置
```

![](/static/uploads/2018/07/97bba3045f1c898642cd6ae5411199d6.png)

### 启动文件svnserve.conf

用于配置启动参数，具体说明如下图所示 ![](/static/uploads/2018/07/5c5f5d4251b858965373d5e80d834773.png)

## 运行SVNSERVE

配置完毕之后就可以启动SVN服务了，svnserve是一个服务端程序，启动后就能向其他人提供SVN服务。 **参数说明：** svnserve -d -r /路径 --listen-port 端口号 -d 表示守护进程 -r 表示SVN服务根路径 --listen-port 可指定端口号

```
# 运行服务端程序
[root@localhost conf]# svnserve -d -r /var/svn/
# 检查是否运行成功
[root@localhost conf]# ps -fe  grep svnserve
root      1586     1  0 02:25 ?        00:00:00 svnserve -d -r /var/svn/
root      1588  1448  0 02:25 pts/0    00:00:00 grep svnserve
```

### 本地仓库测试

假设现在这个服务器上搭建有一个网站，路径是`/var/www/html/`，需要使用svn更新代码。 **【注】**首次使用，需要输入root密码和配置的用户名+密码，root密码直接输入空就行了。

```
[root@localhost html]# svn checkout svn://127.0.0.1/test/
认证领域: <svn://127.0.0.1:3690> My First Repository
“root”的密码: 
认证领域: <svn://127.0.0.1:3690> My First Repository
用户名: Jeffrey
“Jeffrey”的密码: 

-----------------------------------------------------------------------
注意!  你的密码，对于认证域:

   <svn://127.0.0.1:3690> My First Repository

只能明文保存在磁盘上!  如果可能的话，请考虑配置你的系统，让 Subversion
可以保存加密后的密码。请参阅文档以获得详细信息。

你可以通过在“/root/.subversion/servers”中设置选项“store-plaintext-passwords”为“yes”或“no”，
来避免再次出现此警告。
-----------------------------------------------------------------------
保存未加密的密码(yes/no)?yes
[root@localhost html]# svn checkout svn://127.0.0.1/test/
取出版本 0。
```

#### 添加文件

【注】先切换到`/var/www/html/test/`，再进行SVN操作，还有**使用的用户是Jeffrey用户，只有src目录的写权限，所以只能提交src目录，否则报权限错误。**。 【注】测试是，我给`src/index.html`里边写入了"Hello World!"

```
[root@localhost html]# cd test/
[root@localhost test]# mkdir doc/ src
[root@localhost test]# vim src/index.html
[root@localhost test]# svn status
?       doc
?       src
```

#### 提交

commit的时候，-m 后写上理由即可。

```
[root@localhost test]# svn add src/
A         src
A         src/index.html
[root@localhost test]# svn commit -m '初始化提交'
增加           src
增加           src/index.html
传输文件数据.
提交后的版本为 1。
```

#### 实现脚本自动同步最新代码

Linux中可以使用shell+contab定时拉取SVN中的最新代码，也不麻烦。 ![](/static/uploads/2018/07/timg-6-300x219.jpg) 如果没有安装crontab，需要`yum install crontabs`，再`service crond start` `~/svnupdate.sh`内容：

```
#! /bin/sh
# 进入目标目录
cd /var/www/html/test
# 调用svn update即可同步
svn update
```

`crontab -e`内容：

```
*/60 * * * * /bin/shell /root/autogitpull.sh
```

命令概览

```
[root@localhost test]# vim ~/svnupdate.sh
[root@localhost test]# chmod u+x ~/svnupdate.sh
# 测试脚本是否正确
[root@localhost test]# ~/svnupdate.sh 
版本 1。
```

### 远程仓库测试

首先，SVN服务器的IP为`192.168.220.128`，所以Linux或Windows需要在指定地址是`svn://192.168.220.128/test`（创建的test项目） **【注】**如果出现如下错误，请按照采坑实录调整双方防火墙

```
[root@localhost html]# svn checkout svn://192.168.220.128/test/
svn: 无法连接主机“192.168.220.128”: 没有到主机的路由
```

#### Linux同步SVN测试

首先，按照相同的方法安装subversion

```
[root@localhost ~]# /usr/bin/yum install subversion
已加载插件：fastestmirror
设置安装进程
Loading mirror speeds from cached hostfile
```

特别需要注意，如果Linux本机没有开启8690的INPUT和OUTPUT，也会出现问题。 不用开启subserve，直接`svn checkout svn://192.168.220.128/test/src/`即可。

```
[root@localhost html]# svn checkout svn://192.168.220.128/test/src
认证领域: <svn://192.168.220.128:3690> My First Repository
“root”的密码: 
认证领域: <svn://192.168.220.128:3690> My First Repository
用户名: Sally
“Sally”的密码: 

-----------------------------------------------------------------------
注意!  你的密码，对于认证域:

   <svn://192.168.220.128:3690> My First Repository

只能明文保存在磁盘上!  如果可能的话，请考虑配置你的系统，让 Subversion
可以保存加密后的密码。请参阅文档以获得详细信息。

你可以通过在“/root/.subversion/servers”中设置选项“store-plaintext-passwords”为“yes”或“no”，
来避免再次出现此警告。
-----------------------------------------------------------------------
保存未加密的密码(yes/no)?yes
A    src/index.html
取出版本 1。
[root@localhost html]# tail src/index.html 
Hello World!
```

#### Windows同步SVN测试

安装完TortoiseSVN之后，在一个空文件夹里 右键 》 SVN Checkout ![](/static/uploads/2018/07/db0c19e75a1a6ce3e9af4255ab2c2bce.png) 填写SVN地址，点OK再输入密码，即可同步 ![](/static/uploads/2018/07/55053a8c06ae360634e4f70d61d5841a.png) ![](/static/uploads/2018/07/9e6460e9196b5bb28d71f29a2cbfa27e.png)

## 踩坑实录

各位请留步，博主还针对容易出现各种"Connect timed out"、"Connect refused"等问题，进行了一下总结。 ![](/static/uploads/2018/07/timg-2-300x187.jpg)

### 认证失败

如果出现“认证失败”或“Authorization Faild!”，首先考虑`conf/authz`的权限配置问题，特别是\[/\]分组的权限配置问题，先改为\* = rw试一试。

```
[/]
* = rw
```

如果是输入密码错误导致的"Authorization Faild!"可以删除`~/.subversion/auth/`，删除后即可重新输入密码。

```
[root@localhost test]# rm -rf ~/.subversion/auth
```

### iptables配置

svnserve的默认端口是3690，所以需要开启INPUT链 dport 3690和OUTPUT链 sport 3690。 `vim /etc/sysconfig/iptables` 在末尾加入下边三行代码

```
# SVN
-A INPUT -m state --state NEW -m tcp -p tcp --dport 3690 -j ACCEPT
-A OUTPUT -p tcp -m tcp --sport 3690 -m state --state ESTABLISHED -j ACCEPT
```

然后重启防火墙，Centos6：`service iptables restart`，Centos7，如果使用fireware则另行配置。 如果SVN客户端也是Linux的话，需要给客户端Linux加上防火墙配置

```
# SVN
-A INPUT -m state --state NEW -m tcp -p tcp --sport 3690 -j ACCEPT
-A OUTPUT -p tcp -m tcp --dport 3690 -m state --state ESTABLISHED -j ACCEPT
```

### 深坑-云主机的安全组

**深坑警告！！！** ![](/static/uploads/2018/07/timg-1-1-300x258.jpg) 腾讯云、阿里云均设有 安全组功能，类似于在服务器最外层加上了一个iptables，很多时候也能起到一定的防护作用。但是当我们启动了一个服务，但是没有在安全组内放开端口权限的话，，肯定报超时！！！

#### 创建安全组

![](/static/uploads/2018/07/e8c51965110239a6757f2859bd1c01dc.png) 选择**自定义** ![](/static/uploads/2018/07/7aeb6d29649bb8ce25193c7a6be2a10b.png)

#### 配置8690开放

入站规则 》 新建规则中，允许所有的8609访问 ![](/static/uploads/2018/07/07f78c01b1b7264547f26f9ce083b779.png) 出站规则 》 新建规则中，允许所有流量流出即可 ![](/static/uploads/2018/07/103ef20c68667829ebf2d29f1c9ca981.png) 顺手来个表情包 ![](/static/uploads/2018/07/timg-300x138.jpg)

## 总结

配置不坑，防火墙和安全组挺坑的