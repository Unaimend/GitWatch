package main

import "testing"

func TestRealMain(t *testing.T) {
  var ret = realMain([]string{"Program Name", "1"})
  if ret != "Unknown mode" {
    t.Error("Wrong error message (1)")
  }
}
