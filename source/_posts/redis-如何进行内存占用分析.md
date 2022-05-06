---
title: Redis-如何进行内存占用分析?
tags:
  - redis
  - 内存占用
id: '601'
categories:
  - - Linux
  - - 后端开发
date: 2020-07-01 16:26:47
cover: /static/uploads/2020/02/u15834439403268442047fm26gp0.jpg
---



## 目的

博主最近在工作上需要清理一个数年项目的 Redis，内存告警 90%，但是不能单靠看代码和直觉来清理，需要有个直观的参考指标。 并且线上的数据删除需要谨慎，删除的 key 众多，如果删出了问题势必需要快速回滚保证业务正常运行，所以还编写了一个安全删除的脚本，方便线上操作。

## 导出

在做全量的分析之前，首先需要拿到全量的 `.rdb` 文件。 如果是运行在 Linux 服务器中的 `redis-server`，可以直接使用 `dump.rdb`，存的东西太大的话就直接在从库主机上做分析即可，不用全量下载到本地。 如果是云服务商的话，就需要按照其官方文档的引导，获取到 `.rdb`，博主目前所在公司使用了 aws 的全套服务，针对 redis 实例导出 rdb 文件搜集了如下资料。 已经获取到 `.rdb` 文件的可以跳过这一段了，由于里边操作较为复杂为节省同样使用 aws 同学的时间，大概整理一下操作过程

*   首先，创建一个与 redis实例 **同区**的存储桶
*   然后给存储桶配置权限，进入桶的访问控制列表 > 其他AWS账户的访问权限，添加账户输入 `540804c33a284a299d2547575ce1010f2312ef3da9b3a053c8bc45bf233e4353`，权限勾选
    *   **列出对象**
    *   **写入对象**
    *   **读取存储桶权限**
*   然后，在 `ElastiCache控制面板` > `备份` > 选中备份 > 复制
*   写一个标志符，目标s3的位置选刚才建好的桶，点击复制，直到导出完成
*   s3中发现对应的 rdb 文件：memory-clear-1-0001.rdb，下载即可

> 手动备份的官方文档：https://docs.aws.amazon.com/zh\_cn/AmazonElastiCache/latest/red-ug/backups-manual.html 导出备份的官方文档：https://docs.aws.amazon.com/zh\_cn/AmazonElastiCache/latest/red-ug/backups-exporting.html

## 分析

拿到 `.rdb` 文件后，下一步就是分析这个 rdb 文件中各个 key 占用的内存了，网络上也提供到了相关的工具。

### 安装工具并分析

```shell
pip install rdbtools
pip install python-lzf # 加快转储速度
```

进入到 rdb 所在目录，生产内存报告

```shell
rdb -c memory  dump.rdb > dump.csv
```

下一步的分析，可以借助 excel 也可以借助 Mysql 等数据库工具导入 csv 后进行 CSV的内部结构大概是这样的：

database

type

key

size\_in\_bytes

encoding

num\_elements

len\_largest\_element

expiry

0

hash

Funcoloring\_update\_picture\_formal:com.cm.fun.color:Android:2.1.1:v3:large\_preview

459060

hashtable

1

432708

0

hash

ad\_stats:20200624:apps:com.ios.qyb.incolour

201

ziplist

4

19

2020-07-09T00:00:04.621000

0

hash

ad\_stats:20200617:com.cm.fun.color:upltv:inter\_open

206

ziplist

4

19

2020-07-02T00:00:33.171000

0

hash

ad\_stats:20200620:com.ios.qyb.incolour:MAX:reward\_hint

206

ziplist

3

19

2020-07-05T00:01:31.333000

### 关键列

第一列 `database` 表示数据库id，在excel中可以借助筛选功能，得到每个数据库中占用的具体key数量和内存大小。 第二列 `type` 顾名思义，数据类型 第三列 `key` 不多讲 第四列 `size_in_bytes` 就是占用字节数 其他...

### 如何定位问题 key

我一般的操作是，借助 excel 把 `size_in_bytes` / 1024 / 1024，得到 `Mb` 单位的内存占用量，再使用排序功能得到占用最多的列 我这边遇到的情况是大约 20% 的 key 占用了超过 80% 的空间，所以一般不用看太多就可以定位到问题。 除了找那些占用巨大的key，还需要重点查看 **没有过期时间的key**，讲道理一般场景下的 redis 不适合单独持久化数据，容易导致数据来源混乱不好维护。 找到 key 之后，根据 key 的规则去找相关的代码，根据逻辑将需要删除的 key 找出来，没有设置过期时间的改代码设置下过期时间。

### 恢复数据

如果你下载到了本地想要恢复线上的数据做测试或方便的查看内容，那么你可能需要通过 `.rdb` 文件恢复数据 如果安装了 docker，可以把下载的 rdb 重命名为 `dump.rdb`，映射到 `/data` 中，在 rdb 文件目录下可以使用如下指令，运行成功后

