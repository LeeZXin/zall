package pullrequestmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertPullRequest(ctx context.Context, reqDTO InsertPullRequestReqDTO) (PullRequest, error) {
	ret := PullRequest{
		RepoId:   reqDTO.RepoId,
		Target:   reqDTO.Target,
		Head:     reqDTO.Head,
		PrStatus: reqDTO.PrStatus,
		CreateBy: reqDTO.CreateBy,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func ExistsOpenStatusPrByRepoIdAndRef(ctx context.Context, reqDTO ExistsOpenStatusPrByRepoIdAndRefReqDTO) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", reqDTO.RepoId).
		And("target = ?", reqDTO.Target).
		And("head = ?", reqDTO.Head).
		And("pr_status = ?", PrOpenStatus).
		Exist(new(PullRequest))
}

func UpdatePrStatus(ctx context.Context, id int64, oldStatus, newStatus PrStatus) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("pr_status = ?", oldStatus).
		Cols("pr_status").
		Update(&PullRequest{
			PrStatus: newStatus,
		})
	return rows == 1, err
}

func UpdatePrStatusAndCommitId(ctx context.Context, id int64, oldStatus, newStatus PrStatus, targetCommitId, headCommitId string) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("pr_status = ?", oldStatus).
		Cols("pr_status", "target_commit_id", "head_commit_id").
		Update(&PullRequest{
			TargetCommitId: targetCommitId,
			HeadCommitId:   headCommitId,
			PrStatus:       newStatus,
		})
	return rows == 1, err
}

func GetById(ctx context.Context, id int64) (PullRequest, bool, error) {
	var ret PullRequest
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func InsertReview(ctx context.Context, reqDTO InsertReviewReqDTO) error {
	_, err := xormutil.MustGetXormSession(ctx).Insert(&Review{
		PrId:         reqDTO.PrId,
		Reviewer:     reqDTO.Reviewer,
		ReviewMsg:    reqDTO.ReviewMsg,
		ReviewStatus: reqDTO.Status,
	})
	return err
}

func UpdateReview(ctx context.Context, reqDTO UpdateReviewReqDTO) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", reqDTO.Id).
		Cols("review_status").
		Limit(1).
		Update(&Review{
			ReviewStatus: reqDTO.Status,
		})
	return rows == 1, err
}

func ListReview(ctx context.Context, prId int64) ([]Review, error) {
	ret := make([]Review, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		Find(&ret)
	return ret, err
}

func CountReview(ctx context.Context, prId int64, status ReviewStatus) (int64, error) {
	ret, err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		And("review_status = ?", status).
		Count(new(Review))
	return ret, err
}

func GetReview(ctx context.Context, prId int64, reviewer string) (Review, bool, error) {
	var ret Review
	b, err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		And("reviewer = ?", reviewer).
		Get(&ret)
	return ret, b, err
}
