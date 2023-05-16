---
title: Ubuntu 22.04 安装使用绿联千兆网口 AX88179
date: 2023-05-16 06:47:54
tags:
- NAS
categories:
- Linux
cover: ../static/assets/2023-05-16-06-59-00-img_v2_8c473334-bc0f-4e0c-b20a-ee753896e0bh.jpg
---

## 前言

此篇博客的目标是为 Ubuntu 22.04 的服务器扩展一个千兆网口，JD 上踩过其他品牌的坑，最终还是选择了绿联（PS:没给广告费），虽然官网比较简陋提供的驱动也都还有点问题，但好歹是 Linux 下一番折腾还是能用的。

![](../static/assets/2023-05-16-06-59-00-img_v2_8c473334-bc0f-4e0c-b20a-ee753896e0bh.jpg)

## 自带驱动问题

我使用的是 Ubuntu 22.04 的操作系统，内核版本 5.19，里边自带一个 ax88179_179a  的内核模块，使用 `lsmod` 和 `modinfo` 命令可查看相关信息

```bash
$ uname -a
Linux wang-B250ET2 5.19.0-41-generic #42~22.04.1-Ubuntu SMP PREEMPT_DYNAMIC Tue Apr 18 17:40:00 UTC 2 x86_64 x86_64 x86_64 GNU/Linux
$ lsmod | grep "ax88179"
ax88179_178a           36864  0
usbnet                 53248  1 ax88179_178a
mii                    20480  2 usbnet,ax88179_178a
$ modinfo ax88179_178a
filename:       /lib/modules/5.19.0-41-generic/kernel/drivers/net/usb/ax88179_178a.ko
license:        GPL
description:    ASIX AX88179_178A USB 2.0/3.0 Ethernet Devices
author:         David Hollis
srcversion:     BDE6AA432409AEB043AF2AA
alias:          usb:v0711p0179d*dc*dsc*dp*ic*isc*ip*in*
alias:          usb:v2001p4A00d*dc*dsc*dp*ic*isc*ip*in*
alias:          usb:v04E8pA100d*dc*dsc*dp*ic*isc*ip*in*
alias:          usb:v0930p0A13d*dc*dsc*dp*ic*isc*ip*in*
alias:          usb:v17EFp304Bd*dc*dsc*dp*ic*isc*ip*in*
alias:          usb:v0DF6p0072d*dc*dsc*dp*ic*isc*ip*in*
alias:          usb:v0B95p178Ad*dc*dsc*dp*ic*isc*ip*in*
alias:          usb:v0B95p1790d*dc*dsc*dp*ic*isc*ip*in*
depends:        mii,usbnet
retpoline:      Y
name:           ax88179_178a
vermagic:       5.19.0-41-generic SMP preempt mod_unload modversions 
parm:           msg_enable:usbnet msg_enable (int)
parm:           bsize:RX Bulk IN Queue Size (int)
parm:           ifg:RX Bulk IN Inter Frame Gap (int)
parm:           bEEE:EEE advertisement configuration (int)
parm:           bGETH:Green ethernet configuration (int)
parm:           tx_dma_sg:Whether to use the dma_sg feature for tx if supported, "no" by default (bool)
```

插入之后直接就可以识别出来设备，执行 `lsusb` 可以看到一个 `AX88179 Gigabit Ethern` 的设备，但是 `ip addr` 是看不到新增的网口

```bash
$ lsusb 
Bus 002 Device 003: ID 1058:2622 Western Digital Technologies, Inc. Elements SE 2622
Bus 002 Device 004: ID 0b95:1790 ASIX Electronics Corp. AX88179 Gigabit Ethernet
Bus 002 Device 002: ID 05e3:0626 Genesys Logic, Inc. USB3.1 Hub
Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
Bus 001 Device 003: ID 046d:c52b Logitech, Inc. Unifying Receiver
Bus 001 Device 002: ID 062a:4101 MosArt Semiconductor Corp. Wireless Keyboard/Mouse
Bus 001 Device 004: ID 05e3:0610 Genesys Logic, Inc. Hub
Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub
```

