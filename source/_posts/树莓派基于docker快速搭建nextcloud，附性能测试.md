---
title: 树莓派基于Docker快速搭建nextcloud，附性能测试
tags:
  - nextcloud
  - syncting
  - 同步
id: '540'
categories:
  - - Linux
date: 2020-07-02 09:39:37
cover: /static/uploads/2020/07/20150428-cloud-computing.0.1489222360-1200x661.jpg
---



## 前言

一年前兴起买的树莓派3B+，一番折腾最后成为了家里的电视投屏盒子（OSMC系统），我的自建云盘放在云服务器上，50G的空间快被塞满了，于是就想起来把尘封数月的树莓派复活过来，搭建一个 nextcloud，并且把之前的数据迁移过去。

## 正文

### 树莓派刷系统

如果与我一样有把树莓派作为电视盒子的需求的话，这里怒推一波 OSMC，使用了一年多稳定性还挺好的 链接在此：https://osmc.tv/download/ 或者使用官方的 Raspbian ，属于比较纯粹的环境 链接在此：https://www.raspberrypi.org/downloads/ 下载好的系统使用 win32diskimager 刷入内存卡即可，十分简单，这里贴一个安装链接供参考 https://www.jianshu.com/p/a337ccae5d2b

### 环境准备

首先，ssh进入树莓派，如果是 OSMC 系统的话，用户名和密码都是osmc [![](/static/uploads/2020/03/d2db55614f17f5a966cbdf47c7f359fd.png)](/static/uploads/2020/03/d2db55614f17f5a966cbdf47c7f359fd.png) 连接结果如下：

```shell
Connecting to 192.168.3.7:22...
Connection established.
To escape to local shell, press 'Ctrl+Alt+]'.

WARNING! The remote SSH server rejected X11 forwarding request.
Linux jeffrey 4.19.55-6-osmc #1 SMP PREEMPT Sun Nov 3 22:15:28 UTC 2019 armv7l

The programs included with the Debian GNU/Linux system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Debian GNU/Linux comes with ABSOLUTELY NO WARRANTY, to the extent
permitted by applicable law.
Last login: Mon Mar  2 11:00:03 2020 from 192.168.3.10
osmc@jeffrey:~$
```

#### 调整时区

```shell
root@osmc:/home/osmc/rasp-tools/frp_client# timedatectl set-timezone Asia/Shanghai
root@osmc:/home/osmc/rasp-tools/frp_client# date -R
Fri, 03 Apr 2020 12:38:43 +0800
root@osmc:/home/osmc/rasp-tools/frp_client# date
Fri Apr  3 12:38:47 CST 2020
```

#### Docker安装

这里参考了这篇博客：https://docker\_practice.gitee.io/install/raspberry-pi.html 整理下来，直接执行如下两行命令即可

```shell
$ curl -fsSL get.docker.com -o get-docker.sh
$ sudo sh get-docker.sh --mirror Aliyun
```

> 如果是用的 osmc 系统，执行脚本会报如下错误：

```shell
osmc@jeffrey:~/rasp-tools$ sudo sh get-docker.sh --mirror Aliyun
# Executing docker install script, commit: f45d7c11389849ff46a6b4d94e0dd1ffebca32c1

ERROR: Unsupported distribution 'osmc'
```

需要修改下 `get-docker.sh` 中第280行改为，因为 osmc 没有在里边有脚本配置，假装是 `debian` 就行了

> lsb\_dist="raspbian"

安装完毕后输出大概是这样

> Server: Docker Engine - Community Engine: Version: 19.03.6 API version: 1.40 (minimum version 1.12) Go version: go1.12.16 Git commit: 369ce74 Built: Thu Feb 13 01:31:39 2020 OS/Arch: linux/arm Experimental: false containerd: Version: 1.2.10 GitCommit: b34a5c8af56e510852c35414db4c1f4fa6172339 runc: Version: 1.0.0-rc8+dev GitCommit: 3e425f80a8c931f88e6d94a8c831b9d5aa481657 docker-init: Version: 0.18.0 GitCommit: fec3683 If you would like to use Docker as a non-root user, you should now consider adding your user to the "docker" group with something like: sudo usermod -aG docker your-user Remember that you will have to log out and back in for this to take effect! WARNING: Adding a user to the "docker" group will grant the ability to run containers which can be used to obtain root privileges on the docker host. Refer to https://docs.docker.com/engine/security/security/#docker-daemon-attack-surface for more information.

