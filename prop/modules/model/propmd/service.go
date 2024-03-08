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

func ListEtcdNode(ctx context.Context) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	err := xormutil.MustGetXormSession(ctx).Find(&ret)
	return ret, err
}

func GetEtcdNodeByNodeId(ctx context.Context, nodeId string) (EtcdNode, bool, error) {
	ret := EtcdNode{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", nodeId).
		Get(&ret)
	return ret, b, err
}

func InsertEtcdNode(ctx context.Context, reqDTO InsertEtcdNodeReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&EtcdNode{
		NodeId:    reqDTO.NodeId,
		Endpoints: reqDTO.Endpoints,
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
	})
	return err
}

func DeleteEtcdNode(ctx context.Context, nodeId string) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", nodeId).
		Limit(1).
		Delete(new(EtcdNode))
	return err
}

func UpdateEtcdNode(ctx context.Context, reqDTO UpdateEtcdNodeReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("node_id = ?", reqDTO.NodeId).
		Limit(1).
		Cols("endpoints", "username", "password").
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
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func DeletePropContent(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Limit(1).
		Delete(new(PropContent))
	return rows == 1, err
}

func GetPropContentById(ctx context.Context, id int64) (PropContent, bool, error) {
	ret := PropContent{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func GetPropContentByAppIdAndName(ctx context.Context, appId, name string) (PropContent, bool, error) {
	ret := PropContent{}
	b, err := xormutil.MustGetXormSession(ctx).
		Where("app_id = ?", appId).
		And("name = ?", name).
		Get(&ret)
	return ret, b, err
}

func InsertHistory(ctx context.Context, reqDTO InsertHistoryReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&PropHistory{
		ContentId: reqDTO.ContentId,
		Content:   reqDTO.Content,
		Version:   reqDTO.Version,
	})
	return err
}

func DeleteHistory(ctx context.Context, contentId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).Where("content_id = ?", contentId).Delete(new(PropHistory))
	return err
}

func ListPropContent(ctx context.Context, appId string) ([]PropContent, error) {
	ret := make([]PropContent, 0)
	err := xormutil.MustGetXormSession(ctx).Where("app_id = ?", appId).Find(&ret)
	return ret, err
}

func BatchGetEtcdNodes(ctx context.Context, nodeIdList []string) ([]EtcdNode, error) {
	ret := make([]EtcdNode, 0)
	err := xormutil.MustGetXormSession(ctx).In("node_id", nodeIdList).Find(&ret)
	return ret, err
}

func InsertDeploy(ctx context.Context, reqDTO InsertDeployReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&PropDeploy{
		ContentId: reqDTO.ContentId,
		Content:   reqDTO.Content,
		Version:   reqDTO.Version,
		NodeId:    reqDTO.NodeId,
		Endpoints: reqDTO.Endpoints,
		Username:  reqDTO.Username,
		Password:  reqDTO.Password,
	})
	return err
}

func GetHistoryByVersion(ctx context.Context, contentId int64, version string) (PropHistory, bool, error) {
	var ret PropHistory
	b, err := xormutil.MustGetXormSession(ctx).Where("content_id = ?", contentId).And("version = ?", version).Get(&ret)
	return ret, b, err
}

func ListHistory(ctx context.Context, reqDTO ListHistoryReqDTO) ([]PropHistory, error) {
	ret := make([]PropHistory, 0)
	session := xormutil.MustGetXormSession(ctx).Where("content_id = ?", reqDTO.ContentId)
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
	session := xormutil.MustGetXormSession(ctx).Where("content_id = ?", reqDTO.ContentId)
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
	_, err := xormutil.MustGetXormSession(ctx).Insert(&EtcdAuth{
		AppId:    reqDTO.AppId,
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	})
	return err
}

func GetAuthByAppId(ctx context.Context, appId string) (EtcdAuth, bool, error) {
	var ret EtcdAuth
	b, err := xormutil.MustGetXormSession(ctx).Where("app_id = ?", appId).Get(&ret)
	return ret, b, err
}
