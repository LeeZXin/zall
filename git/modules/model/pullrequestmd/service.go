package pullrequestmd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/xorm/xormutil"
	"time"
	"xorm.io/builder"
)

func IsPrTitleValid(title string) bool {
	return len(title) > 0 && len(title) <= 255
}

func InsertPullRequest(ctx context.Context, reqDTO InsertPullRequestReqDTO) (PullRequest, error) {
	ret := PullRequest{
		RepoId:       reqDTO.RepoId,
		Target:       reqDTO.Target,
		TargetType:   reqDTO.TargetType,
		Head:         reqDTO.Head,
		HeadType:     reqDTO.HeadType,
		PrStatus:     reqDTO.PrStatus,
		CreateBy:     reqDTO.CreateBy,
		PrTitle:      reqDTO.Title,
		CommentCount: reqDTO.CommentCount,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
}

func ExistsPrByRepoIdAndRef(ctx context.Context, reqDTO ExistsPrByRepoIdAndRefReqDTO) (bool, error) {
	return xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", reqDTO.RepoId).
		And("target = ?", reqDTO.Target).
		And("head = ?", reqDTO.Head).
		And("pr_status = ?", reqDTO.Status).
		And("target_type = ?", reqDTO.TargetType).
		And("head_type = ?", reqDTO.HeadType).
		Exist(new(PullRequest))
}

func ClosePrStatus(ctx context.Context, id int64, oldStatus PrStatus, account string) (bool, error) {
	now := time.Now()
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("pr_status = ?", oldStatus).
		Cols("pr_status", "close_by", "closed").
		Update(&PullRequest{
			PrStatus: PrClosedStatus,
			Closed:   &now,
			CloseBy:  account,
		})
	return rows == 1, err
}

func MergePrStatus(ctx context.Context, id int64, oldStatus PrStatus, targetCommitId, headCommitId, account string) (bool, error) {
	now := time.Now()
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		And("pr_status = ?", oldStatus).
		Cols("pr_status", "target_commit_id", "head_commit_id", "merged", "merge_by").
		Update(&PullRequest{
			TargetCommitId: targetCommitId,
			HeadCommitId:   headCommitId,
			PrStatus:       PrMergedStatus,
			Merged:         &now,
			MergeBy:        account,
		})
	return rows == 1, err
}

func GetPullRequestById(ctx context.Context, id int64) (PullRequest, bool, error) {
	var ret PullRequest
	b, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Get(&ret)
	return ret, b, err
}

func GetLastPullRequestByRepoIdAndHead(ctx context.Context, repoId int64, heads []string) ([]PullRequest, error) {
	ret := make([]PullRequest, 0)
	err := xormutil.MustGetXormSession(ctx).
		In("id",
			builder.Select("max(id)").
				Where(builder.And(
					builder.Eq{"repo_id": repoId},
					builder.In("head", heads),
				)).
				From(PullRequestTableName).
				GroupBy("head"),
		).Find(&ret)
	return ret, err
}

func InsertReview(ctx context.Context, reqDTO InsertReviewReqDTO) (Review, error) {
	ret := Review{
		PrId:         reqDTO.PrId,
		Reviewer:     reqDTO.Reviewer,
		ReviewStatus: reqDTO.Status,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&ret)
	return ret, err
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
		OrderBy("id desc").
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

func ListPullRequest(ctx context.Context, reqDTO ListPullRequestReqDTO) ([]PullRequest, int64, error) {
	ret := make([]PullRequest, 0)
	session := xormutil.MustGetXormSession(ctx).Where("repo_id = ?", reqDTO.RepoId)
	session.Limit(reqDTO.PageSize, (reqDTO.PageNum-1)*reqDTO.PageSize)
	if reqDTO.SearchKey != "" {
		session.And("pr_title like ?", reqDTO.SearchKey+"%")
	}
	if reqDTO.Status != PrAllStatus {
		session.And("pr_status = ?", reqDTO.Status)
	}
	total, err := session.OrderBy("id desc").FindAndCount(&ret)
	return ret, total, err
}

func BatchInsertTimeline(ctx context.Context, reqDTOs []InsertTimelineReqDTO) error {
	timelines := listutil.MapNe(reqDTOs, func(t InsertTimelineReqDTO) *Timeline {
		return &Timeline{
			PrId:    t.PrId,
			Action:  &t.Action,
			Account: t.Account,
		}
	})
	_, err := xormutil.MustGetXormSession(ctx).
		Insert(timelines)
	return err
}

func ListTimeline(ctx context.Context, prId int64) ([]Timeline, error) {
	ret := make([]Timeline, 0)
	err := xormutil.MustGetXormSession(ctx).
		Where("pr_id = ?", prId).
		OrderBy("id asc").
		Find(&ret)
	return ret, err
}

func GetTimelineById(ctx context.Context, id int64) (Timeline, bool, error) {
	var ret Timeline
	b, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Get(&ret)
	return ret, b, err
}

func DeleteTimelineById(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("id = ?", id).Delete(new(Timeline))
	return rows == 1, err
}

func DeletePullRequestByRepoId(ctx context.Context, repoId int64) error {
	_, err := xormutil.MustGetXormSession(ctx).
		Where("repo_id = ?", repoId).
		Delete(new(PullRequest))
	return err
}

func IncrCommentCount(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Incr("comment_count").
		Cols("comment_count").
		Update(new(PullRequest))
	return rows == 1, err
}

func DecrCommentCount(ctx context.Context, id int64) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).
		Where("id = ?", id).
		Decr("comment_count").
		Cols("comment_count").
		Update(new(PullRequest))
	return rows == 1, err
}
