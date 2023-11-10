package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/itsemre/go-api-k8s/pkg/config"
	"github.com/stretchr/testify/suite"
)

const serverAddress = "localhost:8080"

var cfg *config.Config

// =============================================================================
// TEST SUITE SETUP
// =============================================================================

type errorBody struct {
	Msg string `json:"error"`
}

type ControllerUnitSuite struct {
	suite.Suite
	server struct {
		HTTPServer *http.Server
		Router     *gin.Engine
	}
	ctrl *Controller
}

func (us *ControllerUnitSuite) SetupSuite() {
	cfg = config.NewConfig()
	controller := NewController(cfg)
	us.ctrl = controller

	router := gin.Default()
	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	us.server = struct {
		HTTPServer *http.Server
		Router     *gin.Engine
	}{srv, router}

	us.server.Router.GET(
		"/comics",
		controller.GetComics,
	)
}

func TestControllerUnitSuite(t *testing.T) {
	suite.Run(t, &ControllerUnitSuite{})
}

// =============================================================================
// COMICS ENDPOINT TEST
// =============================================================================

func (us *ControllerUnitSuite) TestComicsEndpoint() {
	testCases := []struct {
		name             string
		query            string
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			"Normal Operation",
			"?start=40&end=42",
			200,
			struct {
				Comics []Comic `json:"comics"`
			}{
				Comics: []Comic{
					{
						Num:        42,
						Title:      "Geico",
						SafeTitle:  "Geico",
						Day:        "1",
						Month:      "1",
						Year:       "2006",
						Transcript: "I just saved a bunch of money on my car insurance by threatening my agent with a golf club.\n{{title text: David did this}}",
						Img:        "https://imgs.xkcd.com/comics/geico.jpg",
						Alt:        "David did this",
						News:       "",
						Link:       "",
					},
					{
						Num:        40,
						Title:      "Light",
						SafeTitle:  "Light",
						Day:        "1",
						Month:      "1",
						Year:       "2006",
						Transcript: "[[A crowd of figures stand around in the dark. One figure is illuminated by a beam of light.]]\nIn a dark and confusing world, you burn brightly. I never feel lost.\n{{Alt-text: Like a beacon.}}",
						Img:        "https://imgs.xkcd.com/comics/light.jpg",
						Alt:        "Like a beacon",
						News:       "",
						Link:       "",
					},
					{
						Num:        41,
						Title:      "Old Drawing",
						SafeTitle:  "Old Drawing",
						Day:        "1",
						Month:      "1",
						Year:       "2006",
						Transcript: "[[A tree holding a chainsaw over a recently cut-down tree.]]\nI found this in one of my high-school notebooks. I think I drew it just to take revenge on people snooping through my stuff.\nCut-down tree: WELL, YOU STUMPED ME...\n{{I don't want to talk about it}}",
						Img:        "https://imgs.xkcd.com/comics/unspeakable_pun.jpg",
						Alt:        "I don't want to talk about it",
						News:       "",
						Link:       "",
					},
				},
			},
		},
		{
			"Missing Query Parameters 1",
			"",
			400,
			&errorBody{
				Msg: "please include the starting and ending comic numbers",
			},
		},
		{
			"Missing Query Parameters 2",
			"?start=3",
			400,
			&errorBody{
				Msg: "please include the starting and ending comic numbers",
			},
		},
		{
			"Missing Query Parameters 3",
			"?end=7",
			400,
			&errorBody{
				Msg: "please include the starting and ending comic numbers",
			},
		},
		{
			"Invalid Query Parameters 1",
			"?start=now&end=42",
			400,
			&errorBody{
				Msg: "please make sure that the starting number is an integer",
			},
		},
		{
			"Invalid Query Parameters 2",
			"?start=1&end=never",
			400,
			&errorBody{
				Msg: "please make sure that the ending number is an integer",
			},
		},
	}

	for i := range testCases {
		test := testCases[i]
		us.Run(test.name, func() {
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/comics%s", test.query), nil)
			us.Nil(err)

			us.server.Router.ServeHTTP(recorder, request)
			us.Equal(test.expectedStatus, recorder.Code)

			var expectedBytes []byte
			if test.expectedResponse == nil {
				expectedBytes = []byte(nil)
			} else {
				expectedBytes, err = json.Marshal(test.expectedResponse)
				us.Nil(err)
			}
			us.Equal(expectedBytes, recorder.Body.Bytes())
		})
	}
}