```shell
$ docker run --name test-redis -it -v `pwd`/dump.rdb:/data/dump.rdb --rm redis
1:C 01 Jul 2020 07:44:21.073 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
1:C 01 Jul 2020 07:44:21.073 # Redis version=6.0.1, bits=64, commit=00000000, modified=0, pid=1, just started
1:C 01 Jul 2020 07:44:21.073 # Warning: no config file specified, using the default config. In order to specify a config file use redis-server /path/to/redis.conf
```

然后进入容器使用 `redis-cli` 即可

```shell
$ docker exec -it test-redis /bin/bash
root@0910da87ae73:/data#
root@0910da87ae73:/data# redis-cli
127.0.0.1:6379>
127.0.0.1:6379> select 2
OK
127.0.0.1:6379[2]> KEYS *
1) "***:****:***"
```

如果本地搭好了环境，替换下 `dump.rdb` 即可，注意修改配置，否则可能加载失败

```ini
appendonly no
dbfilename dump.rdb # your rdb file name
dir /var/lib/redis # your rdb path
```

## 安全删除脚本

为方便回滚，编写了一个小脚本用于安全删除生产环境上的多个 key，执行删除的时候先把 key 按行放到 `keys.txt` 中，执行如下指令可以在删除的同时备份到 `delete-backup.json`，如果出了问题直接回滚就好。 脚本仅支持五大基本类型，`hyperloglog` 类型的数据占用不会太大，也没想法子去兼容，`list` 类型出于可重复性回滚数据的考虑，每次都会先删除记录再回滚，如果线上在 push 列表，有一定的丢数据风险，考虑停机再次导出并合并有关的 list 数据。

### 用法

```shell
python convert.py delete keys.txt  # safe delte and generate delete-backup.json
python convert.py revert delete-backup.json  # rollback by delete-backup.json
```

### 代码

```python
import redis
import sys
import os
import json


def check_backup_file(backup_path):
    if os.path.exists(backup_path):
        raise BaseException('you are already have a backup file ' + backup_path + ', please remove it first')


def op_backup(r, keys, output):
    check_backup_file(output)
    result = {}
    for key in keys:
        t = r.type(key)

        if t == 'string':
            val = r.get(key)
        elif t == 'list':
            val = r.lrange(key, 0, -1)
        elif t == 'set':
            val = list(r.smembers(key))
        elif t == 'zset':
            data = r.zrange(key, 0, -1, withscores=True)
            val = {}
            for v in data:
                val[v[0]] = v[1]
        elif t == "hash":
            val = r.hgetall(key)
        else:
            print('warning: ', key, ', type:', t, ' is not support')
            continue
        result[key] = {
            "type": t,
            "val": val,
        }
    with open(output, 'w') as f:
        json.dump(result, f)
    return result


def op_revert(r, result):
    for key in result:
        item = result[key]
        if item['type'] == 'string':
            r.set(key, item['val'])
        elif item['type'] == 'list':
            r.delete(key)
            r.rpush(key, *item['val'])
        elif item['type'] == 'set':
            r.sadd(key, *item['val'])
        elif item['type'] == 'zset':
            r.zadd(key, item['val'])
        elif item['type'] == "hash":
            r.hset(key, mapping=item["val"])
        else:
            print('revert warning: ', key, ', type:', item.type, ' is not support')
            continue


def op_delete(r, keys, backup_path):
    op_backup(r, keys, backup_path)
    r.delete(*keys)


def read_keys(filename):
    keys = []
    with open(filename) as f:
        line = f.readline()
        while line:
            line = line.strip()  # remove \n
            if not line == '':
                keys.append(line)
            line = f.readline()
    return keys


def read_json(filename):
    with open(filename) as f:
        return json.load(f)


def init_test_data():
    result = {
        'test-str': {'type': 'string', 'val': 'this is a long long string'},
        'test-list': {'type': 'list', 'val': ['str1', 'str2', 'str3']},
        'test-set': {'type': 'set', 'val': ['str1', 'str2', 'str3']},
        'test-zset': {'type': 'zset', 'val': {"key1": 6.0, "key2": 8.0, "key3": 10.0}},
        'test-hash': {'type': 'hash', 'val': {"key1": "123", "key2": "456", "key3": "789"}},
    }

    with open('test-backup.json', 'w') as f:
        json.dump(result, f)
    with open('test-keys.txt', 'w') as f:
        for key in result.keys():
            f.writelines(key + "\n")


def main():
    r = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

    if len(sys.argv) < 3:
        print("python convert.py backupdeleterevert keys.txtkeys.txtbackup.json")
        return False

    op = sys.argv[1]
    file = sys.argv[2]

    if op == 'backup':
        keys = read_keys(file)
        backup_path = 'backup.json'
        op_backup(r, keys, backup_path)
    elif op == 'delete':
        keys = read_keys(file)
        backup_path = 'delete-backup.json'
        check_backup_file(backup_path)
        op_delete(r, keys, backup_path)
    elif op == 'revert':
        result = read_json(file)
        op_revert(r, result)
    else:
        print("not support", op)


if __name__ == '__main__':
    main()
    # init_test_data()

```

## 总结

本文主要介绍了 redis 内存分析的方法和工具，以及如何恢复数据，如何进行分析，以及如何安全的删除，希望能有所帮助。