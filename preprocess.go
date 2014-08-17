package main

import (
    "fmt"
    "github.com/hermanschaaf/enchant"
    "regexp"
        _ "github.com/mattn/go-sqlite3"
    "database/sql"
        "strings"


)

func spellcheck(in <-chan []string, out chan <- []string) {

    // Create an enchant instance
    enchant, err := enchant.NewEnchant()
    if err != nil {
        panic(err)
    }

    enchant.LoadDict("en_US")
    for {
        // Read the next tweet
        tweet, ok := <-in
        if ok {
            for i, w := range tweet {
                if !enchant.Check(w) {
                    suggestions := enchant.Suggest(w)
                    if len(suggestions) >= 1 {
                        tweet[i] = suggestions[0]
                    }
                }
            }
        } else {
            close(out)
            break
        }
    }
}

func expandAbbreviations(in <-chan []string, out chan <- []string) {
    ab := map[string]string {"brb": "be right back",
        "rt": "re-tweet", "abt": "about",
        "bcos":"because", "b": "be", "b4": "before",
        "bfn": "bye for now", "chk": "check",
        "cld": "could", "da": "the",
        "dm": "direct message", "deets": "details",
        "f2f":"face to face", "fab":"fabulous",
        "fav":"favourite", "FF": "follow Friday",
        "ftl":"for the loss",
        "ftw":"for the win", "idk": "I do not know",
        "idc":"i do not care", "kk":"cool cool",
        "nts":"note to self", "plz":"please",
        "tmb":"tweet me back", "u":"you",
        "wos":"was", "woz":"was",
        "wtv":"whatever", "nvm":"never mind",
        "wht":"what", "eva":"ever", "btw":"by the way",
        "fb":"FaceBook", "ftf":"face to face",
        "gr8":"great", "imo": "in my opinion",
        "imho":"in my honest opinion",
        "lmao":"laughing my ass off", "lol":"laugh out loud",
        "np":"no problem", "omg":"oh my god",
        "omfg":"oh my fucking god", "rofl":"rolling on the floor laughing",
        "tmi":"too much information", "ttys":"talk to you soon",
        "ttyl":"talk to you later",
        "wth":"what the hell", "wtf":"what the fuck",
        "cya":"see you", "gotta":"got to",
        "cul8r":"see you later", "dp":"profile picture",
        "ffs":"for fucks sake", "fml":"fuck my life",
        "kmn":"kill me now", "motd":"match of the day",
        "fubar":"fucked up beyond repair", "jk":"joke",
        "pmsl":"pissing myself laughing", "orly":"oh really",
        "plmk":"please let me know",
        "ty":"thank you", "w/e":"whatever", "yolo":"you only live once",
        "4":"four","cldnt":"could not",
        "fave":"favourite", "frnd":"friend", "fwd":"forward",
        "irl":"in real life", "kewl":"cool", "l8":"late", "l8r":"later",
        "l8er":"later", "mil":"million", "njoy":"enjoy",
        "peeps":"people", "ppl":"people", "props":"proper",
        "probs":"probably", "r":"are", "ru":"are you",
        "wbu":"what about you", "thx":"thanks",
        "tx":"thanks", "tyvm":"thank you very much",
        "ur":"your", "yr":"your",
        "lmk":"let me know", "tbh":"too be honest",
        "smh":"shaking my head", "fyi":"for your information",
        "bb":"BlackBerry", "tbt":"throw back Thursday",
        "tho":"though", "whatcha":"what are you",
        "fu":"fuck you", "lil":"little", "dis":"this",
        "urself":"yourself", "c":"see", "luk":"look",
        "cnt":"cannot", "min":"minute", "mins":"minutes",
        "haha":"haha", "hehe":"hehe","2mrw":"tomorrow",
        "2mr":"tomorrow", "cos":"because", "2mrs":"tomorrows", "oke":"ok",
        "ya":"yes", "gnoc":"get naked on camera", "kids":"children", "coz":"because",
        "cont":"continued", "dey":"they", "dere":"there", "nd":"and", "dem":"them",
        "wiv":"with", "tym":"time", "2moz":"tomorrow",
        "bbm":"BlackBerry messenger", "hbu":"how about you",
        "urslf":"yourself", "tonite":"tonight", "bday":"birthday",
        "stfu":"shut the fuck up", "2nite":"tonight", "wuu2":"what you up to",
        "ily":"i love you", "ikr":"i know right", "esp":"especially", "pls":"please", "urs":"yours",
        "obvs":"obviously", "2mora":"tomorrow", "nite":"night", "mom":"mother", "mum":"mother", "mon":"monday",
        "tues":"tuesday", "weds":"wednesday", "thurs":"thursday", "fri":"friday", "ma":"mum",
        "dont":"do not", "wont":"will not", "cant":"cannot", "its":"it is", "didn":"did not",
        "apt":"apartment","sumn":"something", "ifhu":"i fucking hate you",
        "ilya":"i love you", "srsly":"seriously", "im":"i am", "alrdy":"already",
        "hmwrk":"homework", "ano":"i know", "ino":"I know", "gd":"good", "hasn":"has not",
        "aswell":"as well", "gf":"girlfriend","bf":"boyfriend", "wat":"what",
        "lawl":"laugh out loud", "doesn":"does not", "couldn":"could not", "ly":"love you",
        "facebook":"FaceBook", "skool":"school", "wasnt":"was not",
        "rmb":"remember", "c\"mon":"come on", "thnks":"thanks", "gyal":"girl","boi":"boy",
        "lmfao":"laughing my fucking ass off", "ive":"i have",
        "wouldn":"would not", "omgosh":"oh my gosh",
        "wk":"week", "wasup":"what is up", "bu":"about you?", "didnt":"did not",
        "gonna":"going to", "coulda":"could have", "woulda":"would have",
        "dat":"that", "shud":"should","shuld":"should",
        "thnk":"thank", "rmbr":"remember", "skl":"school",
        "txt":"text", "hrs":"hours",
        "sesh":"session", "2b":"to be",
        "fkn":"fucking", "gota":"got to",
        "wud":"would", "cud":"could",
        "nah":"no","bt":"but", "gt":"got", "tht":"that", "tv":"television",
        "mahn":"man", "isnt":"is not", "dnt":"do not", "ima":"i am going to",
        "wid":"with", "fckn":"fucking", "wats":"what is", "app":"application",
        "wif":"with", "cbb":"cannot be bothered", "cba":"cannot be bothered", "hasnt":"has not", "nt":"not", "msn":"messenger",
    }

    for {
        // Read the next tweet
        tweet, ok := <-in
        if ok {
            for i, w := range tweet {
                if r, ok := ab[w]; ok {
                    tweet[i] = r
                }
            }
            out <- tweet
        } else {
            close(out)
            break
        }
    }

}

