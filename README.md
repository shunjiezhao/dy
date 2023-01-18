# DY2023  

2023字节青训营---抖音项目---卷王小分队

项目整体结构参照 _MVC_ 结构来构建

go ：

    common : 通用工具包，一些通用工具，例如Redis操作，文件上传下载，远程调用
    controller ：负责转发请求，对请求进行处理。
    dao : 增删改查，不涉及业务逻辑，只是达到按某个条件获得指定数据的要求
    entity ：实体层，放置一个个实体，及其相应的set、get方法， 定义模型（VO，DTO，PO等）
    routes ： 增加路由文件，用于根据请求url进行转发
    service ： 建立增删改查的业务逻辑
    main.go : 启动类

resources ：

    image : 图片资源png、jpg
    video ： 视频资源mp4
    application.yaml ：数据配置文件，例如mysql配置，redis配置

导入依赖(失败的话可以尝试下面三个命令)：

go mod init first

go get -u github.com/gin-gonic/gin

go get -u github.com/jinzhu/gorm