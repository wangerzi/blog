---
title: Python的类：抽象类、抽象方法及其继承
tags: []
id: '629'
categories:
  - - 后端开发
---



## 目的

从Python类的基础讲起，整理 Python 类中的抽象类和抽象方法知识点

## 正文

### 为什么需要类，如何定义类？

了解过 C++、Java 的同学很容易理解类这个概念，在一个类中包含了类的属性和方法，可以**把业务中的各个实体通过属性和方法表现出来**。 比如我们需要开发一个学生管理系统，你可以借用 python 的模块用面向过程的方式一把梭，也可以将学生这个实体抽象出来一个对象，这个对象包含**很多属性**：学号、入学时间、生日、宿舍号、所在班级等，还包含**很多方法**，比如查看成绩、选课等等。 如上所述的学生类，在Python中可以表现为如下的代码

```python
class Student:
    def __init__(self, stu_no, school_time, birthday, room_no, class_no):
        self.stu_no = stu_no
        self.school_time = school_time
        self.birthday = birthday
        self.room_no = room_no
        self.class_no = class_no
    def show_grade(self):
        print("grade is xxx")
    def choose_lesson(self, lesson_no):
        print("choose lesson", lesson_no)
```

### 类的公/私有方法和属性

### 类的单继承、多继承

#### 多继承中类属性 or 方法查找顺序

#### 继承的弱点

### 面向接口编程，抽象类和抽象方法

## 总结