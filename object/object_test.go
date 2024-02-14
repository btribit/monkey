package object

import "testing"

func TestStringHashKey(t *testing.T) {
	s1 := &String{Value: "Hello World"}
	s2 := &String{Value: "Hello World"}
	s3 := &String{Value: "My name is johnny"}
	s4 := &String{Value: "My name is johnny"}

	if s1.HashKey() != s2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if s3.HashKey() != s4.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if s1.HashKey() == s3.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
