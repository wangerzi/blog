---
title: mysql5.6部分字段无默认值不能insert的问题
tags:
  - Linux
  - Mysql
id: '127'
categories:
  - - Linux
    - Mysql
date: 2018-03-07 21:41:35
cover: /static/uploads/2018/04/bg-1200x661.jpg
---

# 背景

部署生产环境(mysql5.6)的时候，发现本来在测试环境(mysql5.5)好好的程序，到生产环境却报某字段没有默认值，不能添加的错误，经过一番资料查找，终于定位到/etc/my.cnf里配置的sql\_mode身上。

# 处理办法

如果按照严格的规范，没有默认值的字段是不能在加入的时候忽略的，但是程序已经写在那里，修改的话工作量会很大，所以这里我选择了将`/etc/my.cnf`中的sql\_mode那一行去掉。 我试过只去掉 STRICT\_TANS\_TABLES，但是不管用，全部注释后成功 # Recommended in standard MySQL setup # sql\_mode=NO\_ENGINE\_SUBSTITUTION,STRICT\_TRANS\_TABLES # sql-mode=NO\_AUTO\_CREATE\_USER,NO\_ENGINE\_SUBSTITUTION

# 引用

mysql sql\_mode 解决数据库非空无默认值依然可以插入的问题 http://blog.csdn.net/ctrlk/article/details/52742434 MySQL 5.6中的sql\_mode默认设置问题 http://blog.csdn.net/micahriven/article/details/12030981