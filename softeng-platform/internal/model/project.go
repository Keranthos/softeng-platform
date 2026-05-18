package model

type Project struct {
	ProjectID    int      `json:"projectId"`
	ResourceType string   `json:"resourceType"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	TechStack    []string `json:"techStack"`
	LikeCount    int      `json:"likecount"`
	AuthorName   []string `json:"authername"`
	Cover        string   `json:"cover"`
	CreatedAt    string   `json:"createdat"`
	Loves        int      `json:"loves"`
	Collections  int      `json:"collections"`
	Views        int      `json:"views"`
}

type ProjectDetail struct {
	ProjectID    int              `json:"projectId"`
	ResourceType string           `json:"resourceType"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	GithubURL    string           `json:"githubURL"`
	TechStack    []string         `json:"techStack"`
	Category     string           `json:"catagory"`
	Images       []string         `json:"images"`
	Likes        int              `json:"likes"`
	Views        int              `json:"views"`
	Collections  int              `json:"collections"`
	IsLiked      bool             `json:"isliked"`
	IsCollected  bool             `json:"iscollected"`
	Author       []string         `json:"author"`
	CommentCount int              `json:"comment_count"`
	Comments     []ProjectComment `json:"comments"`
	CreatedAt    string           `json:"createdAt"`
}

type ProjectComment struct {
	CommentID   int              `json:"comment_Id"`
	CommentIDP  *int             `json:"commentId"`
	Nickname    string           `json:"nickname"`
	Avatar      string           `json:"avater"`
	Comment     string           `json:"comment"`
	CommentDate string           `json:"commentDate"`
	LoveCount   int              `json:"love_count"`
	IsOwner     *bool            `json:"isowner"`
	IsReply     bool             `json:"isreply"`
	ReplyTotal  int              `json:"reply_total"`
	Replies     []ProjectComment `json:"replies"`
}
