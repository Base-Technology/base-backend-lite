package school

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/imtp"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	schoolInfoFile   = "config/school_info.json"
	schoolConfigFile = "config/school.json"
)

var schools map[string]*School

func InitSchoolGroup() error {
	_, err := os.Stat(schoolConfigFile)
	if err == nil {
		return initSchoolGroup()
	}

	return updateSchoolConfig()

}

func initSchoolGroup() error {
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
		schools[s.Name] = s
	}
	return nil
}

func updateSchoolConfig() error {
	seelog.Infof("%s not found, start to update school config", schoolConfigFile)

	b, err := ioutil.ReadFile(schoolInfoFile)
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(b), &data); err != nil {
		return err
	}

	token, userID, err := imtp.Login(conf.Conf.IMTPConf.AdminPrivateKey)
	if err != nil {
		return err
	}

	schools := []*School{}
	for k := range data {
		v := data[k].([]interface{})
		for _, s := range v {
			school := s.(map[string]interface{})

			groupID, err := imtp.CreateGroup(token, school["学校名称"].(string), userID)
			if err != nil {
				return err
			}
			schools = append(schools, &School{
				ID:      school["序号"].(string),
				Name:    school["学校名称"].(string),
				GroupID: groupID,
			})

			seelog.Infof("create [%s %s] group success", school["序号"].(string), school["学校名称"].(string))
		}
	}
	b, err = json.Marshal(schools)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(schoolConfigFile, b, 0644); err != nil {
		return err
	}

	if err := initSchoolGroup(); err != nil {
		return err
	}

	if err := updateUserSchoolConfig(); err != nil {
		return err
	}

	return nil
}

func updateUserSchoolConfig() error {
	users := []*database.User{}
	if err := database.GetInstance().Model(&database.User{}).Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		// decode private key
		kBytes := user.PrivateKey
		b, err := hexutil.Decode(kBytes)
		if err != nil {
			return err
		}
		k, err := crypto.ToECDSA(b)
		if err != nil {
			return err
		}
		// login
		_, _, err = imtp.Login(kBytes[2:])
		if err != nil {
			return err
		}
		// invite to imtp group
		address := crypto.PubkeyToAddress(k.PublicKey).Hex()
		if err := InviteUserToSchoolGroup(address, user.School); err != nil {
			return err
		}
		seelog.Infof("update user [%s] to school [%s] success", user.Name, user.School)
	}

	return nil
}
