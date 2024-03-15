package propmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"regexp"
)

var (
	validNodeIdPattern      = regexp.MustCompile("[\\w-]{1,32}")
	validContentNamePattern = regexp.MustCompile("[\\w-]+\\.[a-zA-Z]+")
)

func IsNodeIdValid(nodeId string) bool {
	return validNodeIdPattern.MatchString(nodeId)
}

func IsPropContentNameValid(name string) bool {
	return validContentNamePattern.MatchString(name)
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

func InsertPropContent(ctx context.Context, reqDTO InsertPropContentReqDTO) (PropContent, error) {
	ret := PropContent{
		AppId: reqDTO.AppId,
		Name:  reqDTO.Name,
	}
	_, err := xormutil.MustGetXormSession(ctx).
		Table("zprop_prop_content_" + reqDTO.Env).
		Insert(&ret)
	return ret, err
}

func DeletePropContent(ctx context.Context, id int64, env string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Table("zprop_prop_content_" + env).
		Limit(1).
		Delete(new(PropContent))
	return rows == 1, err
}

func GetPropContentById(ctx context.Context, id int64, env string) (PropContent, bool, error) {
	ret := PropContent{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Table("zprop_prop_content_" + env).
		Get(&ret)
	return ret, b, err
}

func GetPropContentByAppIdAndName(ctx context.Context, appId, name, env string) (PropContent, bool, error) {
	ret := PropContent{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("name = ?", name).
		Table("zprop_prop_content_" + env).
		Get(&ret)
	return ret, b, err
}

func IteratePropContent(ctx context.Context, env string, fn func(content *PropContent) error) error {
	return xormutil.MustGetXormSession(ctx).
		Table("zprop_prop_content_"+env).
		Iterate(new(PropContent), func(_ int, bean interface{}) error {
			return fn(bean.(*PropContent))
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
		Table("zprop_prop_history_" + reqDTO.Env).
		Insert(&PropHistory{
			ContentId: reqDTO.ContentId,
			Content:   reqDTO.Content,
			Version:   reqDTO.Version,
			Creator:   reqDTO.Creator,
		})
	return err
}

func DeleteHistory(ctx context.Context, contentId int64, env string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("content_id = ?", contentId).
		Table("zprop_prop_history_" + env).
		Delete(new(PropHistory))
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

func ListPropContent(ctx context.Context, appId, env string) ([]PropContent, error) {
	ret := make([]PropContent, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		Table("zprop_prop_content_" + env).
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

func GetHistoryByVersion(ctx context.Context, contentId int64, version, env string) (PropHistory, bool, error) {
	var ret PropHistory
	b, err := xormutil.MustGetXormSession(ctx).
		Where("content_id = ?", contentId).
		And("version = ?", version).
		Table("zprop_prop_history_" + env).
		Get(&ret)
	return ret, b, err
}

func ListHistory(ctx context.Context, reqDTO ListHistoryReqDTO) ([]PropHistory, error) {
	ret := make([]PropHistory, 0)
	session := xormutil.MustGetXormSession(ctx).
		Where("content_id = ?", reqDTO.ContentId).
		Table("zprop_prop_history_" + reqDTO.Env)
	if reqDTO.Version != "" {
		session.And("version like ?", reqDTO.Version+"%")
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
