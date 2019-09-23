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
- [x] 将生成的书籍、章节信息存入到数据库中进行管理。
- [x] 在图书目录生成一个静态HTML文件，列出图书的名称和目录链接，方便人工查阅。
- [x] 句子圈点功能
- [ ] 基于句子的基本笔记功能（纯文本）
- [ ] 基于句子的基本笔记功能（Markdown编辑器）
- [ ] 提供函数将带ID的章节文件净化生成普通文本文件，方便阅读。（包括章节排序功能）。
- [ ] 提供一个全库输出功能，将所有书籍输出成普通文本文件的集合，建立压缩包。
- [x] 根据ID来存储评注定位。
- [ ] 使用[authboss](https://github.com/volatiletech/authboss)和[bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt)来实现用户登录和权限功能
- [ ] 重新组织代码目录结构，将数据与数据库处理的模块提出来，供server与prepare模块共享

- [ ] 用户登录信息表
- [x] 集成argon2用于保存和检查密码
- [ ] 利用阿里云发送邮件
- [x] JWT后台

- [ ] 设计词典表单，可以参考[ECDICT](https://github.com/skywind3000/ECDICT) 

- [ ] 参考html2article，利用goquery等工具实现简单的网页内容抓取工具，方便创建外链新闻记事。

- [ ] Markdown处理
  - [x] 处理sentid
  - [x] 处理paraid
  - [x] 支持表格
  - [x] 支持代码块
  - [x] 支持ul列表
  - [x] 支持ol列表
  - [ ] 支持缩进块
  - [x] 支持引用blockquote

## 前后端版本兼容

由于前后端版本号管理是分开的所以需要记录对应兼容的版本号。未来加上前端APP后，这个对应就更重要了。

前端WEB | 后端 
--- | ---
0.2.1 | 0.2.0
0.2.0 | 0.2.0
0.1.9 | 0.1.9
0.1.8 | 0.1.8
0.1.8 | 0.1.7
0.1.7 | 0.1.6
0.1.6 | 0.1.5
0.1.6 | 0.1.4
0.1.5 | 0.1.3
0.1.4 | 0.1.2
0.1.3 | 0.1.2
0.1.2 | 0.1.2
0.1.1 | 0.1
0.1 | 0.1


# API设计

后台API全部返回JSON格式数据，URL以`/api/v1/`开头

| URL | Method | 功能 |
| --- | --- | --- |
| [/book/list](docs/book/list.md) | GET | 图书列表 |
| [/book/toc](docs/book/toc.md) | GET | 章节目录|
| [/book/:bookid](docs/book/info.md) | GET | 图书信息 |
| [/book/:bookid/chapter/:chapid](docs/book/chapter.md) | GET | 章节内容 |
