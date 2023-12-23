package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestRealMain(t *testing.T) {
  var ret string = realMain([]string{"Program Name"})
  assert.Equal(t, ret, "Please provide a mode Either `client` or `server`.", "Should be equal")
  
  ret = realMain([]string{"Program Name", "1"})
  assert.Equal(t, ret, "Unknown mode", "Should be equal")



}
