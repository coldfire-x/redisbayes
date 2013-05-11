package redisbayes

import (
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
	expected_res := []string{"fjalsdfj", "fajs", "ldkfj", "23"}

	words := English_tokenizer(test_string)
	for i, word := range expected_res {
		if words[i] != word {
			t.Errorf("tokenizer failed: %s", expected_res)
		}
	}
}

func TestOccurances(t *testing.T) {
	words := []string{"fjalsdfj", "23", "fjalsdfj", "23", "ok"}
	res := Occurances(words)
	expected_res := map[string]uint{
		"23":       2,
		"fjalsdfj": 2,
		"ok":       1}

	for k, v := range expected_res {
		if res[k] != v {
			t.Errorf("Occurances failed: %s", expected_res)
		}
	}
}
