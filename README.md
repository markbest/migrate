## migrate
Golang实现的数据库迁移管理工具

## 配置使用
- 拷贝conf/conf.toml.exmpale为conf/conf.toml，并配置数据库信息
- go build -o bin/migrate migrate.go，编译生成可执行文件
- bin/migrate create [filename]（创建迁移文件）
- bin/migrate up（执行迁移）
- bin/migrate down（回滚迁移）
- bin/migrate status（查看迁移文件的状态）