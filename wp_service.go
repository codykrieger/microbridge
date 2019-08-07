package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type WPService struct {
	config *Config
}

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
