#!/usr/bin/env python3

"""passdrill: typing drills for practicing passphrases
"""

import sys
import hashlib
import base64
import getpass
import os

HASH_FILENAME = 'passdrill.hash'
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


def pbkdf2(salt, octets):
    algorithm = 'sha512'
    rounds = 100_000
    return hashlib.pbkdf2_hmac(algorithm, octets, salt, rounds)


def compute_hash(key_func, salt, text):
    octets = text.encode('utf-8')
    if key_func == 'pbkdf2':
        return pbkdf2(salt, octets)
    else:
        raise ValueError('Unknown key function ' + repr(key_func))


def build_hash(key_func, text):
    salt = os.urandom(32)
    octets = compute_hash(key_func, salt, text)
    header = key_func.encode('utf-8') + b':' + base64.b64encode(salt)
    return header + b':' + base64.b64encode(octets)


def save_hash(argv):
    if len(argv) > 2 or argv[1] != '-s':
        print('ERROR: invalid argument.', HELP)
        sys.exit(1)
    wrapped_hash = build_hash('pbkdf2', prompt())
    with open(HASH_FILENAME, 'wb') as fp:
        fp.write(wrapped_hash)
    print(f'Passphrase hash saved to', HASH_FILENAME)


def unwrap_hash(wrapped_hash):
    key_func, salt, passwd_hash = wrapped_hash.split(b':')
    return (key_func.decode('utf-8'), base64.b64decode(salt),
        base64.b64decode(passwd_hash))


def practice():
    try:
        with open(HASH_FILENAME, 'rb') as fp:
            wrapped_hash = fp.read()
    except FileNotFoundError:
        print('ERROR: passphrase hash file not found.', HELP)
        sys.exit(2)
    key_func, salt, passwd_hash = unwrap_hash(wrapped_hash)
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
        if compute_hash(key_func, salt, response) == passwd_hash:
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