根据提示，可以把当前用户加入到 docker 组，以便这个用户能方便的使用 docker 命令（需要重启）

```shell
sudo usermod -aG docker osmc
```

出于效率考虑，可以创建 `/etc/docker/daemon.json`，并写入如下内容：

```json
{
    "registry-mirrors": ["http://hub-mirror.c.163.com"]
}
```

使用 `systemctl start docker` 启动服务即可

#### docker-compose的安装

docker-compose 是基于 python 的一个 docker 编排工具，能方便的聚合现有docker镜像，使用起来十分方便 可以使用 `apt-get install docker-compose` 的形式安装

#### nextcloud

nextcloud 是一个开源网盘项目，基于 php7，功能非常的强悍，包括手机端、PC端、服务端的程序，在树莓派上只需要安装服务端就可以了，手机端和PC端都只是客户端，找到相关项目下载即可 服务端项目地址：[https://github.com/nextcloud/server](https://github.com/nextcloud/server) 桌面端项目地址：[https://github.com/nextcloud/desktop](https://github.com/nextcloud/desktop) 安卓项目地址：[https://github.com/nextcloud/android](https://github.com/nextcloud/android) IOS项目地址：[https://github.com/nextcloud/ios](https://github.com/nextcloud/ios)

### 运行nextcloud

基于官方关于树莓派运行nextcloud的文档，找到了其中的 docker-compose.yaml，并针对树莓派的特性稍微做了些修改

```yaml
version: '2'

volumes:
  nextcloud:
  db:

services:
  db:
    image: jsurf/rpi-mariadb
    command: --transaction-isolation=READ-COMMITTED --binlog-format=ROW
    restart: always
    volumes:
      - ./db:/var/lib/mysql
    environment:
      - MYSQL_ROOT_PASSWORD=nyist123
      - MYSQL_PASSWORD=nyist123
      - MYSQL_DATABASE=nextcloud
      - MYSQL_USER=nextcloud
  app:
    image: arm32v7/nextcloud:fpm
    volumes:
      - ./data:/var/www/html/data
    restart: always

  web:
    image: arm32v7/nginx
    ports:
      - 8089:80
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    volumes_from:
      - app
    restart: always
  redis:
    image: redis
    volumes_from:
      - app
    restart: always
networks:
  nextcloud:
```

在 `docker-compose.yaml` 目录下使用 `docker-compose up -d` 开始拉取镜像然后后台运行

```shell
root@osmc:/home/osmc/rasp-tools/nextcloud# docker-compose up -d
WARNING: Some networks were defined but are not used by any service: nextcloud
nextcloud_db_1 is up-to-date
nextcloud_app_1 is up-to-date
nextcloud_redis_1 is up-to-date
nextcloud_web_1 is up-to-date
```

进入 `http://ip地址:8089` 即可开始配置数据库等信息，如果出现 504 time out，需要调整Nginx中的超时时间 在刚才的 `docker-compose.yaml` 中指定了数据库主机是 `db`，用户名是 `nextcloud`，数据库是 `nextcloud`，密码是 `nyist123` [![](/static/uploads/2020/04/1c40b3052fc38dda5980a0888c94a133.png)](/static/uploads/2020/04/1c40b3052fc38dda5980a0888c94a133.png) [![](/static/uploads/2020/04/09e37a9efc26ef715406e710b501814e.png)](/static/uploads/2020/04/09e37a9efc26ef715406e710b501814e.png)

### 内网穿透

内网穿透使用的 `frp`，服务端和客户端的配置可以参考下方的这篇博客，也可以用免费的 frp 服务，就是速度和稳定性就没法保证了，或者使用 `ngrok`，本人尝试下来 `ngrok` 有不稳定的情况，偶尔会断掉 [https://www.jianshu.com/p/00c79df1aaf0]("https://www.jianshu.com/p/00c79df1aaf0") 出于效率和体验的考虑，如果是自建 frp 服务器的话，可以用 `nginx` 反向代理到 `frp` 的 `http` 代理端口上，`https` 也在服务端的 `nginx` 上做

#### 配置 frp

个人的 `frpc.ini` 如下：

```ini
[common]
server_addr = 服务器ip
server_port = 服务器服务端口
privilege_token = 密码
[ssh]
type = tcp
local_ip = 127.0.0.1
local_port = 22
remote_port = 32656 # 远程服务器的此端口映射到 22 端口
[nextcloud]
type = http
local_port = 8089 # 转发到树莓派的 8089
custom_domains = 自定义网盘域名
```

服务端的 `frps.ini` 如下：

```ini
[common]
# binde_addr是指定frp内网穿透服务器端监听的IP地址,默认为127.0.0.1，
#如果使用IPv6地址的话，必须用方括号包括起来，比如 “[::1]:80”, “[ipv6-host]:http” or “[ipv6-host%zone]:80”
bind_addr = 0.0.0.0
# bind_port 是frp内网穿透服务器端监听的端口，默认是7000
bind_port = 7000
#auto_token = 验证token
privilege_token = 密码
#token = 密码
#frp内网穿透服务器可以支持虚拟主机的http和https协议，一般我们都用80，可以直接用域名而不用增加端口号，如果使用其它端口，那么客户端也需要配置相同的其他端口。
vhost_http_port = 8189 # http 的统一转发端口
vhost_https_port = 7179 # https 的统一转发端口

dashboard_user = 管理后台用户 
dashboard_pwd = 管理后台密码
# 这个是frp内网穿透服务器的web界面的端口，可以通过http://你的ip:7500查看frp内网穿透服务器端的连接情况，和各个frp内网穿透客户端的连接情况。
dashboard_port = 端口
auth_token = 密码
```

#### 自动重新运行

还写了一个每分钟检测 frpc 是否正常运行的小脚本 `start.sh`

```shell
pid=`ps -fe  grep frpc  grep -v grep  awk '{print $2}'`
if [ -z "$pid" ]; then
    path=/home/osmc/rasp_tools/frp_client
    nohup $path/frpc -c $path/frpc.ini > $path/nohup.out &
    echo "自动重启"
else
    echo "已经启动PID:$pid"
fi
```

然后安装 `apt-get install cron`，每分钟运行该脚本 在`crontab -e`的最后加了一行，博主用的 `root` 加的这一任务

```bash
*  *    * * *   root    bash /home/osmc/rasp-tools/frp_client/start.sh
```

填加完之后做一个小测试，手动运行 `starts.sh`，发现没问题后，停掉 `frpc`，等待时间到达整分钟之后，观察 `fprc` 有没有自动重启

```bash
root@osmc:/home/osmc/rasp-tools/frp_client# bash start.sh 
自动重启
root@osmc:/home/osmc/rasp-tools/frp_client# nohup: redirecting stderr to stdout
root@osmc:/home/osmc/rasp-tools/frp_client# ps -fe  grep frpc
root      1864     1  0 04:27 pts/0    00:00:00 /home/osmc/rasp-tools/frp_client/frpc -c /home/osmc/rasp-tools/frp_client/frpc.ini
root      1875  1343  0 04:27 pts/0    00:00:00 grep frpc
root@osmc:/home/osmc/rasp-tools/frp_client# kill -9 1864
root@osmc:/home/osmc/rasp-tools/frp_client# ps -fe  grep frpc
root      1883  1343  0 04:27 pts/0    00:00:00 grep frpc
# .... 等待一会儿
root@osmc:/home/osmc/rasp-tools/frp_client# ps -fe  grep frpc
root      1984     1  0 04:29 ?        00:00:00 /home/osmc/rasp-tools/frp_client/frpc -c /home/osmc/rasp-tools/frp_client/frpc.ini
root      1985     1  0 04:29 ?        00:00:00 /home/osmc/rasp-tools/frp_client/frpc -c /home/osmc/rasp-tools/frp_client/frpc.ini
root      2010  1343  0 04:29 pts/0    00:00:00 grep frpc
```

#### 配置反向代理

这个时候就可以根据反向代理的域名，比如我的样例里边用的 [http://nextcloud.wj2015.com:8189/]("http://nextcloud.wj2015.com:8189/") 访问项目，但是这样并不友好，包括

*   没有HTTPS
*   端口是8189，不好看

为解决上述两个问题，我在 fprs服务器配置了nginx反向代理，将域名 `nextcloud.wj2015.com` 配好https，再代理到 `raspi.nextcloud.wj2015.com:8189`，具体的 vhost 配置文件如下

```ini
upstream nextcloud-dashboard{
    server 127.0.0.1:8189;
}

server {
    listen 80;
    #listen [::]:80;
    server_name nextcloud.wj2015.com;
    # enforce https
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    #listen [::]:443 ssl http2;
    server_name nextcloud.wj2015.com;

    ssl_certificate ssl/2018-11-09-SSL/1_nextcloud.wj2015.com_bundle.crt;
    ssl_certificate_key ssl/2018-11-09-SSL/2_nextcloud.wj2015.com.key;

    # Add headers to serve security related headers
    # Before enabling Strict-Transport-Security headers please read into this
    # topic first.
    add_header Strict-Transport-Security "max-age=15768000; includeSubDomains; preload;";
    large_client_header_buffers 4 16k;

    location / {
        proxy_set_header   Host         $host:$server_port;
        proxy_set_header   X-Real-IP        $remote_addr;
        proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;

        proxy_pass  http://nextcloud-dashboard;
    }
}

```

访问 `https://nextcloud.wj2015.com` 可能会看到下面的提示 [![](/static/uploads/2020/04/2f500313ada9cdc421f34109902ee7c3.png)](/static/uploads/2020/04/2f500313ada9cdc421f34109902ee7c3.png) 根据提示改一下可信域名即可，或者一开始安装的时候就在此域名下安装（更推荐） 改可信域名需要使用 `docker exec -it nextcloud_app_1 /bin/bash` 进入到 `docker` 容器中更改对应文件，改动之后重新启动容器就需要慎重了；或者映射对应配置文件到真实路径下 最后能看到这个页面就算成功了 [![](/static/uploads/2020/04/b64d24d3f221b7a8528de78af1ce4f22.png)](/static/uploads/2020/04/b64d24d3f221b7a8528de78af1ce4f22.png)

### 速度测试

由于树莓派3b孱弱的性能，OSMC系统持久输出HDMI信号以及 nectcloud 带来的 php-fpm + nginx 结构上的损失，加之内网穿透后速度很慢的网络，又卡容易丢包，不推荐使用。 除非使用更好的树莓派，或者单独只装一个 nextcloud，就像 nextcloud 官方给到的内置 nextcloud 的系统。

### 其他推荐

nextcloud的功能和插件非常之丰富，但是如果只有同步文件的需求的话，可以使用去中心化的 syncting 方便的在不同电脑上同步文件，炒鸡方便，如果有在线下载的需求，可以内网穿透 + nginx 配合HTTP鉴权简单跑一个文件管理页面。 并且 syncting 可以自己搭建中继服务器，速度+++ 参考的其他博主链接：[https://zhuanlan.zhihu.com/p/69267020?from\_voters\_page=true](https://zhuanlan.zhihu.com/p/69267020?from_voters_page=true) 官方网站：[https://syncthing.net/](https://syncthing.net/)

## 总结

树莓派在运行有其他耗性能的服务时，就不适合再用 nextcloud 做网盘了，如果像我一样有运行其他服务诸如 KODI 的同时，还希望能在不同电脑上同步文件，那么就可以考虑同步神器 syncting 了。