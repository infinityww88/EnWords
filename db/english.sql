-- client db

create table if not exists recall(word text not null primary key, recall_score integer default 0, last_time integer default 0);

-- server db

create virtual table if not exists docs using fts5(body, title unindexed, type unindexed, tokenize='porter unicode61');

-- static db

create table if not exists zh_dict (word text not null primary key, meaning text not null, phonic text not null);

create virtual table if not exists oe_dict using fts5(word_exp, usage unindexed, tokenize='porter unicode61');

