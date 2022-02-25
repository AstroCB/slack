package slack

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

type addBookmarkResponseFull struct {
	Bookmark `json:"bookmark"`
	SlackResponse
}

type Bookmark struct {
	ID                  string `json:"id"`
	DateCreated         int    `json:"date_created"`
	DateUpdated         int    `json:"date_updated"`
	Rank                string `json:"rank"`
	LastUpdatedByUserID string `json:"last_updated_by_user_id"`
	LastUpdatedByTeamID string `json:"last_updated_by_team_id"`
	ShortcutID          string `json:"shortcut_id"`
	AppID               string `json:"app_id"`
	BookmarkOpts
}

type BookmarkOpts struct {
	ChannelID string `json:"channel_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Emoji     string `json:"emoji"`     // optional
	EntityID  string `json:"entity_id"` // optional, only applies to `message` and `file` types
	Link      string `json:"link"`      // optional
	ParentID  string `json:"parent"`    // optional
}

func (api *Client) AddBookmark(opts BookmarkOpts) (Bookmark, error) {
	return api.AddBookmarkContext(context.Background(), opts)
}

func (api *Client) AddBookmarkContext(ctx context.Context, opts BookmarkOpts) (Bookmark, error) {
	if (opts.Type == "message" || opts.Type == "file") && opts.EntityID == "" {
		return Bookmark{}, fmt.Errorf("bookmark of type %v requires an entity_id to be specified", opts.Type)
	}

	if opts.Type == "link" && opts.Link == "" {
		return Bookmark{}, errors.New("bookmark of type link requires a link to be specified")
	}

	values := url.Values{
		"channel_id": {opts.ChannelID},
		"title":      {opts.Title},
		"type":       {opts.Type},
	}
	if opts.Emoji != "" {
		values.Set("emoji", opts.Emoji)
	}
	if opts.EntityID != "" {
		values.Set("entity_id", opts.EntityID)
	}
	if opts.Link != "" {
		values.Set("link", opts.Link)
	}
	if opts.ParentID != "" {
		values.Set("parent_id", opts.ParentID)
	}

	response := &addBookmarkResponseFull{}
	err := api.postMethod(ctx, "bookmarks.add", values, response)
	if err != nil {
		return Bookmark{}, err
	}
	if !response.Ok {
		return Bookmark{}, errors.New(response.Error)
	}

	return response.Bookmark, nil
}
