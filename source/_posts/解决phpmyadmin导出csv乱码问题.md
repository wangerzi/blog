---
title: 解决phpMyAdmin导出csv乱码问题
tags:
  - php
  - phpMyAdmin
id: '136'
categories:
  - - 后端开发
date: 2018-03-10 17:13:34
cover: ../../static/uploads/2018/03/f349197440ceed70b72a1276b3afa024-1.png
---

# 背景

今天使用phpMyAdmin导出自定义sql数据时，采用csv格式导出，结果遇到了打开乱码的问题。

# 解决方法

导出方式选择“自定义”，然后在“文件的字符集”中选择“gb2312”，即可成功导出包含中文的结果！ ![](../static/uploads/2018/03/f349197440ceed70b72a1276b3afa024.png)