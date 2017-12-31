# passdrill

Using long, strong passphrases is great, once you overcome two challenges:

* memorize the passphrase;
* learn to type it quickly and reliably.

`passdrill` lets you practice typing a long passphrase in a *safe* environment: your local console.

This repository contains the same program implemented in Python 3 and Go, tested on GNU/Linux, Windows and MacOS.

> **WARNING**: On MacOS, `passdrill.py` uses Python's `getpass` function, and it does not support keyboards with dead keys for typing combining diacritics, such as the tilde in "não" (`getpass` returns "n~ao"). It does handle non-ASCII text when it is pasted to the prompt, so it is clearly a keyboard handling issue.


## Demo

First, run `passdrill -s` to save the hash of a passphrase you want to practice. The passphrase itself is not saved, only a derived key using **scrypt** or PBKDF2 with SHA-512 (see **Implementation Notes**).

>  **NOTE**: Before saving, `passdrill -s` will display the passphrase on your console so that you can confirm that you've typed it correctly. It will never be shown while you practice.

Sample initial session:

```
$ ./passdrill -s
WARNING: the passphrase will be shown so that you can check it!
Type passphrase to hash (it will be echoed): my extra strong secret       
Passphrase to be hashed -> my extra strong secret
Confirm (y/n): y
Passphrase hash saved to passdrill.hash
```

To practice typing the passphrase, just run `passdrill`.

Sample practice session:

```
$ ./passdrill
Type q to end practice.
1:
  wrong	hits=0	misses=1
2:
  OK	hits=1	misses=1
3:
  wrong	hits=1	misses=2
4:
  OK	hits=2	misses=2
5:
  OK	hits=3	misses=2
6:

5 exercises. 60.0% correct.
```

The numbers (e.g `1:`) are the prompts. Nothing is echoed as you type. Entering the letter `q` alone quits the practice.


## Implementation notes

This program is implemented in Python 3 and Go for didactic reasons. The implementations behave identically, except for these two points:

* On MacOS, Python's `getpass` is used to read the passphrase in practice mode. It uses some MacOS API that shows a nice key prompt in the console, but it does not support dead keys for typing diacritics: when I type "não", I get "n~ao".

* When creating a new `passdrill.hash`, the Go version always uses the stronger **scrypt** method. The Python code tries to use **scrypt** if available, otherwise it uses PBKDF2. The hash file data is prefixed with the name of the method used to create it.


### Comparing the implementations

The source code the for Go version is about 40% longer than the Python version:

|     | Python   | Go   | Δ    |
| ---:| --------:| ----:| ----:| 
|lines| 125      | 177  | +41% |
|words| 407      | 559  | +37% |


#### Error checking

For me, the major irritant in Go code is the repetition of `if err != nil {…}` blocks. In a CLI utility like this, almost every error is fatal: there is nothing to do except terminate, with some informative error output. Python gives me that for free. In Go, I wish I could call a function like this:

```go
	exitIf(err, "the input file is malformed")
``` 

Here is a sub-optimal implementation of `exitIf`:

```
func exitIf(e error, msg string) {
	if e != nil {
		logger.Fatalln(msg, e)
	}
}
```

The `Fatal…` functions output the line number when the logger is configured with the `log.Lshortfile`. The problem with `exitIf` is that the line number logged points to the `logger.Fatalln()` call, which is always the same. Ideally, `exitIf` should report the line number where **it** is called. I need to learn the internals of the `log` module.

#### Missing batteries

I could not find equivalents for Python's `input` and `getpass` functions in the Go standard library. After spending some time looking for them in the Go docs, searching the Web and asking around, I decided to:

* implement my own 10-line `input` function (see `passdrill.go`);
* install the `github.com/howeyc/gopass` package as a dependency.

On the other hand, the **scrypt** password derivation algorithm is available from `golang.org/x/crypto/scrypt`, so it was easy to fetch and use, thanks to the `go get` command. But in Python, the `hashlib.scrypt` function is only available for Python 3.6, and even then, it's only provided if the interpreter was compiled with OpenSSL 1.1+ — which was not my case. 

On GNU/Linux (Ubuntu 16.04) I found it easier to pip-install the 3rd-party `scrypt` package by Magnus Hallin ([available on Pypi](https://pypi.python.org/pypi/scrypt/)) than compiling OpenSSL 1.1 and Python 3.6. However, because `scrypt` relies on C code, installing it on some environments is difficult.

#### Ease of distribuition

The issues with **scrypt** described above exemplify a major deficiency of Python for building utilities like `passdrill`: dealing with external dependencies, particularly when they need to be compiled in some other language. I did not invest the time to setup a C compiler to build `scrypt` on Windows. As a work around, if the `scrypt` package cannot be imported, `passdrill.py` displays a warning and falls back to using the less secure `hashlib.pbkdf2_hmac` derivation function from the Python 3.6 standard library.

The problem with installing dependencies does not affect only the Python programmer: it affects end users as well: they also must deal with these issues to run any script I write with external dependencies. 

In contrast, `go build` compiles to static, stand-alone executables that can be easily distributed in binary form. In addition, because the performance of Go code is comparable to C, Go projects don't depend on external C libraries as much as Python projects do.

### Contributors welcome!

I am an experienced Pythonista but a newbie Gopher. If you know how to improve either version, please post an issue or send a pull request. Thanks!

If you'd like to ask questions or give direct feedback, please tweet me: [@ramalhoorg](https://twitter.com/ramalhoorg)
