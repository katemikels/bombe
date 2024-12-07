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

func contains(letters []rune, guess rune) bool {
	for _, l := range letters {
		if l == guess {
			return true
		}
	}
	return false
}

func containsInMap(letters map[rune]rune, guess rune) bool {
	for k, _ := range letters {
		if k == guess {
			return true
		}
	}
	return false
}

func containsInList(letters []rune, guess rune) bool {
	for _, l := range letters {
		if l == guess {
			return true
		}
	}
	return false
}

func containsContradictions(letters map[rune][]rune, guess rune) bool {
	for k, _ := range letters {
		if k == guess {
			return true
		}
	}
	return false
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

func addToContradictions(plugboard map[rune]rune, contradictions map[rune][]rune) map[rune][]rune {
	for k, _ := range plugboard {
		if _, ok := contradictions[k]; ok {
			contradictions[k] = append(contradictions[k], plugboard[k])
		} else {
			contradictions[k] = []rune{k}
		}
	}
	return contradictions
}

func addToPossibilities(rotorPos string, plugboard map[rune]rune, possibilities map[string][]map[rune]rune) map[string][]map[rune]rune {
	// if not in possibilites, add it
	if _, ok := possibilities[rotorPos]; ok {
		possibilities[rotorPos] = make([]map[rune]rune, len(plugboard))
	}
	// loop through plugboard and add it to possibilites at rotorPos
	for k, v := range plugboard {
		possibilities[rotorPos] = append(possibilities[rotorPos], map[rune]rune{k: v})
	}
	return possibilities
}

/*
	func runBombe(paths []string, inputLetter rune, menu map[rune]map[rune]int) []map[rune]rune {
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
		rotorPositionsList := rotorPositions()
		var possibilities []map[rune]rune // list of possible plugboard solutions

		for _, rotorPosition := range rotorPositionsList {
			// - for all guesses (the alphabet) with input letter
			var impossibilities map[rune][]rune

			for _, guess := range alphabet {
				// - for all paths, break at contradictions (remember them for shortcuts later?)
				// set up first plugboard pair guess
				plugboard := make(map[rune]rune)
				plugboard[inputLetter] = guess
				plugboard[guess] = inputLetter
				contradiction := false
				if _, ok := impossibilities[guess]; ok { // if the letter is in impossibilities, don't even think about it
					if contains(impossibilities[inputLetter], guess) {
						contradiction = true
						break
					}
				}

				for _, path := range paths {
					if contradiction {
						break
					}

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
									impossibilities[cribLetter] = append(impossibilities[cribLetter], guess)
									impossibilities[guess] = append(impossibilities[plugOut], cribLetter)
									break // Contradiction!! -- give up on this path... on this guess??
								}
							} else { // assume that we plug to the correct next location
								// letter --plug--> guess --rotors--> rotorOut --plug--> cribLetter... I think?
								if _, plugExists := plugboard[cribLetter]; plugExists {
									fmt.Println("CONTRADITCTION:", letter, "plugs to", plugOut, "and", cribLetter)
									contradiction = true
									impossibilities[cribLetter] = append(impossibilities[cribLetter], guess)
									impossibilities[guess] = append(impossibilities[plugOut], cribLetter)
									break
								} else {
									if impossiblePairings, outExists := impossibilities[cribLetter]; outExists {
										if contains(impossiblePairings, rotorOut) {
											fmt.Println("CONTRADITCTION:", letter, "plugs to", plugOut, "and", cribLetter)
											contradiction = true
											impossibilities[cribLetter] = append(impossibilities[cribLetter], guess)
											impossibilities[guess] = append(impossibilities[plugOut], cribLetter)
											break
										}
									}
								}
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
		}
		return possibilities
	}
*/
func newRunBombe(paths []string, inputLetter rune, menu map[rune]map[rune]int) map[string][]map[rune]rune {
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
	rotorPositionsList := rotorPositions()
	// from []map[rune]rune
	// rotorPositions (as a string) : {possible plugboard}, {another possible plugboard}
	var possibilities map[string][]map[rune]rune // list of possible plugboard solutions

	// for all rotor positions
	for _, rotorPosition := range rotorPositionsList {

		// map of contradictions for the rotor position
		contradictions := make(map[rune][]rune)

		// for all possible guesses in the alphabet
		for _, guess := range alphabet {

			// set up first plugboard pair guess
			plugboard := make(map[rune]rune)
			plugboard[inputLetter] = guess
			plugboard[guess] = inputLetter
			plugboardPossible := true // the plugboard is currently possible;
			// if not possible, give up on guess

			// for all paths
			for _, path := range paths {
				// for each letter in paths
				for i := 0; i < len(path)-1; i++ {
					// variables :)
					currentLetter := rune(path[i])
					//currentLetterConnections := menu[currentLetter]
					nextLetter := rune(path[i+1])
					//nextLetterConnections := menu[nextLetter]
					stepsToRotorPos := menu[currentLetter][nextLetter]

					// step rotors
					for j := 0; j < stepsToRotorPos; j++ {
						//fmt.Println(rotorPosition)
						rotorPosition = stepRotors(rotors, rotorPosition)
						//fmt.Println(rotorPosition)
					}

					// if currentLetter in plugboard, encrypt the plugboard of current letter
					if containsInMap(plugboard, currentLetter) && plugboardPossible {
						plugCurrentLetter := plugboard[currentLetter]
						encryptedPlugCurrentLetter := encryptChar(plugCurrentLetter, rotors, rotorPosition)

						// contradiction checking
						if containsContradictions(contradictions, nextLetter) {
							// if next letter is in our contradictions list, pull out the contradictions
							nextLetterContradictions := contradictions[nextLetter]
							if containsInList(nextLetterContradictions, encryptedPlugCurrentLetter) {
								// adds the whole plugboard to contradictions
								contradictions = addToContradictions(plugboard, contradictions)
								plugboardPossible = false
								// give up on this
								break
							}
						}

						if containsInMap(plugboard, nextLetter) {
							if plugboard[nextLetter] == encryptedPlugCurrentLetter {
								// found a loop that works - check other possible paths
								break
							} else { // otherwise there was some sort of contradiction
								// add to plugboard
								plugboard[nextLetter] = encryptedPlugCurrentLetter
								plugboard[encryptedPlugCurrentLetter] = nextLetter
								// add the whole plugboard to contradictions
								contradictions = addToContradictions(plugboard, contradictions)
								plugboardPossible = false
								// give up on this
								break
							}
						} else { // if not in plugboard, add
							plugboard[nextLetter] = encryptedPlugCurrentLetter
							plugboard[encryptedPlugCurrentLetter] = nextLetter
						}
					}
				}
				// if we have found that the plugboard no longer works, don't keep going
				if plugboardPossible == false {
					break
				}
			}
			// if the plugboard is still possible, add to possibilities
			// TODO: check if this is correct
			if plugboardPossible {
				possibilities := addToPossibilities(string(rotorPosition), plugboard, possibilities)
				fmt.Println(possibilities) // because it was mad at me....
			}
		}
	}
	return possibilities
}

/*
Run Bombe Pseudocode:

for all rotor combos:
	for all rotor positions:
		create contradictions
		for all guesses in alphabet
			create plugboard
			plugboardPossible = true
			for all paths
				for all letters in paths
					// variables
					currentLetter = path[i]
					currentLetterConnections = menu[currentLetter]
					nextLetter = path[i+1]
					nextLetterConnections = menu[nextLetter]
					stepsToRotorPos = offset from currentLetter to nextLetter in menu

					stepRotors
					encryptedLetter = encrypt currentLetter

					// already contradictions
					if nextLetter is in contradictions
						if encryptedLetter is in contradictions
							contradiction: add the contradiction and the entire plugboard to the contradictions list
							plugboardPossible = false
							break

					// no contradictions
					if nextLetter is in the plugboard
						if plugboard[nextLetter] == encryptedLetter
							path completed - no contradictions found, is possible`5
						else (must be a contradiction)		// found new one
							contradiction: add the contradiction and the entire plugboard to the contradictions list
							plugboardPossible = false
					else (not in plugboard)
						add to plugboard

				if not plugboardPossible (in for each letter in path loop)
					break out of path loop because it isn't going to work

			if plugboardPossible (in for all paths loop)
				add to list of all possible plugboards

outside of loops, should have list of all possible plugboards

*/

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
			//possibilities := runBombe(paths, inputLetter, menu) // update as parameters change and given output
			possibilities := newRunBombe(paths, inputLetter, menu)

			for possibility := range possibilities {
				fmt.Println(possibility)
			}
		}

		// TODO: test possibilities to see if they work
		/*
			for each possibility, go through the characters in the cipher text, building up a decrypted string where we replace the letter we found
			if we don't have the answer, then put a ? in the string
		*/
		// TODO: test code in general....
		// TODO: clean code

		start = cribStart + 1
	}

	// print possible solution
	fmt.Println("Decrypted text: <DNE>")

	return
}
