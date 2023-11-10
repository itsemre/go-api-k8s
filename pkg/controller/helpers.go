package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AbortWithError returns a JSON body containing the error message back to the sender prior to
// creating a *gin.Error object to be logged
func AbortWithError(c *gin.Context, status int, err error) *gin.Error {
	error := c.Error(err)
	c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
	return error
}

// isAlphabetical returns true if the inputted character belongs in the alphabet
func isAlphabetical(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z')
}

// findFirstAlphabeticalCharacter returns the first alphabetical character of the inputted string
func findFirstAlphabeticalCharacter(str string) string {
	for i := 1; i < len(str); i++ {
		if isAlphabetical(str[i]) {
			return str[i:]
		}
	}
	return str
}

// getStartEnd returns the start & end query parameters from the request and parses them into integers
func getStartEnd(c *gin.Context) (int, int, error) {
	queryParams := c.Request.URL.Query()

	if !queryParams.Has("start") || !queryParams.Has("end") {
		return 0, 0, errors.New("please include the starting and ending comic numbers")
	}

	start, err := strconv.Atoi(queryParams.Get("start"))
	if err != nil {
		return 0, 0, errors.New("please make sure that the starting number is an integer")
	}
	end, err := strconv.Atoi(queryParams.Get("end"))
	if err != nil {
		return 0, 0, errors.New("please make sure that the ending number is an integer")
	}

	return start, end, nil
}
