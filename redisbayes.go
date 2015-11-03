package redisbayes

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	redis_conn   redis.Conn
	redis_prefix = "bayes:"
	correction   = 0.1
)

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
	// this slice's length should be initialized to 0
	// otherwise, the first element will be the whitespace(empty string)
	res := make([]string, 0)

	for _, word := range words {
		res = append(res, strings.TrimSpace(word))
	}

	return res
}

// compute word occurances
func Occurances(words []string) map[string]uint {
	counts := make(map[string]uint)
	for _, word := range words {
		if _, ok := counts[word]; ok {
			counts[word] += 1
		} else {
			counts[word] = 1
		}
	}

	return counts
}

func Flush() {
	reply, err := redis_conn.Do("SMEMBERS", redis_prefix+"categories")
	if err != nil {
		log.Println(err)
		return
	}

	for _, key := range reply.([]interface{}) {
		redis_conn.Do("DEL", redis_prefix+string(key.([]byte)))
	}

	redis_conn.Do("DEL", redis_prefix+"categories")
}

func Train(categories, text string) {
	redis_conn.Do("SADD", redis_prefix+"categories", categories)

	token_occur := Occurances(English_tokenizer(text))
	for word, count := range token_occur {
		redis_conn.Do("HINCRBY", redis_prefix+categories, word, count)
	}
}

func Untrain(categories, text string) {
	token_occur := Occurances(English_tokenizer(text))
	for word, count := range token_occur {
		reply, err := redis_conn.Do("HGET", redis_prefix+categories, word)
		if err != nil {
			log.Println(err)
			return
		}

		cur, err := strconv.ParseUint(string(reply.([]byte)), 10, 0)
		if cur != 0 {
			inew := cur - uint64(count)
			if inew > 0 {
				redis_conn.Do("HSET", redis_prefix+categories, word, inew)
			} else {
				redis_conn.Do("HDEL", redis_prefix+categories, word)
			}
		}
	}

	if Tally(categories) == 0 {
		redis_conn.Do("DEL", redis_prefix+categories)
		redis_conn.Do("SREM", redis_prefix+"categories", categories)
	}
}

func Classify(text string) string {
	scores := Score(text)
	key := ""
	max := 0.0

	if scores != nil {
		for k, v := range scores {
			if v <= max {
				max = v
				key = k
			}
		}

		return key
	}

	return "I dont know"
}

func Score(text string) map[string]float64 {
	token_occur := Occurances(English_tokenizer(text))
	res := make(map[string]float64)

	reply, err := redis_conn.Do("SMEMBERS", redis_prefix+"categories")
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, category := range reply.([]interface{}) {
		category := string(category.([]byte))
		tally := Tally(category)
		if tally == 0 {
			continue
		}

		res[category] = 0.0
		for word, _ := range token_occur {
			score, err := redis_conn.Do("HGET", redis_prefix+category, word)
			if err != nil {
				log.Println(err)
				return nil
			}

			if score == nil {
				continue
			}

			iVal, err := strconv.ParseFloat(string(score.([]byte)), 64)
			if err != nil {
				log.Println(err)
				return nil
			}

			if iVal == 0.0 {
				iVal = correction
			}
			res[category] += math.Log(iVal / float64(tally))
		}
	}

	return res
}

func Tally(category string) (sum uint64) {
	vals, err := redis_conn.Do("HVALS", redis_prefix+category)
	if err != nil {
		log.Println(err)
		return
	}

	for _, val := range vals.([]interface{}) {
		iVal, err := strconv.ParseUint(string(val.([]byte)), 10, 0)
		if err != nil {
			log.Println(err)
			return
		}

		sum += iVal
	}

	return sum
}
func ConnectToRedis(host, port string) error {
	var err error
	redis_conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

	return err
}
