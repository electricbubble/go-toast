package toast

import (
	"testing"
)

func TestPush(t *testing.T) {
	checkErr(t, Push("test_message"))
	// checkErr(t, Push("test_message", WithSubtitle("test_subtitle")))
	// checkErr(t, Push("test_message", WithSubtitle("test_subtitle"), WithTitle("test_title")))
	// checkErr(t, Push("test_message", WithAudio(Ping)))
	// checkErr(t, Push("test_message", WithSubtitle("test_subtitle"), WithAudio(Ping), WithObjectiveC()))
	checkErr(t, Push("test_message", WithSubtitle("test_subtitle"), WithObjectiveC()))

	// checkErr(t, Push("test_message", WithSubtitle("test_subtitle"), WithFakeBundleID("com.apple.Safari")))
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
