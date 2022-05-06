---
title: JavaScript基本类型总结
tags:
  - JavaScript
  - 数据类型
id: '410'
categories:
  - - 前端开发
date: 2019-09-03 00:14:35
cover: ../../static/uploads/2019/09/20161107111502_84-1200x661.jpg
---



## 背景

最近仔细看了看Javascript的变量相关章节，包括网道阮一峰大神写的 [JavaScript教程](https://wangdoc.com/javascript)，以及被誉为JS圣经的《JavaScript权威指南》，故对变量的基础类型做一个总结

### 变量有什么类型？

JS中的数据类型总体分为两类：原始类型和对象类型

#### 原始类型

包括数字、字符串、布尔值、null、undefined

#### 对象类型

除了原始类型，其他都是对象类型，Object（普通对象）、Function（函数）、Array（数组）

> PS: 虽然字符和数字也能通过对象的方式调用方法，但是这只是JS自动进行的包装对象，比如：

```javascript
'1,2'.split();// '1,2' 会被转换为 String('1,2')
var a = 3.1415926;
var b = a.toFixed(2)// a 会被转换为 Number(a)
var c = new Number(111); c === 111// false 包装对象与原属性不全等
Number(111) === 111 // true, Number(111) 是显示转换的意思，结果是 111，所以全等

var str = 'test';
str.number = 10;// 包装String对象
console.log(str.number);// undefined，注意这个地方，包装对象的生命周期只在那一个语句，下一个语句就没了

var arr = [];
arr.number = 123;
console.log(arr.number);// 123，这里是因为 [] 本来就是对象，与包装对象不同
```

包装对象的详解可以参考如下博客：[https://blog.csdn.net/lhjuejiang/article/details/79623505](https://blog.csdn.net/lhjuejiang/article/details/79623505)

### 类型转换

由于JS是一个弱类型的语言，所以在不同数据类型的数据进行比较或者操作的时候，需要进行类型转换，类型转换分强制转换（也称显式转换）以及自动转换（也称隐式转换）

#### 强制转换

使用 Number() String() Boolean() 等构造函数进行强制转换，也相当于显示新增了包装对象，相比 parseInt()，parseFloat() 等辅助函数更加严格

```javascript
Number('1') // 1
Number(true) // 1
Number(false) // 0
Number('1.2e8') // 120000000 PS:科学计数法
Number('0b1001')// 9 PS:二进制
Number('0xF') // 15 PS:十六进制
Number('0o777') // 511 PS:八进制
Number('1x') // NaN

// valueOf 用于返回对象自身，默认 return this，如果返回对象则调用原对象的 toString() 方法，如果 toString() 依旧返回对象，则报错无法转换
Number({
    valueOf: function() {
        return 123;
    }
}) // 123
Number({
    valueOf: function() {
        return {}
    },
    toString: function() {
        return {}
    }
})// 报错 Uncaught TypeError: Cannot convert object to primitive value

String(true) // 'true'
String(1.) // '1'
String(1.1) // '1.1'
String(-0) // '0'
String(-Infinity) // -Infinity
String(Infinity) // Infinity

Boolean('') // false
Boolean(0) // false
Boolean(null) // false
Boolean(NaN) // false
Boolean(undefined) // false
Boolean([]) // true PS:除了上边的，其他转出来都是 true
Boolean({}) // true
```

#### 自动转换

JS在不同类型数据互相运算时，非布尔值求布尔值，非数值类型使用 + - \* / ，会自动根据场景调用相应的强制转换函数，既能转为字符又能转为数字时，优先数字

##### 自动转Boolean

在需要用到布尔判断时，会导致此转换，比如 (if/else、三目运算符之类的)

```javascript
if (
    '' == false &&
    null == false &&
    undefined == false &&
    0 == false &&
    NaN == false
){
    console.log(111);// true
}
[] ? 111 : 000 // 111
```

##### 自动转数字

涉及 + - \* / 等操作会转换为数字

```javascript
1 + '2' // 3
1 + 's' // NaN
true + 1 // 2
'c' - 1 // NaN
```

##### 自动转字符

字符 + 非字符 会被转换为字符串

```javascript
'abc' + 1 // abc1
'1' + 0 // 10
'5' + undefined // 5undefined
'5' + function(){} // 5function(){}
```

### 如何检查变量类型

JavaScript是一个弱类型语言，我们可以用内置的三种方法判断一个值到底属于什么类型

*   typeof 运算符
*   instanceof 运算符
*   Object.prototype.toString 方法

#### typeof判断变量类型的各种情况

```javascript
typeof 123 // number
typeof '123' // string
typeof false // boolean

f = function() {}
typeof f // function
typeof undefined // undefined
typeof xxxxx // 未定义变量不报错，返回 undefined

typeof {} // object
typeof window // object
typeof [] // object，数组的typeof是object
typeof null // object，null的typeof也是object
```

需要注意的点：数组、null 调用 typeof 都是 object。因为在JS中，数组本质上是一种特殊的对象，而null是因为最开始的JS把null作为对象的一种特殊值，为了向前兼容，就一直 typeof null === 'object' 了

#### instanceof 通过原型链判断数据类型

instanceof 返回一个 Boolean，表示对象的原型链中是否存在某个构造函数（原始定义：对象是否为某个构造函数的实例）。

> 关于原型链，这里简单提及下，一个对象的 prototype 都指向其构造函数的 prototype，其中有一个方法 Object.getPrototypeOf 返回对象的原型，最终的最终都会指向 Object，而 Object.prototype = null，整个原型链到了尽头，原型链的更为详细的介绍可点击：[参考博客](https://wangdoc.com/javascript/oop/prototype.html#%E5%8E%9F%E5%9E%8B%E9%93%BE)

一些实例如下：

```javascript
[] instanceof Array // true 因为原型链里有 Array 这个构造函数
[] instanceof Object // true 因为原型链里有 Object 这个构造函数
Array instanceof Object // true，因为Array是一个特殊的对象
[] instanceof null // 报错，Right-hand side of 'instanceof' is not an object，因为null不是个对象

var func = function() {}
func instanceof Function // true，原型链中有Function
func instanceof Object // true，Function的原型是 Object
Function instanceof Object // true，函数是一个特殊对象
```

#### Object.prototype.toString()

toString() 用于返回一个对象的字符串形式，用于在自动类型转换时得到想要的字符形式。一些样例具体如下：

```javascript
var obj = {};
obj.toString() // [object Object]
var obj2 = {name: 'xxx'};
obj2.toString = function() {
    return 'Hello '+this.name
}
obj2.toString() // Hello xxx
obj2 == 'Hello xxx' // true

var obj3 = []
obj3.toString() // '' 空字符串，因为数组是空的
var obj4 = [1, 2, 3]
obj4.toString() // '1,2,3'
```

由于对象的 toString 可以被重载，如果要用于类型判断的话，需要用到 `Object.prototype.toString.apply(obj)` 的方式，显示格式 \[object 构造函数\]，比如 `[object Object]`, `[object Number]`, `[object String]`,`[objct Function]`，不过有个问题是无法精确判断 字符对象 和 字符串，因为返回都是 `[object String]`

> `apply` 是 `Function` 对象的一个方法，可以更改函数执行上下文，意思是执行 `Object.prototype.toString` 函数，但是执行过程中的 `this` 指向外部注入的 `obj`，如果还是看不明白可以参考博主之前写的另外一篇博客 [JS中Function.apply() 的骚操作](https://blog.wj2015.com/2018/11/20/js%E4%B8%ADfunction-apply-%E7%9A%84%E9%AA%9A%E6%93%8D%E4%BD%9C/)

```javascript
Object.prototype.toString.apply({}) // [object Object]
Object.prototype.toString.apply([]) // [object Array]
Object.prototype.toString.apply('') // [object String]
Object.prototype.toString.apply(String('')) // [object String]
Object.prototype.toString.apply(123) // [object Number]
Object.prototype.toString.apply(NaN) // [object Number]

Object.prototype.toString.apply(String) //[object Function]
Object.prototype.toString.apply(function(){return 111}) //[object Function]

Object.prototype.toString.apply(undefined)// [object Undefined]
Object.prototype.toString.apply(null) // [object Null]
Object.prototype.toString.apply(true) // [object Boolean]
```

## 总结

变量类型是JS的基础之一，也是平时很容易忽视的地方，稍不注意就容易产生莫名其妙的BUG，特此总结