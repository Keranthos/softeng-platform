package model

type Tool struct {
	ResourceID        int       `json:"resourceId"`
	ResourceType      string    `json:"resourceType"`
	ResourceName      string    `json:"resourceName"`
	ResourceLink      string    `json:"resourceLink"`
	Description       string    `json:"description"`
	DescriptionDetail string    `json:"description_detail"`
	Category          string    `json:"catagory"`
	Tags              []string  `json:"tags"`
	Image             []string  `json:"image"`
	Views             int       `json:"views"`
	Collections       int       `json:"collections"`
	Loves             int       `json:"loves"`
	IsCollected       *bool     `json:"iscollected"`
	IsLiked           bool      `json:"isliked"`
	CommentCount      int       `json:"comment_count"`
	Comments          []Comment `json:"comments"`
	CreatedDate       string    `json:"createdDate"`
	Contributors      []string  `json:"contributors"`
}

type ToolReview struct {
	ResourceID   int     `json:"resourceId"`
	ResourceType string  `json:"resourceType"`
	Resource     string  `json:"resource"`
	AuditStatus  string  `json:"auditStatus"`
	SubmitTime   string  `json:"submitTime"`
	AuditTime    *string `json:"auditTime"`
	RejectReason *string `json:"rejectReason"`
}

type ToolPersonal struct {
	ResourceID   int    `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	Resource     string `json:"resource"`
	Image        string `json:"image"`
	Introduce    string `json:"introduce"`
	Contributer  []User `json:"contributer"`
}

type ToolComment struct {
	CommentID   int           `json:"comment_Id"`
	CommetnID   *int          `json:"commetnId"`
	Nickname    string        `json:"nickname"`
	Avatar      string        `json:"avater"`
	Comment     string        `json:"comment"`
	CommentDate string        `json:"commentDate"`
	LoveCount   int           `json:"love_count"`
	IsOwner     *bool         `json:"isowner"`
	IsReply     bool          `json:"isreply"`
	ReplyTotal  int           `json:"reply_total"`
	Replies     []ToolComment `json:"replies"`
}
