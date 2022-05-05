---
title: Thinkphp中关于容器和门面的应用
tags: []
id: '465'
categories:
  - - 后端开发
---

## 前言

众·所·周·知·，Think PHP5的设计与 Laravel 十分相似，内部也使用 容器、interface、trait 来进行代码复用和拆分，在代码结构规划上也可以使用这些特性给代码解耦，便于日后扩展

## 一些必须知道的概念复习

### 容器是什么

简单来说就是类的实例化管理工具，比如我想得到一个类的实例，最简单的方式是使用 `new` 语法，比如如下的代码

```php
class Demo1 {
    protected $user = null;
}
$a = new Demo1();
```

如果我下一次还需要使用这个类的实例，可能需要重复实例化或者全局变量、参数传递等方式，实现起来并不优雅，影响代码可读性，比如：

```php
class UserModel extends \Think\Model {
    protected $table = 'user';
}

class UserLogic {
    public function login($name, $password)
    {
        $model = new UserModel();
        $user = $model->where([['name' => $name, 'password' => $password]])->get();
        if (empty($user)) {
            return false;
        } else {
            return true;
        }
    }
}

class UserController {
    public function login()
    {
        $logic = new UserLogic();
        if ($logic->login(request()->get('name'), request()->get('passwrod'))) {
            return '登陆成功';
        } else {
            return '登陆失败';
        }
    }
    public function register()
    {
        $logic = new UserLogic(); // 重复的实例化
        $logic->register(request()->only(['name', 'password']));
    }
}
```

可以看到所有地方类的实例化都是通过 `new` 关键字实例化，这种行为称为 `直接耦合`，意思是上层类直接依赖底层的类，底层类变化会直接影响到上层类的变化，如果底层类需要 **被替换** ，则需要修改所有 `new` 关键字实例化的代码； 而容器就能配合实现解耦，不需要使用 `new` 关键字，替换底层类也不需要四处找依赖改代码。 下面可以看一下通过容器实现相同逻辑的代码： `php`

### 依赖注入是什么

依赖注入是