---
title: kali玩转暴力破解
tags: []
id: '178'
categories:
  - - 未分类
---

# 背景

上课的时候老师教了一招暴力破解，主要用于硬跑网站用户的用户密码，适用于撞库等情况。

## 准备工作

kali Linux 或者windows下安装Burp suite 1. 安装JRE 1. burp suite 代理过程 代理填写为127.0.0.1:8080，burp 中打开 interept，抓捕报文 选择爆破方法，队列爆破，用户名字典，密码字典