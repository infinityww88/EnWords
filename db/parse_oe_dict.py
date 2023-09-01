import hashlib
import json
import re

def getDigit(s):
    result = hashlib.md5(s.encode("utf-8"))
    return result.digest()

lineDigits = set()

up_chars = list(map(chr, range(ord('A'), ord('Z')+1)))

lineno = 0

STATE_WORD = "word"
STATE_USAGE = "usage"

item = []
state = STATE_WORD

items = []

def getNewItem(item):
    digit = getDigit(item[0][0])
    if digit in lineDigits:
        return
    lineDigits.add(digit)
    exp = item[0][0]
    usage = ""
    if len(item) > 1:
        usage = re.sub(r"^Usage\s+", "", item[1][0])
    items.append([exp, usage])

for line in map(str.strip, open("raw_oe_dict.txt")):
    lineno += 1
    if line == "":
        continue
    if line in up_chars:
        continue
    if state == STATE_WORD:
        if line.startswith("Usage ") and not line.startswith("Usage  n. "):
            raise Exception("get usage when except word at " + str(lineno))
        else:
            item.append((line, lineno))
            state = STATE_USAGE
    elif state == STATE_USAGE:
        if line.startswith("Usage "):
            item.append((line, lineno))
            state = STATE_WORD
            getNewItem(item)
            item.clear()
        else:
            getNewItem(item)
            item.clear()
            item.append((line, lineno))

print(json.dumps(items, indent=4, ensure_ascii=False))