## migrate
golang简单实现的数据库迁移管理工具

## 配置使用
- 拷贝conf/conf.toml.exmpale为conf/conf.toml，并配置数据库信息
- go build -o m migrate.go,编译生成可执行文件m，加入path中
- m create filename(创建迁移文件)
- m up(执行迁移)
- m down(回滚迁移）
- m status(查看迁移文件的状态)

