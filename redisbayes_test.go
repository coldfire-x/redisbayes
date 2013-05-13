package redisbayes

import (
	"github.com/garyburd/redigo/redis"
	"testing"
)

func TestTidy(t *testing.T) {
	test_string := "fjalsdfj $()*#()#*@)&(*&(*^@#*&)!fajs`ldkfj 23"

	if s_out := Tidy(test_string); s_out != "fjalsdfj fajs ldkfj 23" {
		t.Errorf("Tidy failed:\n expected: fjalsdfj fajsldkfj 23\n result:%s\n", s_out)
	}
}

func TestEnglish_tokenizer(t *testing.T) {
	test_string := "fjalsdfj $(.;)*#()#*@)&(*&(*^@#*&)!fajs`ldkfj 23"
	expected_res := []string{"fjalsdfj", "fajs", "ldkfj"}

	words := English_tokenizer(test_string)
	for i, word := range expected_res {
		if words[i] != word {
			t.Errorf("tokenizer failed, expected: %s", expected_res)
			t.Errorf("tokenizer failed, actually: %s, len:%d", words, len(words))
		}
	}
}

func TestOccurances(t *testing.T) {
	words := []string{"fjalsdfj", "23", "fjalsdfj", "23", "ok"}
	res := Occurances(words)
	expected_res := map[string]uint{
		"23":       2,
		"fjalsdfj": 2,
		"ok":       1,
	}

	for k, v := range expected_res {
		if res[k] != v {
			t.Errorf("Occurances failed: %s", expected_res)
		}
	}
}

func TestFlush(t *testing.T) {
	Train("good", "sunshine drugs love sex lobster sloth")
	Flush()

	exists, err := redis.Bool(redis_conn.Do("EXISTS", redis_prefix+"good"))
	if exists || err != nil {
		t.Errorf("Flush failed")
	}
}

func TestClassify(t *testing.T) {
	Train("good", "sunshine drugs love sex lobster sloth")
	Train("bad", "fear death horror government zombie god")

	class := Classify("sloths are so cute i love them")
	if class != "good" {
		t.Errorf("Classify failed, should be good, result: %s", class)
	}

	class = Classify("i fear god and love the government")
	if class != "bad" {
		t.Errorf("Classify failed, should be bad, result: %s", class)
	}
}

func TestUntrain(t *testing.T) {
    Flush()
	Train("good", "sunshine drugs love sex lobster sloth")
    Untrain("good", "sunshine drugs love sex lobster sloth")

	exists, err := redis.Bool(redis_conn.Do("EXISTS", redis_prefix+"good"))
	if exists || err != nil {
		t.Errorf("Untrain failed %s, %s", exists, err)
	}
}

func init() {
	Flush()
}
