package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/mainflux/mainflux-auth-server/domain"
	"github.com/mainflux/mainflux-auth-server/services"
)

func TestCheckCredentials(t *testing.T) {
	var (
		devId  string            = "test-device"
		chanId string            = "test-chan"
		spec   domain.AccessSpec = domain.AccessSpec{[]domain.Scope{domain.Scope{"R", "channel", chanId}, domain.Scope{"RW", "device", devId}}}
	)

	key, _ := services.AddUserKey(user.Id, user.MasterKey, spec)

	cases := []struct {
		body string
		code int
	}{
		{fmt.Sprintf(`{"action":"R","resource":"user","id":"%s","key":"%s"}`, user.Id, user.MasterKey), 200},
		{`malformed body`, 400},
		{fmt.Sprintf(`{"action":"bad","resource":"user","id":"%s","key":"%s"}`, user.Id, user.MasterKey), 400},
		{fmt.Sprintf(`{"action":"R","resource":"bad","id":"%s","key":"%s"}`, user.Id, user.MasterKey), 400},
		{fmt.Sprintf(`{"action":"R","resource":"user","key":"%s"}`, user.MasterKey), 400},
		{fmt.Sprintf(`{"action":"R","resource":"user","id":"%s"}`, user.Id), 400},

		{fmt.Sprintf(`{"action":"R","resource":"user","id":"%s","key":"%s"}`, user.Id, key), 403},
		{fmt.Sprintf(`{"action":"R","resource":"device","id":"%s","key":"%s"}`, devId, key), 200},
		{fmt.Sprintf(`{"action":"W","resource":"device","id":"%s","key":"%s"}`, devId, key), 200},
		{fmt.Sprintf(`{"action":"R","resource":"channel","id":"%s","device":"%s","key":"%s"}`, chanId, devId, key), 200},
		{fmt.Sprintf(`{"action":"R","resource":"channel","id":"%s","key":"%s"}`, chanId, key), 400},
		{fmt.Sprintf(`{"action":"X","resource":"channel","id":"%s","device":"%s","key":"%s"}`, chanId, devId, key), 403},
	}

	url := fmt.Sprintf("%s/access-checks", ts.URL)

	for i, c := range cases {
		b := strings.NewReader(c.body)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}
	}
}
