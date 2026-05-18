package model

type Course struct {
	CourseID     int      `json:"courseId"`
	ResourceType string   `json:"resourceType"`
	Name         string   `json:"name"`
	Teacher      []string `json:"teacher"`
	Category     []string `json:"category"`
	Semester     string   `json:"semester"`
	Credit       int      `json:"credit"`
	Cover        string   `json:"cover"`
	Views        int      `json:"views"`
	Loves        int      `json:"loves"`
	Collections  int      `json:"collections"`
}

type CourseDetail struct {
	CourseID     int              `json:"courseId"`
	ResourceType string           `json:"resourceType"`
	Name         string           `json:"name"`
	Category     string           `json:"catagory"`
	URLForm      []ResourceWeb    `json:"url_form"`
	UploadForm   []ResourceUpload `json:"upload_form"`
	Contributor  []string         `json:"contributor"`
	Collections  int              `json:"collections"`
	Views        int              `json:"views"`
	Likes        int              `json:"likes"`
	IsLiked      bool             `json:"isliked"`
	IsCollected  bool             `json:"iscollected"`
	CommentTotal int              `json:"comment_total"`
	Comments     []CourseComment  `json:"comments"`
	CreatedAt    string           `json:"createdAt"`
}

type CourseComment struct {
	CommentID   int             `json:"comment_Id"`
	CommentIDP  *int            `json:"commentId"`
	Nickname    string          `json:"nickname"`
	Avatar      string          `json:"avater"`
	Comment     string          `json:"comment"`
	CommentDate string          `json:"commentDate"`
	LoveCount   int             `json:"love_count"`
	IsOwner     *bool           `json:"isowner"`
	IsReply     bool            `json:"isreply"`
	ReplyTotal  int             `json:"reply_total"`
	Replies     []CourseComment `json:"replies"`
}

type TeachReview struct {
	ResourceID   int            `json:"resourceId"`
	ResourceType string         `json:"resourceType"`
	Resource1    *ResourceWeb   `json:"resource1"`
	Resource2    ResourceUpload `json:"resource2"`
	AuditStatus  string         `json:"auditStatus"`
	SubmitTime   string         `json:"submitTime"`
	AuditTime    string         `json:"auditTime"`
	RejectReason string         `json:"rejectReason"`
}

type TeachPersonal struct {
	ResourceID   int    `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	Resource     string `json:"resource"`
	Image        string `json:"image"`
	Introduce    string `json:"introduce"`
	Contributer  []User `json:"contributer"`
}
