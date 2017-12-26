// passdrill: typing drills for practicing passphrases

package main

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

const hashAlgorithm = "sha512"
const hashFilename = "passdrill." + hashAlgorithm
const help = "Use -s to save passphrase hash for practice."

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func input(msg string) string {
	response := ""
	fmt.Print(msg)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		response = scanner.Text()
	}
	check(scanner.Err())
	return response
}

func prompt() string {
	fmt.Println("WARNING: the passphrase will be shown so that you can check it!")
	confirmed := ""
	passwd := ""
	for confirmed != "y" {
		passwd = input("Type passphrase to hash (it will be echoed): ")
		if passwd == "" || passwd == "q" {
			fmt.Println("ERROR: the passphrase cannot be empty or 'q'.")
			continue
		}
		fmt.Println("Passphrase to be hashed ->", passwd)
		confirmed = strings.ToLower(input("Confirm (y/n): "))
	}
	return passwd
}

func hashBytes(text string) []byte {
	octets := sha512.Sum512([]byte(text))
	return []byte(base64.StdEncoding.EncodeToString(octets[:]))
}

func saveHash(args []string) {
	if len(os.Args) > 2 || os.Args[1] != "-s" {
		fmt.Println("ERROR: invalid argument.", help)
		os.Exit(1)
	}
	passwdHash := hashBytes(prompt())
	err := ioutil.WriteFile(hashFilename, passwdHash, 0644)
	check(err)
	fmt.Printf("Passphrase %s hash saved to %s\n",
		hashAlgorithm, hashFilename)
}

func practice() {
	passwdHash, err := ioutil.ReadFile(hashFilename)
	if os.IsNotExist(err) {
		fmt.Println("ERROR: passphrase hash file not found.", help)
		os.Exit(1)
	}
	check(err)
	fmt.Println("Type q to end practice.")
	turn := 0
	response := ""
	correct := 0
	for {
		turn++
		fmt.Printf("%d:", turn)
		octets, err := gopass.GetPasswd()
		check(err)
		response = string(octets)
		if response == "" {
			fmt.Println("Type q to quit.")
			turn-- // don't count this response
			continue
		} else if response == "q" {
			turn-- // don't count this response
			break
		}
		answer := "wrong"
		if bytes.Compare(hashBytes(response), passwdHash) == 0 {
			correct++
			answer = "OK"
		}
		fmt.Printf("  %s\thits=%d\tmisses=%d\n", answer, correct, turn-correct)
	}
	if turn > 0 {
		pct := float32(correct) / float32(turn) * 100
		fmt.Printf("\n%d exercises. %0.1f%% correct.\n", turn, pct)
	}

}

func main() {
	if len(os.Args) > 1 {
		saveHash(os.Args)
	} else {
		practice()
	}
}
