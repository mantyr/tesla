package tesla

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClientSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()
	previousAuthURL := AuthURL
	previousURL := BaseURL
	AuthURL = ts.URL + "/oauth/token"
	BaseURL = ts.URL + "/api/1"

	auth := &Auth{
		GrantType:    "password",
		ClientID:     "abc123",
		ClientSecret: "def456",
		Email:        "elon@tesla.com",
		Password:     "go",
	}
	client, err := NewClient(auth)

	Convey("Should set the HTTP headers", t, func() {
		req, _ := http.NewRequest("GET", "http://foo.com", nil)
		client.setHeaders(req)
		So(req.Header.Get("Authorization"), ShouldEqual, "Bearer ghi789")
		So(req.Header.Get("Accept"), ShouldEqual, "application/json")
		So(req.Header.Get("Content-Type"), ShouldEqual, "application/json")
	})
	Convey("Should login and get an access token", t, func() {
		So(err, ShouldBeNil)
		So(client.Token.AccessToken, ShouldEqual, "ghi789")
	})

	AuthURL = previousAuthURL
	BaseURL = previousURL
}

func serveHTTP(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		req.Body.Close()
		Convey("HTTP headers should be present", t, func() {
			So(req.Header["Accept"][0], ShouldEqual, "application/json")
			So(req.Header["Content-Type"][0], ShouldEqual, "application/json")
		})
		switch req.URL.String() {
		case "/oauth/token":
			Convey("Request body should be set correctly", t, func() {
				auth := &Auth{}
				json.Unmarshal(body, auth)
				So(auth.ClientID, ShouldEqual, "abc123")
				So(auth.ClientSecret, ShouldEqual, "def456")
				So(auth.Email, ShouldEqual, "elon@tesla.com")
				So(auth.Password, ShouldEqual, "go")
				So(auth.URL, ShouldEqual, BaseURL)
				So(auth.StreamingURL, ShouldEqual, StreamingURL)
			})
			w.WriteHeader(200)
			w.Write([]byte("{\"access_token\": \"ghi789\"}"))
		case "/api/1/vehicles":
			w.WriteHeader(200)
			w.Write([]byte(VehiclesJSON))
		case "/api/1/vehicles/1234/mobile_enabled":
			w.WriteHeader(200)
			w.Write([]byte(TrueJSON))
		case "/api/1/vehicles/1234/data_request/charge_state":
			w.WriteHeader(200)
			w.Write([]byte(ChargeStateJSON))
		case "/api/1/vehicles/1234/data_request/climate_state":
			w.WriteHeader(200)
			w.Write([]byte(ClimateStateJSON))
		case "/api/1/vehicles/1234/data_request/drive_state":
			w.WriteHeader(200)
			w.Write([]byte(DriveStateJSON))
		case "/api/1/vehicles/1234/data_request/gui_settings":
			w.WriteHeader(200)
			w.Write([]byte(GuiSettingsJSON))
		case "/api/1/vehicles/1234/data_request/vehicle_state":
			w.WriteHeader(200)
			w.Write([]byte(VehicleStateJSON))
		case "/api/1/vehicles/1234/wake_up":
			w.WriteHeader(200)
			w.Write([]byte(WakeupResponseJSON))
		case "/api/1/vehicles/1234/command/charge_standard":
			w.WriteHeader(200)
			w.Write([]byte(ChargeAlreadySetJSON))
		case "/api/1/vehicles/1234/command/charge_start":
			w.WriteHeader(200)
			w.Write([]byte(ChargedJSON))
		case "/api/1/vehicles/1234/command/charge_stop",
			"/api/1/vehicles/1234/command/charge_max_range",
			"/api/1/vehicles/1234/command/charge_port_door_open",
			"/api/1/vehicles/1234/command/flash_lights",
			"/api/1/vehicles/1234/command/honk_horn",
			"/api/1/vehicles/1234/command/auto_conditioning_start",
			"/api/1/vehicles/1234/command/auto_conditioning_stop",
			"/api/1/vehicles/1234/command/door_lock",
			"/api/1/vehicles/1234/command/set_temps?driver_temp=72&passenger_temp=72",
			"/api/1/vehicles/1234/command/remote_start_drive?password=foo":
			w.WriteHeader(200)
			w.Write([]byte(CommandResponseJSON))
		}
	}))
}
