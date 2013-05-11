package redisbayes

import (
	"fmt"
	"testing"
)

func TestTidy(t *testing.T) {
	test_string := "fjalsdfj $()*#()#*@)&(*&(*^@#*&)!fajs`ldkfj 23"

	if s_out := Tidy(test_string); s_out != "fjalsdfj fajs ldkfj 23" {
		t.Errorf("Tidy failed:\n expected: fjalsdfj fajsldkfj 23\n result:%s\n", s_out)
	}
}

func TestEnglish_tokenizer(t *testing.T) {
	test_string := "fjalsdfj $()*#()#*@)&(*&(*^@#*&)!fajs`ldkfj 23"
	expected_res := []string{"fjalsdfj", "fsjs", "ldkfj", "23"}

	words := English_tokenizer(test_string)

	if words != expected_res {
		t.Errorf("tokenizer failed: %s", expected_res)
	}
}
