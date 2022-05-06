---
title: JS中Function.apply() 的骚操作
tags:
  - apply
  - js
id: '311'
categories:
  - - 前端开发
date: 2018-11-20 20:50:02
cover: ../static/uploads/2019/09/20161107111502_84-1200x661.jpg
---



## 基本介绍

### Function.apply()是什么

apply() 方法是 JS中所有函数默认都有的方法，IE5+支持，如果需要 IE5 的支持的话，可以参考博客 [https://www.jb51.net/article/44151.htm](https://www.jb51.net/article/44151.htm) 都知道，在JS中可以通过 `this` 关键字获取当前执行环境，比如：

```javascript
var obj = {
    name: 'Wang',
    init: function(){
        this.name = 'Jeffrey';
    }
};
obj.init()
console.log(obj.name);// 输出Jeffrey
```

而 `apply` 就可以改变这个上下文环境，比如下边这个例子就是将 obj1 的执行环境变成了 obj2 的执行环境，执行 obj1 的 `init` 之后，却将 obj2 的属性改变了。

```javascript
var obj = {
    name: 'Wang',
    init: function(name){
        this.name = name;
    }
};
var obj2 = {}；
obj.init.apply(obj2, ['Test'])
console.log(obj2.name);// 输出 Test
```

`apply()` 有两个参数，第一个参数是上下文环境的对象，第二个参数是函数列表，支持数组形式传递

### apply()和call()的区别

`call()` 方法也可以改变函数执行的上下文环境，但是与 `apply()` 有一定的区别，注意看以下代码：

```javascript
var obj = {
    name: 'Wang',
    desc: '',
    init: function(name, desc){
        this.name = name;
        this.desc = desc;
    }
};
var obj2 = {}；
obj.init.apply(obj2, ['Test', 'Hello World'])
console.log(obj2.name);// 输出 Test
console.log(obj2.desc);// 输出 Hello World
obj.init.call(obj2, 'Test2', 'Hello World!')
console.log(obj2.name);// 输出 Test2
console.log(obj2.desc);// 输出 Hello World!
```

> apply 传递参数的方式是数组，call 传递参数的方式是函数列表

### 所谓骚操作

骚操作，，就是一般想不到，然后偶尔打代码的时候突然灵光一闪，一些骚操作就能达到 简化代码，增加效率等功效。

##### 骚操作：数组合并

假设现在有『数组A』(1, 2)，还有『数组B』(3, 4, 5)，我们希望把 『数组B』 的数据追加到 『数组A』中，了解`contact` 函数的童鞋，可能会写出这样的代码：

```javascript
var A = [1, 2];
var B = [3, 4, 5];

A = A.contact(B);
console.log(A);// 输出 1,2,3,4,5
```

原理是，contact函数能够连接多个数组，在**不改变原有数组的前提**下，将所有列表合并在一起作为返回值返回到调用方。 所以，在这个过程中，会出现一个临时『数组C』，拥有 A+B的长度，对于只想合并 AB 数组的需求来说，是一种内存浪费。 所以，有如下骚操作：

```javascript
var A = [1, 2];
var B = [3, 4, 5];

A.push.apply(A, B);
console.log(A);// 输出 1,2,3,4,5
```

**解析：** 调用A数组的 `push` 函数的 `apply` 函数，将上下文环境设为『数组A』，参数列表设为 『数组B』，由于 `push` 方法支持如下调用： `puth(item1, item2, item3)`，所以就将数组参数转换为参数列表，从而实现数组合并。并且，**支持IE5+** 相同的调用原理还可以用于 `unshift` ，合并 AB 数组并将 B 数组的数据放在前边。

##### 骚操作：求数组中的最大最小值

JS中提供了 `min()` 和 `max()`，他们都支持一个特点：『动态参数』，意思是传入 `min(1,2,3)` 会返回1，传入 `max(1,2,3,4,5)` 会返回5。 所以衍生了如下操作：

```javascript
var A = [1, 2, 3];
var B = [1, 2, 3, 4, 5];

var ans1 = min.apply(null, A);
var ans2 = max.apply(null, B);
console.log(ans1, ans2);// 输出 1 5
```

### 总结

`apply` 和 `call` 均可以改变上下文执行环境，`apply` 可以动态传入数组参数， `call` 可以按照正常调用函数的方式执行函数，各有特点， `apply` 配合数组、对象还可以巧妙的玩出各种骚操作，增加效率和乐趣！