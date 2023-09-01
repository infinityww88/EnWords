import json
from tap import Tap
import sqlite3

oeTable = "oe_dict"

class ArgParser(Tap):
    dbname: str # db name
    def configure(self):
        self.add_argument("dbname")

args = ArgParser().parse_args()

with open("raw_oe_dict.json") as f:
    dict = json.load(f)

conn = sqlite3.connect(args.dbname)
cur = conn.cursor()

cur.execute(f"drop table if exists {oeTable}")
cur.execute(f'''create virtual table if not exists {oeTable} using fts5 (word_exp, usage, token = "porter unicode61 tokenchars '-_'")''')

insertSql = f"insert into {oeTable} values(?, ?)"

for item in dict:
    cur.execute(insertSql, item)

cur.execute("commit")
