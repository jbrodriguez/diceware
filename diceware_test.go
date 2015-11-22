package main

import (
	// "flag"
	// "fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var rolls, expectedDefault, expectedBeale []string

func TestMain(m *testing.M) {
	rolls = []string{
		"64122",
		"56666",
		"65546",
		"25231",
		"41321",
		"23656",
		"21361",
		"61243",
	}

	expectedDefault = []string{
		"whim",
		"the",
		"44",
		"far",
		"lw",
		"dyke",
		"coot",
		"tift",
	}

	expectedBeale = []string{
		"wm",
		"tile",
		"222",
		"faze",
		"madam",
		"ears",
		"coral",
		"toni",
	}

	m.Run()
}

func TestGetPassword(t *testing.T) {
	chosenDefault := getPassword(strings.NewReader(list), rolls)
	chosenBeale := getPassword(strings.NewReader(bealeList), rolls)

	assert.Equal(t, expectedDefault, chosenDefault, "expectedDefault is wrong")
	assert.Equal(t, expectedBeale, chosenBeale, "expectedBeale is wrong")
}

func TestStrength(t *testing.T) {
	strengthDefault := getPasswordStrength(strings.Join(expectedDefault, " "), nil)
	strengthBeale := getPasswordStrength(strings.Join(expectedBeale, " "), nil)

	// Entropy
	assert.Equal(t, 130.241, strengthDefault.Entropy, "strengthDefault.Entropy is wrong")
	assert.Equal(t, 124.698, strengthBeale.Entropy, "strengthBeale.Entropy is wrong")

	// CrackTime
	assert.Equal(t, 8.043448503439465e+34, strengthDefault.CrackTime, "strengthDefault.CrackTime is wrong")
	assert.Equal(t, 1.725027324318945e+33, strengthBeale.CrackTime, "strengthBeale.CrackTime is wrong")

	// CrackTimeDisplay
	assert.Equal(t, "centuries", strengthDefault.CrackTimeDisplay, "strengthDefault.CrackTimeDispla is wrong")
	assert.Equal(t, "centuries", strengthBeale.CrackTimeDisplay, "strengthBeale.CrackTimeDisplay is wrong")

	// Score
	assert.Equal(t, 4, strengthDefault.Score, "strengthDefault.Score is wrong")
	assert.Equal(t, 4, strengthBeale.Score, "strengthBeale.Score is wrong")
}

func TestWordsEntropy(t *testing.T) {
	entropy := getWordsEntropy(8, 1000000000000.)

	assert.Equal(t, 103.36, entropy.Entropy, "entropy.Entropy is wrong")
	assert.Equal(t, 6.683747269421867e+30, entropy.HalfKeySpace, "entropy.HalfKeySpace is wrong")
	assert.Equal(t, 6.683747269421867e+18, entropy.Seconds, "entropy.Seconds is wrong")
	assert.Equal(t, 1.1139578782369778e+17, entropy.Minutes, "entropy.Minutes is wrong")
	assert.Equal(t, 1.8565964637282962e+15, entropy.Hours, "entropy.Hours is wrong")
	assert.Equal(t, 7.735818598867902e+13, entropy.Days, "entropy.Days is wrong")
	assert.Equal(t, 2.1194023558542197e+11, entropy.Years, "entropy.Years is wrong")
	assert.Equal(t, 3.1538725533544936e+09, entropy.HumanLifetimes, "entropy.HumanLifetimes is wrong")
	assert.Equal(t, 2.1194023558542198e+08, entropy.Millenia, "entropy.Millenia is wrong")
	assert.Equal(t, 15.360214203900707, entropy.UniverseLifetimes, "entropy.UniverseLifetimes is wrong")
}
