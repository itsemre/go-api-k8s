package controller

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// =============================================================================
// TEST SUITE SETUP
// =============================================================================

type HelpersUnitSuite struct {
	suite.Suite
}

func TestHelpersUnitSuite(t *testing.T) {
	suite.Run(t, &HelpersUnitSuite{})
}

// =============================================================================
// UNIT TESTS
// =============================================================================

func (us *HelpersUnitSuite) TestIsAlphabetical() {
	testCases := []struct {
		name           string
		input          byte
		expectedOutput bool
	}{
		{
			"Normal Operation Uppercase",
			'A',
			true,
		},
		{
			"Normal Operation Lowercase",
			'b',
			true,
		},
		{
			"Non Alphabetic Character",
			'@',
			false,
		},
	}

	for i := range testCases {
		test := testCases[i]
		us.Run(test.name, func() {
			us.Equal(test.expectedOutput, isAlphabetical(test.input))
		})
	}
}

func (us *HelpersUnitSuite) TestFindFirstAlphabeticalCharacter() {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			"Case 1",
			"12th Avenue",
			"th Avenue",
		},
		{
			"Case 2",
			"@!42&characters",
			"characters",
		},
		{
			"Case 3",
			"❤©♬☂unicode☯∞✿✡♨",
			"unicode☯∞✿✡♨",
		},
		{
			"Case 4",
			"     spaces",
			"spaces",
		},
	}

	for i := range testCases {
		test := testCases[i]
		us.Run(test.name, func() {
			us.Equal(test.expectedOutput, findFirstAlphabeticalCharacter(test.input))
		})
	}
}

func (us *HelpersUnitSuite) TestGetStartEnd() {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			"Case 1",
			"12th Avenue",
			"th Avenue",
		},
		{
			"Case 2",
			"@!42&characters",
			"characters",
		},
		{
			"Case 3",
			"❤©♬☂unicode☯∞✿✡♨",
			"unicode☯∞✿✡♨",
		},
		{
			"Case 4",
			"     spaces",
			"spaces",
		},
	}

	for i := range testCases {
		test := testCases[i]
		us.Run(test.name, func() {
			us.Equal(test.expectedOutput, findFirstAlphabeticalCharacter(test.input))
		})
	}
}
