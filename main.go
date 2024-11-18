package main

import "fmt"

// sketch of how to do the reflector.
// a switch statement would take a lot more code and writing,
// do we want to use runes? Otherwise we'd be doing a lot of
// conversion to strings on line 15 for alph[i]
func reflector(in rune) string {
	in8 := uint8(in)
	alph := "abcdefghijklmnopqrstuvwxyz"
	ukwb := "yruhqsldpxngokmiebfzcwvjat"

	for i := 0; i < 26; i++ {
		if alph[i] == in8 {
			return string(ukwb[i])
		}
	}
	return ""
}

// this is not the most efficient way, but it works
func findCrib(start int, text string, crib string) (int, int) {
	end := start + len(crib) - 1
	for end < len(text) {
		compareStr := text[start:end]
		locationPlausible := true
		for pos := range compareStr {
			if compareStr[pos] == crib[pos] {
				locationPlausible = false
				break // the crib and text can't have any letters in common
			}
		}
		if locationPlausible {
			return start, end
		} else {
			start += 1
			end += 1
		}
	}
	return -1, -1
}

func main() {
	// in the final product, cipherText and crib won't be pre-initialized,
	// remove the surrounding if block
	cipherText := "zjevjibowhpsvdupnvyyzlseqvgfkfxpqtxqoxhydaydprfgtnqxmcsayakszezmaxwpuoxtetffguvszkaikknfhdfgwopiisytteivnlyde"
	//var cipherText string
	if cipherText == "" {
		_, err := fmt.Scanln(&cipherText)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	crib := "christmasdaywithyou"
	// var crib string
	if crib == "" {
		_, err := fmt.Scanln(&cipherText)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	fmt.Println("Cipher Text: ", cipherText)
	fmt.Println("Crib: ", crib)

	plugboardFound := false
	start := 0
	for plugboardFound == false {
		// find the slice of cipher text that could work with the crib
		// unfortunately, in the hardcoded example above, almost every location is a valid crib
		// the actual solution is
		// start = 29, end = 40, but this returns 0 and 11 respectively
		cribStart, cribEnd := findCrib(start, cipherText, crib)
		fmt.Println(cribStart, cribEnd)

		// graph it up
		// on all rotator settings, set up a way to identify loops to determine the plugboard
		// only duplicate letters allowed are when it exists as a key once in each map
		// the arrangement is correct if the new key and val pair is not in the existing key and val sets respectively

		//cipherCribIntersection := make(map[rune]rune)  // plugboard keys set
		//rotatorCribIntersection := make(map[rune]rune) // plugboard vals set

		// if no errors come up and everything makes it into these two maps,
		// then we found the plugboard settings and the rotator arrangement at the same time
		// (previous lines will continue the loop before getting this far)
		plugboardFound = true
	}

	// go through the cipher text and decode using the plugboard/rotator settings/reflector
	// mark letters with ? that we don't know (don't have plugboard settings for those letters)

	// testing reflector
	fmt.Println("a", reflector('a')) // == y
	fmt.Println("h", reflector('h')) // d
	fmt.Println("i", reflector('i')) // p
	fmt.Println("z", reflector('z')) // t
	fmt.Println("p", reflector('p')) // i

	// print solution
	fmt.Println("Decrypted text: <DNE>")

	// manually fill in the letters that are still unsolved -- can't be done with code
	return
}
