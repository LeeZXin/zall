<template>
  <div style="padding:14px">
    <div class="pr-title no-wrap">
      <span class="pr-id">#{{prStore.id}}</span>
      <span>{{prStore.prTitle}}</span>
    </div>
    <div class="desc">
      <PrStatusTag :status="prStore.prStatus" />
      <span style="font-size:14px">
        {{prStore.head}}
        <ArrowLeftOutlined style="color:gray" />
        {{prStore.target}}
      </span>
    </div>
    <div>
      <a-tabs style="width: 100%;" @change="selectTab">
        <a-tab-pane key="timeline" tab="时间轴">
          <div class="timeline">
            <a-timeline style="margin-top:20px;width:calc(100% - 416px)">
              <a-timeline-item v-if="prStore.prStatus === 1">
                <div class="message-card">
                  <div class="card-title">
                    <template v-if="canMergeDetectLoaded">
                      <CheckCircleOutlined v-if="canMergeDetect.canMerge === true" />
                      <WarningOutlined v-else-if="canMergeDetect.canMerge === false" />
                      <span style="margin-left:4px">{{canMergeDetect.canMerge?'该请求能合并':'该请求不能合并'}}</span>
                    </template>
                    <template v-else>
                      <LoadingOutlined />
                    </template>
                  </div>
                  <div class="card-content" v-if="canMergeDetectLoaded">
                    <div
                      v-if="canMergeDetect.isProtectedBranch === true"
                      style="font-size:14px;color:green;margin-bottom:10px"
                    >
                      <FileProtectOutlined />
                      <span style="margin-left:4px">受保护分支规则保护</span>
                    </div>
                    <div v-if="canMergeDetect.canMerge === true">
                      <a-button @click="mergePr" type="primary">
                        <LoadingOutlined v-if="merging" />
                        <span>提交合并请求</span>
                      </a-button>
                    </div>
                    <div v-else class="can-not-merge-reason">
                      <WarningOutlined />
                      <span style="margin-left:4px">{{cannotMergeReason}}</span>
                    </div>
                  </div>
                </div>
              </a-timeline-item>
              <a-timeline-item v-for="item in timelines" v-bind:key="item.id">
                <template v-if="item.action.actionType === 3">
                  <div class="timeline-text">
                    <CloseCircleFilled
                      style="font-size:20px;color:red"
                      v-if="item.action.pr.status === 2"
                    />
                    <CheckCircleFilled style="font-size:20px;color:green" v-else />
                    <span style="font-weight:bold">{{item.account}}</span>
                    <span>{{prStatusMap[item.action.pr.status]}}</span>
                    <span style="color:gray">{{readableTimeComparingNow(item.created)}}</span>
                  </div>
                </template>
                <template v-if="item.action.actionType === 4">
                  <div class="timeline-text">
                    <CheckCircleFilled style="font-size:20px;color:green" />
                    <span style="font-weight:bold">{{item.account}}</span>
                    <span>评审并同意合并</span>
                    <span style="color:gray">{{readableTimeComparingNow(item.created)}}</span>
                  </div>
                </template>
                <div
                  class="message-card"
                  v-else-if="item.action.actionType === 1 || item.action.actionType === 2"
                >
                  <div class="card-title no-wrap" :id="`comment-${item.id}`">
                    <span style="font-weight:bold">{{item.account}}</span>
                    <span>发表评论</span>
                    <span style="color:gray">{{readableTimeComparingNow(item.created)}}</span>
                    <a-popover placement="bottomRight" trigger="hover">
                      <template #content>
                        <ul class="op-list">
                          <li @click="selectReply(item)" v-if="prStore.prStatus === 1">
                            <EditOutlined />
                            <span style="margin-left:4px">回复</span>
                          </li>
                          <li @click="deleteComment(item.id)" v-if="user.account === item.account">
                            <DeleteOutlined />
                            <span style="margin-left:4px">删除</span>
                          </li>
                        </ul>
                      </template>
                      <div
                        class="op-icon"
                        style="float:right"
                        v-if="prStore.prStatus === 1 || user.account === item.account"
                      >
                        <EllipsisOutlined />
                      </div>
                    </a-popover>
                  </div>
                  <div class="card-content" v-if="item.action.actionType === 1">
                    <div class="comment-text">{{item.action.comment.comment}}</div>
                  </div>
                  <div class="card-content" v-else-if="item.action.actionType === 2">
                    <div
                      class="comment-reply no-wrap"
                      @click="scrollToComment(item.action.reply.fromId)"
                    >{{item.action.reply.fromAccount}}:{{item.action.reply.fromComment}}</div>
                    <div class="comment-text">{{item.action.reply.replyComment}}</div>
                  </div>
                </div>
              </a-timeline-item>
              <a-timeline-item v-if="workflowTaskList.length > 0">
                <div class="message-card">
                  <div class="card-title">
                    <NodeIndexOutlined />
                    <span style="padding-left:4px">工作流</span>
                  </div>
                  <ul class="workflow-task-list">
                    <li v-for="item in workflowTaskList" v-bind:key="'task_' + item.id">
                      <div style="display:flex;justify-content:space-between">
                        <span @click="getWorkflowTaskStatus(item)" style="cursor:pointer">
                          <DownOutlined v-if="item.showJobs" />
                          <RightOutlined v-else />
                          <TaskStatusTag :status="item.taskStatus" style="margin-left:10px" />
                          <span>{{item.name}}</span>
                        </span>
                        <span class="task-detail-btn" @click="gotoTaskDetailPage(item)">查看详情</span>
                      </div>
                      <ul
                        class="workflow-job-list"
                        v-if="item.showJobs && item?.jobsList?.length > 0"
                      >
                        <li v-for="job in item.jobsList" v-bind:key="`job_${job.jobName}`">
                          <RunStatus :status="job.status" :hideText="true" />
                          <span style="padding-left:10px">{{job.jobName}}</span>
                        </li>
                      </ul>
                    </li>
                  </ul>
                </div>
              </a-timeline-item>
              <a-timeline-item v-if="prStore.prStatus === 1">
                <div class="message-card">
                  <div class="card-title">
                    <EditOutlined />
                    <span style="padding-left: 4px">编写评论</span>
                  </div>
                  <div class="card-content">
                    <div class="reply-text" v-if="replyItem.replyFrom > 0">
                      <div style="width:calc(100% - 50px)" class="no-wrap">
                        <span>回复:</span>
                        <span>{{replyItem.fromAccount}}:</span>
                        <span>{{replyItem.fromComment}}</span>
                      </div>
                      <span class="cancel-reply" @click="cancelReply">
                        <close-circle-filled />
                      </span>
                    </div>
                    <a-textarea
                      placeholder="评论"
                      :auto-size="{ minRows: 5, maxRows: 10 }"
                      v-model:value="replyItem.replyComment"
                      ref="commentInput"
                    />
                    <div style="text-align:right;margin-top:10px">
                      <a-button type="primary" style="margin-left:10px" @click="addComment">提交</a-button>
                    </div>
                  </div>
                </div>
              </a-timeline-item>
            </a-timeline>
            <div style="margin-top:13px;width:210px;margin-left:6px;" v-if="prStore.prStatus === 1">
              <div class="reviewer-section">
                <div class="title">
                  <UserOutlined />
                  <span style="padding-left:4px">指定评审人</span>
                </div>
                <ul
                  class="reviewer-account-list"
                  v-if="canMergeDetect?.protectedBranchCfg?.reviewerList?.length > 0"
                >
                  <li
                    v-for="item in canMergeDetect.protectedBranchCfg.reviewerList"
                    v-bind:key="item"
                  >{{item}}</li>
                </ul>
                <div class="no-reviewer-msg" v-else>未指定任何评审人</div>
              </div>
              <div style="margin-top:6px">
                <a-button danger @click="closePr" style="width:100%" type="primary">关闭合并请求</a-button>
              </div>
            </div>
          </div>
        </a-tab-pane>
        <a-tab-pane key="diff" tab="代码差异" v-if="prStore.prStatus === 1 || prStore.prStatus === 3">
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
        <a-tab-pane key="review" tab="评审代码" v-if="prStore.prStatus === 1 || prStore.prStatus === 3">
          <div style="padding:10px;width:calc(100% - 416px)">
            <div class="section" style="margin-bottom:10px" v-if="prStore.prStatus === 1">
              <template v-if="canReviewLoaded">
                <div class="section-title">
                  <span>{{canReviewMsg}}</span>
                </div>
                <div class="section-body" v-if="canReviewState.canReview">
                  <a-button type="primary" @click="agreeReview">同意合并</a-button>
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
                  <span>{{dataItem[dataIndex]}}</span>
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
import TaskStatusTag from "@/components/git/WorkflowTaskStatusTag";
import ZTable from "@/components/common/ZTable";
import FileDiffDetail from "@/components/git/FileDiffDetail";
import CommitList from "@/components/git/CommitList";
import ConflictFiles from "@/components/git/ConflictFiles";
import { ref, reactive, createVNode, nextTick, inject, onUnmounted } from "vue";
import { useUserStore } from "@/pinia/userStore";
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
const reviewColumns = ref([
  {
    title: "评审人",
    dataIndex: "reviewer",
    key: "reviewer"
  },
  {
    title: "操作时间",
    dataIndex: "updated",
    key: "updated"
  },
  {
    title: "状态",
    dataIndex: "reviewStatus",
    key: "reviewStatus"
  }
]);
const router = useRouter();
const reviewList = ref([]);
const user = useUserStore();
const reload = inject("gitRepoLayoutReload");
const scrollToElem = inject("gitRepoLayoutScrollToElem");
const scrollToBottom = inject("gitRepoLayoutScrollToBottom");
const cannotMergeReason = ref("");
const canReviewMsg = ref("");
const workflowTaskList = ref([]);
const workflowTaskInterval = ref(null);
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
      if (prStore.value.prStatus === 1 || prStore.value.prStatus === 3) {
        diffRefs();
      }
      break;
    case "review":
      if (prStore.value.prStatus === 1) {
        canReview();
      }
      listReview();
      break;
    case "timeline":
      if (prStore.value.prStatus === 1) {
        detectCanMerge();
      }
      listTimeline();
      break;
  }
};
const listReview = () => {
  listReviewRequest(route.params.prId).then(res => {
    reviewList.value = res.data.map(item => {
      return {
        ...item,
        key: item.id
      };
    });
  });
};
const prStatusMap = {
  1: "创建合并请求",
  2: "关闭合并请求",
  3: "提交合并请求"
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
const canMergeDetect = ref({});
const route = useRoute();
const prStore = ref({});
const repoId = parseInt(route.params.repoId);
const timelines = ref([]);
const commentInput = ref(null);
const replyItem = reactive({
  replyFrom: 0,
  fromComment: "",
  fromAccount: "",
  replyComment: ""
});
const listTimeline = () => {
  listTimelineRequest(prStore.value.id).then(res => {
    timelines.value = res.data;
  });
};
const detectCanMerge = () => {
  canMergeDetectLoaded.value = false;
  nextTick(() => {
    canMergeRequest(prStore.value.id).then(res => {
      if (res.data.statusChange) {
        reload();
        return;
      }
      canMergeDetectLoaded.value = true;
      canMergeDetect.value = res.data;
      if (!res.data.canMerge) {
        if (res.data.gitCommitCount === 0) {
          cannotMergeReason.value = "无任何代码提交";
        } else if (
          res.data.gitConflictFiles &&
          res.data.gitConflictFiles.length > 0
        ) {
          cannotMergeReason.value = "存在冲突文件";
        } else if (
          res.data.protectedBranchCfg?.reviewCountWhenCreatePr &&
          res.data.protectedBranchCfg.reviewCountWhenCreatePr >
            res.data.reviewCount
        ) {
          cannotMergeReason.value = "评审数量不足";
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
  getPullRequestRequest(route.params.prId).then(res => {
    prStore.value.commentCount = res.data.commentCount;
    prStore.value.createBy = res.data.createBy;
    prStore.value.created = res.data.created;
    prStore.value.head = res.data.head;
    prStore.value.headCommitId = res.data.headCommitId;
    prStore.value.headType = res.data.headType;
    prStore.value.id = res.data.id;
    prStore.value.prComment = res.data.prComment;
    prStore.value.prStatus = res.data.prStatus;
    prStore.value.prTitle = res.data.prTitle;
    prStore.value.repoId = res.data.repoId;
    prStore.value.target = res.data.target;
    prStore.value.targetCommitId = res.data.targetCommitId;
    prStore.value.targetType = res.data.targetType;
    listTimeline();
    if (prStore.value.prStatus === 1) {
      detectCanMerge();
    }
    initGetWorkflowTaskInterval();
  });
};
const diffRefs = () => {
  if (prStore.value.id === 0) {
    return;
  }
  let req = {
    repoId,
    targetType: prStore.value.targetType,
    target: prStore.value.target,
    headType: prStore.value.headType,
    head: prStore.value.head
  };
  if (prStore.value.prStatus === 3) {
    req.head = prStore.value.headCommitId;
    req.headType = "commit";
    req.target = prStore.value.targetCommitId;
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
    title: "你确定要关闭?",
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      closePullRequestRequest(route.params.prId).then(res => {
        console.log(res);
        if (res.data.statusChange) {
          reload();
          return;
        }
        message.success("关闭成功");
        prStore.value.prStatus = 2;
        listTimeline();
      });
    },
    onCancel() {}
  });
};
// 删除评论
const deleteComment = id => {
  Modal.confirm({
    title: "你确定要删除?",
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteCommentRequest(id).then(() => {
        message.success("删除成功");
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
    canReviewRequest(prStore.value.id).then(res => {
      if (res.data.statusChange) {
        reload();
        return;
      }
      canReviewLoaded.value = true;
      canReviewState.canReview = res.data.canReview;
      canReviewState.isInReviewerList = res.data.isInReviewerList;
      canReviewState.isProtectedBranch = res.data.isProtectedBranch;
      canReviewState.reviewerList = res.data.reviewerList;
      canReviewState.hasAgree = res.data.hasAgree;
      if (
        res.data.canReview &&
        res.data.isProtectedBranch &&
        res.data.reviewerList.length === 0
      ) {
        canReviewMsg.value = "无指定任何评审人, 你可以评审该合并请求";
      } else if (
        canReviewState.canReview &&
        canReviewState.isProtectedBranch &&
        canReviewState.reviewerList.length > 0 &&
        canReviewState.isInReviewerList
      ) {
        canReviewMsg.value = "你在指定评审名单里, 你可以评审该合并请求";
      } else if (
        canReviewState.canReview &&
        !canReviewState.isProtectedBranch
      ) {
        canReviewMsg.value = "非保护分支, 你可以评审该合并请求";
      } else if (!canReviewState.canReview && canReviewState.hasAgree) {
        canReviewMsg.value = "你已同意合并该请求";
      } else if (!canReviewState.canReview) {
        canReviewMsg.value = "你不可以评审该合并请求";
      }
    });
  });
};
const addComment = () => {
  if (!prCommentRegexp.test(replyItem.replyComment)) {
    message.warn("评论格式不合法");
    return;
  }
  addCommentRequest({
    prId: prStore.value.id,
    hasReply: replyItem.replyFrom > 0,
    comment: replyItem.replyComment,
    replyFrom: replyItem.replyFrom
  }).then(() => {
    message.success("提交成功");
    replyItem.replyFrom = 0;
    replyItem.replyComment = "";
    listTimeline();
  });
};
const mergePr = () => {
  Modal.confirm({
    title: "你确定要提交合并吗?",
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      merging.value = true;
      mergePullRequestRequest(prStore.value.id)
        .then(res => {
          if (!res.data.statusChange) {
            merging.value = false;
            message.success("合并成功");
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
  agreeReviewRequest(prStore.value.id).then(res => {
    if (res.data.statusChange) {
      reload();
      return;
    }
    message.success("操作成功");
    canReview();
    listReview();
  });
};
const getWorkflowTasks = () => {
  listTaskByPrIdRequest(route.params.prId).then(res => {
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
  if (prStore.value.prStatus === 3 && !workflowTaskInterval.value) {
    getWorkflowTasks();
    workflowTaskInterval.value = setInterval(() => getWorkflowTasks(), 5000);
  }
};
const gotoTaskDetailPage = item => {
  router.push(
    `/gitRepo/${route.params.repoId}/workflow/${item.workflowId}/${item.id}/steps`
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
  font-size: 18px;
  margin-bottom: 10px;
}
.pr-title > span + span {
  margin-left: 4px;
}
.desc {
  font-size: 14px;
  margin-bottom: 10px;
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
  border-bottom: 1px solid #d9d9d9;
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
.workflow-task-list {
  width: 100%;
}
.workflow-task-list > li {
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
.workflow-job-list {
  width: 100%;
  background-color: #f7f7f7;
}
.workflow-job-list > li {
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
.reviewer-account-list {
  width: 100%;
}
.reviewer-account-list > li {
  width: 100%;
  font-size: 14px;
  padding: 10px 20px;
}
.reviewer-account-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.task-detail-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>