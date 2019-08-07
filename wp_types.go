package main

import (
	"time"
)

type Category struct {
	CategoryID          string `xml:"categoryId"`
	ParentID            string `xml:"parentId"`
	Name                string `xml:"categoryName"`
	Description         string `xml:"description"`
	CategoryDescription string `xml:"categoryDescription"`
}

type CustomField struct {
	ID    string `xml:"id"`
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

type Enclosure struct {
	URL    string `xml:"url"`
	Length int    `xml:"length"`
	Type   string `xml:"type"`
}

type Term struct {
	ID   string `xml:"term_id"`
	Name string `xml:"name"`
	Slug string `xml:"slug"`
}

type PostThumbnail struct {
	AttachmentID   string    `xml:"attachment_id"`
	DateCreatedGMT time.Time `xml:"date_created_gmt"`
	ParentID       int       `xml:"parent_id"`
	Link           string    `xml:"link"`
	Title          string    `xml:"title"`
	Caption        string    `xml:"caption"`
	Description    string    `xml:"description"`
}

type Post struct {
	PostID      string    `xml:"post_id"`
	Title       string    `xml:"post_title"`
	DateCreated time.Time `xml:"post_date"`
	// DateCreatedGMT  time.Time      `xml:"post_date_gmt"`
	DateModified time.Time `xml:"post_modified"`
	// DateModifiedGMT time.Time      `xml:"post_modified_gmt"`
	Status        string `xml:"post_status"`
	Type          string `xml:"post_type"`
	Format        string `xml:"post_format"`
	Password      string `xml:"post_password"`
	Name          string `xml:"post_name"` // note: url-safe slug
	Author        string `xml:"post_author"`
	Excerpt       string `xml:"post_excerpt"`
	Content       string `xml:"post_content"`
	Parent        string `xml:"post_parent"`
	MIMEType      string `xml:"post_mime_type"`
	Link          string `xml:"link"`
	GUID          string `xml:"guid"`
	MenuOrder     int    `xml:"menu_order"`
	CommentStatus string `xml:"comment_status"`
	PingStatus    string `xml:"ping_status"`
	Sticky        bool   `xml:"sticky"`
	// PostThumbnail PostThumbnail `xml:"post_thumbnail"`

	Terms        []Term        `xml:"terms"`
	CustomFields []CustomField `xml:"custom_fields"`
	Enclosure    Enclosure     `xml:"enclosure"`
}

type Tag struct {
	ID    int    `xml:"tag_id"`
	Name  string `xml:"name"`
	Slug  string `xml:"slug"`
	Count int    `xml:"count"`
}

type User struct {
	UserID      string    `xml:"user_id"`
	Username    string    `xml:"username"`
	FirstName   string    `xml:"first_name"`
	LastName    string    `xml:"last_name"`
	Bio         string    `xml:"bio"`
	Email       string    `xml:"email"`
	Nickname    string    `xml:"nickname"`
	Nicename    string    `xml:"nicename"`
	URL         string    `xml:"url"`
	DisplayName string    `xml:"display_name"`
	Registered  time.Time `xml:"registered"`
}