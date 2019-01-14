package controllers

//用户表结构体
type User struct {
	UserId         int64  `json:"user_id"`
	UserLogin      string `json:"user_login"`
	UserPass       string `json:"user_pass"`
	UserNicename   string `json:"user_nicename"`
	UserEmail      string `json:"user_email"`
	UserRegistered string `json:"user_registered"`
	UserStatus     string `json:"user_status"`
}

//分类项表结构体
type Term struct {
	TermId   int64  `json:"term_id"`
	TermName string `json:"term_name"`
}

//分类法表结构体
type TermTaxonomy struct {
	TermTaxonomyId int64  `json:"term_taxonomy_id"`
	Taxonomy       string `json:"taxonomy"`
	Description    string `json:"description"`
	TermParentId   string `json:"term_parent_id"`
}

//文章结构体
type Post struct {
	PostId        int64  `json:"post_id"`
	PostAuthor    string `json:"post_author"`
	PostContent   string `json:"post_content"`
	PostTitle     string `json:"post_title"`
	PostDate      string `json:"post_date"`
	PostPreImgUrl string `json:"post_pre_img_url"`
	CommentCount  string `json:"comment_count"`
	PostExcerpt   string `json:"post_excerpt"`
	PostStatus    int    `json:"post_status"`
	PostModified  string `json:"post_modified"`
	CommentStatus int    `json:"comment_status"`
}

//评论结构体
type Comment struct {
	CommentId          int64  `json:"comment_id"`
	CommentAuthor      string `json:"comment_author"`
	CommentAuthorEmail string `json:"comment_author_email"`
	CommentAuthorIP    string `json:"comment_author_ip"`
	CommentDate        string `json:"comment_date"`
	CommentContent     string `json:"comment_content"`
	CommentApproved    int64  `json:"comment_approved"`
}

//链接
type Link struct {
	LinkId          int64  `json:"link_id"`
	LinkUrl         string `json:"link_url"`
	LinkName        string `json:"link_name"`
	LinkImageUrl    string `json:"link_image_url"`
	LinkDescription string `json:"link_description"`
	LinkVisible     int    `json:"link_visible"`
	LinkRating      int    `json:"link_rating"`
	LinkUpdated     string `json:"link_updated"`
}

