---
title: PHP高并发的注意事项以及数据库计数的优化
tags: []
id: '104'
categories:
  - - 后端开发
date: 2018-02-22 23:40:08
---



# 摘要

本文主要介绍了在对数据库更新数据，特别是使用ORM进行数据库操作时需要注意的事项，并且简单介绍了以下高频率计数的实现优化。

# 更新数据的注意事项

*   UPDATE 直接更改数据，更新的同时需要在条件中检查原始数据。 以微信扫码绑定为例：

```sql
CREATE TABLE wechat_bind(
    `wb_id`     INT(11) PRIMARY KEY COMMENT '主键',
    `openid`    CHAR(60) NOT NULL COMMENT '用户openid'
    `scene`     INT(11) NOT NULL COMMENT '场景值'
);
```

在这个案例里，在创建二维码时插入这样一条记录，OPENID为空表示二维码可使用，不为空则视为已绑定。 一般的解决思路是： 当扫描二维码进入时，通过scene查询到openid， `SELECT wechat_bind WHERE scene='xxx'` 如果openid为空，则执行绑定 `UPDATE wechat_bind SET openid='xxxxxx' WHERE wb_id=1` 但这样的思路在高并发的时候会出现绑定覆盖的问题，比如两个人几乎同时扫码，第一个人扫到，查询出openid为空，执行update，在update执行完毕之前，第二个扫码的进程开始，select出来的openid依旧为空，所以同样会执行update语句；由于update执行的时候会出现排他锁，一步一步执行，第二条update将会覆盖以一次的update从而造成更新异常。

## 解决方案

在更新数据的时候，为被更新的数据增加一条筛选条件，如： `update wechat_bind SET openid='xxxx' WHERE wb_id=1 AND openid = ''` 现在，如果出现并发的情况，两个人同时扫码，并巧合的形成了update队列，但在这种情况，第二条update执行结果将为空，不会出现更新覆盖的现象。

# 有最大绑定人数的更新注意事项

*   假设有这这样一种情况，一张二维码可以绑定多个用户，所以，我会这样设计数据表

```sql
CREATE TABLE wechat_qrcode(
    `wqid`  INT(11) PRIMARY KEY,
    `binded`    INT(11) NOT NULL COMMENT '已绑定人数',
    `maxbind`   INT(11) NOT NULL COMMENT '最大绑定人数',
    `scene`     INT(11) NOT NULL COMMENT '场景值'
);
```

在这个场景中，要实现最大绑定人数不超过数据库中的值，我们一般的做法是通过场景值(scene)找到最大绑定人数（maxbind）和已绑定人数(binded)，如果binded < maxbind 则执行绑定，binded加一，反之提示绑定失败 `SELECT * FROM wechat_qrcode WHERE scene=2` `UPDATE wechat_qrcode SET binded=binded+1 WHERE scene=2` 同上一种更新异常的情况，当两个用户同时扫码的时候，用户1的update可能已经将binded变成了maxbind，但取数据的时候，用户2并不知道binded已经达到maxbind了，依然执行binded+1，导致绑定用户大于最大绑定用户，解决方法也类似，将更新语句加一个筛选条件就好。 `UPDATE wecaht_qrcode SET binded=bind+1 WHERE scene=2 AND binded < maxbind` 注意：这里是小于，要是判断条件为<=则binded+1就可能大于maxbind了。

## 关于高频率计数的优化

在高并发状态下，如果借助自带并发访问队列控制的工具（redis、memecache、mysql等），系统的记录效率会受到影响，之前机缘巧合在某博客下看到了一个比较巧妙，并且效率挺高的做法。 背景：一个联网游戏，后台采用PHP编写，每次访问游戏API，API访问量+1，记得原博客所述访问量在1W/s左右 如果维护一个数据表，每次访问update相关字段会导致mysql压力增大，并且每次显示都需要访问此数据库，所以这不是一个合适的解决方案。 最后的解决方案是，使用一个文件来存储网站访问量，每次访问，都往文件后追加一个1，这样既不会造成丢失，并且数据库/内存缓存压力也不会太大 `<?php // 做统计 $fp = fopen('total.txt', 'a'); fwrite($fp, '1'); fclose($fp);`

# 小结

面对高并发，有时候并不需要使用事务等复杂的操作，可能有时候只是需要限制一下更新条件就能更好更方便的避免高并发数据更新错误的情况。