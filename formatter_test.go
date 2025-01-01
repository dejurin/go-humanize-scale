package hscale

import (
	"strconv"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestHumanizeScale(t *testing.T) {
	scales := []Scale{
		{Value: "1000000000", Name: "billion"},
		{Value: "1000000", Name: "million"},
		{Value: "1000", Name: "thousand"},
	}

	fallback := func(n string) string {
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return n
		}
		p := message.NewPrinter(language.English)
		return p.Sprintf("%d", i)
	}

	tests := []struct {
		number   string
		expected string
	}{
		{"1234000000", "1.234 billion"},
		{"1200000", "1.2 million"},
		{"1230000", "1.23 million"},
		{"1234000", "1.234 million"},
		{"100000000", "100 million"},
		{"1000000.1", "1000000.1"},
		{"10000000", "10 million"},
		{"15000", "15 thousand"},
		{"100000", "100 thousand"},
		{"10000", "10 thousand"},
		{"150000000", "150 million"},
		{"15000000", "15 million"},
		{"101000", "101 thousand"},
		{"100100", "100,100"},
		{"2340000", "2.34 million"},
		{"11000000", "11 million"},
		{"99500", "99,500"},
		{"999999", "999,999"},
		{"9999", "9,999"},
		{"2300000", "2.3 million"},
		{"2000000", "2 million"},
		{"2345678", "2,345,678"},        // fallback
		{"12500", "12,500"},             // fallback
		{"15100", "15,100"},             // fallback
		{"100001", "100,001"},           // fallback
		{"1234000001", "1,234,000,001"}, // fallback
		{"1234500", "1,234,500"},        // fallback
		{"1234500000", "1,234,500,000"}, // fallback
		{"1234200", "1,234,200"},        // fallback
		{"1000", "1,000"},               // fallback
		{"1", "1"},                      // fallback
		{"4000", "4,000"},               // fallback
		{"1234500", "1,234,500"},        // fallback
	}

	for _, tt := range tests {
		res, err := Formatter(tt.number, "10000", scales, fallback)
		if err != nil {
			t.Errorf("number %q => unexpected error: %v", tt.number, err)
			continue
		}
		if res != tt.expected {
			t.Errorf("number %q => got %q, want %q", tt.number, res, tt.expected)
		}
	}
}

func TestIndianHumanizeScale(t *testing.T) {
	scales := []Scale{
		{Value: "1000000000", Name: "billion"},
		{Value: "10000000", Name: "crore"}, // Indian numbering
		{Value: "1000000", Name: "million"},
		{Value: "100000", Name: "lakh"}, // Indian numbering
		{Value: "1000", Name: "thousand"},
	}

	fallback := func(n string) string {
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return n
		}
		p := message.NewPrinter(language.English)
		return p.Sprintf("%d", i)
	}

	tests := []struct {
		number   string
		expected string
	}{
		{"1234000000", "1.234 billion"},
		{"1200000", "1.2 million"},
		{"1230000", "1.23 million"},
		{"1234000", "1.234 million"},
		{"100000000", "10 crore"},
		{"1000000.1", "1000000.1"},
		{"10000000", "1 crore"},
		{"15000", "15 thousand"},
		{"100000", "1 lakh"},
		{"10000", "10 thousand"},
		{"150000000", "15 crore"},
		{"15000000", "1.5 crore"},
		{"101000", "1.01 lakh"},
		{"100100", "100,100"},
		{"2340000", "2.34 million"},
		{"11000000", "1.1 crore"},
		{"99500", "99,500"},
		{"999999", "999,999"},
		{"9999", "9,999"},
		{"2300000", "2.3 million"},
		{"2000000", "2 million"},
		{"2345678", "2,345,678"},        // fallback
		{"12500", "12,500"},             // fallback
		{"15100", "15,100"},             // fallback
		{"100001", "100,001"},           // fallback
		{"1234000001", "1,234,000,001"}, // fallback
		{"1234500", "1,234,500"},        // fallback
		{"1234500000", "1,234,500,000"}, // fallback
		{"1234200", "1,234,200"},        // fallback
		{"1000", "1,000"},               // fallback
		{"1", "1"},                      // fallback
		{"4000", "4,000"},               // fallback
		{"1234500", "1,234,500"},        // fallback
	}

	for _, tt := range tests {
		res, err := Formatter(tt.number, "10000", scales, fallback)
		if err != nil {
			t.Errorf("number %q => unexpected error: %v", tt.number, err)
			continue
		}
		if res != tt.expected {
			t.Errorf("number %q => got %q, want %q", tt.number, res, tt.expected)
		}
	}
}
