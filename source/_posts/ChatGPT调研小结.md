---
title: ChatGPT 调研小结
date: 2023-06-02 21:52:07
tags:
  - chatgpt
categories:
- 机器学习
cover: ../static/assets/38a6112a28460a6c5f12cb07b41ff8e659450ae4.jpg
---

## 前言

博主其实并不是专门从事 LLM 方向的程序员，但 23年初 ChatGPT 横空出世的时候着实震惊了我一把，博主也算第一时间注册了账号进行体验；

事实证明这很明智，这两个月发现新账号越来越难注册了，界面也越来越不稳定，于是乎我又化身为一个资(bei)源(po)的付费用户，直到今天。

与 GPT 交流的时间长了，然后就一直很好奇一些问题：

- ChatGPT 是怎么演变而来的？它的内部结构是什么样子？

- ChatGPT 有智慧么？

- 为什么它有时候会胡言乱语？ 

- 为什么我跟它说话，它会一个字一个字往外蹦？而不是一个词一个词的回答我？

本文会从怎么使用聊起，探索它现在有什么有趣的玩法，以及它是怎么形成、怎么工作的；不过碍于篇幅，这篇文章并不会过多的啰嗦，只是一个使用小结，大家可以轻松愉快的看完它。

## 怎么跟它对话

想要玩 ChatGPT，自然就得先访问到他们的服务嘛，那这里就分两条路了：

- 注册 ChatGPT，直接在 web 界面上使用

- 使用国内转出来的“盗版” ChatGPT

