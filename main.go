package main

import (
	"fmt"
	"slices"
)

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// TODO confirm these are the right turnovers??? From OtherBombCode/rotors.txt. Also assuming moving B->R is what causes turnover in rotor I
var rotorInfo = map[string][]string{
	"I":   {"EKMFLGDQVZNTOWYHXUSPAIBRCJ", "R"},
	"II":  {"AJDKSIRUXBLHWTMCQGZNPYFVOE", "F"},
	"III": {"BDFHJLCPRTXVZNYEIWGAKMUSQO", "W"},
	"IV":  {"ESOVPZJAYQUIRHXLNFTGKDCMWB", "K"},
	"V":   {"VZBRGITYUPSDNHLXAWMJQOFECK", "A"},
}

type rotorStruct struct {
	name          string
	letters       string
	turnover      string
	wiring        map[rune]rune
	wiringInverse map[rune]rune
}

// practically to get where the letter is on a rotator. "Where is `pos` in my `str`?
func index(pos rune, str string) rune {
	for i, c := range str {
		if c == pos {
			return rune(i)
		}
	}
	return -1
}

func makeWiring(r *rotorStruct) {
	wiring := make(map[rune]rune)
	wiringInverse := make(map[rune]rune)

	for i, letter := range r.letters {
		in := rune(i)
		out := index(letter, alphabet)
		wiring[in] = out
		wiringInverse[out] = in
	}
	r.wiring = wiring
	r.wiringInverse = wiringInverse
	return
}

// create an array of all possible rotor positions
func rotorPositions() []string {
	positions := make([]string, 0)
	for i := 0; i < len(alphabet); i++ {
		for j := 0; j < len(alphabet); j++ {
			for k := 0; k < len(alphabet); k++ {
				position := string(alphabet[i]) + string(alphabet[j]) + string(alphabet[k])
				positions = append(positions, position)
			}
		}
	}
	return positions
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
	return -1, -1 // no possible crib found
}

