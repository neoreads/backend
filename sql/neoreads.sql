drop schema public cascade;
create schema public;


--- begin of utils

--- ID Gen Trigger, refer to https://blog.andyet.com/2016/02/23/generating-shortids-in-postgres/

-- Create a trigger function that takes no arguments.
-- Trigger functions automatically have OLD, NEW records
-- and TG_TABLE_NAME as well as others.
CREATE OR REPLACE FUNCTION gen_unique_id()
RETURNS TRIGGER AS $$

 -- Declare the variables we'll be using.
DECLARE
  key TEXT;
  qry TEXT;
  found TEXT;
BEGIN

  -- generate the first part of a query as a string with safely
  -- escaped table name, using || to concat the parts
  qry := 'SELECT id FROM ' || quote_ident(TG_TABLE_NAME) || ' WHERE id=';

  -- This loop will probably only run once per call until we've generated
  -- millions of ids.
  LOOP

    -- Generate our string bytes and re-encode as a base64 string.
    key := encode(gen_random_bytes(TG_ARGV[0]::int*3/4), 'base64');

    -- Base64 encoding contains 2 URL unsafe characters by default.
    -- The URL-safe version has these replacements.
    key := replace(key, '/', '_'); -- url safe replacement
    key := replace(key, '+', '-'); -- url safe replacement

    -- Concat the generated key (safely quoted) with the generated query
    -- and run it.
    -- SELECT id FROM "test" WHERE id='blahblah' INTO found
    -- Now "found" will be the duplicated id or NULL.
    EXECUTE qry || quote_literal(key) INTO found;

    -- Check to see if found is NULL.
    -- If we checked to see if found = NULL it would always be FALSE
    -- because (NULL = NULL) is always FALSE.
    IF found IS NULL THEN

      -- If we didn't find a collision then leave the LOOP.
      EXIT;
    END IF;

    -- We haven't EXITed yet, so return to the top of the LOOP
    -- and try again.
  END LOOP;

  -- NEW and OLD are available in TRIGGER PROCEDURES.
  -- NEW is the mutated row that will actually be INSERTed.
  -- We're replacing id, regardless of what it was before
  -- with our key variable.
  NEW.id = key;

  -- The RECORD returned here is what will actually be INSERTed,
  -- or what the next trigger will get if there is one.
  RETURN NEW;
END;
$$ language 'plpgsql';


---- end of utils

---- begin of tables

-- table for book info
drop table if exists books;
create table books (
    id char(8),
    lang char(2) NOT NULL DEFAULT 'zh', -- language: en, zh, etc. ISO_639_1; TODO: might need ISO_639_2 in the future;
    title varchar(200),
    cover char(8),
    intro text,
    CONSTRAINT books_pk PRIMARY KEY (id)
);

-- table for person info
drop table if exists people;
create table people (
    id char(8) PRIMARY KEY,
    lastname varchar(100) NOT NULL DEFAULT '', -- last name
    firstname varchar(100) NOT NULL DEFAULT '', -- first name
    fullname varchar(100) NOT NULL DEFAULT '', -- full name
    othernames varchar(200) NOT NULL DEFAULT '',
    intro text NOT NULL DEFAULT '',
    avatar char(8) NOT NULL DEFAULT '' -- avatar photo id
);

DROP TRIGGER if exists trigger_people_genid on people;
CREATE TRIGGER trigger_people_genid BEFORE INSERT ON people FOR EACH ROW EXECUTE PROCEDURE gen_unique_id(8);


-- TODO: rename to books_authors
drop table if exists books_people;
create table books_people(
    bookid char(8),
    pid char(8)
);

drop table if exists books_collaborators;
create table books_collaborators(
    bookid char(8) NOT NULL,
    kind smallint NOT NULL DEFAULT 0, -- kind: 0: initiator; 1: watcher; 2: contributor; 3: translator
    pid char(8) NOT NULL
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
    ntype smallint, -- note type: 0: mark; 1: note; 2: phonetics; 3: reference;  4: translation;
    ptype smallint, -- position type: 0: word; 1: sentence; 2: paragraph; 3: article; 4: collection;
    pid char(8), -- person id, refers to table people
    colid char(8),
    artid char(8),
    paraid char(4),
    sentid char(4),
    startpos smallint DEFAULT 0,
    endpos smallint DEFAULT 0,
    content text DEFAULT '', -- for simple notes, this is markdown content; for complex notes like dictionary, this field is empty, and out reference table is required
    CONSTRAINT notes_key PRIMARY KEY (id)
);

