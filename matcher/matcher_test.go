package matcher

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// Given an inexisting file, it should return an error
	if _, err := New("/does/not/exist"); err == nil {
		t.Error("No error for inexisting file")
	}
	// Given a bogus regex
	tmpfile := mktemp(t, `foo(`)
	defer os.Remove(tmpfile)
	if _, err := New(tmpfile); err == nil {
		t.Error("No error for an invalid regex")
	}
}

func TestMatches(t *testing.T) {
	// Should fail for a missing file
	if _, err := New("/does/not/exist"); err == nil {
		t.Error("No error for inexisting file")
	}

	tmpfile := mktemp(t, "^test\n")
	defer os.Remove(tmpfile)
	m, err := New(tmpfile)
	if err != nil {
		t.Fatal(err)
	}
	// Should match for a matching regex
	if !m.Matches("test this thing") {
		t.Error("Expected '^test' to match 'test this thing', but it did not")
	}
	// Should not match for a non-matching regex
	if m.Matches("don't test this thing") {
		t.Error("Expected '^test' not to match 'don't test this thing', but it did")
	}

}

func mktemp(t *testing.T, contents string) string {
	tmpfile, err := ioutil.TempFile("", "journalcheck_matcher_test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmpfile.WriteString(contents)
	if err != nil {
		os.Remove(tmpfile.Name())
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		os.Remove(tmpfile.Name())
		t.Fatal(err)
	}
	return tmpfile.Name()
}
