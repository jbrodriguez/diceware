diceware
========

**tl;dr** diceware is an implementation of the [diceware passphrase](http://world.std.com/~reinhold/diceware.html), created by Arnold G. Reinhold

[![Build Status](https://drone.io/github.com/jbrodriguez/diceware/status.png)](https://drone.io/github.com/jbrodriguez/diceware/latest)

## Screen output
```
$ ./diceware -e -s -w 8
diceware v0.1.0-1.9fdd687
Copyright (c) 2015 Juan B. Rodriguez
Web: http://jbrodriguez.io
Github: https://github.com/jbrodriguez/diceware

Using the default word list ...

You chose to create a password of 8 words, which is ~103 bits of entropy (~12.92 bits/word and ~5.16/symbol).
That many words equates to an effective keyspace of ~6,683,747,269,421,867,000,000,000,000,000 possible keys,
(7776^words)/2 which is 50% of the total key space since on average a key will be found in that reduced keyspace.
Conservatively assuming an adversary can guess passwords at the rate of a trillion (1,000,000,000,000) keys/second,
as Edward Snowden claims, an exhaustive brute-force search on half the total keyspace might take:

~6,683,747,269,421,867,000 seconds
~111,395,787,823,697,780 minutes
~1,856,596,463,728,296 hours
~77,358,185,988,679 days
~211,940,235,585 years
~3,153,872,553 x avg. lifespan
~211,940,236 millenia
~15 x age Universe

Learn more about calculating entropy at http://world.std.com/~reinhold/dicewarefaq.html#calculatingentropy

Roll 5 dices 8 times and enter the result after each prompt (press Ctrl+C to stop at any moment):

Roll #1: 64122
Roll #2: 56666
Roll #3: 65546
Roll #4: 25231
Roll #5: 41321
Roll #6: 23656
Roll #7: 21361
Roll #8: 61243

Your password is:
whim the 44 far lw dyke coot tift


Entropy: 130.241000
Crack Time: 80434485034394646894966093304037376.000000 seconds
Alternate Crack Time: centuries
Score: 4

Key for score:
0 = too guessable: risky password
1 = very guessable: protection from throttled online attacks
2 = somewhat guessable: protection from unthrottled online attacks
3 = safely unguessable: moderate protection from offline slow-hash scenario
4 = very unguessable: strong protection from offline slow-hash scenario
```

## Rationale
There's already a great diceware golang version in https://github.com/natefinch/diceware.

But I happen to own a set of Bicycle dice, which although are not the same as Casino dice, they're still quite usable.

So I created a utility that allows me to enter the result of my rolling the dice, and output the  password.

Additionally, I added some stats:
- Some theoretical entropy stats based on number of words in the password (taken from https://grempe.github.io/diceware/)
- A password strength meter based on [dropbox's zxcvbn implementation](https://github.com/dropbox/zxcvbn) (and golang's version in https://github.com/nbutton23/zxcvbn-go)

## Installation
Download the binary that corresponds to your platform.

## Development
It's assumed you have GO15VENDOREXPERIMENT set to 1.

The recommended way is to install [gom](https://github.com/mattn/gom), then

```
$ git clone https://github.com/jbrodriguez/diceware
$ cd diceware
$ gom install
$ make build
$ ./diceware
```

Alternatively (but you lose vendor dependency management), you can

```
$ go get https://github.com/nbutton23/zxcvbn-go
$ go get github.com/stretchr/testify(for testing only)
$ git clone https://github.com/jbrodriguez/diceware
$ cd diceware
$ make build
$ ./diceware
```

## Usage
```
$ ./diceware -h
Usage of ./diceware:
  -b    use alternate word list from Alan Beale
  -e    display entropy
  -s    display password strength
  -w int
        number of words to generate (default 6)
```

## Credits
- [Arnold G. Reinhold](http://world.std.com/~reinhold/diceware.html)
- [Glenn Rempe](https://github.com/grempe/diceware)
- [Nate Finch](https://github.com/natefinch/diceware)
- [Nathan Button](https://github.com/natefinch/diceware)

## LICENSE
MIT
