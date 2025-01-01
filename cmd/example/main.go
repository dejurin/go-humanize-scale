package main

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	hscale "github.com/dejurin/go-humanize-scale"
)

func main() {
	scales := []hscale.Scale{
		{Value: "1000000000", Name: "billion"},
		{Value: "10000000", Name: "crore"}, // Indian numbering
		{Value: "1000000", Name: "million"},
		{Value: "1000", Name: "thousand"},
	}

	fallbackFunc := func(n string) string {
		// Here we use "golang.org/x/text" for comma formatting:
		p := message.NewPrinter(language.English)
		return p.Sprintf("%s", n)
	}

	// Example usage:
	result, err := hscale.Formatter("100001", "10000", scales, fallbackFunc)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("1230000 =>", result) // "1.23 million"
	}

	result2, err := hscale.Formatter("1234500", "10000", scales, fallbackFunc)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("1234500 =>", result2) // "1,234,500" (fallback)
	}
}
