package imtp

type LoginRequest struct {
	SenderAddress string `json:"senderAddress"`
	Signature     string `json:"signature"`
	Network       int    `json:"network"`
}

type LoginResponse struct {
	Token  interface{} `json:"token"`
	UserID string      `json:"userID"`
}

type CreateGroupRequest struct {
	GroupType    int    `json:"groupType"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	OwnerUserID  string `json:"ownerUserID"`
	MemberList   []*struct {
		RoleLevel int    `json:"roleLevel"`
		UserID    string `json:"userID"`
	} `json:"memberList"`
	OperationID string `json:"operationID"`
}

type CreateGroupResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		GroupID string `json:"groupID"`
	} `json:"data"`
}

type InviteUserToGroupRequest struct {
	GroupID           string   `json:"groupID"`
	InvitedUserIDList []string `json:"invitedUserIDList"`
	Reason            string   `json:"reason"`
	OperationID       string   `json:"operationID"`
}

type InviteUserToGroupResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    []struct {
		Result int    `json:"result"`
		UserID string `json:"userID"`
	} `json:"data"`
}
