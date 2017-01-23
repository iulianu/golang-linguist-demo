package main

import (
	"testing"
	"reflect"
)

func TestParsePhrasePairStream(t *testing.T) {
	text := []byte("PP 2 2 l'homme $ the man & blah PP 2 2 the woman $ la femme & lorem PP 1 1 Prost $ cheers & noroc")
	pairs, err := ParsePhrasePairStream(text)
	if err != nil {
		t.Fatalf("Parsing failed with: %s", err)
	}
	expected := []PhrasePair {
		{	source:Phrase{tokenCount:2, tokens:"l'homme"},
			target:Phrase{tokenCount:2, tokens:"the man"},
			memo: "blah"},
		{	source:Phrase{tokenCount:2, tokens:"the woman"},
			target:Phrase{tokenCount:2, tokens:"la femme"},
			memo: "lorem"},
		{	source:Phrase{tokenCount:1, tokens:"Prost"},
			target:Phrase{tokenCount:1, tokens:"cheers"},
			memo: "noroc"},
	}
	if ! reflect.DeepEqual(pairs, expected) {
		t.Fatalf("Parsing failed, expected %s, got %s", expected, pairs)
	}
}

func TestDetectsSyntaxError(t *testing.T) {
	text := []byte("PP brokenString")
	_, err := ParsePhrasePairStream(text)
	if( err.Error() != "Phrase pair syntax error" ) {
		t.Fatalf("Parsing should have failed with Phrase pair syntax error, failed with %s", err)
	}
}

func TestDetectsBrokenEncoding(t *testing.T) {
	text := []byte("PP 2 2 l'")
	text = append(text, 0xFA)
	_, err := ParsePhrasePairStream(text)
	if( err.Error() != "Phrase pair syntax error" ) {
		t.Fatalf("Parsing should have failed with Phrase pair syntax error, failed with %s", err)
	}
}