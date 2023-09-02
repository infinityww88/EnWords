import re, sys

for line in map(str.strip, open("words.txt")):
    line = line.lower()
    if len(line) <= 2:
        #print(line, file=sys.stderr)
        continue
    if ' ' in line and not re.search(r"\s+past of.*$", line):
        #print(line, file=sys.stderr)
        continue
    if "'" in line:
        print(line, file=sys.stderr)
        continue;
    line = re.sub(r"\s+past of.*$", "", line)
    print(line)
    