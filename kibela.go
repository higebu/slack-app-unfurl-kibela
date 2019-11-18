package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

const KibelaTimeFormat = "2006-01-02T15:04:05.000-07:00"

var (
	team         = ""
	endpointBase = "https://%s.kibe.la/api/v1"
)

var query = `query ($path: String!) {
  note: noteFromPath(path: $path) {
    author {
      id
      account
      avatarImage {
        url
      }
      url
    }
    id
    title
    url
    publishedAt
    summary: contentSummaryHtml
  }
}`

type Response struct {
	Note Note `json:"note"`
}

type Note struct {
	Author      Author `json:"author"`
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
	Summary     string `json:"summary"`
}

type Author struct {
	ID          string      `json:"id"`
	Account     string      `json:"account"`
	AvatarImage AvatarImage `json:"avatarImage"`
	URL         string      `json:"url"`
}

type AvatarImage struct {
	URL string `json:"url"`
}

type KibelaClient struct {
	c     *graphql.Client
	token string
}

func NewKibelaClient(team, token string) *KibelaClient {
	client := graphql.NewClient(fmt.Sprintf(endpointBase, team))
	return &KibelaClient{c: client, token: token}
}

func (c *KibelaClient) NoteFromPath(ctx context.Context, path string) (*Response, error) {
	req := graphql.NewRequest(query)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Var("path", path)
	var res Response
	if err := c.c.Run(ctx, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
