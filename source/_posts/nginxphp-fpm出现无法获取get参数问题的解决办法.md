---
title: nginx+php-fpm出现无法获取GET参数问题的解决办法
tags:
  - nginx
  - php-fpm
  - 参数获取
id: '121'
categories:
  - - Linux
date: 2018-03-07 10:28:58
cover: ../static/uploads/2018/03/TIM%E6%88%AA%E5%9B%BE20180307102409.png
---

# 背景

今天在配置nginx服务器的过程中，发现php能正常访问，但是$\_GET始终为空，就此开始了解决问题之旅。

# 方法

网上的方法大多是：

> 将/etc/nginx/conf/conf.d/virtual.conf中对应虚拟机的try\_files改变为： try\_files $uri $uri/ /index.php?$query\_string;

原：

```null
location / {
            try_files $uri $uri/ /index.php?s=$args;
            # try_files $uri $uri/;
    }
```

目标：

```null
location / {
            try_files $uri $uri/ /index.php?$query_string;
            # try_files $uri $uri/;
    }
```

但是，亲测无效

> 最后，通过与功能完善的nginx配置文件进行对比，发现我在/etc/nginx/fastcgi\_params中配置顺序有问题。 经测试，SCRIPT\_FILENAME参数只能在最前边，放最后就会导致$\_GET为空的情况。

正确配置：

```null
fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
fastcgi_param  QUERY_STRING       $query_string;
fastcgi_param  REQUEST_METHOD     $request_method;
fastcgi_param  CONTENT_TYPE       $content_type;
fastcgi_param  CONTENT_LENGTH     $content_length;

fastcgi_param  SCRIPT_NAME        $fastcgi_script_name;
fastcgi_param  REQUEST_URI        $request_uri;
fastcgi_param  DOCUMENT_URI       $document_uri;
fastcgi_param  DOCUMENT_ROOT      $document_root;
fastcgi_param  SERVER_PROTOCOL    $server_protocol;
fastcgi_param  REQUEST_SCHEME     $scheme;
fastcgi_param  HTTPS              $https if_not_empty;

fastcgi_param  GATEWAY_INTERFACE  CGI/1.1;
fastcgi_param  SERVER_SOFTWARE    nginx/$nginx_version;

fastcgi_param  REMOTE_ADDR        $remote_addr;
fastcgi_param  REMOTE_PORT        $remote_port;
fastcgi_param  SERVER_ADDR        $server_addr;
fastcgi_param  SERVER_PORT        $server_port;
fastcgi_param  SERVER_NAME        $server_name;

# PHP only, required if PHP was built with --enable-force-cgi-redirect
fastcgi_param  REDIRECT_STATUS    200;
# fastcgi_param PATH_INFO                $fastcgi_script_name;
```

# 总结

之所以出现这个问题就是因为当初配置php-fpm参数的时候，听从网上的答案，在fastcgi\_params追加了两行配置，结果浪费我1个多小时的时间，谨记教训！！！！