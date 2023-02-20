# 个人博客

从 Wordpress 迁移至 hexo，最终还是走向了静态博客的道路，地址： [blog.wj2015.com](https://blog.wj2015.com/)

## 历程

Wordpress 已经是五年前在个人服务器上搭建的了，当时用的服务器还是 centos 6，借着迁移服务器的风，把可以静态化的服务统一放到腾讯云 COS 上，配合其 CDN 分发。

博客以 hexo + neoxeo + githubAction + githubPage + netlify + 腾讯云 CDN 搭建

## 博客管理

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


