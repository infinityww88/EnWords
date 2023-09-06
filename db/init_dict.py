import sqlite3
import json

zhDictSql = '''create table if not exists zh_dict (word text not null primary key, meaning text not null, phonic text not null)'''

oeDictSql = '''create virtual table if not exists oe_dict using fts5(word_exp, usage unindexed, tokenize='porter unicode61')'''

conn = sqlite3.connect("dict.db")
cur = conn.cursor()
cur.execute("drop table if exists zh_dict")
cur.execute("drop table if exists oe_dict")
cur.execute(zhDictSql)
cur.execute(oeDictSql)

i = 0
oeDict = json.load(open("oe_dict.json"))
for item in oeDict:
    cur.execute("insert into oe_dict values (?, ?)", (item[0], item[1]))
    print("oe", i)
    i += 1

i = 0
zhDict = json.load(open("word_exp.json"))
for item in zhDict:
    cur.execute("insert into zh_dict values (?, ?, ?)", (item[0], item[1], item[2]))
    print("zh", i)
    i += 1

cur.execute("commit")