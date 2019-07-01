-- table for book info
drop table if exists book;
create table book (
    id char(10),
    title varchar(200),
    authors varchar(100)
);

-- table for person info
drop table if exists person;
create table person (
    id char(10),
    fname varchar(100), -- first name
    mname varchar(100), -- middle name
    sname varchar(100) -- surname
);

-- n to n relation book <-> author
create table book_author;
create table book_author (
    book_id char(10),
    author_id char(10)
);

drop table if exists chapter;
-- table for chapter info
create table chapter (
    id char(8),
    title varchar(200),
    path varchar(200) -- location of the chapter's file
);