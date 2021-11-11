package toast

import (
	"testing"
)

func TestPush(t *testing.T) {
	checkErr(t, Push("test_message"))
	checkErr(t, Push("test_message", WithAppID("test_AppID")))
	checkErr(t, Push("test_message", WithAppID("test_AppID"), WithTitle("test_title")))
	checkErr(t, Push("test_message", WithAudio(Default)))
	checkErr(t, Push("test_message", WithAudio(Default), WithProtocolAction("click me")))
	checkErr(t, Push("test_message", WithProtocolAction("Open Maps", "bingmaps:?q=beijing")))
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
