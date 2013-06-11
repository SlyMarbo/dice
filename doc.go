/*
Package dice is a simple dice-rolling library with automatic collection of roll statistics.

Example use:
		package main

		import "github.com/SlyMarbo/dice"

		func main() {
			// Perform a single roll and store just the roll.
			roll, err := dice.RollSimple("1d6 +2")

			// Perform a single roll and store the roll and its statistics.
			result, err := dice.Roll("1d6 +2")

			// Perform a series of rolls and store the full results.
			results, err := dice.RollAll("1d6 +2, D12 -4, 18d100")
		}
*/
package dice
