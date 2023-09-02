import json, re
import sys
from unidecode import unidecode

items = json.load(open("raw_oe_dict.json"))
words = set()

for item in items:
    sentence: str = item[0]
    if re.search("\sAbbr.\s", sentence):
        print(t, file=sys.stderr)
        continue
    word = re.split(r'\s[A-Za-z—\-.]+\.\s', sentence)[0]
    t = word.split(" ")[0]
    if t.startswith("-") or t.endswith("-"):
        print(t, file=sys.stderr)
    else:
        word = word.strip()
        word = re.sub(r'\d+', "", word)
        word = re.sub(r'\s*\(.*\)', "", word)
        word = re.sub(r'\s+see.*\.$', "", word)
        words.add(unidecode(word))
words = sorted(words)
for w in words:
    print(w)
