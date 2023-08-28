create table if not exists words(wid integer primary key autoincrement, word text not null, is_phrase integer default 0, phonetic_symbol text, meaning text);

create table if not exists sentences(sid integer primary key autoincrement, sentence text not null);

create table if not exists word_sentences(wsid integer primary key autoincrement, wid int not null, sid int not null);

create table if not exists notes(nid integer primary key autoincrement, note text not null);

create table if not exists word_notes(wnid integer primary key autoincrement, wid int not null, nid int not null);
