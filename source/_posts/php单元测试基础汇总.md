---
title: phpunit单元测试基础汇总
tags:
  - php
  - phpunit
  - 单元测试
id: '379'
categories:
  - - 后端开发
date: 2019-07-31 15:17:03
---



## 什么是单元测试

> 单元测试（unit testing），是指对软件中的最小可测试单元进行检查和验证。单元就是人为规定的最小的被测功能模块。单元测试是在软件开发过程中要进行的最低级别的测试活动，软件的独立单元将在与程序的其他部分相隔离的情况下进行测试。

在php里边，最小单元可以指一个函数、或者类，需要验证的就是每个函数，每个类的功能与我们预想的一致。

## 单元测试有什么意义

*   可以**减少一些细节错误的发生**，比如应该报错的情况没有报错，入参、结果是否与需求对应上等。
*   **便于日后修改维护**，实际工作中存在不少情况是做出了一版功能，但是上线后需要对里边的细节进行调整，有单元测试的话改起来会更加放心，并且完善单元测试的过程也是进一步理解需求的过程。
*   更容易**发现平时无法走到的异常分支**，而这个分支的处理逻辑可能人工测试需要经历很多步骤才能走到，省时间

最近在工作中也尝试着为开发中的功能写单元测试，切实意识到了单元测试的好处，需求里边有一个比较复杂的时间推算逻辑，最开始自认为各种情况考虑周全然后劈里啪啦写完，不过运行了事先写好的单元测试之后，依旧发现了几个隐藏比较深问题（**再自信也得过一遍测试啊**）。 修复问题后提测的过程中遇到了需求变更，不少关键代码需要改动，正常这种情况自测的话会很费劲，因为需要数据库找各种各样情况的数据去跑接口，然后数据对不上改完还得重新跑接口自测。但是这次先把单元测试规定正确后，放心大胆的按照自己的想法改造代码，经历了 改代码 > 跑测试 > 改代码 > 跑测试的循环后，快速交付了需求。

## 单元测试的一些概念

之前也接触过php、python、JS之类的语言，对这些语言的单元测试也有一定了解，下边先看一下单元测试中通用的一些概念。

### 断言

