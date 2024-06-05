package propertymd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
)

var (
	validNodeIdPattern   = regexp.MustCompile("[\\w-]{1,32}")
	validFileNamePattern = regexp.MustCompile("[\\w-]+\\.[a-zA-Z]+")
)

func IsNodeIdValid(nodeId string) bool {
	return validNodeIdPattern.MatchString(nodeId)
}

func IsFileNameValid(name string) bool {
	return validFileNamePattern.MatchString(name)
}

func ListEtcdNode(ctx context.Context, env string) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	err := xormutil.MustGetXormSession(ctx).
		Table("zprop_etcd_node_" + env).
		Find(&ret)
	return ret, err
}

func GetEtcdNodeByNodeId(ctx context.Context, nodeId, env string) (EtcdNode, bool, error) {
	ret := EtcdNode{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", nodeId).
		Table("zprop_etcd_node_" + env).
		Get(&ret)
	return ret, b, err
}

func InsertEtcdNode(ctx context.Context, reqDTO InsertEtcdNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zprop_etcd_node_" + reqDTO.Env).
		Insert(&EtcdNode{
			NodeId:    reqDTO.NodeId,
			Endpoints: reqDTO.Endpoints,
			Username:  reqDTO.Username,
			Password:  reqDTO.Password,
		})
	return err
}

func DeleteEtcdNode(ctx context.Context, nodeId, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", nodeId).
		Table("zprop_etcd_node_" + env).
		Limit(1).
		Delete(new(EtcdNode))
	return err
}

func UpdateEtcdNode(ctx context.Context, reqDTO UpdateEtcdNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", reqDTO.NodeId).
		Limit(1).
		Cols("endpoints", "username", "password").
		Table("zprop_etcd_node_" + reqDTO.Env).
		Update(&EtcdNode{
			Endpoints: reqDTO.Endpoints,
			Username:  reqDTO.Username,
			Password:  reqDTO.Password,
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

func DeletePropContent(ctx context.Context, id int64, env string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Table("zprop_prop_content_" + env).
		Limit(1).
		Delete(new(File))
	return rows == 1, err
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

func IteratePropContent(ctx context.Context, env string, fn func(content *File) error) error {
	return xormutil.MustGetXormSession(ctx).
		Table("zprop_prop_content_"+env).
		Iterate(new(File), func(_ int, bean interface{}) error {
			return fn(bean.(*File))
		})
}

func IterateDeletedDeployByNodeId(ctx context.Context, nodeId, env string, fn func(deploy *PropDeploy) error) error {
	return xormutil.MustGetXormSession(ctx).Where("deleted = 1").
		And("node_id = ?", nodeId).
		Table("zprop_prop_deploy_"+env).
		Iterate(new(PropDeploy), func(_ int, bean interface{}) error {
			return fn(bean.(*PropDeploy))
		})
}

func InsertHistory(ctx context.Context, reqDTO InsertHistoryReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(&History{
			FileId:      reqDTO.FileId,
			Content:     reqDTO.Content,
			Version:     reqDTO.Version,
			LastVersion: reqDTO.LastVersion,
			Creator:     reqDTO.Creator,
			Env:         reqDTO.Env,
		})
	return err
}

func DeleteHistory(ctx context.Context, contentId int64, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("content_id = ?", contentId).
		Table("zprop_prop_history_" + env).
		Delete(new(History))
	return err
}

func DeleteDeploy(ctx context.Context, contentId int64, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("content_id = ?", contentId).
		And("deleted = 0").
		Table("zprop_prop_deploy_" + env).
		Cols("deleted").
		Update(&PropDeploy{
			Deleted: true,
		})
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

func BatchGetEtcdNodes(ctx context.Context, nodeIdList []string, env string) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("node_id", nodeIdList).
		Table("zprop_etcd_node_" + env).
		Find(&ret)
	return ret, err
}

func InsertDeploy(ctx context.Context, reqDTO InsertDeployReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zprop_prop_deploy_" + reqDTO.Env).
		Insert(&PropDeploy{
			ContentId:    reqDTO.ContentId,
			Content:      reqDTO.Content,
			Version:      reqDTO.Version,
			NodeId:       reqDTO.NodeId,
			ContentAppId: reqDTO.ContentAppId,
			ContentName:  reqDTO.ContentName,
			Endpoints:    reqDTO.Endpoints,
			Username:     reqDTO.Username,
			Password:     reqDTO.Password,
			Creator:      reqDTO.Creator,
			Deleted:      false,
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

func ListDeploy(ctx context.Context, reqDTO ListDeployReqDTO) ([]PropDeploy, error) {
	ret := make([]PropDeploy, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("content_id = ?", reqDTO.ContentId).
		Table("zprop_prop_deploy_" + reqDTO.Env).
		And("deleted = 0")
	if reqDTO.Version != "" {
		session.And("version like ?", reqDTO.Version+"%")
	}
	if reqDTO.NodeId != "" {
		session.And("node_id = ?", reqDTO.NodeId)
	}
	if reqDTO.Cursor > 0 {
		session.And("id < ?", reqDTO.Cursor)
	}
	if reqDTO.Limit > 0 {
		session.Limit(reqDTO.Limit)
	}
	err := session.OrderBy("id desc").Find(&ret)
	return ret, err
}

func InsertAuth(ctx context.Context, reqDTO InsertAuthReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zprop_etcd_auth_" + reqDTO.Env).
		Insert(&EtcdAuth{
			AppId:    reqDTO.AppId,
			Username: reqDTO.Username,
			Password: reqDTO.Password,
		})
	return err
}

func GetAuthByAppId(ctx context.Context, appId, env string) (EtcdAuth, bool, error) {
	var ret EtcdAuth
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Table("zprop_etcd_auth_" + env).
		Get(&ret)
	return ret, b, err
}

func GetLatestDeployByNodeId(ctx context.Context, contentId int64, nodeId, env string) (PropDeploy, bool, error) {
	ret := PropDeploy{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("content_id = ?", contentId).
		Table("zprop_etcd_deploy_"+env).
		And("node_id = ?", nodeId).
		And("deleted = 0").
		OrderBy("id desc").
		Get(&ret)
	return ret, b, err
}
