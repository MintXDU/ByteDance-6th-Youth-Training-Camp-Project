# 开发项目

## 需要做的
1. 数据目前是以对象的形式存在内存中的，需要把在内存中的数据存在数据库中

## 项目结构
/controller 控制层
/service    业务层
/dao        数据层
/public     静态资源

理想的单向调用链：控制层 => 业务层 => 数据层

### GORM
#### 安装
```
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```
#### 连接
见文件 /service/mysql 中的函数 Connection()

其他服务想要连接数据库可调用该函数