func replaceFeatures(in <-chan []string, out chan <- []string) {
    urlPattern := regexp.MustCompile("((www\\.[:alpha:]+)|(https?://[:alpha:]+)|(http?://[:alpha:]+))")
    idPattern := regexp.MustCompile("@[^:alpha:]+")
    for {
        // Read the next tweet
        tweet, ok := <-in
        if ok {
            for i, w := range tweet {
                if urlPattern.MatchString(w) {
                    tweet[i] = "URL"
                } else if idPattern.MatchString(w) {
                    tweet[i] = "@USERID"
                }
            }
            out <- tweet
        } else {
            close(out)
            break
        }
    }
}

func readTweets(out chan <- []string) {
    // Open database
    db, err := sql.Open("sqlite3", "../emotionannotate/spam.sqlite")
    if err != nil {
        panic(err)
    }

    // Read from the input table
    sql := "SELECT document FROM input"
    rows, err := db.Query(sql)
    if err != nil {
        panic(err)
    }

    defer rows.Close()
    for rows.Next() {
        var doc string
        rows.Scan(&doc)
        words := strings.Split(doc, " ")
        out <- words
    }
}


func main() {

    replaceChan := make(chan []string, 16)
    abbrevChan := make(chan []string, 16)
    spellCheckChan := make(chan []string, 16)

    go readTweets(replaceChan)
    go replaceFeatures(replaceChan, abbrevChan)
    go expandAbbreviations(abbrevChan, spellCheckChan)

    for i := range spellCheckChan {
        fmt.Println(i)
    }


}
