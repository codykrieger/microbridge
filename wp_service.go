package main

import (
	"fmt"
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
			Username:    s.config.Username,
			DisplayName: s.config.Username,
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
		Category{
			CategoryID: "1",
			ParentID:   "0",
			Name:       "foo",
		},
		Category{
			CategoryID: "2",
			ParentID:   "0",
			Name:       "bar",
		},
		Category{
			CategoryID: "3",
			ParentID:   "2",
			Name:       "baz",
		},
	}

	return nil
}

type NewCategoryArgs struct {
	BlogID   string
	Username string
	Password string
	Category struct {
		Name string `xml:"name"`
	}
}

type NewCategoryReply struct {
	CategoryID int
}

func (s *WPService) NewCategory(req *http.Request, args *NewCategoryArgs, reply *NewCategoryReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
	}).Info("---> wp.NewCategory")

	// args.Category.Name

	reply.CategoryID = 999

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
		"fields": args.Fields,
	}).Info("---> wp.GetPosts")

	if args.Filter.PostType == "post" {
		now := time.Now()

		reply.Posts = []Post{
			Post{
				PostID:        "1",
				Title:         "Title!",
				Date:          now,
				DateModified:  now,
				Status:        "publish",
				Type:          "post",
				Format:        "standard",
				Name:          "title",
				Author:        "1",
				Content:       "content one!",
				Parent:        "0",
				MIMEType:      "text/plain",
				Link:          s.config.PostsURL + "/title",
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
				Date:          now,
				DateModified:  now,
				Status:        "draft",
				Type:          "post",
				Format:        "standard",
				Name:          "title-2",
				Author:        "1",
				Content:       "content two!",
				Parent:        "0",
				MIMEType:      "text/plain",
				Link:          s.config.PostsURL + "/title-2",
				GUID:          "cb146c10-0294-4dc4-a578-84dd1b98d3c0",
				MenuOrder:     0,
				CommentStatus: "closed",
				PingStatus:    "closed",
				Sticky:        false,
				// PostThumbnail:   PostThumbnail{},

				Terms: []Term{
					Term{ID: "1", Name: "footag", Taxonomy: "post_tag"},
					Term{ID: "2", Name: "foo", Taxonomy: "category"},
				},
				CustomFields: []CustomField{},
			},
		}
	}

	return nil
}

type EditPostArgs struct {
	BlogID   string
	Username string
	Password string
	PostID   string
	Content  Post
}

type EditPostReply struct {
	Success bool
}

func (s *WPService) EditPost(req *http.Request, args *EditPostArgs, reply *EditPostReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
		"pid": args.PostID,
	}).Info("---> wp.EditPost")

	// args.Content.Title
	// args.Content.Status
	// args.Content.Content
	// args.Content.Date
	// args.Content.Terms (categories)
	// args.Content.TermsNames (tags)
	// args.Content.Name (?)
	// args.Content.Enclosure (optional; don't need this yet)

	reply.Success = true

	return nil
}

type NewPostArgs struct {
	BlogID   string
	Username string
	Password string
	Content  Post
}

type NewPostReply struct {
	PostID string
}

func (s *WPService) NewPost(req *http.Request, args *NewPostArgs, reply *NewPostReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
	}).Info("---> wp.NewPost")

	// args.Content.Title
	// args.Content.Status
	// args.Content.Content
	// args.Content.Date (optional, apparently)
	// args.Content.Terms (categories)
	// args.Content.TermsNames (tags)
	// args.Content.Name (?)
	// args.Content.Enclosure (optional; don't need this yet)

	reply.PostID = "999"

	return nil
}

type GetPostArgs struct {
	BlogID   string
	Username string
	Password string
	PostID   string
	Fields   []string
}

type GetPostReply struct {
	Post Post
}

func (s *WPService) GetPost(req *http.Request, args *GetPostArgs, reply *GetPostReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
		"pid": args.PostID,
	}).Info("---> wp.GetPost")

	// FIXME: Huh. Well, this doesn't seem to result in an error being returned to the client.
	// ಠ_ಠ
	return fmt.Errorf("not implemented")
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
		Tag{
			ID:   1,
			Name: "footag",
		},
		Tag{
			ID:   2,
			Name: "bartag",
		},
	}

	return nil
}
