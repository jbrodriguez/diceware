package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/nbutton23/zxcvbn-go/matching"
	"github.com/nbutton23/zxcvbn-go/scoring"
	zgm "github.com/nbutton23/zxcvbn-go/utils/math"
)

var Version string

const diceCount = 5

type Entropy struct {
	Entropy           float64
	KeySpace          float64
	HalfKeySpace      float64
	Seconds           float64
	Minutes           float64
	Hours             float64
	Days              float64
	Years             float64
	Millenia          float64
	HumanLifetimes    float64
	UniverseLifetimes float64
}

func main() {
	var length int
	var beale bool
	var entropy bool
	var strength bool

	flag.IntVar(&length, "w", 6, "number of words to generate")
	// flag.BoolVar(&extra, "e", false, "generate an extra random character")
	flag.BoolVar(&beale, "b", false, "use alternate word list from Alan Beale")
	flag.BoolVar(&entropy, "e", false, "display entropy")
	flag.BoolVar(&strength, "s", false, "display password strength")
	flag.Parse()

	fmt.Printf("diceware v%s\n", Version)
	fmt.Println("Copyright (c) 2015 Juan B. Rodriguez")
	fmt.Println("Web: http://jbrodriguez.io")
	fmt.Println("Github: https://github.com/jbrodriguez/diceware")
	fmt.Println()

	wordList := "default"
	if beale {
		wordList = "alternate Alan Beale's"
	}

	fmt.Printf("Using the %s word list ...\n", wordList)

	if entropy {
		displayEntropy(length)
	}

	fmt.Println()
	fmt.Printf("Roll 5 dices %d times and enter the result after each prompt (press Ctrl+C to stop at any moment):\n\n", length)

	rolls := make([]string, length)
	count := 1

	for {
		if count > length {
			break
		}

		var ln string
		fmt.Printf("Roll #%d: ", count)
		fmt.Scanf("%s", &ln)

		if len(ln) != diceCount {
			fmt.Printf("You must input exactly %d digits. Please try again\n", diceCount)
			continue
		}

		if !digitsAreOk(ln, diceCount) {
			fmt.Printf("Each digit must be between 1 and 6. Please try again\n")
			continue
		}

		rolls[count-1] = ln
		count++
	}

	var chosen []string
	if beale {
		chosen = getPassword(strings.NewReader(bealeList), rolls)
	} else {
		chosen = getPassword(strings.NewReader(list), rolls)
	}

	fmt.Println("\nYour password is:")
	fmt.Printf("%s\n", strings.Join(chosen, " "))

	if strength {
		fmt.Println()
		displayStrength(strings.Join(chosen, " "))
	}

}

func digitsAreOk(line string, count int) (ok bool) {
	for i := 0; i < count; i++ {
		digit, err := strconv.Atoi(string(line[i]))
		ok = err == nil && digit >= 1 && digit <= 6
		if !ok {
			break
		}
	}

	return ok
}

func getPassword(rd io.Reader, rolls []string) (chosen []string) {
	r := csv.NewReader(rd)
	r.Comma = '\t'
	r.FieldsPerRecord = 2

	words := make(map[string]string, 7776)
	rec, err := r.Read()
	for err == nil {
		words[rec[0]] = rec[1]
		rec, err = r.Read()
	}

	if err != io.EOF {
		log.Fatalln(err)
	}

	if len(words) != 7776 {
		log.Fatalf("Expected exactly 7776 words in tsv wordlist, but found %d", len(words))
	}

	chosen = make([]string, len(rolls))
	for k, v := range rolls {
		chosen[k] = words[v]
	}

	return chosen
}

