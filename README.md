# passdrill

Using long, strong passphrases is great, once you overcome two challenges:

* memorize the passphrase;
* learn to type it fase;

`passdrill` lets you practice typing a long passphrase in a *safe* environment: your local console.

## Demo

First, run `passdrill -s` to save the hash of a passphrase you want to practice. The passphrase itself is not saved, only its SHA-512 hash.

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

The numbers (e.g `1:`) are the prompts. Nothing is echoed as you type. Typing just `q` quits the practice.