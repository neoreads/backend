# NeoReads后台

NeoReads应用的后台程序：

- 编程语言：Golang
- 数据库：PostgreSQL
- 操作系统：Ubuntu 18.04 LTS
- 数据格式: JSON

## 近期计划 | TODO

- [x] 对史记文件进行预处理，生成章节文件。
- [ ] 生成BookID, ChapID, ParaID和SentID，并写入到章节文件中。（可以和预处理过程结合到一起，直接生成带ID的章节文件。
- [ ] 将生成的ID存入到数据库中进行管理。
- [ ] 提供函数将带ID的章节文件净化生成普通文本文件，方便阅读。（包括章节排序功能）。
- [ ] 根据ID来存储评注定位。

# API设计

后台API全部返回JSON格式数据，URL以`/api/v1/`开头

| URL | Method | 功能 |
| --- | --- | --- |
| [/book/list](docs/book/list.md) | GET | 图书列表 |
| [/book/toc](docs/book/toc.md) | GET | 章节目录|
| [/book/:bookid](docs/book/info.md) | GET | 图书信息 |
| [/book/:bookid/chapter/:chapid](docs/book/chapter.md) | GET | 章节内容 |
