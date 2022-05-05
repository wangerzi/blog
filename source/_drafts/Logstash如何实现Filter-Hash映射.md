---
title: Logstash如何实现Filter Hash映射
tags: []
id: '632'
categories:
  - - 未分类
---



## 前言

Logstash 是 ELK 中重要的一环，时不时会遇到需要固定 key 映射到字符串的场景，但是 Logstash 的 Filtr 中不支持动态映射对象属性。