package simplequery

import (
	"fmt"
	"net/url"
)

func ExampleSimpleUseCase() {
	urlQ, _ := url.ParseQuery("set")
	q := FromQuery(urlQ)

	fmt.Println("Flag `set` is", q.Get("set").Bool())
	fmt.Println("Flag `unset` is", q.Get("unset").Bool())

	// Output:
	// Flag `set` is true
	// Flag `unset` is false
}

func ExampleDefaultValues() {
	urlQ, _ := url.ParseQuery("name=vasya")
	q := FromQuery(urlQ)

	fmt.Println("Value of `name` is", q.Get("name").String())
	// Need the trailing `|` to keep the whitespaces at the end of the line
	// visible
	fmt.Println("Value of `middlename` is", q.Get("middlename").String(), "|")
	fmt.Println("Value of `lastname` is", q.Get("lastname").String("<unknown>"))
	fmt.Println("")

	// Output:
	// Value of `name` is vasya
	// Value of `middlename` is  |
	// Value of `lastname` is <unknown>
}

func ExampleParsingValues() {
	urlQ, _ := url.ParseQuery("updated_at=21")
	q := FromQuery(urlQ)

	fmt.Println("String value of `updated_at` is", q.Get("updated_at").String())
	fmt.Println("Int value of `updated_at` is", q.Get("updated_at").Int64())
	fmt.Println("Time value of `updated_at` is", q.Get("updated_at").Time())

	// Output:
	// String value of `updated_at` is 21
	// Int value of `updated_at` is 21
	// Time value of `updated_at` is 1970-01-01 00:00:21 +0000 UTC
}
