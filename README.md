# 通用的后端架构
---

1. common
   - 公共方法
   - 通用响应方法
     - 成功 200
     - 访问拒绝 403
     - 未找到 404
     - server error 500
2. conf
   - 加载配置方法
   - 使用yaml的方式解析配置
   - mysql配置加载-主从分离
3. controller
   - 控制器
   - 登录
   - 用户信息
   - 用户权限信息
4. middleware
   - 中间件
   - 跨域解决
5. model
   - orm模型
   - 用户权限表
   - 用户表
   - 用户权限关联表
6. router
   - 路由
7. service
   - 业务
   - 登录
   - 用户信息
   - 用户权限信息
8. utils
   - 配置值
   - 响应码
     - 成功 200
     - 访问拒绝 403
     - 未找到 404
     - server error 500