func getWordsEntropy(numWords int, guessesPerSec float64) Entropy {
	entropy := Entropy{}

	var (
		avgHumanLifespanInYears = 67.2         // https://en.wikipedia.org/wiki/Life_expectancy
		ageOfUniverseInYears    = 13798000000. // https://en.wikipedia.org/wiki/Age_of_the_universe
	)

	entropy.Entropy = 12.92 * float64(numWords)

	// https://xkcd.com/936/
	// https://security.stackexchange.com/questions/62832/is-the-oft-cited-xkcd-scheme-no-longer-good-advice/62881#62881
	// https://hashcat.net/forum/thread-2580.html
	entropy.KeySpace = math.Pow(7776., float64(numWords))

	// Divide the keySpace in half. On average it is expected that an
	// exhaustive search of only half the keySpace will result in success.
	entropy.HalfKeySpace = entropy.KeySpace / 2

	// // "Assume that your adversary is capable of a trillion guesses per second" - Snowden
	// // http://www.nytimes.com/2013/08/18/magazine/laura-poitras-snowden.html?pagewanted=all&_r=0
	// guessesPerSec = 1000000000000.

	entropy.Seconds = entropy.HalfKeySpace / guessesPerSec
	entropy.Minutes = entropy.Seconds / 60
	entropy.Hours = entropy.Minutes / 60
	entropy.Days = entropy.Hours / 24
	entropy.Years = entropy.Days / 365
	entropy.Millenia = entropy.Years / 1000
	entropy.HumanLifetimes = entropy.Years / avgHumanLifespanInYears
	entropy.UniverseLifetimes = entropy.Years / ageOfUniverseInYears

	return entropy
}

func getPasswordStrength(password string, userInputs []string) scoring.MinEntropyMatch {
	start := time.Now()
	matches := matching.Omnimatch(password, userInputs)
	result := scoring.MinimumEntropyMatchSequence(password, matches)
	end := time.Now()

	calcTime := end.Nanosecond() - start.Nanosecond()
	result.CalcTime = zgm.Round(float64(calcTime)*time.Nanosecond.Seconds(), .5, 3)
	return result
}

func displayEntropy(numWords int) {
	entropy := getWordsEntropy(numWords, 1000000000000.)

	fmt.Printf(`
You chose to create a password of %d words, which is ~%s bits of entropy (~12.92 bits/word and ~5.16/symbol).
That many words equates to an effective keyspace of ~%s possible keys,
(7776^words)/2 which is 50%% of the total key space since on average a key will be found in that reduced keyspace.
Conservatively assuming an adversary can guess passwords at the rate of a trillion (1,000,000,000,000) keys/second,
as Edward Snowden claims, an exhaustive brute-force search on half the total keyspace might take:

~%s seconds
~%s minutes
~%s hours
~%s days
~%s years
~%s x avg. lifespan
~%s millenia
~%s x age Universe

Learn more about calculating entropy at http://world.std.com/~reinhold/dicewarefaq.html#calculatingentropy	
`,
		numWords,
		Commaf(trunc(entropy.Entropy)),
		Commaf(trunc(entropy.HalfKeySpace)),
		Commaf(trunc(entropy.Seconds)),
		Commaf(trunc(entropy.Minutes)),
		Commaf(trunc(entropy.Hours)),
		Commaf(trunc(entropy.Days)),
		Commaf(trunc(entropy.Years)),
		Commaf(trunc(entropy.HumanLifetimes)),
		Commaf(trunc(entropy.Millenia)),
		Commaf(trunc(entropy.UniverseLifetimes)),
	)

}

func displayStrength(password string) {
	strength := getPasswordStrength(password, nil)

	fmt.Printf(`
Entropy: %f
Crack Time: %f seconds
Alternate Crack Time: %s
Score: %d

Key for score:
0 = too guessable: risky password
1 = very guessable: protection from throttled online attacks
2 = somewhat guessable: protection from unthrottled online attacks
3 = safely unguessable: moderate protection from offline slow-hash scenario
4 = very unguessable: strong protection from offline slow-hash scenario
`,
		strength.Entropy,
		strength.CrackTime,
		strength.CrackTimeDisplay,
		strength.Score,
	)
}

// b, err := json.MarshalIndent(strength, "", " ")
// if err != nil {
// 	fmt.Println(err)
// 	os.Exit(1)
// }

// fmt.Println(string(b))

func Commaf(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{','}

	parts := strings.Split(strconv.FormatFloat(v, 'f', -1, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func trunc(v float64) float64 {
	if v < 1 {
		return v
	}

	return rint(v)
}

func rint(x float64) float64 {
	v, frac := math.Modf(x)
	if x > 0.0 {
		if frac > 0.5 || (frac == 0.5 && uint64(v)%2 != 0) {
			v += 1.0
		}
	} else {
		if frac < -0.5 || (frac == -0.5 && uint64(v)%2 != 0) {
			v -= 1.0
		}
	}

	return v
}
