package propertymd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
)

func IsFileNameValid(name string) bool {
	return regexp.MustCompile("[\\w-]+\\.[a-zA-Z]+").MatchString(name)
}

func IsPropertySourceNameValid(name string) bool {
	return len(name) > 0 && len(name) <= 32
}

func ListEtcdNode(ctx context.Context, reqDTO ListEtcdNodeReqDTO) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", reqDTO.AppId).
		And("env = ?", reqDTO.Env)
	if len(reqDTO.Cols) > 0 {
		session.Cols(reqDTO.Cols...)
	}
	err := session.Find(&ret)
	return ret, err
}

func InsertEtcdNode(ctx context.Context, reqDTO InsertEtcdNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&EtcdNode{
			AppId:     reqDTO.AppId,
			Name:      reqDTO.Name,
			Env:       reqDTO.Env,
			Endpoints: reqDTO.Endpoints,
			Username:  reqDTO.Username,
			Password:  reqDTO.Password,
		})
	return err
}

func DeleteEtcdNodeById(ctx context.Context, id int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(EtcdNode))
	return err
}

func DeleteEtcdNodeByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(EtcdNode))
	return err
}

func UpdateEtcdNode(ctx context.Context, reqDTO UpdateEtcdNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("endpoints", "username", "password", "name").
		Update(&EtcdNode{
			Endpoints: reqDTO.Endpoints,
			Username:  reqDTO.Username,
			Password:  reqDTO.Password,
			Name:      reqDTO.Name,
		})
	return rows == 1, err
}

func InsertFile(ctx context.Context, reqDTO InsertFileReqDTO) (File, error) {
	ret := File{
		AppId: reqDTO.AppId,
		Name:  reqDTO.Name,
		Env:   reqDTO.Env,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func DeleteFileById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Delete(new(File))
	return rows == 1, err
}

func DeleteFileByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(File))
	return err
}

func GetFileById(ctx context.Context, id int64) (File, bool, error) {
	ret := File{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ExistFile(ctx context.Context, appId, name, env string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("name = ?", name).
		And("env = ?", env).
		Exist(new(File))
}

func InsertHistory(ctx context.Context, reqDTO InsertHistoryReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&History{
			FileId:      reqDTO.FileId,
			Content:     reqDTO.Content,
			Version:     reqDTO.Version,
			LastVersion: reqDTO.LastVersion,
			Creator:     reqDTO.Creator,
		})
	return err
}

func DeleteHistoryByFileId(ctx context.Context, fileId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("file_id = ?", fileId).
		Delete(new(History))
	return err
}

func DeleteDeployByFileId(ctx context.Context, fileId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("file_id = ?", fileId).
		Delete(new(Deploy))
	return err
}

func DeleteDeployByAppId(ctx context.Context, appId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Delete(new(Deploy))
	return err
}

func ListFile(ctx context.Context, appId, env string) ([]File, error) {
	ret := make([]File, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("env = ?", env).
		Find(&ret)
	return ret, err
}

func BatchGetEtcdNodes(ctx context.Context, nodeIdList []int64) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id", nodeIdList).
		Find(&ret)
	return ret, err
}

func InsertDeploy(ctx context.Context, reqDTO InsertDeployReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&Deploy{
			HistoryId: reqDTO.HistoryId,
			NodeName:  reqDTO.NodeName,
			FileId:    reqDTO.FileId,
			AppId:     reqDTO.AppId,
			Endpoints: reqDTO.Endpoints,
			Username:  reqDTO.Username,
			Password:  reqDTO.Password,
			Creator:   reqDTO.Creator,
		})
	return err
}

func GetHistoryByVersion(ctx context.Context, fileId int64, version string) (History, bool, error) {
	var ret History
	b, err := xormutil.MustGetXormSession(ctx).
		Where("file_id = ?", fileId).
		And("version = ?", version).
		Get(&ret)
	return ret, b, err
}

func GetHistoryById(ctx context.Context, id int64) (History, bool, error) {
	var ret History
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func ExistHistoryByVersion(ctx context.Context, fileId int64, version string) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("file_id = ?", fileId).
		And("version = ?", version).
		Get(new(History))
}

func PageHistory(ctx context.Context, reqDTO PageHistoryReqDTO) ([]History, int64, error) {
	ret := make([]History, 0)
	total, err := xormutil.MustGetXormSession(ctx).
		Where("file_id = ?", reqDTO.FileId).
		Desc("id").
		Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize).
		FindAndCount(&ret)
	return ret, total, err
}

func ListDeployByHistoryId(ctx context.Context, historyId int64) ([]Deploy, error) {
	ret := make([]Deploy, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("history_id = ?", historyId).
		OrderBy("id desc").
		Find(&ret)
	return ret, err
}

func GetEtcdNodeById(ctx context.Context, nodeId int64) (EtcdNode, bool, error) {
	var ret EtcdNode
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", nodeId).Get(&ret)
	return ret, b, err
}
