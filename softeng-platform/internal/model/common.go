package model

type ResourceWeb struct {
	ResourceIntro string `json:"resource_intro"`
	ResourceURL   string `json:"resource_url"`
	ResourceID    int    `json:"resource_id"`
}

type ResourceUpload struct {
	ResourceIntro  string `json:"resource_intro"`
	ResourceUpload string `json:"resource_upload"`
	ResourceID     int    `json:"resource_id"`
}

type ResourceReview struct {
	ResourceID   int    `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	Resource     string `json:"resource"`
	AuditStatus  string `json:"auditStatus"`
	SubmitTime   string `json:"submitTime"`
	AuditTime    string `json:"auditTime"`
	RejectReason string `json:"rejectReason"`
}

type ResourcePersonal struct {
	ResourceID   int    `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	Resource     string `json:"resource"`
	Image        string `json:"image"`
	Introduce    string `json:"introduce"`
	Contributer  []User `json:"contributer"`
}

type Maneuver struct {
	ResourceID   int    `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	NewStatus    string `json:"newstatus"`
	OldStatus    string `json:"oldestatus"`
	OperateTime  string `json:"operateTime"`
	Operator     string `json:"operator"`
}

type Submit struct {
	Submitor     string   `json:"submitor"`
	SubmitDate   string   `json:"submitDate"`
	ResourceID   int      `json:"reourceId"`
	ResourceType string   `json:"resourceType"`
	Category     string   `json:"catagory"`
	Link         string   `json:"link"`
	File         string   `json:"file"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	ResourceName string   `json:"resourcename"`
}

type Comment struct {
	CommentID   int       `json:"comment_Id"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avater"`
	Comment     string    `json:"comment"`
	CommentDate string    `json:"commentDate"`
	LoveCount   int       `json:"love_count"`
	ReplyTotal  int       `json:"reply_total"`
	Replies     []Comment `json:"replies"`
}
