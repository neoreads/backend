drop schema public cascade;
create schema public;

-- table for book info
drop table if exists books;
create table books (
    id char(8),
    title varchar(200),
    authors varchar(100),
    CONSTRAINT books_pk PRIMARY KEY (id)
);

-- table for person info
drop table if exists people;
create table people (
    id char(8) PRIMARY KEY,
    lastname varchar(100), -- last name
    firstname varchar(100), -- first name
    fullname varchar(100) -- full name
);

-- n to n relation book <-> author
drop table if exists book_author;
create table book_author (
    book_id char(8),
    author_id char(8)
);

create table chapters (
    id char(4),
    "order" int,
    bookid char(8),
    title varchar(200),
    CONSTRAINT book_chapter_key PRIMARY KEY (id, bookid)
);

-- notes
drop table if exists notes;
create table notes (
    id char(8),
    "time" timestamp NOT NULL DEFAULT NOW(),
    ntype smallint, -- note type: 0: mark; 1: note; 2: annotation; 3: comment; 4: reference;  5: dict;
    ptype smallint, -- position type: 0: word; 1: sentence; 2: paragraph; 3: chapter; 4: book;
    pid char(8), -- person id, refers to table people
    bookid char(8),
    chapid char(4),
    paraid char(4),
    sentid char(4),
    wordid char(4), -- TODO: this pos may change to startOffset & endOffset
    content text DEFAULT '', -- for simple notes, this is markdown content; for complex notes like dictionary, this field is empty, and out reference table is required
    CONSTRAINT notes_key PRIMARY KEY (id)
);

-- comment
drop table if exists comments;
create table comments (
    id char(8) PRIMARY KEY
    nid cahr(8), -- note id, refers to note table
);

-- dict
drop table if exists dict;

create table dict (
    id SERIAL PRIMARY KEY,
    lang char(2), -- quick access to languages(lang)
    "langid" smallint, -- refer to language table
    word varchar(100), -- e.g.: '天', 'Sky', 'Ciel'
);

-- sentence dict
drop table if exists sentdict;
create table sentdict (
    id SERIAL PRIMARY KEY,
    lang char(2),
    "langid" smallint,
    sent varchar(1000)
);

-- languages
drop table if exists languages;
create table languages (
    id SERIAL PRIMARY KEY,
    lang char(2), -- language: en, zh, etc. ISO_639_1; TODO: might need ISO_639_2 in the future;
    scode varchar(4), -- script code (for writing system), ISO_15924, e.g. Hans for Simplified Chinese
    sno smallint, -- script number, ISO_15924
);

-- user
drop table if exists users;
create table users (
    id SERIAL PRIMARY KEY,
    username varchar(12),
    email varchar(40),
    pid char(8) REFERENCES people(id), -- person id
    pwd varchar(100)
);

CREATE OR REPLACE VIEW users_people AS
 SELECT u.id AS uid,
    p.id AS pid,
    u.username,
    p.firstname,
    p.lastname,
    p.fullname
   FROM users u,
    people p
  WHERE u.pid = p.id;

-- test data ---

-- insert into books VALUES ('00000001', '史记', '司马迁');
-- insert into people VALUES ('00000001', '司马', '迁', '司马迁');
-- insert into book_author VALUES ('00000001', '00000001');
-- insert into chapters VALUES ('0001', 1, '00000001', '五帝本纪第一');
-- insert into chapters VALUES ('0002', 2, '00000001', '夏本纪第二');
-- insert into chapters VALUES ('0003', 3, '00000001', '殷本纪第三');