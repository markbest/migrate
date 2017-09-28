## migrate
Golang实现的数据库迁移管理工具

## 配置使用
- 拷贝conf/conf.toml.exmpale为conf/conf.toml，并配置数据库信息
- go build -o migrate migrate.go，编译生成可执行文件migrate并加入path中
- migrate create filename(创建迁移文件)
- migrate up(执行迁移)
- migrate down(回滚迁移）
- migrate status(查看迁移文件的状态)

