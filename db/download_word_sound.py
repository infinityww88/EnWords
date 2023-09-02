import subprocess as sp
import sys, atexit, time

if len(sys.argv) < 3:
    print("usage: total_words done_words")
    exit(0)

def getWords(filename):
    with open(filename) as f:
        return set(map(str.strip, f))

oeWords = getWords(sys.argv[1])
doneWords = getWords(sys.argv[2])

remainWords = oeWords - doneWords

df = open(sys.argv[2], "a")

@atexit.register
def clearup():
    df.close()

for w in remainWords:
    url = f"https://dict.youdao.com/dictvoice?audio={w}&type=2"
    cmd = f"curl '{url}' > word_sounds/{w}.mp3"
    for i in range(3):
        time.sleep(5)
        retcode = sp.call(cmd, shell=True, stderr=sp.DEVNULL)
        if retcode == 0:
            print(f"ok {w}")
            print(w, file=df)
            df.flush()
            break
        else:
            print(f"failed {w}")