想要更加细致的了解断言的话，这里推荐一篇博客：[https://www.jianshu.com/p/9b8c88deed6a](https://www.jianshu.com/p/9b8c88deed6a) 在软件测试特别是在单元测试时,必用的一个功能就是“断言”（Assert)，顾名思义，编写程序时，常会做出一定的假设，那断言就是用来捕获假设的异常。 下边举个栗子： 一个简单的函数 add() 拥有两个参数，功能是返回两个参数的和，当我需要验证这个函数的正确性的时候就需要模拟两个入参并 **判断函数的返回值是否为两个入参之和**，判断返回值是否准确这个过程即为断言。

```php
function add($a, $b)
{
    return $a + $b;
}
```

### 基境

每一个单元测试方法都是一个独立的个体，每次单元测试完毕，需要将数据恢复到正确的状态中，不至于被其他测试方法给影响。 在phpunit中，给出的 TestCase 基类即有两个方法，`setUp` 和 `setDown` 分别用于为每个单元测试创建测试对象和清理测试对象

### 数据供给器

对同一类情况进行测试，通常可以用数据供给器传入不同入参和相应的预期返回值。 测试方法可以接受任意参数。这些参数由数据供给器方法提供。在phpunit中使用 @dataProvider 标注来指定使用哪个数据供给器方法。

## php如何集成单元测试

PHP的单元测试依赖一个测试框架：phpunit（官方文档：[https://phpunit.readthedocs.io/zh\_CN/latest/index.html](https://phpunit.readthedocs.io/zh_CN/latest/index.html)）

### 如何安装

可以通过phar的方式安装

```shell
$  wget https://phar.phpunit.de/phpunit-7.0.phar
$  chmod +x phpunit-7.0.phar
$  sudo mv phpunit-7.0.phar /usr/local/bin/phpunit
$  phpunit --version
```

也可以通过 composer 进行统一管理

```shell
$ composer require phpunit/phpunit
```

在 `composer.json` 中会出现如下依赖

```json
{
    "require": {
        "phpunit/phpunit": "^7.5"
    }
}
```

并且会出现 `vendor/bin/phpunit` 文件，直接运行即可

### 如何编写单元测试

所有类需要继承 `PHPUnit\Framework\TestCase`，`setUp` 函数用于初始化测试对象，`setDown` 函数用于清理测试对象，更多规范 更具体写法可以查看底部的 `举个栗子`

### phpunit常用断言方法

更多断言方法详见 phpunit 官方文档，基本都能顾名思义。

断言函数

作用

示例

assertEquals(\\$except, \\$value)

断言相等

$this->assertEquals(2, 1 + 1)

assertEmpty($value)

断言为空

$this->assertEmpty(\[\])

assertNotEmpty($value)

断言不为空

$this->assertNotEmpty(\[1, 2, 3\])

assertTrue($value)

断言为真

$this->assertTrue(1 === 1)

assertFalse($value)

断言为假

$this->assertFalse(1 === '1')

expectException(\\Exception $e)

断言本次测试会出现特定异常

$this->expectException(\\Exception::class); throw new \\Exception('测试', 1002);

expectExceptionCode($code)

断言异常状态码

$this->expectExceptionCode(1002); throw new \\Exception('测试', 1002);

expectExceptionMessage($msg)

断言异常信息

$this->expectExceptionMessage('测试'); throw new \\Exception('测试', 1002);

expectOutputString($msg)

断言输出

$this->expectOutputString('Hello');echo "Hello";

getActualOutput()

获取实际输出

echo "Hello";\\$result = \\$this->getActualOutput();

### 如何运行单元测试

```bash
# 运行全部测试
phpunit
# 运行某个分组的单元测试
phpunit --group GroupA
# 运行指定测试类的所有测试用例
phpunit tests/xxxxTest.php
# 运行所有测试类中满足filter条件的方法
phpunit --filter xxxFunc
# 运行某个测试类中满足filter条件的
```

### phpunit.xml 是什么

phpunit.xml 是一个XML格式的配置文件，能够配置单元测试中的一些默认行为，比如环境变量、启动文件、日志记录等，官方文档如下[https://phpunit.readthedocs.io/zh\_CN/latest/configuration.html](https://phpunit.readthedocs.io/zh_CN/latest/configuration.html) 一个样例配置如下所示：

```markup
<?xml version="1.0" encoding="UTF-8"?>
<!--phpunit标签是配置中的核心，这里配置了启动文件 "./tests/bootstrap.php"-->
<phpunit backupGlobals="false"
         backupStaticAttributes="false"
         bootstrap="./tests/bootstrap.php"
         colors="true"
         convertErrorsToExceptions="true"
         convertNoticesToExceptions="true"
         convertWarningsToExceptions="true"
         processIsolation="false"
         stopOnFailure="false">
    <!--测试套件：非常多的测试用例放在一起即可成为测试套件，执行时会扫描包含的所有 *Test.php文件-->
    <testsuites>
        <testsuite name="Unit">
            <directory suffix="Test.php">./tests/Unit</directory>
        </testsuite>
    </testsuites>
    <filter>
        <!--这里配置了白名单，只有这里边的代码会被统计覆盖率-->
        <whitelist processUncoveredFilesFromWhitelist="true">
            <directory suffix=".php">./app/library</directory>
            <directory suffix=".php">./app/models</directory>
        </whitelist>
    </filter>
    <!--这里配置了PHP的环境变量-->
    <php>
        <server name="APP_ENV" value="product"/>
        <server name="BCRYPT_ROUNDS" value="4"/>
        <server name="CACHE_DRIVER" value="array"/>
        <server name="MAIL_DRIVER" value="array"/>
        <server name="QUEUE_CONNECTION" value="sync"/>
        <server name="SESSION_DRIVER" value="array"/>
    </php>
    <!--这里是日志记录，把覆盖率信息保存到 ./tests/codeCoverage-->
    <logging>
        <log type="coverage-html" target="./tests/codeCoverage"/>
    </logging>
</phpunit>

```

### 如何查看代码覆盖率

执行 phpunit 之后，根据 `<logging>` 中的配置，会自动生成代码覆盖率信息至 `./tests/codeCoverage/`，打开其中 `index.html` 即可查看覆盖率信息。

## 举个栗子

以一个简单的斐波拉契数列计算函数为例

> 斐波那契数列由0和1开始，之后的斐波那契系数就是由之前的两数相加而得出。

#### 输入输出分析

根据函数特点，我们可以通过验证已知情况和特殊情况的方式去验证，经过分析结果如下

#### 正常输入的已知情况：

入参

预期返回

描述

0

0

规则

1

1

规则

2

1

0 + 1 = 1

3

2

1 + 1 = 2

4

3

1 + 2 = 3

5

5

2 + 3 = 5

6

8

3 + 5 = 8

...

...

...

12

144

55 + 89 = 144

#### 异常输入的情况处理

处理为0，或者抛出异常均可

入参

预期返回

描述

\-1

0

非正常输入处理为0

''

0

非正常输入处理为0

1.1

0

非正常输入处理为0

'文字'

0

非正常输入处理为0

#### 编写测试类

`tests/FunctionTest.php`

```php
use PHPUnit\Framework\TestCase;
class FunctionsTest extends TestCase
{
    /**
     * @dataProvider fibon_normal_provider
     * @param $input
     * @param $except
     */
    public function test_fibon_normal($input, $except)
    {
        $this->assertEquals($except, fibon($input));
    }

    public function fibon_normal_provider()
    {
        return [
            [0, 0],
            [1, 1],
            [2, 1],
            [3, 2],
            [4, 3],
            [5, 5],
            [6, 8],
            [12, 144],
        ];
    }

    /**
     * @dataProvider fibon_error_provider
     * @param $input
     * @param $except
     */
    public function test_fibon_error($input, $except)
    {
        $this->assertEquals($except, fibon($input));
    }

    public function fibon_error_provider()
    {
        return [
            [-1, 0],
            [1.1, 0],
            ['', 0],
            ['文字', 0],
        ];
    }
}
```

#### 函数功能实现

（PS:此法效率很差，约莫是$O(n^2)$的复杂度，仅用于此处演示） `functions.php`

```php
function fibon($a)
{
    if (!is_int($a)) {
        return 0;
    }
    if ($a <= 0) {
        return 0;
    } elseif ($a == 1) {
        return 1;
    } else {
        return fibon($a - 1) + fibon($a - 2);
    }
}
```

#### 运行结果

```shell
vagrant@homestead:~/code/bmtrip/platoReceivable$ phpunit tests/FunctionsTest.php --filter test_fibon
PHPUnit 7.5.14 by Sebastian Bergmann and contributors.

............                                                      12 / 12 (100%)

Time: 5.77 seconds, Memory: 26.00 MB

OK (12 tests, 12 assertions)

Generating code coverage report in HTML format ... done
```