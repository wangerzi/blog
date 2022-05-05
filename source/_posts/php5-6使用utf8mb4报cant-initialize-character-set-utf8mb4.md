---
title: PHP5.6使用utf8mb4报Can't initialize character set utf8mb4
tags:
  - Linux
  - mysql
  - nginx
  - php
id: '125'
categories:
  - - Linux
  - - 后端开发
date: 2018-03-07 21:16:59
---

# 背景

今天部署环境的时候，使用了mysql数据库的utf8mb4;因为之前就知道mysql版本太低不能使用utf8mb4，所以早早就准备好了mysql5.6。 结果部署好之后发现报错：Can't initialize character set utf8mb4 (path: /usr/share/mysql/charsets/) in ...

# 解决方法：

经过一番查找，终于在stackoverflow上找到一篇帖子，解决了这个问题 原帖： https://stackoverflow.com/questions/33834191/php-pdoexception-sqlstatehy000-2019-cant-initialize-character-set-utf8mb4

```null
yum erase php56w-mysql
yum install php56w-mysqlnd
```