# NeoReads后台

NeoReads应用的后台程序：

- 编程语言：Golang
- 数据库：PostgreSQL
- 操作系统：Ubuntu 18.04 LTS
- 数据格式: JSON


# API设计

| URL | Method | 功能 |
| --- | --- | --- |
| [/book/list](docs/book/list.md) | GET | 图书列表 |
| [/book/{id}/cover](docs/book/cover.md) | GET | 图书封面 |
| [/book/{id}/read](docs/book/read.md) | GET | 图书阅读界面 |
