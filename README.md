# NeoReads后台

NeoReads应用的后台程序：

- 编程语言：Golang
- 数据库：PostgreSQL
- 操作系统：Ubuntu 18.04 LTS
- 数据格式: JSON


# API设计

后台API全部返回JSON格式数据，URL以`/api/v1/`开头

| URL | Method | 功能 |
| --- | --- | --- |
| [/book/list](docs/book/list.md) | GET | 图书列表 |
| [/book/toc](docs/book/toc.md) | GET | 章节目录|
| [/book/:bookid](docs/book/info.md) | GET | 图书信息 |
| [/book/:bookid/chapter/:chapid](docs/book/chapter.md) | GET | 章节内容 |
