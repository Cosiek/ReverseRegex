package reverse_regex_test

import (
	"testing"

	reverseRegex "github.com/Cosiek/ReverseRegex"
)

func TestReverseRegex(t *testing.T) {
	// no work case
	// basic case
	rRx := reverseRegex.NewReverseRegex(`/products/(?P<id>\d)/edit`)
	result := rRx.GetReversedString("15")
	expected := "/products/15/edit"
	if result != expected {
		t.Fatalf("Basic case failed\nexpected: %s\ngot: %s", expected, result)
	}

	// something with escaped ()
	rRx = reverseRegex.NewReverseRegex(`/products/(?P<id>\d+)/e\(d\)it`)
	rRx.GetReversedString("15")
	result = rRx.GetReversedString("15")
	expected = "/products/15/e(d)it"
	if result != expected {
		t.Fatalf("Escaping '()' failed\nexpected: %s\ngot: %s", expected, result)
	}

	// multiple groups
	rRx = reverseRegex.NewReverseRegex(`/article/(?P<id>\d)-(?P<slug>.*)`)
	result = rRx.GetReversedString("15", "title-or-something")
	expected = "/article/15-title-or-something"
	if result != expected {
		t.Fatalf("Escaping '()' failed\nexpected: %s\ngot: %s", expected, result)
	}

	// TODO: groups with modifiers - +*?{0-3}
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d)+/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d)*/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d)?/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d){0}/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d){1}/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d){2}/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d){0-1}/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d){0-2}/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d){1-2}/edit`)
	_ = reverseRegex.NewReverseRegex(`/products/(?P<id>\d){2-3}/edit`)
}
