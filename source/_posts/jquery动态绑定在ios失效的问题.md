---
title: jquery动态绑定在IOS失效的问题
tags:
  - ios
  - 动态绑定
  - 踩坑
id: '202'
categories:
  - - 前端开发
date: 2018-04-08 21:59:01
cover: /static/uploads/2018/03/TIM%E6%88%AA%E5%9B%BE20180321022207.png
---



# 背景

在一个微信项目中，由于某些元素是动态添加的，所以使用了jquery的on('click', 'xx', function(){});来进行动态绑定点击事件，但是在IOS下，此法失效，经过一番资料查找，终于找到了解决方案。

# 问题复现

html代码：

```markup
<html>
    <head>
        <title>测试</title>
    </head>
    <body>
        <ul>
            <li class="elem">操作1</li>
            <li class="elem">操作1</li>
            <li class="elem">操作1</li>
            <li class="elem">操作1</li>
            <li class="elem">操作1</li>
            <li class="elem">操作1</li>
        </ul>
    </body>
    <!--引入jquery-->
    <script src="js/jquery.js"></script>
    <script>
        // 给.elem增加点击事件
        $('body').on('click', '.elem', function(){
            alert('已点击');
        });
    </script>
</html>
```

> 使用电脑调试工具以及android均没有问题，但是在IOS的微信浏览器/safari浏览器均无法触发点击事件。

# 解决方案

经查，在IOS中，点击对象需要拥有"cursor:pointer;"样式才能触发点击事件，具体原因据说是无法找到DOM。 所以，上边的代码，只需要加一个样式，即可生效。

```css
.elem {cursor:pointer;}
```