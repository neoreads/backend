drop schema public cascade;
create schema public;

-- table for book info
drop table if exists books;
create table books (
    id char(8),
    title varchar(200),
    authors varchar(100)
);

-- table for person info
drop table if exists people;
create table people (
    id char(8),
    surname varchar(100), -- surname
    name varchar(101), -- name
    fname varchar(100) -- full name
);

-- n to n relation book <-> author
drop table if exists book_author;
create table book_author (
    book_id char(8),
    author_id char(8)
);

create table chapters (
    id char(3),
    "order" int,
    bookid char(8),
    title varchar(200),
    CONSTRAINT book_chapter_key PRIMARY KEY (id, bookid)
);

-- test data ---
insert into books VALUES ('00000001', '史记', '司马迁');
insert into people VALUES ('00000001', '司马', '迁', '司马迁');
insert into book_author VALUES ('00000001', '00000001');
insert into chapters VALUES ('001', 1, '00000001', '五帝本纪第一');
insert into chapters VALUES ('002', 2, '00000001', '夏本纪第二');
insert into chapters VALUES ('003', 3, '00000001', '殷本纪第三');