package main

import (
	"fmt"
	"slices"
)

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// TODO confirm these are the right turnovers??? From OtherBombCode/rotors.txt.
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

// Go doesn't have a function for this
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
func rotorPositions() [][]rune {
	positions := make([][]rune, 0)
	for i := 0; i < len(alphabet); i++ {
		for j := 0; j < len(alphabet); j++ {
			for k := 0; k < len(alphabet); k++ {
				position := []rune{rune(alphabet[i]), rune(alphabet[j]), rune(alphabet[k])}
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
	// pos: the current letter offset on each rotor. Ex: "ABC" means left rotor will start on A

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
func encryptChar(char rune, rotors []rotorStruct, pos []rune) rune {
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
		// I gotta be honest, I took this idea from enigma.py
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
		//charIdx = (charIdx - offsets[i]) % 26 // carried from enigma.py... I don't know what it's supposed to do
	}

	char = rune(alphabet[charIdx]) // back to ASCII
	return char
}

// DFS to find paths in the menu
func searchForPaths(letter rune, menu map[rune]map[rune]int, current rune, path []rune, paths *[]string) {
	// for all the letters associated with the current letter
	for newLetter := range menu[current] {
		// if the new letter is the same as the start (and the path is longer than two) then we found a loop
		if newLetter == letter && len(path) > 2 {
			// copy
			pathCopy := make([]rune, len(path))
			copy(pathCopy, path)
			pathCopy = append(pathCopy, letter)
			// add to list
			*paths = append(*paths, string(pathCopy))
			continue
		}

		// if already in the path and not start, move on - not helpful
		if slices.Contains(path, newLetter) {
			continue
		}

		// copy
		pathCopy := make([]rune, len(path))
		copy(pathCopy, path)
		pathCopy = append(pathCopy, newLetter)

		// recursive call
		searchForPaths(letter, menu, newLetter, pathCopy, paths)
	}
}

func runBombe(paths []string, inputLetter rune, menu map[rune]map[rune]int) {
	var rotors []rotorStruct
	rotorNames := []string{"I", "IV", "III"} // TODO hardcoded rotors?
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
	rotorPositionsList := rotorPositions()
	var possibilities []map[rune]rune // list of possible plugboard solutions

	for _, rotorPosition := range rotorPositionsList {
		// - for all guesses (the alphabet) with input letter
		for _, guess := range alphabet {
			// - for all paths, break at contradictions (remember them for shortcuts later?)
			// set up first plugboard pair guess
			plugboard := make(map[rune]rune)
			plugboard[inputLetter] = guess
			plugboard[guess] = inputLetter

			contradiction := false
			for _, path := range paths {
				// - no contradictions, is a possibility (check other paths?)
				// add all possibilities to list of all possibilities for all possible cribs

				// for each letter in the path
				for i := 0; i < len(path)-1; i++ {
					letter := rune(path[i])
					cribLetter := rune(path[i+1])
					cribPosition := menu[letter][cribLetter]
					fmt.Println("crib position: ", cribPosition)

					// - send through enigma rotators/reflector -> write the rotators to step through to the index
					// step rotators to cribPosition
					for j := 0; j < cribPosition; j++ {
						// example code for stepping rotors
						fmt.Println(rotorPosition)
						rotorPosition = stepRotors(rotors, rotorPosition)
						fmt.Println(rotorPosition)
					}

					// check to see if the letter we are looking at is in the plugboard, else we can't do anything?
					rotorIn, exists := plugboard[letter]
					//var rotorInput rune
					if exists {
						fmt.Println("from plugboard, what is going into rotors: ", rotorIn)
						// I thiiiiiiiinnkkkk the logic is right?? I'm not really sure.
						rotorOut := encryptChar(rotorIn, rotors, rotorPosition)
						fmt.Println("letter from rotors for plugboard: ", rotorOut)

						if plugOut, ok := plugboard[rotorOut]; ok { // `exists` was giving multiple declaration warnings
							fmt.Println("encryption result: ", plugOut)
							// check for contradictions
							if plugOut != cribLetter {
								fmt.Println("CONTRADITCTION:", letter, "plugs to", plugOut, "and", cribLetter)
								contradiction = true
								break // Contradiction!! -- give up on this path... on this rotor position?? // TODO change how this breaks?
							}
						} else { // assume that we plug to the correct next location
							// letter --plug--> guess --rotors--> rotorOut --plug--> cribLetter... I think?
							plugboard[rotorOut] = cribLetter
							plugboard[cribLetter] = rotorOut
							fmt.Println("encryption result: ", cribLetter)
						}
					}
				}
			}
			// add all possibilities to list of all possibilities for all possible cribs
			if !contradiction {
				possibilities = append(possibilities, plugboard)
			}
		}
	} // return possibilities?
}

func main() {
	// in the final product, cipherText and crib won't be pre-initialized,
	// remove the surrounding if block
	cipherText := "ZJEVJIBOWHPSVDUPNVYYZLSEQVGFKFXPQTXQOXHYDAYDPRFGTNQXMCSAYAKSZEZMAXWPUOXTETFFGUVSZKAIKKNFHDFGWOPIISYTTEIVNLYDE"
	//var cipherText string
	if cipherText == "" {
		_, err := fmt.Scanln(&cipherText)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	crib := "CHRISTMASDAYWITHYOU"
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

		// paths found
		if len(paths) != 0 {
			fmt.Println("paths: ", paths)
			fmt.Println(paths[0][0]) // first letter of first path
			inputLetter := rune(paths[0][0])
			runBombe(paths, inputLetter, menu) // update as parameters change and given output
		}

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