可能是设备启动时出了问题，那么下一步就是看 `dmesg`，从输出中很容易就找到了下面这个错误

> ax88179_178a 2-5.4:1.0 (unnamed net_device) (uninitialized): Failed to read reg index 0x0040: -32

## 安装新驱动

既然内核自带的驱动不好使，这个时候我就考虑重新安装驱动了

我首先尝试的是绿联官网的千兆网卡驱动：[绿联Type-C千兆网卡驱动下载](https://www.lulian.cn/download/4-cn.html)

下载压缩包并解压其中的 Linux 驱动 .tar.bz2 文件

```bash
$ tar -jxvf AX88179_178A_LINUX_DRIVER_v1.20.0_SOURCE.tar.bz2 
AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.h
AX88179_178A_Linux_Driver_v1.20.0_source/Makefile
AX88179_178A_Linux_Driver_v1.20.0_source/readme
AX88179_178A_Linux_Driver_v1.20.0_source/
AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c
```

然后安装编译环境

```bash
sudo apt-get install make gcc
```

执行编译，如果你内核版本跟我一样(5.13+)，大概会碰到这样的错误：`error: ‘usbnet_get_stats64’ undeclared here (not in a function); did you mean ‘usbnet_cdc_status’?`

```bash
$ make
make -C /lib/modules/5.19.0-41-generic/build M=/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source modules
make[1]: 进入目录“/usr/src/linux-headers-5.19.0-41-generic”
warning: the compiler differs from the one used to build the kernel
  The kernel was built by: x86_64-linux-gnu-gcc (Ubuntu 11.3.0-1ubuntu1~22.04.1) 11.3.0
  You are using:           gcc (Ubuntu 11.3.0-1ubuntu1~22.04) 11.3.0
  CC [M]  /home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.o
In file included from ./include/linux/string.h:253,
                 from ./include/linux/bitmap.h:11,
                 from ./include/linux/cpumask.h:12,
                 from ./arch/x86/include/asm/cpumask.h:5,
                 from ./arch/x86/include/asm/msr.h:11,
                 from ./arch/x86/include/asm/processor.h:22,
                 from ./arch/x86/include/asm/timex.h:5,
                 from ./include/linux/timex.h:67,
                 from ./include/linux/time32.h:13,
                 from ./include/linux/time.h:60,
                 from ./include/linux/stat.h:19,
                 from ./include/linux/module.h:13,
                 from /home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:30:
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c: In function ‘ax88179_set_mac_addr’:
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1000:19: warning: passing argument 1 of ‘__builtin_memcpy’ discards ‘const’ qualifier from pointer target type [-Wdiscarded-qualifiers]
 1000 |         memcpy(net->dev_addr, addr->sa_data, ETH_ALEN);
      |                ~~~^~~~~~~~~~
./include/linux/fortify-string.h:379:27: note: in definition of macro ‘__fortify_memcpy_chk’
  379 |         __underlying_##op(p, q, __fortify_size);                        \
      |                           ^
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1000:9: note: in expansion of macro ‘memcpy’
 1000 |         memcpy(net->dev_addr, addr->sa_data, ETH_ALEN);
      |         ^~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1000:19: note: expected ‘void *’ but argument is of type ‘const unsigned char *’
 1000 |         memcpy(net->dev_addr, addr->sa_data, ETH_ALEN);
      |                ~~~^~~~~~~~~~
./include/linux/fortify-string.h:379:27: note: in definition of macro ‘__fortify_memcpy_chk’
  379 |         __underlying_##op(p, q, __fortify_size);                        \
      |                           ^
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1000:9: note: in expansion of macro ‘memcpy’
 1000 |         memcpy(net->dev_addr, addr->sa_data, ETH_ALEN);
      |         ^~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1004:47: warning: passing argument 6 of ‘ax88179_write_cmd’ discards ‘const’ qualifier from pointer target type [-Wdiscarded-qualifiers]
 1004 |                                  ETH_ALEN, net->dev_addr);
      |                                            ~~~^~~~~~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:216:46: note: expected ‘void *’ but argument is of type ‘const unsigned char *’
  216 |                              u16 size, void *data)
      |                                        ~~~~~~^~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c: At top level:
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1030:35: error: ‘usbnet_get_stats64’ undeclared here (not in a function); did you mean ‘usbnet_cdc_status’?
 1030 |         .ndo_get_stats64        = usbnet_get_stats64,
      |                                   ^~~~~~~~~~~~~~~~~~
      |                                   usbnet_cdc_status
In file included from ./include/linux/string.h:253,
                 from ./include/linux/bitmap.h:11,
                 from ./include/linux/cpumask.h:12,
                 from ./arch/x86/include/asm/cpumask.h:5,
                 from ./arch/x86/include/asm/msr.h:11,
                 from ./arch/x86/include/asm/processor.h:22,
                 from ./arch/x86/include/asm/timex.h:5,
                 from ./include/linux/timex.h:67,
                 from ./include/linux/time32.h:13,
                 from ./include/linux/time.h:60,
                 from ./include/linux/stat.h:19,
                 from ./include/linux/module.h:13,
                 from /home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:30:
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c: In function ‘access_eeprom_mac’:
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1436:32: warning: passing argument 1 of ‘__builtin_memcpy’ discards ‘const’ qualifier from pointer target type [-Wdiscarded-qualifiers]
 1436 |                 memcpy(dev->net->dev_addr, buf, ETH_ALEN);
      |                        ~~~~~~~~^~~~~~~~~~
./include/linux/fortify-string.h:379:27: note: in definition of macro ‘__fortify_memcpy_chk’
  379 |         __underlying_##op(p, q, __fortify_size);                        \
      |                           ^
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1436:17: note: in expansion of macro ‘memcpy’
 1436 |                 memcpy(dev->net->dev_addr, buf, ETH_ALEN);
      |                 ^~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1436:32: note: expected ‘void *’ but argument is of type ‘const unsigned char *’
 1436 |                 memcpy(dev->net->dev_addr, buf, ETH_ALEN);
      |                        ~~~~~~~~^~~~~~~~~~
./include/linux/fortify-string.h:379:27: note: in definition of macro ‘__fortify_memcpy_chk’
  379 |         __underlying_##op(p, q, __fortify_size);                        \
      |                           ^
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1436:17: note: in expansion of macro ‘memcpy’
 1436 |                 memcpy(dev->net->dev_addr, buf, ETH_ALEN);
      |                 ^~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c: In function ‘ax88179_get_mac’:
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1498:54: warning: passing argument 2 of ‘access_eeprom_mac’ discards ‘const’ qualifier from pointer target type [-Wdiscarded-qualifiers]
 1498 |                 ret = access_eeprom_mac(dev, dev->net->dev_addr, 0x0, 1);
      |                                              ~~~~~~~~^~~~~~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1399:54: note: expected ‘u8 *’ {aka ‘unsigned char *’} but argument is of type ‘const unsigned char *’
 1399 | static int access_eeprom_mac(struct usbnet *dev, u8 *buf, u8 offset, int wflag)
      |                                                  ~~~~^~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:1535:45: warning: passing argument 6 of ‘ax88179_write_cmd’ discards ‘const’ qualifier from pointer target type [-Wdiscarded-qualifiers]
 1535 |                           ETH_ALEN, dev->net->dev_addr);
      |                                     ~~~~~~~~^~~~~~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:216:46: note: expected ‘void *’ but argument is of type ‘const unsigned char *’
  216 |                              u16 size, void *data)
      |                                        ~~~~~~^~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c: In function ‘ax88179_reset’:
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:2104:45: warning: passing argument 6 of ‘ax88179_write_cmd’ discards ‘const’ qualifier from pointer target type [-Wdiscarded-qualifiers]
 2104 |                           ETH_ALEN, dev->net->dev_addr);
      |                                     ~~~~~~~~^~~~~~~~~~
/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.c:216:46: note: expected ‘void *’ but argument is of type ‘const unsigned char *’
  216 |                              u16 size, void *data)
      |                                        ~~~~~~^~~~
make[2]: *** [scripts/Makefile.build:257：/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/ax88179_178a.o] 错误 1
make[1]: *** [Makefile:1850：/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source] 错误 2
make[1]: 离开目录“/usr/src/linux-headers-5.19.0-41-generic”
make: *** [Makefile:30：default] 错误 2
```

于是乎，我拿着错误在网上搜了搜，尝试了很多解决方案均无法正常使用网卡，所以我求助了绿联的专业客服，但他们给我的回复是：“仅支持 5.13 以下内核，建议降低内核版本试试，更新的内核驱动修复日期不详”

![](../static/assets/2023-05-16-07-17-23-img_v2_22e39f24-bbc3-41f8-b007-d48d1a28f89h.jpg)

不过内核版本直接影响系统稳定，肯定是不能随便降的，但是作为一个折腾服务器的老网络人，还是没有放弃，继续寻找解决方案，最终让我发现了这么一个 ISSUE: [Compile errors with kernel 5.16.10 // Getting your imrovements upstream into the Linux kernel · Issue #1 · nothingstopsme/AX88179_178A_Linux_Driver · GitHub](https://github.com/nothingstopsme/AX88179_178A_Linux_Driver/issues/1)

作者说已经解决了这个 BUG，于是乎我就下载了这个开源的驱动仓库并尝试编译，成功！

```bash
$ git clone https://github.com/nothingstopsme/AX88179_178A_Linux_Driver
$ cd AX88179_178A_Linux_Driver/source/
$ make
make -C /lib/modules/5.19.0-41-generic/build M=/home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/AX88179_178A_Linux_Driver/source modules
make[1]: 进入目录“/usr/src/linux-headers-5.19.0-41-generic”
warning: the compiler differs from the one used to build the kernel
  The kernel was built by: x86_64-linux-gnu-gcc (Ubuntu 11.3.0-1ubuntu1~22.04.1) 11.3.0
  You are using:           gcc (Ubuntu 11.3.0-1ubuntu1~22.04) 11.3.0
  CC [M]  /home/wang/AX88179_178A_Linux_Driver_v1.20.0_source/AX88179_178A_Linux_Driver/source/ax88179_178a.o
# 做 make install 之前，可以
$ make install
```

编译完成后会出现编译产物 `ax88179_178a.ko`，使用 `make install` 安装此模块，如果担心玩坏了，可以先备份一下 `/lib/modules/5.19.0-41-generic/kernel/drivers/net/usb/ax88179_178a.ko`

```bash
$ ls -lah
总计 1.4M
drwxrwxr-x 2 wang wang 4.0K  5月 16 07:20 .
drwxrwxr-x 4 wang wang 4.0K  5月 16 07:20 ..
-rw-rw-r-- 1 wang wang  68K  5月 16 07:20 ax88179_178a.c
-rw-rw-r-- 1 wang wang  12K  5月 16 07:20 ax88179_178a.h
-rw-rw-r-- 1 wang wang 574K  5月 16 07:20 ax88179_178a.ko
-rw-rw-r-- 1 wang wang  497  5月 16 07:20 .ax88179_178a.ko.cmd
-rw-rw-r-- 1 wang wang  100  5月 16 07:20 ax88179_178a.mod
-rw-rw-r-- 1 wang wang 3.7K  5月 16 07:20 ax88179_178a.mod.c
-rw-rw-r-- 1 wang wang  364  5月 16 07:20 .ax88179_178a.mod.cmd
-rw-rw-r-- 1 wang wang  54K  5月 16 07:20 ax88179_178a.mod.o
-rw-rw-r-- 1 wang wang  32K  5月 16 07:20 .ax88179_178a.mod.o.cmd
-rw-rw-r-- 1 wang wang 522K  5月 16 07:20 ax88179_178a.o
-rw-rw-r-- 1 wang wang  59K  5月 16 07:20 .ax88179_178a.o.cmd
-rw-rw-r-- 1 wang wang 1.2K  5月 16 07:20 Makefile
-rw-rw-r-- 1 wang wang  101  5月 16 07:20 modules.order
-rw-rw-r-- 1 wang wang  343  5月 16 07:20 .modules.order.cmd
-rw-rw-r-- 1 wang wang    0  5月 16 07:20 Module.symvers
-rw-rw-r-- 1 wang wang  382  5月 16 07:20 .Module.symvers.cmd
-rw-rw-r-- 1 wang wang 3.5K  5月 16 07:20 readme
# 执行此操作前可以先备份一下
$ make install
# 安装完后重启
$ sudo reboot
```

重启前 `dmesg` 检查下，没有再次出现之前的报错就算问题解决了，重启后就能看到新识别出来的网卡

```bash
$ ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: enp2s0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    link/ether b8:97:5a:ff:45:cc brd ff:ff:ff:ff:ff:ff
3: enx000ec68ee90d: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    link/ether 00:0e:c6:8e:e9:0d brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.16/24 brd 192.168.1.255 scope global dynamic noprefixroute enx000ec68ee90d
       valid_lft 56815sec preferred_lft 56815sec
    inet6 240e:388:9f06:100:f4d:af8b:e099:f2d3/64 scope global temporary dynamic 
       valid_lft 259180sec preferred_lft 56494sec
    inet6 240e:388:9f06:100:3497:293c:4393:a48d/64 scope global dynamic mngtmpaddr noprefixroute 
       valid_lft 259180sec preferred_lft 172780sec
    inet6 fe80::8b25:81fe:35fb:b655/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever
```


翻阅了一下代码，发现是内核 API 变动导致的，作者使用 #IF 做了内核版本的判断，就解决了这个 BUG

```c
#if LINUX_VERSION_CODE > KERNEL_VERSION(2, 6, 29)
static const struct net_device_ops ax88179_netdev_ops = {
    .ndo_open        = usbnet_open,
    .ndo_stop        = usbnet_stop,
    .ndo_start_xmit        = usbnet_start_xmit,
    .ndo_tx_timeout        = usbnet_tx_timeout,
    .ndo_change_mtu        = ax88179_change_mtu,
    .ndo_do_ioctl        = ax88179_ioctl,
    .ndo_set_mac_address    = ax88179_set_mac_addr,
    .ndo_validate_addr    = eth_validate_addr,
#if LINUX_VERSION_CODE <= KERNEL_VERSION(3, 2, 0)
    .ndo_set_multicast_list    = ax88179_set_multicast,
#else
    .ndo_set_rx_mode    = ax88179_set_multicast,
#endif
#if LINUX_VERSION_CODE >= KERNEL_VERSION(2, 6, 39)
    .ndo_set_features    = ax88179_set_features,
#endif
#if LINUX_VERSION_CODE >= KERNEL_VERSION(5, 11, 0)
    .ndo_get_stats64    = dev_get_tstats64,
#elif LINUX_VERSION_CODE >= KERNEL_VERSION(4, 12, 0)
    .ndo_get_stats64    = usbnet_get_stats64,
#endif
};
#endif
```

随后笔者为绿联客服反馈了相关解决方案，不过在笔者写这篇文章的时候，绿联官网上千兆网卡驱动的最后更新时间还是 "2022/07/29"，如果你也有类似的需求，请注意不要下错了。

![](../static/assets/2023-05-16-07-17-33-img_v2_f66d4ec8-5f33-402f-9451-6400aa70a08h.jpg)

## 总结

所谓官方驱动不一定是真”官方“，碰到问题还是需要多在网上搜索尝试，希望能帮到你。
