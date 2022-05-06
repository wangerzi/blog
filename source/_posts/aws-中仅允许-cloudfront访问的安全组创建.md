---
title: AWS 中仅允许 CloudFront访问的安全组创建
tags:
  - aws
  - cloudfront
  - lambda
id: '636'
categories:
  - - 后端开发
    - Serverless
date: 2020-11-22 13:21:01
cover: /static/uploads/2020/11/timg.jpg
---



## 目的

生产环境中出于安全和并发的考虑，一般会用 CDN 和 LB 为后端服务/资源进行加速和并发防护，在 AWS 中，CDN 服务是 CloudFront，为了从根源上避免直接访问到源，需要通过安全组的方式组织非 CloudFront 边缘节点的访问。

## 思路

根据官方博客 [How to Automatically Update Your Security Groups for Amazon CloudFront and AWS WAF by Using AWS Lambda](https://aws.amazon.com/cn/blogs/security/how-to-automatically-update-your-security-groups-for-amazon-cloudfront-and-aws-waf-by-using-aws-lambda/)中所属步骤，创建好 Lambda 函数和相关安全组，并监听 IP 更新的 SNS 即可。 对应代码仓库：[aws-cloudfront-samples](https://github.com/awslabs/aws-cloudfront-samples) 原理就是每次在AWS更新 CloudFront 的IP时都会触发 Lambda 函数的执行，Lambda 里边获取 [IP ranges](https://ip-ranges.amazonaws.com/ip-ranges.json)，找出里边 CloudFront 的 Global IP 和 Region IP，保存到对应的安全组去。

> 注意：如果服务器在多个区或同一个区不同的 VPC，需要每个区都执行安全组创建+lambda，同区的不同 VPC 只要创建好安全组即可。 注意：如果运行过程中提示你每个安全组规则数量不可超过 50 导致无法更新成功，需要在 AWS 控制台申请扩大单个安全组规则最大限制，100~200即可。

## 操作过程

### 创建安全组

[安全组控制台](https://us-west-2.console.aws.amazon.com/ec2/v2/home?region=us-west-2#SecurityGroups:) 首先把需要CloudFront IP 访问控制的安全组创建起来，由于 Global 和 Region 是分两个组存的，HTTP 和 HTTPS 也是分两个组存的，所以我们需要创建四个安全组，名称和标签如下表所示。

安全组名称

描述

协议

类型

Tag:Name

Tag:Protocol

Tag:AutoUpdate

HTTP\_CLOUDFRONT\_IP\_ONLY

HTTP for cloudfront global ip.

HTTP

Global

cloudfront\_g

http

true

HTTP\_CLOUDFRONT\_IP\_ONLY\_R

HTTP for cloudfront region ip.

HTTP

Region

cloudfront\_r

http

true

HTTPS\_CLOUDFRONT\_IP\_ONLY

HTTPS for cloudfront global ip.

HTTPS

Global

cloudfront\_g

https

true

HTTPS\_CLOUDFRONT\_IP\_ONLY\_R

HTTPS for cloudfront region ip.

HTTPS

Region

cloudfront\_r

https

true

> 注：所有的 EC2 和 LB 都需要选择 HTTP\_CLOUDFRONT\_IP\_ONLY 和 HTTP\_CLOUDFRONT\_IP\_ONLY\_R 这两个安全组。

[![](/static/uploads/2020/11/wp_editor_md_7092ef83875ab140cb81c7e1d9fed2d8.jpg)](/static/uploads/2020/11/wp_editor_md_7092ef83875ab140cb81c7e1d9fed2d8.jpg) 创建完大概就是这样： [![](/static/uploads/2020/11/wp_editor_md_55e62d5672bde84b42fba64c5b992b0c.jpg)](/static/uploads/2020/11/wp_editor_md_55e62d5672bde84b42fba64c5b992b0c.jpg)

### 创建Lambda执行角色

[IAM角色](https://console.aws.amazon.com/iam/home?region=us-west-1#/roles) 首先进入 IAM 角色控制台，点击『创建角色』，选择『Lambda』，点击『下一步』 [![](/static/uploads/2020/11/wp_editor_md_f1d3d67dca44e87a06e1e5311d580ccb.jpg)](/static/uploads/2020/11/wp_editor_md_f1d3d67dca44e87a06e1e5311d580ccb.jpg) 再选择『创建策略』，点击『JSON』标签，将如下内容粘贴进去

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeSecurityGroups",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:RevokeSecurityGroupIngress"
      ],
      "Resource": "*"
    }
  ]
}
```

[![](/static/uploads/2020/11/wp_editor_md_98594b8a74ef2f31faaf96482ea86672.jpg)](/static/uploads/2020/11/wp_editor_md_98594b8a74ef2f31faaf96482ea86672.jpg) 命好名字，创建策略即可。 [![](/static/uploads/2020/11/wp_editor_md_005357534d0d74fcb6e2386836423f97.jpg)](/static/uploads/2020/11/wp_editor_md_005357534d0d74fcb6e2386836423f97.jpg) 随后返回继续创建角色，搜索到对应的策略后勾选并点击『下一步：标签』 [![](/static/uploads/2020/11/wp_editor_md_3db02e43119ec89211db531ff29314c4.jpg)](/static/uploads/2020/11/wp_editor_md_3db02e43119ec89211db531ff29314c4.jpg) 直到最后，点击创建角色即可，**这个角色会用来执行 Lambda 函数**。 [![](/static/uploads/2020/11/wp_editor_md_3db02e43119ec89211db531ff29314c4.jpg)](/static/uploads/2020/11/wp_editor_md_3db02e43119ec89211db531ff29314c4.jpg)

### 创建函数

[Lambda函数控制台](https://us-west-2.console.aws.amazon.com/lambda/home?region=us-west-2#/functions) 点击『创建函数』，定义好名称选择好角色，点击『创建函数』 [![](/static/uploads/2020/11/wp_editor_md_043d0764d073f967d9a6b1244a3ef411.jpg)](/static/uploads/2020/11/wp_editor_md_043d0764d073f967d9a6b1244a3ef411.jpg) 将参考了 [官方代码](https://github.com/aws-samples/aws-cloudfront-samples/blob/master/update_security_groups_lambda/update_security_groups.py) 的改进版代码拷贝到『函数代码』区域，在右上角点击『保存』

#### 代码改进

##### 第10行

```python
INGRESS_PORTS = { 'http' : 80, 'https': 443, 'example': 8080}
```

改为

```python
INGRESS_PORTS = { 'http' : 80, 'https': 443}
```

[![](/static/uploads/2020/11/wp_editor_md_b81cb7f4b41fda58c33c74aa7a1f8d2c.jpg)](/static/uploads/2020/11/wp_editor_md_b81cb7f4b41fda58c33c74aa7a1f8d2c.jpg) 修改函数的执行时间设置为1分钟，否则容易超时 [![](/static/uploads/2020/11/wp_editor_md_2b9bec1d31d33d7d265e7d560cf9670c.jpg)](/static/uploads/2020/11/wp_editor_md_2b9bec1d31d33d7d265e7d560cf9670c.jpg) 可以点击运行测试用例尝试下是否会更新安全组，如果出现 md5 错误，手动修改下测试数据即可

### 自动更新

点击『添加触发器』，选择SNS，主题输入框输入 `arn:aws:sns:us-east-1:806199016981:AmazonIpSpaceChanged`，点击『添加』 [![](/static/uploads/2020/11/wp_editor_md_d3e3ab1240d6a3bf5e5ff0d8634402c0.jpg)](/static/uploads/2020/11/wp_editor_md_d3e3ab1240d6a3bf5e5ff0d8634402c0.jpg)

### 更新 EC2 或 ELB 的安全组策略

进入到 EC2 控制台或者 ELB 的控制台，根据需要更新安全组为刚才生成的 `HTTP_CLOUDFRONT_IP_ONLY`, `HTTP_CLOUDFRONT_IP_ONLY_R` 安全组即可。

## 总结

整体思路就是通过一个 lambda 函数 **监听 AWS 官方通知队列** ，这个队列一旦有消息就会执行函数体，函数内部 **拉去最新的 cloudfront ip 列表** ，然后 **自动生成规则** 更新带有对应 TAG 的安全组