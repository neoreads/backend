package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/neoreads/backend/util"
)

func testExtractSentID() {
	var sentEndPat = regexp.MustCompile("^.*\\[\\#....\\]$")

	s := "abcd. [#abcd]"

	fmt.Printf("match: %v\n", sentEndPat.MatchString(s))

}

func testOlist() {

	var olistStartPat = regexp.MustCompile(`^\s*[0-9]+\. `)
	idx := olistStartPat.FindStringIndex(" 1. ABC")
	fmt.Printf("%v \n", idx)
}

func testSents() {
	t := "And there is an end. [#sdfs]I want to know. Do you know? If you think about it: it will be ok. {#abcd}"

	xs := util.SplitSents(t)

	fmt.Printf("\n==============================\n\n")
	for i := range xs {
		fmt.Printf("{%v}\n", xs[i])
	}

	gen := util.NewN64Generator(4)
	ts := util.ParseSents(t, gen)

	fmt.Printf("\n==============================\n\n")
	for i := range ts {
		fmt.Printf("{%#v}\n", ts[i])
	}

	testExtractSentID()

}
func testMain() {

	var paraEndPat = regexp.MustCompile("^.*\\{\\#....\\}$")
	log.Printf("match:%v\n", paraEndPat.MatchString("{#sdfd}"))
	/*
			s := `#Chapter Title

		This Chapter covers:

		- First section
		- Second section

		And there is an end. I want to know. If you think about it: it will be ok.
		This is the real end.

		Where is the money? I haven't seen any one.
		`
	*/

	/*
			s1 := `#Chapter Title [#xf7X] {#16VI}

		This Chapter covers: [#8A6L] {#yNoq}

		- First section [#blBy] {#fqEi}
		- Second section [#_HLU] {#STRU}

		And there is an end. [#rfBc]I want to know. [#FTZi]If you think about it: [#dia3]it will be ok. [#NjuA]This is the real end. [#_pl9] {#RnEh}

		Where is the money? [#8WDH]I haven't seen any one. [#3uZs] {#GAF2}

		~~~python
		print("hello")
		print("world")
		~~~
		`
	*/
	s1 := `
~~~python
def x(y):
	print (y)

	# this is comment
	return y
~~~
`
	md := util.ApplyIDs(s1)

	fmt.Printf("\n==============================\n\n")
	fmt.Printf("[%v]\n", md)

}

func main() {
	testMain()
}
