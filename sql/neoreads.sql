-- table for book info
create table book (
    id char(10),
    title varchar(200),
    authors varchar(100)
)

-- table for author info
create table author (
    id char(10),
    fname varchar(100),
    mname varchar(100),
    sname varchar(100)
)

-- n to n relation book <-> author
create table book_author (
    book_id char(10)
    author_id char(10)
)

-- table for chapter info
create table chapter (
    id char(10),
    title varchar(200),
    file varchar(200), -- location of the chapter's file
    loc bigint -- location of the chapter, if it is part of a file
)
