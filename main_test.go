package main

import (
	"fmt"
	"testing"
)

func TestGetToken(t *testing.T) {
	getToken("")
}

func TestGetUserInfo(t *testing.T) {
	userInfo := getUser("")
	fmt.Println(userInfo)
}
