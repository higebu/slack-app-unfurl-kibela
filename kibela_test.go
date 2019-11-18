package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	testKibelaTeam  = ""
	testKibelaToken = ""
)

func getTestEnv(t *testing.T) {
	testKibelaTeam = os.Getenv("KIBELA_TEAM")
	if testKibelaTeam == "" {
		t.Fatal("KIBELA_TEAM is empty")
	}
	testKibelaToken = os.Getenv("KIBELA_TOKEN")
	if testKibelaToken == "" {
		t.Fatal("KIBELA_TOKEN is empty")
	}
}

func TestNoteFromPath(t *testing.T) {
	getTestEnv(t)
	kibela := NewKibelaClient(testKibelaTeam, testKibelaToken)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	path := fmt.Sprintf("https://%s.kibe.la/notes/1", testKibelaTeam)
	res, err := kibela.NoteFromPath(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if res.Note.URL != path {
		t.Fatalf("got %s want %s", res.Note.URL, path)
	}
	t.Logf("%#v", res)
}
