# NeoReads后台

NeoReads应用的后台程序：

- 编程语言：Golang
- 数据库：PostgreSQL
- 操作系统：Ubuntu 18.04 LTS
- 数据格式: JSON

## 近期计划 | TODO

- [x] 对史记文件进行预处理，生成章节文件。
- [x] 生成BookID, ChapID, ParaID和SentID，并写入到章节文件中。（可以和预处理过程结合到一起，直接生成带ID的章节文件。
- [x] 修改章节文件格式，支持Markdown
- [x] 将史记处理转化为基本Markdown格式
- [ ] 将生成的书籍、章节信息存入到数据库中进行管理。
- [ ] 在图书目录生成一个静态HTML文件，列出图书的名称和目录链接，方便人工查阅。
- [ ] 提供函数将带ID的章节文件净化生成普通文本文件，方便阅读。（包括章节排序功能）。
- [ ] 提供一个全库输出功能，将所有书籍输出成普通文本文件的集合，建立压缩包。
- [ ] 根据ID来存储评注定位。
- [ ] 使用[authboss](https://github.com/volatiletech/authboss)和[bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt)来实现用户登录和权限功能
- [ ] 使用[makrdown-it](https://github.com/markdown-it/markdown-it)/[showdown](https://github.com/showdownjs/showdown)或[remark](https://github.com/remarkjs/remark)来实现其前端Markdown解析


# API设计

后台API全部返回JSON格式数据，URL以`/api/v1/`开头

| URL | Method | 功能 |
| --- | --- | --- |
| [/book/list](docs/book/list.md) | GET | 图书列表 |
| [/book/toc](docs/book/toc.md) | GET | 章节目录|
| [/book/:bookid](docs/book/info.md) | GET | 图书信息 |
| [/book/:bookid/chapter/:chapid](docs/book/chapter.md) | GET | 章节内容 |
