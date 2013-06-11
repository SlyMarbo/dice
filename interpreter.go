package dice

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

// Ensure pseudo-random rolls.
func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	ParseFailure = errors.New("Error: Roll string could not be parsed.")
	DieSizeError = errors.New("Error: Die size is invalid. A die must have at least 2 sides.")
	interpreter  = regexp.MustCompile(`(\d*)[dD](\d+)\s?([\+\-]\d+)?`)
)

const (
	ALL_MATCHES = -1
)

// Result represents a single dice roll.
// It contains the relevant statistics for
// the roll, such as the minimum, maximum,
// and average results of the given roll
// string.
type Result struct {
	Roll   int
	Min    int
	Max    int
	Avg    float64
	String string
}

// Results represents a set of one or more
// dice rolls. It contains the relevant
// statistics for the rolls, such as the
// minimum, maximum, and average results
// of the combined rolls.
type Results struct {
	Rolls []*Result
	Min   int
	Max   int
	Avg   float64
}

func newResults(n int) *Results {
	out := new(Results)
	out.Rolls = make([]*Result, n)
	return out
}

// Roll takes a roll string, parses it, and
// returns the roll result and an error.
//
// Some example roll strings follow:
//		"1d6"
//		"D100"
//		"4d7 -18"
func SimpleRoll(roll string) (int, error) {
	results, err := doRoll(roll, 1)
	if err != nil {
		return 0, err
	}
	if results != nil {
		return results.Rolls[0].Roll, nil
	}
	return 0, nil
}

// Roll takes a roll string, parses it, and
// returns the roll result and an error.
//
// Some example roll strings follow:
//		"1d6"
//		"D100"
//		"4d7 -18"
func Roll(roll string) (result *Result, err error) {
	results, err := doRoll(roll, 1)
	if err != nil {
		return nil, err
	}
	if results != nil {
		return results.Rolls[0], nil
	}
	return nil, nil
}

// RollAll takes a roll string, parses it,
// and returns the roll results and an error.
//
// The roll string may contain multiple
// rolls, separated by any non-numberical
// character.
//
// Some example roll strings follow:
//
//		"1d6"
//		"D100"
//		"4d7 -18"
//		"1d6 +3, 2d2"
func RollAll(roll string) (results *Results, err error) {
	return doRoll(roll, ALL_MATCHES)
}

// doRoll performs the actual rolling.
func doRoll(roll string, n int) (results *Results, err error) {
	matches := interpreter.FindAllStringSubmatch(roll, n)
	if matches == nil {
		return nil, ParseFailure
	}

	// Create the output results and set the minimum
	// very large so that non-zero results can still
	// become the minimum.
	results = newResults(len(matches))
	results.Min = 1<<31 - 1

	// Iterate through roll strings.
	for i, match := range matches {
		result := new(Result)

		// num is the number of dice to roll.
		num := 1
		if match[1] != "" {
			num, err = strconv.Atoi(match[1])
			if err != nil {
				return nil, err
			}
		}

		// size is the number of faces on the dice.
		size, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}
		if size < 2 {
			return nil, DieSizeError
		}

		// mod is the roll modifier.
		mod := 0
		if match[3] != "" {
			mod, err = strconv.Atoi(match[3])
			if err != nil {
				return nil, err
			}
		}

		// Store the roll string.
		result.String = match[0]

		// Check whether dice are rolled.
		if num > 0 {
			result.Min = (num * 1) + mod
			result.Max = (num * size) + mod
			result.Avg = (float64(num) * (float64(size+1) / 2.0)) + float64(mod)

			result.Roll = mod
			for i := 0; i < num; i++ {
				result.Roll += rand.Intn(size) + 1
			}
		} else {
			result.Roll = mod
			result.Min = mod
			result.Max = mod
			result.Avg = float64(mod)
		}

		// Update overall statistics.
		if result.Min < results.Min {
			results.Min = result.Min
		}
		if result.Max > results.Max {
			results.Max = result.Max
		}
		results.Avg += float64(result.Avg)

		// Store the result.
		results.Rolls[i] = result
	}

	// Finalise the overall average.
	if results.Avg != 0 {
		results.Avg /= float64(len(results.Rolls))
	}

	return results, nil
}
