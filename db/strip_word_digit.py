import json, re

dict = json.load(open("raw_oe_dict.json"))
ret = []
for item in dict:
    exp, usage = item
    parts = exp.split(" ")
    parts[0] = re.sub(r"\d+", "", parts[0])
    exp = ' '.join(parts)
    ret.append([exp, usage])

print(json.dumps(ret, indent=4))
