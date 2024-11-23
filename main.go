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

func makeMenu(crib string, cipherCrib string) map[string]map[string]int {
	menu := make(map[string]map[string]int)
	fmt.Println(cipherCrib)

	for i, cipherLetterInt := range cipherCrib {
		cipherLetter := string(uint8(cipherLetterInt))
		plainLetter := string(crib[i])

		// the cipher letter is not in the menu, add it linked to the plain letter
		if _, ok := menu[cipherLetter]; !ok {
			menu[cipherLetter] = map[string]int{plainLetter: i}
		} else {
			// cipher letter is in the menu, and the plain letter is not linked, add it!
			if _, ok := menu[cipherLetter][plainLetter]; !ok {
				menu[cipherLetter][plainLetter] = i
			}
		}
		// the plain letter is not in the menu, link it to cipher letter
		if _, ok := menu[plainLetter]; !ok {
			menu[plainLetter] = map[string]int{cipherLetter: i}
		} else {
			// plain letter is in the menu, and the cipher letter is not linked, add it!
			if _, ok := menu[plainLetter][cipherLetter]; !ok {
				menu[plainLetter][cipherLetter] = i
			}
		}

		// if the combo is already in the menu, we don't add anything or change the index vals
	}
	return menu
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

	start := 0
	for start + len(crib) - 1 < len(cipherText) {
		// find the slice of cipher text that could work with the crib
		// unfortunately, in the hardcoded example above, almost every location is a valid crib
		// the actual solution is
		// start = 29, end = 40, but this returns 0 and 11 respectively
		cribStart, cribEnd := findCrib(start, cipherText, crib)
		cipherCrib := cipherText[cribStart : cribEnd+1] // +1 because findCrib() returned inclusive values?
		menu := makeMenu(crib, cipherCrib)
		fmt.Println(menu) // so that the go compiler doesn't get mad at us not using menu

		// call function -> createMenu()
		// - decide on paths
		// - decide on input letter (start of loop path)
		// runBombe() - returns possible plugboards with associated rotator settings
		// - check all rotator positions (which ones, what order, starting order)
		// - for all guesses (the alphabet) with input letter
		// - for all paths, break at contradictions (remember them for shortcuts later?)
		// - send through enigma rotators/reflector -> write the rotators to step through to the index
		// - no contradictions, is a possibility (check other paths?)
		// add all possibilities to list of all possibilities for all possible cribs
		// after the loop, checkCipherText() to check our possibilities against the entire message
		// prints out decrypted ciphertext
		// command line utility????

		// if no errors come up and everything makes it into these two maps,
		// then we found the plugboard settings and the rotator arrangement at the same time
		// (previous lines will continue the loop before getting this far)
		start = cribStart + 1
	}

	// go through the cipher text and decode using the plugboard/rotator settings/reflector
	// mark letters with ? that we don't know (don't have plugboard settings for those letters)

	// print possible solution
	fmt.Println("Decrypted text: <DNE>")

	return
}
