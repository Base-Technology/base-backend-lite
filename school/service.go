package school

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/Base-Technology/base-backend-lite/imtp"
	"github.com/pkg/errors"
)

const (
	schoolConfigFile = "config/school.json"
)

var schools map[string]*School

func InitSchoolGroup() error {
	b, err := ioutil.ReadFile(schoolConfigFile)
	if err != nil {
		return err
	}
	schoolList := []*School{}
	if err := json.Unmarshal(b, &schoolList); err != nil {
		return err
	}
	schools = make(map[string]*School)
	for _, s := range schoolList {
		schools[s.ID] = s
	}
	return nil
}

func InviteUserToSchoolGroup(address, schoolID string) error {
	school, ok := schools[schoolID]
	if !ok {
		return errors.Errorf("school [%s] not found", schoolID)
	}
	userID := imtp.GetUserIDFromAddress(address)
	token, _, err := imtp.Login(conf.Conf.IMTPConf.AdminPrivateKey)
	if err != nil {
		return err
	}
	if err := imtp.InviteUserToGroup(token, school.GroupID, userID); err != nil {
		return err
	}
	return nil
}
