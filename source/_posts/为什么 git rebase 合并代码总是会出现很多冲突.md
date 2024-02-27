---
title: 为什么 git rebase 合并代码总是很多冲突
tags:
  - git
  - 多人协作
categories:
  - - Linux
  - - 后端开发
date: 2024-02-27 15:09:47
cover: ../static/uploads/2019/07/git%E5%B7%A5%E4%BD%9C%E6%B5%81.png
---

# 为什么 git rebase 合并代码总是很多冲突

## 前言

与 git 打了几年交道，合并代码是多人协同的时候最常见的操作，但合并代码有两种方式： MERGE 和 REBASE，多数人都知道 rebase 可以让分支提交变得更新性，使用 merge 合并分支总是会有很多分叉。但同时 rebase 也会让我们解决非常多的提交，不禁思考为什么会发生这样的现象呢？

## merge 的本质

**git merge** 是将两个分支的改动合并到一起。如果两个分支都有新的提交，Git 会创建一个新的提交，这个提交有两个父提交。这种方式保留了完整的提交历史和分支结构，但可能会导致历史记录复杂。

```bash
A---B---C feature
     /         \
D---E---F---G master
```

在上图中，如果我们在 `master`​ 分支上执行 `git merge feature`​，会得到下面的结果：

```bash
A---B---C feature
     /         \
D---E---F---G---H master
```

​`H`​ 是新的合并提交，它的父提交是 `G`​ 和 `C`​。

## rebase 的本质

**git rebase** 是将一系列提交应用到另一个基础之上。如果两个分支都有新的提交，Git 会逐个取出分支上的提交，然后在目标分支上重新应用。这种方式会创建一条线性的提交历史，使得历史记录更清晰，但是会丢失一些历史信息。

```bash
A---B---C feature
     /
D---E---F---G master
```

在上图中，如果我们在 `feature`​ 分支上执行 `git rebase master`​，会得到下面的结果：

```bash
             A'--B'--C' feature
             /
D---E---F---G master
```

​`A'`​, `B'`​, `C'`​ 是原来 `feature`​ 分支上的提交 `A`​, `B`​, `C`​ 在 `master`​ 分支上重新应用的结果。

至于为什么 `rebase`​ 时经常需要手动处理冲突，这是因为 `rebase`​ 是逐个应用提交，所以如果有多个提交修改了同一部分代码，就需要多次解决冲突。而 `merge`​ 只需要解决一次冲突，因为它是一次性合并所有改动。

## 总结

本文比较浅显，只是直接了当的解释了一个常见困惑，本质上还是 commit 合并的原理不同。

## 参考资料

详细的 git rebase 操作指南：https://www.codercto.com/a/45325.html

git merge 操作指南：https://blog.csdn.net/u010665216/article/details/129885666