第一条路呢，网上写的注册教程铺天盖地，随便找了一个，仅供参考：[ChatGPT注册教程（完整指南） · xiaoming2028/FreePAC Wiki · GitHub](https://github.com/xiaoming2028/FreePAC/wiki/ChatGPT%E6%B3%A8%E5%86%8C%E6%95%99%E7%A8%8B%EF%BC%88%E5%AE%8C%E6%95%B4%E6%8C%87%E5%8D%97%EF%BC%89)

重点就是：

- 解决网络问题

- 搞一个谷歌账号

- 搞一个SIM接码平台

- 官网注册

![](../static/assets/2023-06-02-22-21-11-image.png)

那个 sms activate 里边的号码很多都被玩坏了，尝试新注册过很多次都碰到风控了，但我同事就能一次成功，所以碰到了就多试几次，功夫不负有心人，再者小黄鱼各种平台都有卖账号的，想尝试的都能想到办法。

第二条路，国内外可用的资源就比较多了，基本上就是注册账号就能免费玩上那么一阵子，需要持久使用的话多个站点切换就行，重点是免费

搜集到一些可用的网站，可直接体验：

- https://ministers.ai/home

- https://chat.forefront.ai/

- [Chat &#8211; ChatGPT Online](https://chatgptonline.ai/chat/)

- [Use ChatGPT Online Free - GPT-4 Login](https://gpt4login.com/use-chatgpt-online-free/)

## 一些好玩的东西

### 角色扮演

[GitHub - SillyTavern/SillyTavern: LLM Frontend for Power Users.](https://github.com/SillyTavern/SillyTavern)

[Welcome to Character Book! - Character.AI](https://book.character.ai/character-book/welcome-to-character-book)

![](../static/assets/2023-06-02-22-28-38-image.png)

![](../static/assets/2023-06-02-22-31-09-image.png)

![](../static/assets/2023-06-02-22-31-30-image.png)

通过 SillyTavern 这个工具，可以比较方便的定制自己喜欢的角色、通过 StableDiffusion 扩展描绘场景和角色，甚至让 ChatGPT 扮演文字探险游戏的 GM 来一场刺激的地牢探险

![](../static/assets/2023-06-02-22-32-40-img_v2_97a8c31f-3e69-42bd-8552-45a808fd520h.jpg)

### 写小说

[仅用自然语言，让ChatGPT输出连贯的长篇小说！苏黎世联邦理工大学提出RecurrentGPT](https://mp.weixin.qq.com/s/26S-y9ZhTxi3KHWL6v62Ew)

本质是使用 prompt 引导 GPT 思考情节，并对每次新增的内容不断的做总结（压缩），思考后续场景时把前面次总结 + 总体的大纲带上，模拟人类长期记忆过程。

开源项目地址：[https://github.com/aiwaves-cn/RecurrentGPT](https://github.com/aiwaves-cn/RecurrentGPT)
作者开放体验地址：[交互式小说生成](https://www.aiwaves.org/interactivefiction)

![](../static/assets/2023-06-03-11-12-34-a081c51f5c5ffaece100657672c016e.jpg)

### 自建 GPT 服务

[GitHub - Yidadaa/ChatGPT-Next-Web: One-Click to deploy well-designed ChatGPT web UI on Vercel. 一键拥有你自己的 ChatGPT 网页服务。](https://github.com/Yidadaa/ChatGPT-Next-Web)

[GitHub - GaiZhenbiao/ChuanhuChatGPT: GUI for ChatGPT API and many LLMs](https://github.com/GaiZhenbiao/ChuanhuChatGPT)

这类服务最大的意义是，可以让身边感兴趣但折腾能力不足的人无缝上手，避免去注册账号、搞定网络、找网站；而且最大的一点：稳定。

> 不过这类服务如果自己部署的话，在注册账号免费的 API 额度过期后，需要绑定海外信用卡才可以继续使用下去。这个成本会比较高，操作也比较麻烦，有需要自行检索即可。

我尝试部署过很多个这样的服务，最终留下了 ChatGPT Next，面具和预设真的很方便，分享给别人后，聊天数据只会持久化到浏览器，隐私困扰少。

![](../static/assets/2023-06-02-22-42-18-image.png)

### LangChain

[GitHub - binary-husky/gpt_academic: 为ChatGPT/GLM提供图形交互界面，特别优化论文阅读润色体验，模块化设计支持自定义快捷按钮&函数插件，支持代码块表格显示，Tex公式双显示，支持Python和C++等项目剖析&自译解功能，PDF/LaTex论文翻译&总结功能，支持并行问询多种LLM模型，支持清华chatglm等本地模型。兼容复旦MOSS, llama, rwkv, 盘古, newbing, claude等](https://github.com/binary-husky/gpt_academic)

[GitHub - mayooear/gpt4-pdf-chatbot-langchain: GPT4 & LangChain Chatbot for large PDF docs](https://github.com/mayooear/gpt4-pdf-chatbot-langchain)

[GitHub - jerryjliu/llama_index: LlamaIndex (GPT Index) is a data framework for your LLM applications](https://github.com/jerryjliu/llama_index)

基于 LangChain 可以让 GPT 做更多的事情，比如与一篇 PDF 文档聊天、读论文等等

> PS:  这张图是 gpt_academic 项目截图，底层并非使用 langchain，只有一个理解 pdf 功能与 langchain 的 combine 总结功能类似。

![thumbnailwebp](https://static.flomoapp.com/file/2023-05-28/73871/2b8019867a4c451bb4e60500383a57b9.png?OSSAccessKeyId=LTAI4G9PcaGksWVKCPrE1TVL&Expires=1685889191&Signature=SfYEb7j9z69Bf3RdJmccpAWYtUI%3D)

[GitHub - FlowiseAI/Flowise: Drag & drop UI to build your customized LLM flow using LangchainJS](https://github.com/FlowiseAI/Flowise)

[GitHub - hwchase17/langchain: ⚡ Building applications with LLMs through composability ⚡](https://github.com/hwchase17/langchain/)

更技术一点，就是可以使用工作流的方式把 LangChain 组件组合起来，为自己服务

![](https://github.com/FlowiseAI/Flowise/raw/main/images/flowise.gif?raw=true)

## 一些底层原理

作为技术工作者，自然也需要对其原理做一下简单的了解，以明确其能力边界和特性，文字版教程推荐：[人人都能看懂的 ChatGPT 原理课 - JioNLP - 掘金小册](https://juejin.cn/book/7225592349563289600?suid=1679709499036206&source=android)

![](../static/assets/2023-06-02-22-44-41-img_v2_4970535a-5101-4393-b020-ac914870caah.jpg)

视频版教程推荐两个 B站的教程：

[训练原理科普](https://www.bilibili.com/video/BV1yV4y1k7Tc/?spm_id_from=333.880.my_history.page.click)

[ChatGPT会产生意识么？](https://www.bilibili.com/video/BV1em4y187CU/?spm_id_from=333.880.my_history.page.click&vd_source=d3df27ba79e25e88d0242b745a857c29)

看完相关教程后，博主还整理了一个思维导图分享，供大家参考：

![](../static/assets/f5a5a4f733061a73f9062f1de907a546c1cd5fe9.Jpeg)

还有更清楚的 PDF 版本：[ChatGPT原理-230602193139.pdf](/static/assets/ChatGPT原理-230602193139.pdf)

## 总结

ChatGPT 所代表的大语言模型潜力还远远没有被挖掘出来，上面提到的都只是冰山一角，博主自我感觉连门槛都没有踏入。

博主相信随着 LLM 潜力的进一步挖掘，基于 ChatGPT 也一定会涌现出来无数高效率工具，它会如20年前的互联网一样，一步步成为我们工作生活中密不可分的一部分（参见博主写的另一篇文章：[基于ChatGPT恐怖能力的高效率工具涌现意味着什么 - 九层台](https://blog.wj2015.com/2023/04/05/%E5%9F%BA%E4%BA%8Echatgpt%E6%81%90%E6%80%96%E8%83%BD%E5%8A%9B%E7%9A%84%E9%AB%98%E6%95%88%E7%8E%87%E5%B7%A5%E5%85%B7%E6%B6%8C%E7%8E%B0%E6%84%8F%E5%91%B3%E7%9D%80%E4%BB%80%E4%B9%88/)）。
