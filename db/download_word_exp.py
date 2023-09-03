import pyquery as py
import sys, atexit, time, json

dictPageUrl = "https://www.youdao.com/result?word={word}&lang=en"

if len(sys.argv) < 4:
    print("usage: word_exp total_words done_words")
    exit(0)

def getWords(filename):
    with open(filename) as f:
        return set(map(str.strip, f))

oeWords = getWords(sys.argv[1])
doneWords = getWords(sys.argv[2])

remainWords = oeWords - doneWords

df = open(sys.argv[2], "a")
ef = open(sys.argv[3], "a")

@atexit.register
def clearup():
    df.close()
    ef.close()

def download(word):
    url = dictPageUrl.format(word=word)
    doc = py.PyQuery(url=url)

    t = doc(".catalogue_author .trans-container .word-exp")
    meaning = '\n'.join(map(lambda e: e.text(), t.items()))

    t = doc(".trans-container .phonetic")
    phonic = ""
    if len(t) > 0:
        t = py.PyQuery(t[-1]).text()
        phonic = t.replace("/", "").strip()

    return [word, meaning, phonic]

for w in remainWords:
    for i in range(3):
        time.sleep(1)
        try:
            ret = download(w)
            print(json.dumps(ret, ensure_ascii=False) + ",", file=ef)
            ef.flush()
            print(w, file=df)
            df.flush()
            print(f"ok {w}")
            break
        except Exception as e:
            print(f"failed {w}", e)

            
