---
title: 令最新JS-XLSX支持样式的改造方法
tags:
  - excel
  - js-xlsx
  - style
  - 导出样式
  - 改造
id: '360'
categories:
  - - 前端开发
date: 2019-05-01 15:19:24
---

## 背景

最近五一小长假，统一处理下之前开源的 [Layui导出插件 lay-excel](https://github.com/wangerzi/layui-excel) 反馈的部分问题，这个插件核心使用的是经过改造的 protobi/js-xlsx，支持设置样式但是不支持诸如 导出文件压缩、边距设置等功能，还存在很多BUG，效率也不高。 为解决这些问题，博主开始着手改造最新 JS-XLSX 让其支持样式设置，断点调试再加上代码对比，最后改造成功。 改造完毕的 xlsx.js 在此，可自行引入： [https://github.com/wangerzi/layui-excel/blob/master/src/xlsx.js](https://github.com/wangerzi/layui-excel/blob/master/src/xlsx.js) 记得配合 jszip.js 使用，否则报错： [https://github.com/wangerzi/layui-excel/blob/master/src/jszip.js](https://github.com/wangerzi/layui-excel/blob/master/src/jszip.js)

## 开始前的准备

> SheetJs/js-xlsx是目前依旧维护的最新版本，protobi/js-xlsx 是大约两年前的支持设置样式的版本，很多比较实用的功能都没有。

开源项目名称

地址

用于

[SheetJS / js-xlsx](https://github.com/SheetJS/js-xlsx)

[https://github.com/SheetJS/js-xlsx](https://github.com/SheetJS/js-xlsx)

导出的基础逻辑

[protobi / js-xlsx](https://github.com/protobi/js-xlsx)

[https://github.com/protobi/js-xlsx](https://github.com/protobi/js-xlsx)

可以设置样式，用于补全样式功能

## 改造思路

博主做了个插件预览放在了服务器上，如果想了解具体的改造原理，可以跟着思路在 [http://excel.wj2015.com](http://excel.wj2015.com) 上自行调试一下

### 断点调试，找有关样式设置的代码

*   点击批量设置样式按钮，在调用行处打断点

```javascript
var wbout = XLSX.write(wb, {bookType: type, type: 'binary', cellStyles: true});
```

![](../static/uploads/2019/05/59cb58df4b8239af744e7dedb4c202dd.png)

*   发现xlsx导出进入了此函数

![](../static/uploads/2019/05/0d060dd000f14d4248d9a0b690d158dc.png)

*   进一步调试，发现了差异

![](../static/uploads/2019/05/ec1dbef1c7bb6294902479f9476f6dad.png)

#### 根据支持样式的版本，补全缺失代码

*   补全样式代码

> 进入最新版[https://github.com/SheetJS/js-xlsx/blob/master/dist/xlsx.js](https://github.com/SheetJS/js-xlsx/blob/master/dist/xlsx.js)，搜索 `StyleBuilder` 发现没有这个类，但是在 [https://github.com/protobi/js-xlsx/blob/master/dist/xlsx.js](https://github.com/protobi/js-xlsx/blob/master/dist/xlsx.js) 发现了这个类，与其依赖一起复制过来

![](../static/uploads/2019/05/53ca8f3c054d618a5c642f00b64d2a70.png)

*   复制到目标文件暴露全局变量之前

![](../static/uploads/2019/05/9d9e1b38659099dbf158141ef8946282.png)

*   查找 `style_builder` 这个变量在 `protobi/js-xlsx` 中有哪些地方使用到了，依次复制到 `SheetJS/protobi` 中

![](../static/uploads/2019/05/a236429e076ed4b49cdaf83b0a9917eb.png) ![](../static/uploads/2019/05/19a56922a11ed2b160636e0a739abb6e.png) PS:几年过去了，这里的逻辑基本没变，所以放心复制好了

*   最后，测试一下导出，发现可以设置样式了

![](../static/uploads/2019/05/31876adf7ddaf845e9a6be32dcb9cd4f.png)