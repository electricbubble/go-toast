package toast

import (
	"testing"
)

func TestPush(t *testing.T) {
	checkErr(t, Push("test_message"))
	checkErr(t, Push("test_message", WithSubtitle("test_subtitle")))
	checkErr(t, Push("test_message", WithSubtitle("test_subtitle"), WithTitle("test_title")))
	checkErr(t, Push("test_message", WithAudio(Ping)))
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
