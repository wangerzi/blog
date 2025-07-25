# 个人博客

从 Wordpress 迁移至 hexo，最终还是走向了静态博客的道路。

## 历程

Wordpress 已经是五年前在个人服务器上搭建的了，当时用的服务器还是 centos 6，借着迁移服务器的风，把可以静态化的服务统一放到腾讯云 COS 上，配合其 CDN 分发。

博客以 hexo + neoxeo + githubAction + githubPage + netlify + 腾讯云 CDN 搭建

> GitHub README.md 中的 link 默认为 norefer，会触发防盗链导致无法打开博客，如需浏览请右键复制链接后粘贴到浏览器访问

历程整理：

- [测试环境升级整理实录](https://blog.wj2015.com/2022/05/11/%E6%B5%8B%E8%AF%95%E7%8E%AF%E5%A2%83%E5%8D%87%E7%BA%A7%E6%95%B4%E7%90%86%E5%AE%9E%E5%BD%95/)
- [hexo博客迁移](https://blog.wj2015.com/2022/05/08/hexo%E5%8D%9A%E5%AE%A2%E8%BF%81%E7%A7%BB/)

## 博客管理

本地环境：

```shell
npm i
npm run server
```
访问 http://localhost:4000/

新建博客：

```shell
npx hexo new "博客的标题"
```

marktxt 配置：

- 要选择复制文件到相对路径

- 打开 blog 项目的根目录

- 配置相对路径为 sources/static/assets

![](source/static/assets/2023-02-20-23-29-08-image.png)

然后提交图片前记得压缩一下，不然很大

封面可以如如下网站找：

[Unsplash](https://unsplash.com/s/photos/cc0)

[Pexels](https://www.pexels.com/zh-cn/)

## 问题

如果碰到推不上去，报如下异常时

```
% git push
Enumerating objects: 38, done.
Counting objects: 100% (38/38), done.
Delta compression using up to 8 threads
Compressing objects: 100% (32/32), done.
error: RPC failed; HTTP 400 curl 22 The requested URL returned error: 400
send-pack: unexpected disconnect while reading sideband packet
Writing objects: 100% (32/32), 12.92 MiB | 12.31 MiB/s, done.
Total 32 (delta 6), reused 0 (delta 0), pack-reused 0
fatal: the remote end hung up unexpectedly
Everything up-to-date
```

可以尝试如下命令

```bash
git config http.postBuffer 524288000  # 设置为 500 MB
git config http.maxRequestBuffer 100M
```