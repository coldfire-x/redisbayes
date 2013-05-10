package main

import (
	//"math"
    "fmt"
    "log"
    "strings"
    "github.com/kylelemons/go-gypsy/yaml"
    //"github.com/garyburd/redigo/redis"
)

type Text string
type Word string

// replace \_.,<>:;~+|\[\]?`"!@#$%^&*()\s chars with whitespace
// re.sub(r'[\_.,<>:;~+|\[\]?`"!@#$%^&*()\s]', ' ' 
func Tidy(s Text) Text {
        return ""
}

// tidy the input text, ignore those text composed with less than 2 chars 
func English_tokenizer(s Text) Text {
        return ""
}

// compute word occurances
func Occurances(w Word) map[Word]uint {
        return nil
}

func main() {
    // load config file
    cfg_filename := "config.yaml"
    config, err := yaml.ReadFile(cfg_filename)
    if err != nil {
            log.Fatalf("readfile(%s): %s", cfg_filename, err)
    }

    // get english ignore entire string
    english_ignore, err := config.Get("english_ignore")
    if err != nil {
            log.Fatalf("%s parse error: %s\n", english_ignore, err)
    }

    // get each separated words
    english_ignore_words_list := strings.Fields(english_ignore)
    for _, value := range english_ignore_words_list {
        log.Println(value)
    }

    fmt.Println("ok")
}
