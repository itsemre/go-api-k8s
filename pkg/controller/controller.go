package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itsemre/go-api-k8s/pkg/config"
)

// Controller is the struct implementing the corresponding gin.HandlerFunc fxns
type Controller struct {
	Cfg *config.Config
}

// NewController returns a pointer to a new Controller instance
func NewController(conf *config.Config) *Controller {
	return &Controller{
		Cfg: conf,
	}
}

// Comic is a representation of comic book metadata
type Comic struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Transcript string `json:"transcript"`
	Img        string `json:"img"`
	Alt        string `json:"alt"`
	News       string `json:"news"`
	Link       string `json:"link"`
}

// GetComics a gin handler function that returns a list of comics that were published on
// odd months, and whose number is in between the 'start' and 'end' query parameters, sorted
// alphabetically by the title.
func (ctrl *Controller) GetComics(c *gin.Context) {
	var (
		comics []Comic
		comic  Comic
	)

	// Extract query parameters
	start, end, err := getStartEnd(c)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, err)
		return
	}

	client := &http.Client{}

	// Iterate through the comics range
	for i := start; i <= end; i++ {
		// Get the metadata of the current comic book
		req, err := http.NewRequest("GET", fmt.Sprintf("https://xkcd.com/%d/info.0.json", i), nil)
		if err != nil {
			AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			AbortWithError(c, http.StatusInternalServerError, err)
			return
		}
		err = json.Unmarshal(bodyText, &comic)
		if err != nil {
			AbortWithError(c, http.StatusInternalServerError, err)
			return
		}

		// If the published month is odd, add the comic to the list
		if m, _ := strconv.Atoi(comic.Month); m%2 != 0 {
			comics = append(comics, comic)
		}
	}

	// Sort list alphabetically by the title
	sort.Slice(comics, func(i, j int) bool {
		title1 := comics[i].Title
		title2 := comics[j].Title

		// If the first character is non-alphabetic get the first alphabetic character
		if !isAlphabetical(title1[0]) {
			title1 = findFirstAlphabeticalCharacter(title1)
		}
		if !isAlphabetical(title2[0]) {
			title2 = findFirstAlphabeticalCharacter(title2)
		}

		return title1 < title2
	})

	c.JSON(http.StatusOK, gin.H{"comics": comics})
}

// Health is a simple handler that allows us to check the status of our API
func (ctrl *Controller) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
