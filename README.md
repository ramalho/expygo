# passdrill

Using long, strong passphrases is great, once you overcome two challenges:

* memorize the passphrase;
* learn to type it quickly and reliably.

`passdrill` lets you practice typing a long passphrase in a *safe* environment: your local console.

This repository contains the same program implemented in Python 3 and Go, tested on GNU/Linux, Windows and MacOS.

> **WARNING**: On MacOS, `passdrill.py` uses Python's `getpass` function, and it does not support keyboards with dead keys for typing combining diacritics, such as the tilde in "não" (`getpass` returns "n~ao"). It does handle non-ASCII text when it is pasted to the prompt, so it is clearly a keyboard handling issue.


## Demo

First, run `passdrill -s` to save the hash of a passphrase you want to practice. The passphrase itself is not saved, only its SHA-512 digest.

>  **NOTE**: Before saving, `passdrill -s` will display the passphrase on your console so that you can confirm that you've typed it correctly. It will never be shown while you practice.

Sample initial session:

```
$ ./passdrill -s
WARNING: the passphrase will be shown so that you can check it!
Type passphrase to hash (it will be echoed): my extra strong secret       
Passphrase to be hashed -> my extra strong secret
Confirm (y/n): y
Passphrase sha512 hash saved to passdrill.sha512
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


## About the code

This program is implemented in Python 3 and Go for didactic reasons. The implementations behave identically as far as I can tell, except for this:

* On MacOS, Python's `getpass` is used to read the passphrase in practice mode. It uses some MacOS API that shows a nice key prompt in the console, but it does not support dead keys for typing diacritics: when I type "não", I get "n~ao".


### Comparing the implementations

The source code the for Go version is longer than the Python version:

|     | Python   | Go   | Δ    |
| ---:| --------:| ----:| ----:| 
|lines| 82       | 117  | +43% |
|words| 247      | 330  | +34% |

The Python version uses only packages from the Python standard library.

I could not find equivalents for Python's `input` and `getpass` functions in the Go standard library. After wasting some time looking for them in the Go docs, I had to do some research and ask around. I decided to:

* implement my own 10-line `input` function (see `passdrill.go`);
* install the `github.com/howeyc/gopass` package as a dependency.


### Contributors welcome!

I am an experienced Pythonista but a newbie Gopher. If you know how to improve either version, please post an issue or send a pull request. Thanks!
