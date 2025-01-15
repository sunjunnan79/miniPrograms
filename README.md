# MIniprograms项目

## 概述:

这个是用来规避小程序审核的服务,运行时需要启动两个一样的服务

- miniprograms
- miniprograms_changer

基本使用方式为

1. 使用账号密码创建项目
2. 审核时修改项目状态
3. 先获取当前状态以隐藏无法通过审核的页面
4. 当审核通过后修改项目状态
5. 此时再访问就可以显示被隐藏的页面了





## API:

提供两个接口(身份验证直接使用账号密码)：

直接看desc/swagger文档



## 技术栈:

基于golang1.23开发

- 数据库：sqlite(并发写效率极低，但是并发读没问题，运行需要开启cgo)，极其适合这种轻量级场景
- 缓存：sync.map(同样低频写,高频读)
- 框架：gin+gorm
- MVC架构(其实所有业务逻辑都在main函数里面算不上mvc...)



## 部署:

使用docker进行二阶段打包

使用docker-compose进行部署，在deploy文件夹下存放了部署相关文档,使用的时候需要修改配置文件为你需要的配置文件



需要部署在两个域名下，为了不影响正在使用的服务，更换时前端应当切换到另一个域名

- miniprograms.muxixyz.com
- miniprogramschanger.muxixyz.com
