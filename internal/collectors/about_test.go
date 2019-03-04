package collectors

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var (
	mockAboutData = `{ "data": { "accounts_active": 99, "subscribers": 1 } }`
)

func setup() {
	httpmock.Activate()

	httpmock.RegisterResponder("GET", "https://api.reddit.com/r/pass/about.json",
		httpmock.NewStringResponder(200, mockAboutData),
	)
	httpmock.RegisterResponder("GET", "https://api.reddit.com/r/fail/about.json",
		httpmock.NewStringResponder(469, "Error"),
	)
}

func TestGet(t *testing.T) {
	setup()

	defer httpmock.DeactivateAndReset()

	// pass
	c := NewAboutSubredditCollector("pass", new(http.Client))
	info, _ := c.get()
	assert.Equal(t, float64(99), info.AccountsActive)
	assert.Equal(t, float64(1), info.Subscribers)

	//fail
	c = NewAboutSubredditCollector("fail", new(http.Client))
	info, err := c.get()
	assert.NotNil(t, err)
	assert.Nil(t, info)
}

func TestCollect(t *testing.T) {
	setup()

	defer httpmock.DeactivateAndReset()
	c := NewAboutSubredditCollector("pass", new(http.Client))
	prom.MustRegister(c)

	testAboutData, _ := os.Open("about_test_data.txt")
	assert.Nil(t, testutil.GatherAndCompare(prom.DefaultGatherer, testAboutData, "subreddit_active_users", "subreddit_subscriber_users"))
}
