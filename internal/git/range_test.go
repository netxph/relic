package git

import "testing"

func TestParseRange_SingleHash(t *testing.T) {
	r, err := ParseRange("abc1234")
	if err != nil {
		t.Fatal(err)
	}
	if r.From != "abc1234" || r.To != "HEAD" {
		t.Errorf("got %+v", r)
	}
}

func TestParseRange_ExplicitRange(t *testing.T) {
	r, err := ParseRange("abc1234..def5678")
	if err != nil {
		t.Fatal(err)
	}
	if r.From != "abc1234" || r.To != "def5678" {
		t.Errorf("got %+v", r)
	}
}

func TestParseRange_TooShortHash(t *testing.T) {
	_, err := ParseRange("abc12")
	if err == nil {
		t.Fatal("expected error for short hash")
	}
}

func TestParseRange_EmptyInput(t *testing.T) {
	_, err := ParseRange("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseRange_TooShortFrom(t *testing.T) {
	_, err := ParseRange("abc12..def5678")
	if err == nil {
		t.Fatal("expected error for short from hash")
	}
}
