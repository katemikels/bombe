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
		charIdx = rotor.wiringInverse[charIdx]
		//charIdx = (charIdx + offsets[i]) % 26 // carried from enigma.py... I don't know what it's supposed to do
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

func newRunBombe(paths []string, inputLetter rune, menu map[rune]map[rune]int) map[string][]map[rune]rune {
	var rotors []rotorStruct
	rotorNames := []string{"I", "II", "III"}
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
					nextLetter := rune(path[i+1])

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

				fmt.Println("here: ", possibilities) // because it was mad at me....
			}

		}
	}
	return possibilities
}

func enigma(plaintext string, rotors []rotorStruct, initialRotors []rune, plugboard map[rune]rune) string {
	cipherText := ""

	// starting at initial rotor position
	rotorPosition := make([]rune, len(initialRotors))
	copy(rotorPosition, initialRotors)
	// step rotors until initial position is reached
	for string(rotorPosition) != string(initialRotors) {
		rotorPosition = stepRotors(rotors, rotorPosition)
	}

	// for all letters in the plain text
	for _, currentLetter := range plaintext {

		// through the plugboard
		plugCurrentLetter, ok := plugboard[currentLetter]

		if !ok {
			// if not in the plugboard -> mapped to itself
			plugCurrentLetter = currentLetter
		}

		// encrypt the character
		encryptedCurrentLetter := encryptChar(plugCurrentLetter, rotors, rotorPosition)

		// back through the plugboard
		encryptedPlugCurrentLetter, ok := plugboard[encryptedCurrentLetter]
		if !ok {
			encryptedPlugCurrentLetter = encryptedCurrentLetter
		}

		// add encrypted letter to the cipher text
		cipherText += string(encryptedPlugCurrentLetter)

	}
	return cipherText
}

func main() {

	// set up rotor structs for encryption
	var rotors []rotorStruct
	rotorNames := []string{"I", "II", "III"}
	for _, name := range rotorNames {
		var r rotorStruct
		r.name = name
		r.letters = rotorInfo[name][0]
		r.turnover = rotorInfo[name][1]
		makeWiring(&r)
		rotors = append(rotors, r)
	}
	// each line is one mapping out of the tem
	plugboard := map[rune]rune{
		//K/M, Z/G, P/D, A/C, X/O, E/R, B/N, J/H, I/S, T/F
		'K': 'M', 'M': 'K',
		'Z': 'G', 'G': 'Z',
		'P': 'D', 'D': 'P',
		'A': 'C', 'C': 'A',
		'X': 'O', 'O': 'X',
		'E': 'R', 'R': 'E',
		'B': 'N', 'N': 'B',
		'J': 'H', 'H': 'J',
		'I': 'S', 'S': 'I',
		'T': 'F', 'F': 'T',
	}

	// call our own enigma function to create our cipher text
	cipherText := enigma("DEARSWEETIEIMSOEXCITEDTOSPENDCHRISTMASDAYWITHYOUDOYOUHAVEIDEASFORWHATWESHOULDBUYLITTLEJIMMYMISSYOUXOXOBARBARA", rotors, []rune{'J', 'G', 'H'}, plugboard)

	crib := "CHRISTMASDAYWITHYOU"

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
		cipherCrib := cipherText[cribStart : cribEnd+1]

		menu := createMenu(crib, cipherCrib)

		// decide on paths
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
				fmt.Println("possibility: ", possibility)
			}
		}

		start = cribStart + 1
	}

	// print possible solution
	fmt.Println("Decrypted text: <DNE>")

	return
}
