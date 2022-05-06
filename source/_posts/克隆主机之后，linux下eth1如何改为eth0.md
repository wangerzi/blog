---
title: 克隆主机之后，Linux下eth1如何改为eth0
tags:
  - eth0
  - Vmware
  - 网卡
id: '93'
categories:
  - - Linux
date: 2017-12-24 14:55:52
cover: /static/uploads/linux.jpg
---



# 写在前面

在日常虚拟机的使用中，很可能会使用虚拟机软件的克隆功能，但是克隆之后的主机，网卡的MAC信息跟被克隆主机的mac信息是不一样的。 由于克隆主机MAC变化，在启动网卡(`ifup eth0`)时就会产生`Device eth0 does not seem to be present, delaying initialization.`错误，就像这样： ![](/static/uploads/2017/12/aff14d0b7df83bd7a47131344fc99ccb.png)，那么如何解决这样的问题就是这篇博客的主要内容。

# 产生问题原因

之所以虚拟机会让克隆主机的MAC改变，是因为在DHCP协议中或者说在内网的通信中，都是通过ARP协议将IP转换为MAC或将MAC转换为IP的，换句话说，局域网通信就是依靠MAC地址通信的，两台主机MAC地址相同则会导致 `/etc/sysconfig/network-script/ifcfg-eth0`中的网卡信息是被克隆主机的MAC，新主机的网卡已经不是那个MAC了。 查看`/etc/udev/rules.d/70-persistent-net.rules`可查看本机所有的mac信息，在里面可以看到，被克隆主机的MAC信息是eth0，新产生的MAC信息被识别为了eth1。 ![](/static/uploads/2017/12/62b7fc002ba5cb2ae723017d10ec1deb.png)

# 解决方案

## 快速解决方案

将`/etc/sysconfig/network-script/ifcfg-eth0`中的MAC换成新机的MAC，并将eth0换成eth1。 ![](/static/uploads/2017/12/f0b18c413e720760c59df6e63463ff95.png) `ifup eth0`即可发现eth1被启动 ![](/static/uploads/2017/12/0208660422e60d18fc30e73bd1291167.png) 但是这种解决方法只能解决上网问题，如果说我们需要让克隆主机的网卡eth0表示第一个网卡，eth1表示第二个网卡，这种方案是行不通的，并且很容易把网卡搞混淆。

## 升级版解决方案

先将`/etc/udev/rules.d/70-persistent-net.rules`中eth0删掉，并将其中的eth1改为eth0，之后将`/etc/sysconfig/network-script/ifcfg-eth0`的MAC改为上面文件中的MAC。 ![](/static/uploads/2017/12/4610d13b397c4944701f40871f741523.png) 最后，重启Linux之后，ifup eth0就可以发现IP获取正常了！ PS:这里试过只重启network服务，但是不能达到预想的效果，所以就直接重启Linux了。