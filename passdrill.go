// passdrill: typing drills for practicing passphrases

package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"golang.org/x/crypto/pbkdf2"
)

const hashFilename = "passdrill.hash"
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

func myPbkdf2(salt, content []byte) []byte {
	algorithm := sha512.New
	rounds := 200000
	keyLen := 64
	return pbkdf2.Key(content, salt, rounds, keyLen, algorithm)
}

func computeHash(keyFunc string, salt []byte, text string) []byte {
	if keyFunc == "pbkdf2" {
		return myPbkdf2(salt, []byte(text))
	}
	panic("Unknown key function " + keyFunc)
}

func buildHash(keyFunc, text string) []byte {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	check(err)
	octets := computeHash(keyFunc, salt, text)
	headerStr := keyFunc + ":" + base64.StdEncoding.EncodeToString(salt) +
		":" + base64.StdEncoding.EncodeToString(octets[:])
	return []byte(headerStr)
}

func saveHash(args []string) {
	if len(os.Args) > 2 || os.Args[1] != "-s" {
		fmt.Println("ERROR: invalid argument.", help)
		os.Exit(1)
	}
	wrappedHash := buildHash("pbkdf2", prompt())
	err := ioutil.WriteFile(hashFilename, wrappedHash, 0600)
	check(err)
	fmt.Println("Passphrase hash saved to", hashFilename)
}

func unwrapHash(wrappedHash []byte) (string, []byte, []byte) {
	fields := strings.Split(string(wrappedHash), ":")
	keyFunc := fields[0]
	salt, err := base64.StdEncoding.DecodeString(fields[1])
	check(err)
	passwdHash, err := base64.StdEncoding.DecodeString(fields[2])
	check(err)
	return keyFunc, salt, passwdHash
}

func practice() {
	wrappedHash, err := ioutil.ReadFile(hashFilename)
	if os.IsNotExist(err) {
		fmt.Println("ERROR: passphrase hash file not found.", help)
		os.Exit(1)
	}
	check(err)
	keyFunc, salt, passwdHash := unwrapHash(wrappedHash)
	fmt.Println("Type q to end practice.")
	turn := 0
	correct := 0
	for {
		turn++
		fmt.Printf("%d:", turn)
		octets, err := gopass.GetPasswd()
		check(err)
		response := string(octets)
		if response == "" {
			fmt.Println("Type q to quit.")
			turn-- // don't count this response
			continue
		} else if response == "q" {
			turn-- // don't count this response
			break
		}
		answer := "wrong"
		if bytes.Compare(computeHash(keyFunc, salt, response), passwdHash) == 0 {
			correct++
			answer = "OK"
		}
		fmt.Printf("  %s\thits=%d\tmisses=%d\n", answer, correct, turn-correct)
	}
	if turn > 0 {
		pct := float64(correct) / float64(turn) * 100
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
