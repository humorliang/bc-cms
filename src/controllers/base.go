package controllers

//用户表结构体
type User struct {
	UserId         int64  `json:"user_id,omitempty"`  // 转化为json时：omitempty 忽略空值  - 忽略此字段
	UserLogin      string `json:"user_login,omitempty"`
	UserPass       string `json:"user_pass,omitempty"`
	UserNicename   string `json:"user_nicename,omitempty"`
	UserEmail      string `json:"user_email,omitempty"`
	UserRegistered string `json:"user_registered,omitempty"`
	UserStatus     string `json:"user_status,omitempty"`
}

//分类项表结构体
type Term struct {
	TermId   int64  `json:"term_id,omitempty"`
	TermName string `json:"term_name,omitempty"`
}

//分类法表结构体
type TermTaxonomy struct {
	TermTaxonomyId int64  `json:"term_taxonomy_id,omitempty"`
	Taxonomy       string `json:"taxonomy,omitempty"`
	Description    string `json:"description,omitempty"`
	TermParentId   string `json:"term_parent_id,omitempty"`
}

//文章结构体
type Post struct {
	PostId        int64  `json:"post_id,omitempty"`
	PostAuthor    string `json:"post_author,omitempty"`
	PostContent   string `json:"post_content,omitempty"`
	PostTitle     string `json:"post_title,omitempty"`
	PostDate      string `json:"post_date,omitempty"`
	PostPreImgUrl string `json:"post_pre_img_url,omitempty"`
	CommentCount  string `json:"comment_count,omitempty"`
	PostExcerpt   string `json:"post_excerpt,omitempty"`
	PostStatus    int    `json:"post_status,omitempty"`
	PostModified  string `json:"post_modified,omitempty"`
	CommentStatus int    `json:"comment_status,omitempty"`
}

//评论结构体
type Comment struct {
	CommentId          int64  `json:"comment_id,omitempty"`
	CommentAuthor      string `json:"comment_author,omitempty"`
	CommentAuthorEmail string `json:"comment_author_email,omitempty"`
	CommentAuthorIP    string `json:"comment_author_ip,omitempty"`
	CommentDate        string `json:"comment_date,omitempty"`
	CommentContent     string `json:"comment_content,omitempty"`
	CommentApproved    int64  `json:"comment_approved,omitempty"`
}

//链接
type Link struct {
	LinkId          int64  `json:"link_id,omitempty"`
	LinkUrl         string `json:"link_url,omitempty"`
	LinkName        string `json:"link_name,omitempty"`
	LinkImageUrl    string `json:"link_image_url,omitempty"`
	LinkDescription string `json:"link_description,omitempty"`
	LinkVisible     int    `json:"link_visible,omitempty"`
	LinkRating      int    `json:"link_rating,omitempty"`
	LinkUpdated     string `json:"link_updated,omitempty"`
}