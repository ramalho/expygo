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

The source code the for Go version is 32% longer than the Python version:

|     | Python   | Go   | Δ    |
| ---:| --------:| ----:| ----:| 
|lines| 125      | 165  | +32% |
|words| 407      | 533  | +31% |

I could not find equivalents for Python's `input` and `getpass` functions in the Go standard library. After spending some time looking for them in the Go docs, searching the Web and asking around, I decided to:

* implement my own 10-line `input` function (see `passdrill.go`);
* install the `github.com/howeyc/gopass` package as a dependency.

On the other hand, the **script** password derivation algorithm is available from `golang.org/x/crypto/scrypt`, so it was easy to fetch and use, thanks to the `go get` command. But in Python, the `hashlib.scrypt` function is only available for Python 3.6, and even then, it's only provided if the interpreter was compiled with OpenSSL 1.1+ — which was not my case. 

On GNU/Linux (Ubuntu 16.04) I found it easier to pip-install the 3rd-party `scrypt` package by Magnus Hallin ([available on Pypi](https://pypi.python.org/pypi/scrypt/)) than compiling OpenSSL 1.1 and Python 3.6. However, because `scrypt` relies on C code, installing it on some environments is difficult.

This highlights a major deficiency of Python for building utilities like `passdrill`: dealing with external dependencies, particularly when they need to be compiled in some other language. I did not invest the time to setup a C compiler to build `scrypt` on Windows. As a work around, if the `scrypt` package cannot be imported, `passdrill.py` displays a warning and falls back to using the less secure `hashlib.pbkdf2_hmac` derivation function from the Python 3.6 standard library. 

In contrast, Go compiles to static, stand-alone executables that can be easily compiled and distributed in binary form. In addition, because the performance of Go code is comparable to C, Go projects don't depend on external C libraries as much as Python projects do.


### Contributors welcome!

I am an experienced Pythonista but a newbie Gopher. If you know how to improve either version, please post an issue or send a pull request. Thanks!
