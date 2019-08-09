package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/codykrieger/microbridge/micropub"
	"github.com/codykrieger/microbridge/xmlrpc"
	log "github.com/sirupsen/logrus"
)

type WPService struct {
	config *Config
}

func (s *WPService) checkAuth(username, password string) error {
	if username == "" || password == "" {
		return xmlrpc.ErrForbidden
	}

	client := micropub.NewClient(s.config.MicropubEndpoint, password)
	config, err := client.GetConfig()
	if err != nil {
		return err
	}

	if len(config.Destination) == 0 {
		// FIXME: Micro.blog almost exclusively returns 200 OK, even when the
		// bearer token is flat-out wrong. So we use the presence of a
		// Destination in the config object to determine whether the included
		// bearer token successfully authenticated the client.
		log.Error("micropub config contains no destinations; assuming authentication failure")
		return xmlrpc.ErrForbidden
	}

	return nil
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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

	reply.Users = []User{
		User{
			UserID:      "1",
			Username:    "You",
			DisplayName: "You",
		},
	}

	return nil
}

type GetAuthorsArgs struct {
	BlogID   string
	Username string
	Password string
}

type GetAuthorsReply struct {
	Users []User
}

func (s *WPService) GetAuthors(req *http.Request, args *GetAuthorsArgs, reply *GetAuthorsReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
	}).Info("---> wp.GetAuthors")

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

	client := micropub.NewClient(s.config.MicropubEndpoint, args.Password)

	categories, err := client.GetCategories()
	if err != nil {
		return err
	}

	reply.Categories = []Category{}

	for i, v := range categories {
		reply.Categories = append(reply.Categories, Category{
			CategoryID: fmt.Sprintf("%d", i),
			Name:       v,
		})
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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

	if args.Filter.PostType != "post" {
		return nil
	}

	client := micropub.NewClient(s.config.MicropubEndpoint, args.Password)

	posts, err := client.GetPosts()
	if err != nil {
		return err
	}

	reply.Posts = []Post{}

	for _, v := range posts {
		status := v.Properties.PostStatus[0]
		wpStatus := ""
		switch status {
		case "published":
			wpStatus = "publish"
		case "draft":
			wpStatus = "draft"
		default:
			log.Warnf("unknown micropub post status '%s'", status)
			wpStatus = "publish"
		}

		dateString := v.Properties.Published[0]
		date, err := time.Parse(time.RFC3339, dateString) //"2006-01-02T15:04:05-07:00", dateString)
		if err != nil {
			return err
		}

		date = date.Local()
		log.WithField("d", date).Infof("parsed date (%s)", dateString)

		reply.Posts = append(reply.Posts, Post{
			PostID:        fmt.Sprintf("%d", v.Properties.UID[0]),
			Title:         v.Properties.Name[0],
			Date:          date,
			DateModified:  date,
			Status:        wpStatus,
			Type:          "post",
			Format:        "standard",
			Name:          "",
			Author:        "1",
			Content:       v.Properties.Content[0],
			Parent:        "0",
			MIMEType:      "text/plain",
			Link:          v.Properties.URL[0],
			CommentStatus: "closed",
			PingStatus:    "closed",
			Sticky:        false,
			Terms:         []Term{},
			CustomFields:  []CustomField{},
		})
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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

	return xmlrpc.ErrNotImplemented
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

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

	return nil
}

type NewMediaObjectArgs struct {
	BlogID   string
	Username string
	Password string
	Object   struct {
		Name string `xml:"name"`
		Bits string `xml:"bits"`
		Type string `xml:"type"`
	}
}

type NewMediaObjectReply struct {
}

func (s *WPService) NewMediaObject(req *http.Request, args *NewMediaObjectArgs, reply *NewMediaObjectReply) error {
	log.WithFields(log.Fields{
		"bid": args.BlogID,
		"u":   args.Username,
	}).Info("---> metaWeblog.newMediaObject")

	if err := s.checkAuth(args.Username, args.Password); err != nil {
		return err
	}

	log.Infof("object: %s; type: %s", args.Object.Name, args.Object.Type)

	return nil
}
