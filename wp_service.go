package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type WPService struct{}

type GetUsersArgs struct {
	BlogID   string
	Username string
	Password string
	Filter   struct {
		Who string `xml:"who"`
	}
}

type GetUsersReply struct {
	Users []User
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

func (s *WPService) GetUsers(req *http.Request, args *GetUsersArgs, reply *GetUsersReply) error {
	log.WithFields(log.Fields{
		"bid":    args.BlogID,
		"u":      args.Username,
		"filter": args.Filter,
	}).Info("---> wp.GetUsers")

	reply.Users = []User{
		User{
			UserID:      "1",
			Username:    "cjk",
			FirstName:   "Cody",
			LastName:    "Krieger",
			Bio:         "",
			Email:       "cody@krieger.io",
			Nickname:    "cody-nickname",
			Nicename:    "cody-nicename",
			URL:         "http://localhost:4567/users/cjk",
			DisplayName: "cody-displayname",
		},
	}

	return nil
}

type GetCategoriesArgs struct {
	BlogID   string
	Username string
	Password string
}

type GetCategoriesReply struct {
	Categories []Category
}

type Category struct {
	CategoryID          string `xml:"categoryId"`
	ParentID            string `xml:"parentId"`
	Name                string `xml:"categoryName"`
	Description         string `xml:"description"`
	CategoryDescription string `xml:"categoryDescription"`
	HTMLURL             string `xml:"htmlUrl"`
	RSSURL              string `xml:"rssUrl"`
}

func (s *WPService) GetCategories(req *http.Request, args *GetCategoriesArgs, reply *GetCategoriesReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
	}).Info("---> wp.GetCategories")

	reply.Categories = []Category{
		// Category{
		//     CategoryID:          "1",
		//     Name:                "foo(1)",
		//     Description:         "foo desc",
		//     CategoryDescription: "foo(2)",
		// },
		// Category{
		//     CategoryID:          "2",
		//     Name:                "bar(1)",
		//     Description:         "bar desc",
		//     CategoryDescription: "bar(2)",
		// },
	}

	return nil
}

type GetPostsArgs struct {
	BlogID   string
	Username string
	Password string
	Filter   struct {
		PostType   string `xml:"post_type"`
		PostStatus string `xml:"post_status"`
		Number     int    `xml:"number"`
		Offset     int    `xml:"offset"`
		OrderBy    string `xml:"orderby"`
		Order      string `xml:"order"`
	}
	Fields []string
}

type GetPostsReply struct {
	Posts []Post
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

func (s *WPService) GetPosts(req *http.Request, args *GetPostsArgs, reply *GetPostsReply) error {
	log.WithFields(log.Fields{
		"bid":    args.BlogID,
		"u":      args.Username,
		"filter": args.Filter,
		"Fields": args.Fields,
	}).Info("---> wp.GetPosts")

	if args.Filter.PostType == "post" {
		now := time.Now()

		reply.Posts = []Post{
			Post{
				PostID:        "1",
				Title:         "Title!",
				DateCreated:   now,
				DateModified:  now,
				Status:        "publish",
				Type:          "post",
				Format:        "standard",
				Name:          "title",
				Author:        "1",
				Excerpt:       "excerpt one",
				Content:       "content one!",
				Parent:        "0",
				MIMEType:      "text/plain",
				Link:          "http://localhost:4567/posts/title",
				GUID:          "b87c926c-377a-4a73-9609-fa1edd1f891e",
				MenuOrder:     0,
				CommentStatus: "closed",
				PingStatus:    "closed",
				Sticky:        false,
				// PostThumbnail:   PostThumbnail{},

				Terms:        []Term{},
				CustomFields: []CustomField{},
			},
			Post{
				PostID:        "2",
				Title:         "Title 2!!",
				DateCreated:   now,
				DateModified:  now,
				Status:        "draft",
				Type:          "post",
				Format:        "standard",
				Name:          "title-2",
				Author:        "1",
				Excerpt:       "excerpt two",
				Content:       "content two!",
				Parent:        "0",
				MIMEType:      "text/plain",
				Link:          "http://localhost:4567/posts/title-2",
				GUID:          "cb146c10-0294-4dc4-a578-84dd1b98d3c0",
				MenuOrder:     0,
				CommentStatus: "closed",
				PingStatus:    "closed",
				Sticky:        false,
				// PostThumbnail:   PostThumbnail{},

				Terms:        []Term{},
				CustomFields: []CustomField{},
			},
		}
	}

	return nil
}

type GetTagsArgs struct {
	BlogID   string
	Username string
	Password string
}

type GetTagsReply struct {
	Tags []Tag
}

type Tag struct {
	ID      int    `xml:"tag_id"`
	Name    string `xml:"name"`
	Slug    string `xml:"slug"`
	Count   int    `xml:"count"`
	HTMLURL string `xml:"html_url"`
	RSSURL  string `xml:"rss_url"`
}

func (s *WPService) GetTags(req *http.Request, args *GetTagsArgs, reply *GetTagsReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
	}).Info("---> wp.GetTags")

	reply.Tags = []Tag{
		// Tag{
		//     ID:      1,
		//     Name:    "Tag A!",
		//     Slug:    "tag-a",
		//     Count:   0,
		//     HTMLURL: "http://localhost:4567/tags/tag-a",
		//     RSSURL:  "http://localhost:4567/tags/tag-a",
		// },
		// Tag{
		//     ID:      2,
		//     Name:    "Tag B!",
		//     Slug:    "tag-b",
		//     Count:   0,
		//     HTMLURL: "http://localhost:4567/tags/tag-b",
		//     RSSURL:  "http://localhost:4567/tags/tag-b",
		// },
	}

	return nil
}
