package imtp

import (
	"fmt"
	"net/http"

	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/Base-Technology/base-backend-lite/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

func Login(privateKey string) (string, string, error) {
	k, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", "", err
	}
	sign, err := utils.SignMessage("hello", k)
	if err != nil {
		return "", "", err
	}

	request := &LoginRequest{
		SenderAddress: crypto.PubkeyToAddress(k.PublicKey).Hex(),
		Signature:     sign,
		Network:       1,
	}
	response := &LoginResponse{}
	if err := utils.SendHttpRequest(fmt.Sprintf("%s%s", conf.Conf.IMTPConf.APPServer, "/api/v1/login"), http.MethodPost, nil, request, response); err != nil {
		return "", "", err
	}
	return response.Token, response.UserID, nil
}

func CreateGroup(token, groupName, ownerUserID string) (string, error) {
	request := &CreateGroupRequest{
		GroupName:   groupName,
		OwnerUserID: ownerUserID,
		MemberList: []*struct {
			RoleLevel int    "json:\"roleLevel\""
			UserID    string "json:\"userID\""
		}{{3, ownerUserID}},
		OperationID: "CreateGroup",
	}
	header := make(map[string]string)
	header["token"] = token
	response := &CreateGroupResponse{}
	if err := utils.SendHttpRequest(fmt.Sprintf("%s%s", conf.Conf.IMTPConf.APIServer, "/group/create_group"), http.MethodPost, header, request, response); err != nil {
		return "", err
	}
	return response.Data.GroupID, nil
}

func InviteUserToGroup(token, groupID, userID string) error {
	request := &InviteUserToGroupRequest{
		GroupID:           groupID,
		InvitedUserIDList: []string{userID},
		OperationID:       "InviteUserToGroup",
	}
	header := make(map[string]string)
	header["token"] = token
	response := &InviteUserToGroupResponse{}
	if err := utils.SendHttpRequest(fmt.Sprintf("%s%s", conf.Conf.IMTPConf.APIServer, "/group/invite_user_to_group"), http.MethodPost, header, request, response); err != nil {
		return err
	}
	return nil
}
