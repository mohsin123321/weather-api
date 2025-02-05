package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"weather-api/internal/server"
	mock_adapters "weather-api/tests/mocks/adapters"
	mock_persistence "weather-api/tests/mocks/persistence"

	"go.uber.org/mock/gomock"
)

// tester is a struct that holds the mocks and the server instance.
type tester struct {
	OpenWeatherMapAdapter *mock_adapters.MockOpenWeatherMapAdapter
	CacheMap              *mock_persistence.MockCacheMap
	Server                *server.Server
}

// errorMessage is a struct that holds the error message.
type errorMessage struct {
	Message string `json:"message"`
}

// setup is a helper function to setup the tester struct.
func setup(t *testing.T) *tester {
	t.Parallel()
	t.Helper()

	ctrl := gomock.NewController(t)
	mocksAdapters := mock_adapters.NewMockOpenWeatherMapAdapter(ctrl)
	mocksPersistence := mock_persistence.NewMockCacheMap(ctrl)

	return &tester{
		OpenWeatherMapAdapter: mocksAdapters,
		CacheMap:              mocksPersistence,
		Server:                server.NewServer("8080", mocksPersistence, mocksAdapters),
	}
}

// mockReqGenerator is a helper function to generate a mock request and response recorder.
func mockReqGenerator(url string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	return req, w
}

// unmarshalResponseBody is a helper function to unmarshal the response body into a given struct.
func unmarshalResponseBody(result interface{}, resp *http.Response, t *testing.T) {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %v", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}

func TestHomePage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tester := setup(t)

	t.Run("Failed template rendering", func(t *testing.T) {
		dir, err := os.Getwd()
		if err != nil {
			panic("Failed to get working directory: " + err.Error())
		}

		// set to wrong directory to fail the template rendering
		if err := os.Chdir(".."); err != nil {
			panic("Failed to set working directory: " + err.Error())
		}

		req, w := mockReqGenerator("/")
		tester.Server.HomePage(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %d", resp.StatusCode)
		}

		t.Cleanup(func() {
			// set back to the original directory
			if err := os.Chdir(dir); err != nil {
				panic("Failed to set working directory: " + err.Error())
			}
		})
	})

	t.Run("Successful template rendering", func(t *testing.T) {
		req, w := mockReqGenerator("/")
		tester.Server.HomePage(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}
	})
}

func TestMain(m *testing.M) {
	// Change working directory to the project root
	// useful for the template rendering for tests
	if err := os.Chdir("../.."); err != nil {
		panic("Failed to set working directory: " + err.Error())
	}

	// Run the tests and exit with the result code
	os.Exit(m.Run())
}