func createMenu(crib string, cipherCrib string) map[rune]map[rune]int {
	menu := make(map[rune]map[rune]int)
	//fmt.Println(cipherCrib)

	for i, cipherLetterInt := range cipherCrib {
		cipherLetter := cipherLetterInt
		plainLetter := rune(crib[i])

		// the cipher letter is not in the menu, add it linked to the plain letter
		if _, ok := menu[cipherLetter]; !ok {
			menu[cipherLetter] = map[rune]int{plainLetter: i}
		} else {
			// cipher letter is in the menu, and the plain letter is not linked, add it!
			if _, ok := menu[cipherLetter][plainLetter]; !ok {
				menu[cipherLetter][plainLetter] = i
			}
		}
		// the plain letter is not in the menu, link it to cipher letter
		if _, ok := menu[plainLetter]; !ok {
			menu[plainLetter] = map[rune]int{cipherLetter: i}
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

func stepRotors(rotors []rotorStruct, pos []rune) []rune {
	// rotors: list of 3 rotor structs, the left, middle, and right rotor names respectively
	// pos: the current letter visible on each rotor. Ex: "ABC" means left rotor is on A

	// use inverse to loop right rotor first (backwards from the written order)
	var rotorsInverse []rotorStruct
	for i := len(rotors) - 1; i >= 0; i-- {
		rotorsInverse = append(rotorsInverse, rotors[i])
	}

	var posInverse []rune
	for i := len(pos) - 1; i >= 0; i-- {
		posInverse = append(posInverse, pos[i])
	}

	// stepping rotors only advances the right rotor, unless it crosses a turning point
	for i, r := range rotorsInverse {
		currentPos := index(posInverse[i], r.letters)
		updatePos := (currentPos + 1) % 26
		posInverse[i] = rune(r.letters[updatePos])

		if rune(r.turnover[0]) != posInverse[i] {
			break
		}
	}
	var newPos []rune
	newPos = append(newPos, posInverse[2], posInverse[1], posInverse[0])
	return newPos
}

// this makes me cringe. but yes. it does work.
func encryptChar(char rune, rotors []rotorStruct, pos []rune) (rune, []rune) {
	// TODO this does not use the plugboard in the encryption

	// before a character goes through, the rotators step
	pos = stepRotors(rotors, pos)

	// use inverse to loop right rotor first (backwards from the written order)
	var rotorsInverse []rotorStruct
	for i := len(rotors) - 1; i >= 0; i-- {
		rotorsInverse = append(rotorsInverse, rotors[i])
	}

	var posInverse []rune
	for i := len(pos) - 1; i >= 0; i-- {
		posInverse = append(posInverse, pos[i])
	}

	// get the rotor offsets -- how many clicks forward has the rotor gone?
	// How much offset is the rotor starting at after the beginning of the string?
	var offsets []rune
	for i := 0; i < len(pos); i++ {
		offset := index(pos[i], rotors[i].letters)
		offsets = append(offsets, offset)
	}

	var offsetsInverse []rune
	for i := len(offsets) - 1; i >= 0; i-- {
		offsetsInverse = append(offsetsInverse, offsets[i])
	}

	// value of the char as an offset, not ASCII value
	charIdx := char - 65

	// through rotors right -> left
	for i, rotor := range rotorsInverse {
		// I gotta be honest, I copied this idea straight from enigma.py
		// I never would have thought it out like this without it.
		charIdx = (charIdx + offsetsInverse[i]) % 26 // adjust for rotor offset
		charIdx = rotor.wiring[charIdx]              // wiring gives us the alphabet idx of the char on the rotor
		//charIdx = (charIdx - offsetsInverse[i]) % 26 // why does enigma.py have this line??
	}

	// reflector
	ukwb := "YRUHQSLDPXNGOKMIEBFZCWVJAT"     // reflector string
	reflectorChar := rune(ukwb[charIdx])     // the character in the reflector at our index
	charIdx = index(reflectorChar, alphabet) // where that character is in the alphabet (the wiring)

	// opposite direction through the rotors
	for i, rotor := range rotors {
		charIdx = (charIdx + offsets[i]) % 26 // adjust for rotor offset
		charIdx = rotor.wiring[charIdx]
		//charIdx = (charIdx - offsets[i]) % 26 // back out of rotor offset... why???
	}

	char = rune(alphabet[charIdx]) // back to ASCII
	return char, pos
}

// DFS to find paths in the menu
func searchForPaths(letter rune, menu map[rune]map[rune]int, current rune, path []rune, paths *[]string) {
	// for all the letters associated with the current letter
	// TODO rename `new` bc new is a keyword
	for new := range menu[current] {
		// if the new letter is the same as the start (and the path is longer than two) then we found a loop
		if new == letter && len(path) > 2 {
			// copy
			pathCopy := make([]rune, len(path))
			copy(pathCopy, path)
			pathCopy = append(pathCopy, letter)
			// add to list
			*paths = append(*paths, string(pathCopy))
			continue
		}

		// if already in the path and not start, move on - not helpful
		if slices.Contains(path, new) {
			continue
		}

		// copy
		pathCopy := make([]rune, len(path))
		copy(pathCopy, path)
		pathCopy = append(pathCopy, new) // Changed from letter to new

		// recursive call
		searchForPaths(letter, menu, new, pathCopy, paths)
	}
}

func runBombe(paths []string, inputLetter string, menu map[rune]map[rune]int) {
	// hardcoded rotors?
	var rotors []rotorStruct
	rotorNames := []string{"I", "IV", "III"}
	for _, name := range rotorNames {
		var r rotorStruct
		r.name = name
		r.letters = rotorInfo[name][0]
		r.turnover = rotorInfo[name][1]
		makeWiring(&r)
		rotors = append(rotors, r)
	}

	// - check all rotator positions (which ones, what order, starting order)
	// this does not currently consider the different possible rotors
	rotorPositions := rotorPositions()

	for _, rotor := range rotorPositions {
		// - for all guesses (the alphabet) with input letter
		for _, guess := range alphabet {
			// - for all paths, break at contradictions (remember them for shortcuts later?)
			for _, path := range paths {
				// - send through enigma rotators/reflector -> write the rotators to step through to the index

				// example code for stepping rotors
				var position = []rune{'B', 'A', 'I'}
				fmt.Println(position)
				//newPosition := stepRotors(rotors, position)
				//fmt.Println(newPosition) // to make newPosition used

				// example code for encrypting char -- I thiiiiiiiinnkkkk the logic is right?? I'm not really sure. But if my understanding is right, then it works.
				char := 'C'
				char, position = encryptChar(char, rotors, position)

				// - no contradictions, is a possibility (check other paths?)
				// add all possibilities to list of all possibilities for all possible cribs

				// keep things happy :)
				fmt.Println(rotor)
				fmt.Println(guess)
				fmt.Println(path)
			}
		}
	}
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

	start := 0 //solution: start = 29, end = 47 (inclusive)
	for start+len(crib) < len(cipherText) {
		// find a possible crib and create the corresponding menu
		cribStart, cribEnd := findCrib(start, cipherText, crib)
		if cribStart == -1 || cribEnd == -1 {
			// no cribs left -- stop and go through all the possibilities collected
			break
		}
		cipherCrib := cipherText[cribStart : cribEnd+1] // +1 because findCrib() returned inclusive values?

		menu := createMenu(crib, cipherCrib)
		//fmt.Println(menu) // so that the go compiler doesn't get mad at us not using menu

		// - decide on paths
		var paths []string
		for letter := range menu {
			path := []rune{letter}
			searchForPaths(letter, menu, letter, path, &paths)
		}

		fmt.Println("paths: ", paths)

		// paths found
		if len(paths) != 0 {
			// - decide on input letter (start of loop path)
			//runBombe(paths, inputLetter, menu) // update as parameters change and given output
			runBombe(paths, "A", menu)
		}

		// after the loop, checkCipherText() to check our possibilities against the entire message
		// prints out decrypted ciphertext
		// command line utility????

		// if no errors come up and everything makes it into these two maps,
		// then we found the plugboard settings and the rotator arrangement at the same time
		// (previous lines will continue the loop before getting this far)
		start = cribStart + 1
		break
	}

	// go through the cipher text and decode using the plugboard/rotator settings/reflector
	// mark letters with ? that we don't know (don't have plugboard settings for those letters)

	// print possible solution
	fmt.Println("Decrypted text: <DNE>")

	return
}
