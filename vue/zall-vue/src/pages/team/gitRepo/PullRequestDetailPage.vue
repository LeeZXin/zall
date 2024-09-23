<template>
  <div style="padding:10px">
    <div class="pr-title no-wrap">
      <span class="pr-id">#{{prStore.prIndex}}</span>
      <span>{{prStore.prTitle}}</span>
    </div>
    <div class="desc">
      <PrStatusTag :status="prStore.prStatus" />
      <span style="font-size:14px">
        <span style="margin-right:6px">{{prStore.head}}</span>
        <ArrowLeftOutlined style="color:gray" />
        <span style="margin-left:6px">{{prStore.target}}</span>
      </span>
    </div>
    <div>
      <a-tabs style="width: 100%;" @change="selectTab">
        <a-tab-pane key="timeline" :tab="t('pullRequest.timeline')">
          <div class="timeline">
            <a-timeline style="margin-top:20px;width:calc(100% - 416px)">
              <a-timeline-item v-if="prStore.prStatus === 1">
                <div class="message-card">
                  <div class="card-title">
                    <template v-if="canMergeDetectLoaded">
                      <template v-if="mergeDetect.canMerge">
                        <CheckCircleOutlined />
                        <span style="margin-left:8px">{{t('pullRequest.thisPrCanMerge')}}</span>
                      </template>
                      <template v-else>
                        <WarningOutlined />
                        <span style="margin-left:8px">{{t('pullRequest.thisPrCantMerge')}}</span>
                      </template>
                    </template>
                    <template v-else>
                      <LoadingOutlined />
                    </template>
                  </div>
                  <div class="card-content" v-if="canMergeDetectLoaded">
                    <div
                      v-if="mergeDetect.isProtectedBranch"
                      style="font-size:14px;color:green;margin-bottom:10px"
                    >
                      <FileProtectOutlined />
                      <span style="margin-left:8px">{{t('pullRequest.protectedByProtectedBranch')}}</span>
                    </div>
                    <div v-if="mergeDetect.canMerge">
                      <a-button
                        @click="mergePr"
                        type="primary"
                        v-if="repoStore.perm?.canSubmitPullRequest"
                      >
                        <LoadingOutlined v-if="merging" />
                        <span>{{t('pullRequest.mergePr')}}</span>
                      </a-button>
                    </div>
                    <div v-else class="can-not-merge-reason" style="margin-top:4px">
                      <WarningOutlined />
                      <span style="margin-left:8px">{{cannotMergeReason}}</span>
                    </div>
                  </div>
                </div>
              </a-timeline-item>
              <a-timeline-item v-for="item in timelines" v-bind:key="item.id">
                <template v-if="item.action.actionType === 3">
                  <!-- 创建、合并、关闭合并请求 -->
                  <div class="timeline-text no-wrap">
                    <CloseCircleFilled
                      style="color:red;margin-right:6px"
                      v-if="item.action?.pr?.status === 2"
                    />
                    <CheckCircleFilled style="color:green;margin-right:6px" v-else />
                    <ZAvatar
                      :url="item.avatarUrl"
                      :name="item.name"
                      :account="item.account"
                      :showName="true"
                    />
                    <span>{{t(prStatusMap[item.action?.pr?.status])}}</span>
                    <span style="color:gray">{{readableTimeComparingNow(item.created)}}</span>
                  </div>
                </template>
                <template v-if="item.action.actionType === 4">
                  <!-- 评审并同意合并 -->
                  <div class="timeline-text no-wrap">
                    <CheckCircleFilled style="color:green;margin-right:6px" />
                    <ZAvatar
                      :url="item.avatarUrl"
                      :name="item.name"
                      :account="item.account"
                      :showName="true"
                    />
                    <span>{{t('pullRequest.reviewedAndAgreedMerge')}}</span>
                    <span style="color:gray">{{readableTimeComparingNow(item.created)}}</span>
                  </div>
                </template>
                <div
                  class="message-card"
                  v-else-if="item.action.actionType === 1 || item.action.actionType === 2"
                >
                  <!-- 评论或回复评论 -->
                  <div class="card-title no-wrap flex-between" :id="`comment-${item.id}`">
                    <div class="inline-flex-center">
                      <ZAvatar
                        :url="item.avatarUrl"
                        :name="item.name"
                        :account="item.account"
                        :showName="true"
                      />
                      <span>{{t('pullRequest.madeComment')}}</span>
                      <span
                        style="color:gray;margin-left:6px"
                      >{{readableTimeComparingNow(item.created)}}</span>
                    </div>
                    <a-popover
                      placement="bottomRight"
                      trigger="hover"
                      v-if="prStore.prStatus === 1"
                    >
                      <template #content>
                        <ul class="op-list">
                          <li @click="selectReply(item)">
                            <EditOutlined />
                            <span style="margin-left:6px">{{t('pullRequest.reply')}}</span>
                          </li>
                          <li @click="deleteComment(item.id)" v-if="user.account === item.account">
                            <DeleteOutlined />
                            <span style="margin-left:6px">{{t('pullRequest.delete')}}</span>
                          </li>
                        </ul>
                      </template>
                      <div class="op-icon">
                        <EllipsisOutlined />
                      </div>
                    </a-popover>
                  </div>
                  <div class="card-content">
                    <template v-if="item.action.actionType === 1">
                      <div class="comment-text">{{item.action?.comment?.comment}}</div>
                    </template>
                    <template v-else-if="item.action.actionType === 2">
                      <div
                        class="comment-reply no-wrap"
                        @click="scrollToComment(item.action.reply.fromId)"
                      >{{item.action.reply.fromAccount}}:{{item.action.reply.fromComment}}</div>
                      <div class="comment-text">{{item.action?.reply?.replyComment}}</div>
                    </template>
                  </div>
                </div>
              </a-timeline-item>
              <a-timeline-item v-if="workflowTaskList.length > 0">
                <div class="message-card">
                  <div class="card-title">
                    <NodeIndexOutlined />
                    <span style="padding-left:4px">{{t('pullRequest.workflow')}}</span>
                  </div>
                  <ul class="workflow-task-ul">
                    <li v-for="item in workflowTaskList" v-bind:key="'task_' + item.id">
                      <div>
                        <span style="cursor:pointer" @click="getWorkflowTaskStatus(item)">
                          <DownOutlined v-if="item.showJobs" />
                          <RightOutlined v-else />
                        </span>
                        <TaskStatusTag :status="item.taskStatus" style="margin-left:10px" />
                        <span
                          class="task-detail-btn"
                          @click="gotoTaskDetailPage(item)"
                        >{{item.name}}</span>
                      </div>
                      <ul class="workflow-job-ul" v-if="item.showJobs && item.jobsList?.length > 0">
                        <li v-for="job in item.jobsList" v-bind:key="`job_${job.jobName}`">
                          <RunStatus :status="job.status" :hideText="true" />
                          <span style="padding-left:10px">{{job.jobName}}</span>
                        </li>
                      </ul>
                    </li>
                  </ul>
                </div>
              </a-timeline-item>
              <a-timeline-item
                v-if="prStore.prStatus === 1 && repoStore.perm?.canAddCommentInPullRequest"
              >
                <div class="message-card">
                  <div class="card-title">
                    <EditOutlined />
                    <span style="padding-left: 6px">{{t('pullRequest.writeComment')}}</span>
                  </div>
                  <div class="card-content">
                    <div class="reply-text" v-if="replyItem.replyFrom > 0">
                      <div style="width:calc(100% - 50px)" class="no-wrap">
                        <span>{{t('pullRequest.reply')}}:</span>
                        <span>{{replyItem.fromAccount}}:</span>
                        <span>{{replyItem.fromComment}}</span>
                      </div>
                      <span class="cancel-reply" @click="cancelReply">
                        <close-circle-filled />
                      </span>
                    </div>
                    <div class="flex-center">
                      <a-textarea
                        :auto-size="{ minRows: 5, maxRows: 10 }"
                        v-model:value="replyItem.replyComment"
                        style="width:100%"
                        show-count
                        :maxlength="1024"
                        ref="commentInput"
                      />
                    </div>
                    <div style="margin-top:10px;" class="flex-right">
                      <a-button
                        type="primary"
                        style="margin-left:10px"
                        @click="addComment"
                      >{{t('pullRequest.submit')}}</a-button>
                    </div>
                  </div>
                </div>
              </a-timeline-item>
            </a-timeline>
            <div style="margin-top:13px;width:210px;margin-left:6px;" v-if="prStore.prStatus === 1">
              <div class="reviewer-section">
                <div class="title">
                  <UserOutlined />
                  <span style="padding-left:6px">{{t('pullRequest.assignedReviewer')}}</span>
                </div>
                <ul
                  class="reviewer-account-ul"
                  v-if="mergeDetect.protectedBranchCfg?.reviewerList?.length > 0"
                >
                  <li
                    v-for="item in mergeDetect.protectedBranchCfg.reviewerList"
                    v-bind:key="item.account"
                  >
                    <ZAvatar
                      :url="item.avatarUrl"
                      :name="item.name"
                      :account="item.account"
                      :showName="true"
                    />
                  </li>
                </ul>
                <div class="no-reviewer-msg" v-else>{{t('pullRequest.noAssignedReviewer')}}</div>
              </div>
              <div style="margin-top:6px" v-if="repoStore.perm?.canSubmitPullRequest">
                <a-button
                  danger
                  @click="closePr"
                  style="width:100%"
                  type="primary"
                >{{t('pullRequest.closePr')}}</a-button>
              </div>
            </div>
          </div>
        </a-tab-pane>
        <a-tab-pane
          key="diff"
          :tab="t('pullRequest.codeDiff')"
          v-if="prStore.prStatus === 1 || prStore.prStatus === 3"
        >
          <div style="width:calc(100% - 416px)">
            <ConflictFiles v-if="conflictFiles.length > 0" :conflictFiles="conflictFiles" />
            <CommitList :commits="commits" :diffNumsStats="diffNumsStats" />
            <div>
              <FileDiffDetail
                v-for="item in fileDetails"
                :stat="item"
                :repoId="repoId"
                v-bind:key="item.filePath"
                :target="headTargetCommitId.target"
                :head="headTargetCommitId.head"
              />
            </div>
          </div>
        </a-tab-pane>
        <a-tab-pane
          key="review"
          :tab="t('pullRequest.codeReview')"
          v-if="prStore.prStatus === 1 || prStore.prStatus === 3"
        >
          <div style="padding:10px;width:calc(100% - 416px)">
            <div class="section" style="margin-bottom:10px" v-if="prStore.prStatus === 1">
              <template v-if="canReviewLoaded">
                <div class="section-title">
                  <span>{{canReviewMsg}}</span>
                </div>
                <div class="section-body" v-if="canReviewState.canReview">
                  <a-button type="primary" @click="agreeReview">{{t('pullRequest.agreeMerge')}}</a-button>
                </div>
              </template>
              <div style="padding:14px" v-else>
                <LoadingOutlined />
              </div>
            </div>
            <ZTable
              :columns="reviewColumns"
              :dataSource="reviewList"
              style="margin-top:0"
              v-if="reviewList.length > 0"
            >
              <template #bodyCell="{dataIndex, dataItem}">
                <template v-if="dataIndex === 'updated'">
                  <span>{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
                </template>
                <template v-else>
                  <ZAvatar
                    :url="dataItem.avatarUrl"
                    :name="dataItem.name"
                    :account="dataItem.account"
                    :showName="true"
                  />
                </template>
              </template>
            </ZTable>
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import TaskStatusTag from "@/components/git/WorkflowTaskStatusTag";
import ZTable from "@/components/common/ZTable";
import FileDiffDetail from "@/components/git/FileDiffDetail";
import CommitList from "@/components/git/CommitList";
import ConflictFiles from "@/components/git/ConflictFiles";
import { ref, reactive, createVNode, nextTick, inject, onUnmounted } from "vue";
import { useRepoStore } from "@/pinia/repoStore";
import {
  getPullRequestRequest,
  closePullRequestRequest,
  listTimelineRequest,
  addCommentRequest,
  deleteCommentRequest,
  canMergeRequest,
  mergePullRequestRequest,
  canReviewRequest,
  agreeReviewRequest,
  listReviewRequest
} from "@/api/git/prApi";
import {
  listTaskByPrIdRequest,
  getTaskStatusRequest
} from "@/api/git/workflowApi";
import { useRoute, useRouter } from "vue-router";
import PrStatusTag from "@/components/git/PrStatusTag";
import { diffRefsRequest } from "@/api/git/repoApi";
import { message, Modal } from "ant-design-vue";
import {
  ExclamationCircleOutlined,
  LoadingOutlined,
  CheckCircleOutlined,
  FileProtectOutlined,
  CloseCircleFilled,
  CheckCircleFilled,
  WarningOutlined,
  EllipsisOutlined,
  EditOutlined,
  DeleteOutlined,
  ArrowLeftOutlined,
  NodeIndexOutlined,
  UserOutlined,
  RightOutlined,
  DownOutlined
} from "@ant-design/icons-vue";
import { readableTimeComparingNow } from "@/utils/time";
import { prCommentRegexp } from "@/utils/regexp";
import RunStatus from "@/components/git/WorkflowRunStatus";
import { useI18n } from "vue-i18n";
import { useUserStore } from "@/pinia/userStore";
const user = useUserStore();
const { t } = useI18n();
const reviewColumns = [
  {
    i18nTitle: "pullRequest.reviewer",
    dataIndex: "reviewer",
    key: "reviewer"
  },
  {
    i18nTitle: "pullRequest.agreeTime",
    dataIndex: "updated",
    key: "updated"
  }
];
const router = useRouter();
const reviewList = ref([]);
const reload = inject("gitRepoLayoutReload");
const scrollToElem = inject("gitRepoLayoutScrollToElem");
const scrollToBottom = inject("gitRepoLayoutScrollToBottom");
const cannotMergeReason = ref("");
const canReviewMsg = ref("");
const workflowTaskList = ref([]);
const workflowTaskInterval = ref(null);
const repoStore = useRepoStore();
// 冲突文件
const conflictFiles = ref([]);
// 提交列表
const commits = ref([]);
// 差异信息
const diffNumsStats = reactive({
  deleteNums: 0,
  insertNums: 0,
  fileChangeNums: 0
});
// 目标/源的commitId
const headTargetCommitId = reactive({
  head: "",
  target: ""
});
const canReviewLoaded = ref(false);
const selectTab = key => {
  switch (key) {
    case "diff":
      if (prStore.prStatus === 1 || prStore.prStatus === 3) {
        diffRefs();
      }
      break;
    case "review":
      if (prStore.prStatus === 1) {
        canReview();
      }
      listReview();
      break;
    case "timeline":
      if (prStore.prStatus === 1) {
        detectCanMerge();
      }
      listTimeline();
      break;
  }
};
// 获取审批列表
const listReview = () => {
  listReviewRequest(prStore.id).then(res => {
    reviewList.value = res.data.map(item => {
      return {
        ...item,
        key: item.id
      };
    });
  });
};
const prStatusMap = {
  1: "pullRequest.createPrStatus",
  2: "pullRequest.closePrStatus",
  3: "pullRequest.mergePrStatus"
};
const canReviewState = reactive({
  canReview: false,
  isInReviewerList: false,
  isProtectedBranch: false,
  reviewerList: [],
  hasAgree: false
});
const merging = ref(false);
const canMerge = ref(false);
const canMergeDetectLoaded = ref(false);
const mergeDetect = ref({});
const route = useRoute();
const prStore = reactive({
  commentCount: 0,
  createBy: "",
  created: "",
  head: "",
  headCommitId: "",
  headType: "",
  id: 0,
  prStatus: 0,
  prTitle: "",
  repoId: 0,
  target: "",
  targetCommitId: "",
  targetType: "",
  prIndex: 0
});
const repoId = parseInt(route.params.repoId);
const prIndex = parseInt(route.params.index);
const timelines = ref([]);
const commentInput = ref(null);
const replyItem = reactive({
  replyFrom: 0,
  fromComment: "",
  fromAccount: "",
  replyComment: ""
});
// 获取时间轴列表
const listTimeline = () => {
  listTimelineRequest(prStore.id).then(res => {
    timelines.value = res.data;
  });
};
// 检查是否可以合并
const detectCanMerge = () => {
  canMergeDetectLoaded.value = false;
  nextTick(() => {
    canMergeRequest(prStore.id).then(res => {
      if (res.data.statusChange) {
        reload();
        return;
      }
      canMergeDetectLoaded.value = true;
      mergeDetect.value = res.data;
      if (!res.data.canMerge) {
        if (res.data.gitCommitCount === 0) {
          cannotMergeReason.value = t("pullRequest.noCommits");
        } else if (
          res.data.gitConflictFiles &&
          res.data.gitConflictFiles.length > 0
        ) {
          cannotMergeReason.value = t("pullRequest.existConflictFiles");
        } else if (
          res.data.protectedBranchCfg?.reviewCountWhenCreatePr &&
          res.data.protectedBranchCfg.reviewCountWhenCreatePr >
            res.data.reviewCount
        ) {
          cannotMergeReason.value = t("pullRequest.reviewNotEnough");
        }
      }
    });
  });
};
const clearGetWorkflowTaskInterval = () => {
  if (workflowTaskInterval.value) {
    clearInterval(workflowTaskInterval.value);
    workflowTaskInterval.value = null;
  }
};
const getPullRequest = () => {
  getPullRequestRequest({
    repoId,
    index: prIndex
  }).then(res => {
    prStore.commentCount = res.data.commentCount;
    prStore.createBy = res.data.createBy;
    prStore.created = res.data.created;
    prStore.head = res.data.head;
    prStore.headCommitId = res.data.headCommitId;
    prStore.headType = res.data.headType;
    prStore.id = res.data.id;
    prStore.prStatus = res.data.prStatus;
    prStore.prTitle = res.data.prTitle;
    prStore.repoId = res.data.repoId;
    prStore.target = res.data.target;
    prStore.targetCommitId = res.data.targetCommitId;
    prStore.targetType = res.data.targetType;
    prStore.prIndex = res.data.prIndex;
    listTimeline();
    if (prStore.prStatus === 1) {
      detectCanMerge();
    }
    initGetWorkflowTaskInterval();
  });
};
const diffRefs = () => {
  if (prStore.id === 0) {
    return;
  }
  let req = {
    repoId,
    targetType: prStore.targetType,
    target: prStore.target,
    headType: prStore.headType,
    head: prStore.head
  };
  if (prStore.prStatus === 3) {
    req.head = prStore.headCommitId;
    req.headType = "commit";
    req.target = prStore.targetCommitId;
    req.targetType = "commit";
  }
  diffRefsRequest(req).then(res => {
    canMerge.value = res.data.canMerge;
    commits.value = res.data.commits;
    fileDetails.value = [];
    // 强制重新渲染
    nextTick(() => {
      conflictFiles.value = res.data.conflictFiles;
      if (res.data.diffNumsStats) {
        diffNumsStats.insertNums = res.data.diffNumsStats.insertNums;
        diffNumsStats.deleteNums = res.data.diffNumsStats.deleteNums;
        diffNumsStats.fileChangeNums = res.data.diffNumsStats.fileChangeNums;
        fileDetails.value = res.data.diffNumsStats.stats;
      } else {
        diffNumsStats.insertNums = 0;
        diffNumsStats.deleteNums = 0;
        diffNumsStats.fileChangeNums = 0;
        fileDetails.value = [];
      }
      if (res.data.headCommit) {
        headTargetCommitId.head = res.data.headCommit.commitId;
      }
      if (res.data.targetCommit) {
        headTargetCommitId.target = res.data.targetCommit.commitId;
      }
    });
  });
};
const fileDetails = ref([]);
// 关闭合并请求
const closePr = () => {
  Modal.confirm({
    title: `${t("pullRequest.confirmClosePr")}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      closePullRequestRequest(prStore.id).then(res => {
        if (res.data.statusChange) {
          reload();
          return;
        }
        message.success(t("operationSuccess"));
        prStore.prStatus = 2;
        listTimeline();
      });
    },
    onCancel() {}
  });
};
// 删除评论
const deleteComment = id => {
  Modal.confirm({
    title: `${t("pullRequest.confirmDeleteComment")}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteCommentRequest(id).then(() => {
        message.success(t("operationSuccess"));
        listTimeline();
      });
    },
    onCancel() {}
  });
};
const selectReply = item => {
  replyItem.replyFrom = item.id;
  replyItem.fromAccount = item.account;
  if (item.action.actionType === 1) {
    replyItem.fromComment = item.action.comment.comment;
  } else if (item.action.actionType === 2) {
    replyItem.fromComment = item.action.reply.replyComment;
  }
  scrollToBottom();
  if (commentInput.value) {
    commentInput.value.focus();
  }
};
const scrollToComment = id => {
  scrollToElem(`comment-${id}`);
};
const cancelReply = () => {
  replyItem.replyFrom = 0;
};
const canReview = () => {
  canReviewLoaded.value = false;
  nextTick(() => {
    canReviewRequest(prStore.id).then(res => {
      if (res.data.statusChange) {
        reload();
        return;
      }
      canReviewLoaded.value = true;
      canReviewState.canReview = res.data.canReview;
      if (
        res.data.canReview &&
        res.data.isProtectedBranch &&
        res.data.reviewerList.length === 0
      ) {
        canReviewMsg.value = t("pullRequest.noAssignedReviewerAndCanReviewPr");
      } else if (
        res.data.canReview &&
        res.data.isProtectedBranch &&
        res.data.reviewerList?.length > 0 &&
        res.data.isInReviewerList
      ) {
        canReviewMsg.value = t("pullRequest.inReviewerListAndCanReviewPr");
      } else if (res.data.canReview && !res.data.isProtectedBranch) {
        canReviewMsg.value = t("pullRequest.notProtectedBranchAndCanReviewPr");
      } else if (!res.data.canReview && res.data.hasAgree) {
        canReviewMsg.value = t("pullRequest.hasAgreePr");
      } else if (!res.data.canReview) {
        canReviewMsg.value = t("pullRequest.cantReviewPr");
      }
    });
  });
};
const addComment = () => {
  if (!prCommentRegexp.test(replyItem.replyComment)) {
    message.warn(t("pullRequest.commentFormatErr"));
    return;
  }
  addCommentRequest({
    prId: prStore.id,
    hasReply: replyItem.replyFrom > 0,
    comment: replyItem.replyComment,
    replyFrom: replyItem.replyFrom
  }).then(() => {
    message.success(t("operationSuccess"));
    replyItem.replyFrom = 0;
    replyItem.replyComment = "";
    listTimeline();
  });
};
const mergePr = () => {
  Modal.confirm({
    title: `${t("pullRequest.confirmMergePr")}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      merging.value = true;
      mergePullRequestRequest(prStore.id)
        .then(res => {
          if (!res.data.statusChange) {
            merging.value = false;
            message.success(t("operationSuccess"));
          }
          reload();
        })
        .catch(() => {
          merging.value = false;
        });
    },
    onCancel() {}
  });
};
const agreeReview = () => {
  agreeReviewRequest(prStore.id).then(res => {
    if (res.data.statusChange) {
      reload();
      return;
    }
    message.success(t("operationSuccess"));
    canReview();
    listReview();
  });
};
const getWorkflowTasks = () => {
  listTaskByPrIdRequest(prStore.id).then(res => {
    if (workflowTaskList.value.length === 0) {
      workflowTaskList.value = res.data.map(item => {
        return {
          ...item,
          showJobs: false,
          jobsList: []
        };
      });
    } else {
      let taskList = workflowTaskList.value;
      res.data.forEach(item => {
        // 判断是否存在
        let task = taskList.find(t => {
          return t.id === item.id;
        });
        if (task) {
          if (task.taskStatus !== item.taskStatus) {
            task.taskStatus = item.taskStatus;
          }
        } else {
          workflowTaskList.value.push({
            ...item,
            showJobs: false,
            jobsList: [],
            jobsInterval: null
          });
        }
      });
    }
  });
};
const doGetWorkflowTaskStatus = item => {
  getTaskStatusRequest(item.id).then(res => {
    if (["running", "queue"].includes(res.data.status)) {
      if (!item.jobsInterval) {
        item.jobsInterval = setInterval(() => {
          doGetWorkflowTaskStatus(item);
        }, 5000);
      }
    } else if (item.jobsInterval) {
      clearInterval(item.jobsInterval);
      item.jobsInterval = null;
    }
    let jobs = res.data.jobStatus;
    jobs = jobs ? jobs : [];
    item.jobsList = jobs.map(job => {
      return {
        ...job
      };
    });
  });
};
const getWorkflowTaskStatus = item => {
  if (item.showJobs) {
    item.showJobs = false;
    if (item.jobsInterval) {
      clearInterval(item.jobsInterval);
      item.jobsInterval = null;
    }
  } else {
    doGetWorkflowTaskStatus(item);
    item.showJobs = true;
  }
};

const initGetWorkflowTaskInterval = () => {
  if (prStore.prStatus === 3 && !workflowTaskInterval.value) {
    getWorkflowTasks();
    workflowTaskInterval.value = setInterval(() => getWorkflowTasks(), 5000);
  }
};
const gotoTaskDetailPage = item => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/${item.workflowId}/${item.id}/steps`
  );
};
onUnmounted(() => {
  clearGetWorkflowTaskInterval();
  workflowTaskList.value.forEach(item => {
    if (item.jobsInterval) {
      clearInterval(item.jobsInterval);
      item.jobsInterval = null;
    }
  });
});
getPullRequest();
</script>
<style scoped>
.pr-id {
  color: gray;
}
.pr-title {
  font-size: 16px;
  margin-bottom: 10px;
}
.pr-title > span + span {
  margin-left: 4px;
}
.desc {
  font-size: 14px;
  margin-bottom: 16px;
}
.desc > span + span {
  margin-left: 4px;
}
.create-by {
  font-weight: bold;
}
.message-card {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  position: relative;
  background-color: white;
}
.message-card > .text {
  padding: 14px;
}
.message-card::before {
  content: " ";
  position: absolute;
  width: 10px;
  height: 10px;
  top: 10px;
  left: -6px;
  background-color: white;
  border-top: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  transform: rotate(-45deg);
}
.card-title {
  padding: 14px;
  font-size: 14px;
}
.card-title > span + span,
.message-card > .text > span + span {
  margin-left: 4px;
}
.timeline-text {
  line-height: 32px;
  font-size: 14px;
  display: flex;
  align-items: center;
}
.timeline-text > span + span {
  padding-left: 4px;
}
.card-content {
  border-top: 1px solid #d9d9d9;
  padding: 14px;
}
.comment-reply {
  font-size: 14px;
  color: gray;
  position: relative;
  padding: 0 10px;
  cursor: pointer;
}
.comment-reply::before {
  content: " ";
  width: 4px;
  height: 22px;
  position: absolute;
  top: 0;
  left: 0;
  border-radius: 2px;
  background-color: gray;
}
.comment-text {
  font-size: 14px;
}
.timeline {
  width: 100%;
  display: flex;
}
.reply-text {
  background-color: #d9d9d9;
  border-radius: 4px;
  width: 100%;
  line-height: 32px;
  font-size: 14px;
  padding: 0px 20px;
  color: gray;
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
}
.reply-text > span + span {
  margin-left: 4px;
}
.cancel-reply {
  line-height: 32px;
  float: right;
  cursor: pointer;
}
.cancel-reply:hover {
  color: gray;
}
.can-not-merge-reason {
  font-size: 14px;
  color: darkred;
}
.can-not-merge-reason + .can-not-merge-reason {
  margin-top: 10px;
}
.workflow-task-ul {
  border-top: 1px solid #d9d9d9;
  width: 100%;
}
.workflow-task-ul > li {
  width: 100%;
  line-height: 42px;
  font-size: 14px;
  white-space: nowrap;
  word-break: break-all;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 10px 10px 10px;
  border-radius: 4px;
}
.workflow-job-ul {
  width: 100%;
  background-color: #f7f7f7;
  border-radius: 4px;
}
.workflow-job-ul > li {
  width: 100%;
  line-height: 42px;
  font-size: 14px;
  white-space: nowrap;
  word-break: break-all;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 50px;
}
.reviewer-section {
  width: 100%;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}
.reviewer-section > .title {
  font-size: 14px;
  padding: 10px 20px;
  border-bottom: 1px solid #d9d9d9;
}
.no-reviewer-msg {
  font-size: 14px;
  text-align: center;
  padding: 10px;
}
.reviewer-account-ul {
  width: 100%;
}
.reviewer-account-ul > li {
  width: 100%;
  font-size: 14px;
  height: 40px;
  padding-top: 10px;
  padding-left: 20px;
  padding-right: 20px;
}
.reviewer-account-ul > li + li {
  border-top: 1px solid #d9d9d9;
}
.task-detail-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>