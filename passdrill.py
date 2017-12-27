#!/usr/bin/env python3

"""passdrill: typing drills for practicing passphrases
"""

import sys
import hashlib
import base64
import getpass

HASH_ALGORITHM = 'sha512'
HASH_FILENAME = 'passdrill.' + HASH_ALGORITHM
HELP = 'Use -s to save passphrase hash for practice.'


def prompt():
    print('WARNING: the passphrase will be shown so that you can check it!')
    confirmed = ''
    while confirmed != 'y':
        passwd = input('Type passphrase to hash (it will be echoed): ')
        if passwd == '' or passwd == 'q':
            print('ERROR: the passphrase cannot be empty or "q".')
            continue
        print(f'Passphrase to be hashed -> {passwd}')
        confirmed = input('Confirm (y/n): ').lower()
    return passwd


def hash_bytes(text):
    octets = hashlib.sha512(text.encode('utf-8')).digest()
    return base64.b64encode(octets)


def save_hash(argv):
    if len(argv) > 2 or argv[1] != '-s':
        print('ERROR: invalid argument.', HELP)
        sys.exit(1)
    passwd_hash = hash_bytes(prompt())
    with open(HASH_FILENAME, 'wb') as fp:
        fp.write(passwd_hash)
    print(f'Passphrase {HASH_ALGORITHM} hash saved to',
          HASH_FILENAME)


def practice():
    try:
        with open(HASH_FILENAME, 'rb') as fp:
            passwd_hash = fp.read()
    except FileNotFoundError:
        print('ERROR: passphrase hash file not found.', HELP)
        sys.exit(2)
    print('Type q to end practice.')
    turn = 0
    correct = 0
    while True:
        turn += 1
        response = getpass.getpass(f'{turn}:')
        if response == '':
            print('Type q to quit.')
            turn -= 1  # don't count this response
            continue
        elif response == 'q':
            turn -= 1  # don't count this response
            break
        if hash_bytes(response) == passwd_hash:
            correct += 1
            answer = 'OK'
        else:
            answer = 'wrong'
        print(f'  {answer}\thits={correct}\tmisses={turn-correct}')

    if turn:
        pct = correct / turn * 100
        print(f'\n{turn} exercises. {pct:0.1f}% correct.')


if __name__ == '__main__':
    if len(sys.argv) > 1:
        save_hash(sys.argv)
    else:
        practice()
