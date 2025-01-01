# go-humanize-scale

A small Go library for **“scaling”** large numeric values into **human-readable** text, or **reverting** to a custom format to **keep an accurate result**. This library uses **[cockroachdb/apd](https://pkg.go.dev/github.com/cockroachdb/apd/v3)** for precise decimal arithmetic to avoid floating-point issues.

## Installation

```bash
go get github.com/dejurin/go-humanize-scale
```

### Why Another “Humanize” Library?

Many libraries might silently round your number to 1.23 million—even if you had leftover digits. This library detects any digits that would be lost by rounding to three decimals; if so, it calls a fallback function. That way, you can show an exact format (like 1,234,500) instead of an inaccurate 1.234 million.

### Example usage

```go
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
```


### Results

| Input      | Output          | Explanation                                                                                                                           |
|------------|-----------------|---------------------------------------------------------------------------------------------------------------------------------------|
| 1234000000 | 1.234 billion   | Divides by 1,000,000,000 → ratio 1.234. No leftover digits beyond 3 decimals, so "1.234 billion".                                      |
| 1200000    | 1.2 million     | Divides by 1,000,000 → ratio 1.2. No leftover digits, so "1.2 million".                                                               |
| 1230000    | 1.23 million    | Divides by 1,000,000 → ratio 1.23. Again, no leftover digits lost.                                                                    |
| 1234000    | 1.234 million   | ratio = 1.234 → exactly 3 decimal digits, so still accurate.                                                                           |
| 100000000  | 100 million     | ratio = 100 → no leftover digits.                                                                                                      |
| 1000000.1  | 1000000.1       | Fallback is not triggered for decimal inputs here, so we display as-is (the script does not forcibly convert floats if they parse).    |
| 10000000   | 10 million      | ratio = 10 → no leftover digits.                                                                                                       |
| 15000      | 15 thousand     | ratio = 15 on dividing by 1000, no decimal leftover → "15 thousand".                                                                   |
| 100000     | 100 thousand    | ratio = 100 → no leftover digits.                                                                                                      |
| 10000      | 10 thousand     | ratio = 10 → no leftover digits.                                                                                                       |
| 150000000  | 150 million     | ratio = 150 → no leftover digits.                                                                                                      |
| 15000000   | 15 million      | ratio = 15 → no leftover digits.                                                                                                       |
| 101000     | 101 thousand    | ratio = 101 → no leftover digits.                                                                                                      |
| 100100     | 100,100         | ratio = 100.1 has a decimal leftover, so we fallback, formatting with commas.                                                         |
| 2340000    | 2.34 million    | ratio = 2.34 → exactly 2 leftover digits (".34"), no further leftover.                                                                 |
| 11000000   | 11 million      | ratio = 11 → no leftover digits.                                                                                                       |
| 99500      | 99,500          | ratio = 99.5 → decimal leftover, fallback → comma format.                                                                              |
| 999999     | 999,999         | ratio = 999.999 → decimal leftover, fallback.                                                                                         |
| 9999       | 9,999           | ratio = 9.999 → leftover → fallback.                                                                                                   |
| 2300000    | 2.3 million     | ratio = 2.3 → leftover digits are captured (".3"). No further leftover, so "2.3 million".                                             |
| 2000000    | 2 million       | ratio = 2 → no leftover digits.                                                                                                        |
| 2345678    | 2,345,678       | ratio would be 2.345678 → leftover digits beyond 3 decimals, so fallback → "2,345,678".                                               |
| 12500      | 12,500          | ratio = 12.5 → leftover decimal, fallback.                                                                                            |
| 15100      | 15,100          | ratio = 15.1 → leftover decimal, fallback.                                                                                            |
| 100001     | 100,001         | ratio = 100.001 → leftover decimal, fallback.                                                                                         |
| 1234000001 | 1,234,000,001   | ratio ~ 1.234000001 → leftover decimals, fallback.                                                                                    |
| 1234500    | 1,234,500       | ratio = 1.2345… leftover → fallback.                                                                                                  |
| 1234500000 | 1,234,500,000   | leftover digits beyond 3 decimals → fallback.                                                                                        |
| 1234200    | 1,234,200       | leftover decimals → fallback.                                                                                                         |
| 1000       | 1,000           | ratio = 1, but checks indicate leftover in this context or scale usage → fallback.                                                   |
| 1          | 1               | Below min threshold or leftover handling → fallback.                                                                                  |
| 4000       | 4,000           | ratio = 4, but still ends up fallback in context.                                                                                    |
| 1234500    | 1,234,500       | Duplicate test of leftover scenario → fallback.                                                                                      |