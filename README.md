# DY2023  

2023字节青训营---抖音项目---卷王小分队

项目整体结构参照 `easy_note` 结构来构建


resources ：

    image : 图片资源png、jpg
    video ： 视频资源mp4
    application.yaml ：数据配置文件，例如mysql配置，redis配置
## 项目目录
```
├── go.sum
├── idl
│   └── user.proto
├── kitex_gen  --kitex生成的rpc 接口文件
├── pkg        
│   ├── configs  -- 一些配置
│   ├── constants -- 常量
│   ├── errno -- 错误码
│   └── middleware  -- 中间件
│       └──  jwt.go -- jwt 鉴权
└── service 
    ├── api
    │   ├── handlers
    │   │   └── user
    │   ├── router                -- 调用各个用于初始化router的函数
    │   ├── rpc                   -- http 用到的 rpc 服务
    │   │   ├── mock
    │   │   └── user              -- 用户rpc服务
    │   └── tests                 --测试整条链路的调用
    └── user -- 用户微服务
        ├── main.go
        ├── handler.go      -- rpc 服务 调用 service 层提供的服务
        ├── model           -- dao 层
        │   └── db
        ├── pack            -- 将 dao 层对象打包
        ├── script
        └── service         -- service 层
 

```
