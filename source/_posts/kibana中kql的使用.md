---
title: Kibana中KQL的使用
tags: []
id: '620'
categories:
  - - Linux
    - ElasticSearch
date: 2020-08-26 21:17:33
cover: /static/uploads/2020/08/Monasterio_Khor_Virap_Armenia_2016-10-01_DD_25.jpg
---



## 前言

当我们需要查看 ElasticSearch 中存放的数据时，通常会使用 Kibana 这个可视化工具。 但是 Kibana 中的 Discover 页面默认只会展示最近收到的数据，当我们需要查询符合某个条件数据时，就需要用到 KQL(Kibana Query Language) 了，没有接触之前觉得高大上，看完之后才发现设计如此简单。 英文还可以的同学可以直接看官方文档，结合数据实例进行讲解，也比较清晰 官方文档：[https://www.elastic.co/guide/en/kibana/current/kuery-query.html#kuery-query](https://www.elastic.co/guide/en/kibana/current/kuery-query.html#kuery-query)

## 筛选语法

将所有涉及到的语法铺展开来，首先准备好官网文档中的数据如下：

```json
{
  "grocery_name": "Elastic Eats",
  "items": [
    {
      "name": "banana",
      "stock": "12",
      "category": "fruit"
    },
    {
      "name": "peach",
      "stock": "10",
      "category": "fruit"
    },
   {
      "name": "peach test",
      "stock": "10",
      "category": "fruit"
    },
    {
      "name": "carrot",
      "stock": "9",
      "category": "vegetable"
    },
    {
      "name": "broccoli",
      "stock": "5",
      "category": "vegetable"
    }
  ]
}
```

需要根据实际情况做尝试的时候进入 Kibana 的 Discover页面在输入框中填入筛选即可。 [![](/static/uploads/2020/08/bf38129fb9a6a1d593feb32d28b50582.png)](/static/uploads/2020/08/bf38129fb9a6a1d593feb32d28b50582.png)

> 生产环境截图不便还请谅解 😁

### 简单查询

简单查询就是 关键字匹配、字符串包含等，比如说如下语句会找出 name 字段是 banana 的所有数据：

```kql
name: banana
```

但是如果 name 包含 `peach` 和 `peach test`，然后下面两个语句查出来会是两个结果。

```kql
name: peach test
```

上述查询会将 name 是 `peach` 和 name 是 `peach test` 的都给查出来

```kql
name: "peach test"
```

上述查询只会将 `peach test` 查出来，因为如果不加引号会自动关键字分词，将包含该关键字的所有数据匹配出来。

### 条件运算符

条件运算符就是 > >= < <=，在 KQL 里边都支持，使用也很简单，比如如下语句表示 age 字段大于等于 10。

```kql
age >= 10
```

### 逻辑运算符

查询语言自然少不了逻辑运算符 与或非，在 KQL 中代表了 and or not and 的用法：

```kql
age >= 10 and age < 100
```

上述语句表示查询出 age 在 10 到 100 的左开右闭区间中的所有数据。 or 的用法：

```kql
name: "Jeff" or name: "Kitty"
```

上述语句表示筛选出 name 包含 `Jeff` 或者 `Kitty` 关键字的所有数据。 not 的用法：

```kql
not age >= 10
```

上述语句表示筛选出 age 小于 10 的所有数据。 其中 and 的优先级比 or 的高

```kql
age < 100 or name: wang and age >= 10
```

and 优先级高会先结合，所以意思是 满足 name 是wang age >= 10 或者 age < 100。 当然也可以通过小括号来改变优先级，比如：

```kql
(age < 100 or name: wang) and age >= 10
```

意思是 age >=10 并且这条数据的 name是wang或者age < 100

#### 同一字段运算符简写

可以用括号将多个逻辑运算符和条件合并到一起

```kql
age = 10 or age = 100
# 等价于
age: ( 10 or 100)
```

### 通配符

通配符可以用于查找出存在某个key的数据

```kql
name: *
```

表示查找出所有带 name 字段的数据

```kql
system: win*
```

可以匹配到 system: win7，system: win10 等。

### 字段嵌套查询

首先准备一个多层的数据，比如下面的这几条数据。

```json
{
  "level1": [
    {
      "level2": [
        {
          "prop1": "foo",
          "prop2": "bar"
        },
        {
          "prop1": "baz",
          "prop2": "qux"
        }
      ]
    }
  ]
}
```

比如想筛选 level1.level2.prop1 是 `foo` 或者是 `baz`的，可以这样写：

```kql
level1.level2 { prop1: "foo" or prop1: "baz" }
```

## 总结

KQL是一个比较简单筛选数据的查询语言，包括条件、逻辑、多层查询等用法，能辅助报表的制作和实时日志的筛选。