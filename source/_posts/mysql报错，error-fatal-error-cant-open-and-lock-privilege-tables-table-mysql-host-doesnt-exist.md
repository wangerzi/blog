---
title: >-
  Mysql报错，[ERROR] Fatal error: Can't open and lock privilege tables: Table
  'mysql.host' doesn't exist
tags:
  - Mysql
id: '232'
categories:
  - - Linux
    - Mysql
date: 2018-04-29 13:59:31
cover: /static/uploads/2018/04/bg-1200x661.jpg
---

如果在使用mysql过程中，调用了mysqld\_safe更改了数据文件存放路径，但是新路径里边又没有需要的数据库文件，然后启动mysql时就可能出现这个报错。 `service mysqld start`

> mysqld\_safe --user mysql --datadir=/usr/local/data --datadir指定新的数据文件存放路径

如果有之间的数据文件，那么直接mv 原数据文件地址 现数据文件地址即可。 没有的话，直接执行mysql\_install\_db，重新生成数据库权限表。

> mysql\_install\_db 重新生成数据库权限表。

`# /usr/mysql_install_db` 最后，重启mysqld，还可以执行mysql\_secure\_installation初始化权限信息 `service mysqld restart`