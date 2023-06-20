package school

import (
	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/Base-Technology/base-backend-lite/imtp"
	"github.com/pkg/errors"
)

func InviteUserToSchoolGroup(address, name string) error {
	school, ok := schools[name]
	if !ok {
		return errors.Errorf("school [%s] not found", name)
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

func GetGroupIDByName(name string) (string, error) {
	school, ok := schools[name]
	if !ok {
		return "", errors.Errorf("school [%s] not found", name)
	}
	return school.GroupID, nil
}
