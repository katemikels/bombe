package main

import (
	"fmt"
)

// bc go doesn't have a queue implementation. Taken from https://www.geeksforgeeks.org/queue-in-go-language/
func enqueue(queue []rune, element rune) []rune {
	queue = append(queue, element) // Simply append to enqueue.
	fmt.Println("Enqueued:", element)
	return queue
}

func dequeue(queue []rune) (rune, []rune) {
	element := queue[0] // The first element is the one to be dequeued.
	if len(queue) == 1 {
		var tmp = []rune{}
		return element, tmp

	}

	return element, queue[1:] // Slice off the element once it is dequeued.
}

// sketch of how to do the reflector.
// a switch statement would take a lot more code and writing,
// do we want to use runes? Otherwise we'd be doing a lot of
// conversion to strings on line 15 for alph[i]
func reflector(in rune) rune {
	alph := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ukwb := "YRUHQSLDPXNGOKMIEBFZCWVJAT"

	for i := 0; i < 26; i++ {
		if rune(alph[i]) == in {
			return rune(ukwb[i])
		}
	}
	return '0' // not found
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
	fmt.Println(cipherCrib)

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

func findPaths(menu map[rune]map[rune]int) [][]rune {
	var paths [][]rune
	//var pathsByRune = make(map[rune][][]rune) // so we're not enqueuing/dequeuing on messy data structures?
	//queue := make([]rune, 0)
	//firstElement := true
	//for k, v := range menu {
	//	if firstElement { // initialize the queue
	//		firstElement = false
	//		enqueue(queue, k)
	//	} else {
	//
	//	}
	//}
	return paths
}

// TODO if you can reverse a string in go, you can turn this into a loop.
// using indices would make for shorter code? But would also be less
// readable. This code (I think?) is pretty clear in what it does.
func stepRotors(rotors []string, pos []rune) []rune {
	// rotors: list of 3 strings, the left, middle, and right rotor names respectively
	// pos: the current letter visible on each rotor. Ex: "ABC" means left rotor is on A

	// TODO confirm these are the right turnovers??? From OtherBombCode/rotors.txt. Also assuming moving B->R is what causes turnover in rotor I
	rotorInfo := map[string][]string{
		"I":   {"EKMFLGDQVZNTOWYHXUSPAIBRCJ", "R"},
		"II":  {"AJDKSIRUXBLHWTMCQGZNPYFVOE", "F"},
		"III": {"BDFHJLCPRTXVZNYEIWGAKMUSQO", "W"},
		"IV":  {"ESOVPZJAYQUIRHXLNFTGKDCMWB", "K"},
		"V":   {"VZBRGITYUPSDNHLXAWMJQOFECK", "A"},
	}

	// friendly names!
	leftPos := pos[0]
	middlePos := pos[1]
	rightPos := pos[2]

	leftRotor := rotors[0]
	middleRotor := rotors[1]
	rightRotor := rotors[2]

	// stepping rotors only advances the right rotor, unless it crosses a turning point
	for i, c := range rotorInfo[rightRotor][0] {
		if c == rightPos {
			rightPos = rune(rotorInfo[rightRotor][0][(i+1)%26])
			break
		}
	}
	// if right hit the turnover point
	if rightPos == rune(rotorInfo[rightRotor][1][0]) {
		for i, c := range rotorInfo[middleRotor][0] {
			if c == middlePos {
				middlePos = rune(rotorInfo[middleRotor][0][(i+1)%26])
				break
			}
		}
	}
	// if middle hit a turnover point
	if middlePos == rune(rotorInfo[middleRotor][1][0]) {
		for i, c := range rotorInfo[leftRotor][0] {
			if c == leftPos {
				leftPos = rune(rotorInfo[leftRotor][0][(i+1)%26])
				break
			}
		}
	}
	var newPos []rune
	newPos = append(newPos, leftPos, middlePos, rightPos)
	return newPos
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

// this makes me cringe. but yes. it does work.
func encryptChar(char rune, rotors []string, pos []rune) (rune, []rune) {
	// TODO this does not use the plugboard in the encryption

	// before a character goes through, the rotators step
	pos = stepRotors(rotors, pos)

	//charIdx := char - 65

	return char, pos

	//// friendly names!
	//leftRotor := rotors[0]
	//middleRotor := rotors[1]
	//rightRotor := rotors[2]
	//
	//leftPos := index(pos[0], rotorInfo[leftRotor][0])
	//middlePos := index(pos[1], rotorInfo[middleRotor][0])
	//rightPos := index(pos[2], rotorInfo[rightRotor][0])
	//
	//charIdx := char - 65
	//// right rotor
	////char = rune(rotorInfo[rightRotor][0][(charIdx+rightPos)%26])
	////charIdx = index(char, rotorInfo[rightRotor][0])
	////fmt.Printf("after right: %c\n", char)
	//charIdx = (charIdx + rightPos) % 26
	//
	//// middle rotor
	////char = rune(rotorInfo[middleRotor][0][(charIdx+middlePos)%26])
	////charIdx = index(char, rotorInfo[middleRotor][0])
	////fmt.Printf("after middle: %c\n", char)
	//charIdx = (charIdx + middlePos) % 26
	//
	//// left rotor
	////char = rune(rotorInfo[leftRotor][0][(charIdx+leftPos)%26])
	////charIdx = index(char, rotorInfo[leftRotor][0])
	////fmt.Printf("after left: %c\n", char)
	//charIdx = (charIdx + leftPos) % 26
	//
	//// reflector
	//char = reflector(char)
	//charIdx = index(char, "YRUHQSLDPXNGOKMIEBFZCWVJAT") // hardcoded for the reflector
	////fmt.Printf("reflector: %c\n", char)
	//
	//// left rotor
	////char = rune(rotorInfo[leftRotor][0][(charIdx+leftPos)%26])
	////charIdx = index(char, rotorInfo[leftRotor][0])
	////fmt.Printf("left: %c\n", char)
	//charIdx = (charIdx + leftPos) % 26
	//
	//// middle rotor
	////char = rune(rotorInfo[middleRotor][0][(charIdx+middlePos)%26])
	////charIdx = index(char, rotorInfo[middleRotor][0])
	////fmt.Printf("middle: %c\n", char)
	//charIdx = (charIdx + middlePos) % 26
	//
	//// right rotor
	////char = rune(rotorInfo[rightRotor][0][(charIdx+rightPos)%26])
	////fmt.Printf("right: %c\n", char)
	//charIdx = (charIdx + rightPos) % 26
	//char = charIdx + 65
	//fmt.Printf("final: %c\n", char)
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

		// TODO find paths -- this isn't implemented yet
		paths := findPaths(menu)
		fmt.Println(paths)

		// runBombe() - returns possible plugboards with associated rotator settings
		// - check all rotator positions (which ones, what order, starting order)

		// example code for stepping rotors. Works!
		var position = []rune{'B', 'A', 'J'}
		fmt.Println(position)
		newPosition := stepRotors([]string{"I", "IV", "III"}, position)
		fmt.Println(newPosition) // to make newPosition used

		// example code for encrypting char
		char := 'C'
		char, position = encryptChar(char, []string{"I", "IV", "III"}, position)

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
		break
	}

	// go through the cipher text and decode using the plugboard/rotator settings/reflector
	// mark letters with ? that we don't know (don't have plugboard settings for those letters)

	// print possible solution
	fmt.Println("Decrypted text: <DNE>")

	return
}
