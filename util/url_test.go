package util

import (
	"net/http"
	"net/url"
	"testing"
)

func TestFullUrlFrom(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://big-data-warahouse.xyz?hello=world", nil)

	expected := "http://big-data-warahouse.xyz?overriden=true#"
	actual := fullUrlFrom(req, url.Values{"overriden": {"true"}})
	if expected != actual {
		t.Fatalf("fullUrlFrom should returns original URL suffixed with # and overriden query params")
	}
}

func TestFullUrlFromRequest(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://big-data-warahouse.xyz?hello=world", nil)

	expected := "http://big-data-warahouse.xyz?hello=world#"
	actual := FullUrlFromRequest(req)

	if expected != actual {
		t.Fatalf("FullUrlFromRequest should returns original URL suffixed with # and original query params")
	}
}

func TestNextFullUrlFromRequest_WithoutPageQueryArg(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://big-data-warahouse.xyz?hello=world", nil)

	expected := "http://big-data-warahouse.xyz?hello=world&page=2#"
	actual := NextFullUrlFromRequest(req)

	if expected != actual {
		t.Fatalf("NextFullUrlFromRequest should returns original URL suffixed with # and original query params and page set to 2")
	}
}

func TestNextFullUrlFromRequest_WithPageQueryArg(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://big-data-warahouse.xyz?hello=world&page=12", nil)

	expected := "http://big-data-warahouse.xyz?hello=world&page=13#"
	actual := NextFullUrlFromRequest(req)

	if expected != actual {
		t.Fatalf("NextFullUrlFromRequest should returns original URL suffixed with # and original query params and page set to page+1")
	}
}

func TestNextPreviousFullUrlFromRequest_WithoutPageQueryArg(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://big-data-warahouse.xyz?hello=world", nil)

	expected := ""
	actual := PreviousFullUrlFromRequest(req)

	if expected != actual {
		t.Fatalf("PreviousFullUrlFromRequest should returns nothing page query argument is missing")
	}
}

func TestNextPreviousFullUrlFromRequest_WithPageQueryArg(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://big-data-warahouse.xyz?hello=world&page=12", nil)

	expected := "http://big-data-warahouse.xyz?hello=world&page=11#"
	actual := PreviousFullUrlFromRequest(req)

	if expected != actual {
		t.Fatalf("PreviousFullUrlFromRequest should returns original URL suffixed with # and original query params and page set to page-1")
	}
}
