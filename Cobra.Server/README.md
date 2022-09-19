## 可以直接查看文档
https://pv4fr4bbm0.feishu.cn/docx/doxcnM2Y69KvBKS8hf4msfDkpgd


## 前言
本文主要是对Go语言具体应用做一个，前面已经读了关于Go的基础知识书籍， 希望通过具体开源项目学习下别人写法，进而对Go有更深刻了解。
本文主要学习的是：https://github.com/aceld/zinx  
官网地址: http://zinx.me  


简单介绍：Zinx主要用于客户端 和 服务端通过TCP建立链接，然后对客户端触达的消息进行处理。zinx已经在很多企业进行开发使用，具有一定的价值。

**为什么选择Zinx源码看？**
1. 从star量/近期代码提交量看不算太低
2. 对TCP通信部分有过较多接触
3. 当然最重要的是项目可以run起来

**能学到什么？**
1. 可以对看过的知识进行巩固。遇到不懂的用法，反过来在看基础知识，相互加深。
2. 更重要的是了解源码作者设计的思想。源码作者是怎么设计并发，怎么设计类，怎么架构的等等。

**为什么会新开一个分支？**

https://github.com/CodingStone/learn_zinx/tree/main/Cobra.Server
我觉得代码是看不会的，需要多写。  
因此即使是大部分文件结构相同，也会熬着手敲一遍，多看几遍，每次都会逐行修复注释。  
最终确保理解项目架构。

## 架构图
见文档 https://pv4fr4bbm0.feishu.cn/docx/doxcnM2Y69KvBKS8hf4msfDkpgd

## 个人感悟

### 框架优点：
1. 各个模块职责清晰，单独功能都拆出来。
2. 框架可插拔性高，例如 ConnectStart、PreHandler、PostHandler、ConnectEnd等hook都留了。
3. 里面有些类是单例的，或者每次都需要实例的。作者思路很清晰（在类架构图中可以体现到）
### 框架缺点：
1. 给最终用户写的Router，也就是消息处理函数，设计不合理。 传入的request对象权限太大了，可以通过request获得最高级别的单例 server。   我想在大型团队合作中，这块会出极大风险，导致整个服务不可用。 
2. Connect对象设计不够抽象。 本质上来说MsgHander不关心是TCP/UDP 设置上层的HTTP等，当前设计的过于狭隘。 例如：我想在某些场景下由TCP降级为UDP发送，或者websocket链接降级为http请求。本质上来说当前做不到。 
   再简单的说：MessageHandler消息处理单元，不需要关注消息是通过什么方式给我的，维持与客户端关系那是Connection对象的事情。