-- comment
drop table if exists comments;
create table comments (
    id char(8) PRIMARY KEY,
    nid cahr(8), -- note id, refers to note table
);

-- articles
drop table if exists articles;
create table articles (
    id char(8) PRIMARY KEY,
    kind smallint NOT NULL DEFAULT 0, -- article type: 0: chapter, 1: blog, 2: poem, 3: emark 
    addtime timestamp NOT NULL DEFAULT NOW(),
    modtime timestamp NOT NULL DEFAULT NOW(),
    title varchar(255) NOT NULL,
    content text -- may move to disk as markdown file in the future
);

drop table if exists article_people;
-- one article may have multiple authors
create table article_people (
    aid char(8), -- article id
    pid char(8), -- person id
    CONSTRAINT article_people_pk PRIMARY KEY (pid, aid)
);

-- collections
drop table if exists collections;
create table collections (
    id char(8) PRIMARY KEY,
    kind smallint NOT NULL DEFAULT 0, -- collection type: 0: book, 1: collection
    addtime timestamp NOT NULL DEFAULT NOW(),
    modtime timestamp NOT NULL DEFAULT NOW(),
    title varchar(255) NOT NULL,
    intro text
);

drop table if exists collections_people;
-- one collection may have multiple authors
create table collections_people (
    colid char(8), -- collection id
    pid char(8) -- person id
);

-- n to n relationship
drop table if exists collections_articles;
create table collections_articles (
    colid char(8),
    artid char(8)
);


-- tags
drop table if exists tags;
create table tags (
    id char(8) PRIMARY KEY,
    kind smallint NOT NULL DEFAULT 0, -- tag type: 0: topic, 1: event, 2: people, 3: place, 4: time, 5: emotion
    role smallint NOT NULL DEFAULT 0, -- used for: 0: books, 1: articles, 2: news, 3: people, 4: notes
    tag varchar(200)
);

DROP TRIGGER if exists trigger_tags_genid on tags;
CREATE TRIGGER trigger_tags_genid BEFORE INSERT ON tags FOR EACH ROW EXECUTE PROCEDURE gen_unique_id(8);


drop table if exists people_tags;
create table people_tags (
    pid char(8),
    tid char(8),
    CONSTRAINT people_tags_pk PRIMARY KEY (pid, tid)
);

-- news
drop table if exists news;
create table news (
    id char(8) PRIMARY KEY,
    kind smallint, -- news type: 0: external-link, 1: markdown post, 2: image/gif/video
    addtime timestamp DEFAULT now(),
    modtime timestamp DEFAULT now(),
    link text,
    source varchar(200),
    title text,
    summary text,
    content text
);

-- news_tags relation
drop table if exists news_tags;
create table news_tags (
    newsid char(8),
    tagid char(8),
    PRIMARY KEY (newsid, tagid)
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


-- stars --

-- 用来记录每人每日每个类型剩余可用的喜爱值
drop table if exists stars_quota;
create table stars_quota (
    kind smallint not null default 0, -- kind: 0:chapter, 1:blog, 2:poem, 3: emark
    userid int,
    remaining smallint
);

-- 用来记录每人每个实体的总喜爱值
drop table if exists stars;
create table stars (
    kind smallint not null default 0, -- kind: 0:chapter, 1:blog, 2:poem, 3: emark
    uid int,
    eid varchar(20),
    num smallint
);

-- test data ---

-- insert into books VALUES ('00000001', '史记', '司马迁');
-- insert into people VALUES ('00000001', '司马', '迁', '司马迁');
-- insert into book_author VALUES ('00000001', '00000001');
-- insert into chapters VALUES ('0001', 1, '00000001', '五帝本纪第一');
-- insert into chapters VALUES ('0002', 2, '00000001', '夏本纪第二');
-- insert into chapters VALUES ('0003', 3, '00000001', '殷本纪第三');

