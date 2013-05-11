package redisbayes

import (
	//"math"
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
	"regexp"
	"strings"
	//"github.com/garyburd/redigo/redis"
)

var english_ignore_words_map = make(map[string]int)

// replace \_.,<>:;~+|\[\]?`"!@#$%^&*()\s chars with whitespace
// re.sub(r'[\_.,<>:;~+|\[\]?`"!@#$%^&*()\s]', ' ' 
func Tidy(s string) string {
	reg, err := regexp.Compile("[\\_.,:;~+|\\[\\]?`\"!@#$%^&*()\\s]+")
	if err != nil {
		log.Fatal(err)
	}

	text_in_lower := strings.ToLower(s)
	safe := reg.ReplaceAllLiteralString(text_in_lower, " ")

	return safe
}

// tidy the input text, ignore those text composed with less than 2 chars 
func English_tokenizer(s string) []string {
	words := strings.Fields(Tidy(s))

	for index, word := range words {
		strings.TrimSpace(word)
		_, omit := english_ignore_words_map[word]
		if omit || len(word) < 2 {
			words = words[:index+copy(words[index:], words[index+1:])]
		}
	}

	return words
}

// compute word occurances
func Occurances(w string) map[string]uint {
	return nil
}

// init function, load the configs
// fill english_ignore_words_map
func init() {
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
	for _, word := range english_ignore_words_list {
		word = strings.TrimSpace(word)
		english_ignore_words_map[word] = 1
	}
}
