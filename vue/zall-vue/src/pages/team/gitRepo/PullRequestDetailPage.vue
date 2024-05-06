<template>
  <div style="padding:14px">
    <ZNaviBack name="合并请求列表" :url="`/gitRepo/${route.params.repoId}/pullRequest/list`" />
    <div class="title no-wrap">
      <span class="pr-id">
        <span>#</span>
        <span>{{prStore.id}}</span>
      </span>
      <span>{{prStore.prTitle}}</span>
    </div>
    <div class="desc">
      <PrStatusTag :status="prStore.prStatus" />
      <span class="create-by">{{prStore.createBy}}</span>
      <span>wants to merge into {{prStore.head}} from {{prStore.target}}</span>
    </div>
    <div>
      <a-tabs style="width: 100%;" @change="selectTab">
        <a-tab-pane key="timeline" tab="时间轴">
          <div class="timeline">
            <a-timeline style="margin-top:20px;width:calc(100% - 316px)">
              <a-timeline-item v-if="prStore.prStatus === 1">
                <div class="message-card">
                  <div class="card-title">
                    <template v-if="canMergeDetectLoaded">
                      <span v-if="canMergeDetect.canMerge === true">
                        <CheckCircleOutlined />
                        <span style="margin-left:4px">该请求能合并</span>
                      </span>
                      <span v-else-if="canMergeDetect.canMerge === false">
                        <WarningOutlined />
                        <span style="margin-left:4px">该请求不能合并</span>
                      </span>
                    </template>
                    <template v-else>
                      <LoadingOutlined />
                    </template>
                  </div>
                  <div class="card-content" v-if="canMergeDetect.canMerge === true">
                    <a-button @click="mergePr" type="primary">
                      <LoadingOutlined v-if="merging" />
                      <span>提交合并请求</span>
                    </a-button>
                  </div>
                </div>
              </a-timeline-item>
              <a-timeline-item v-for="item in timelines" v-bind:key="item.id">
                <div class="message-card">
                  <template v-if="item.action.actionType === 3">
                    <div class="text">
                      <span>{{item.account}}</span>
                      <span v-if="item.action.pr.status === 1">创建合并请求</span>
                      <span v-else-if="item.action.pr.status === 2">关闭合并请求</span>
                      <span v-else-if="item.action.pr.status === 3">提交合并请求</span>
                      <span>#{{item.prId}}</span>
                      <span>{{readableTimeComparingNow(item.created)}}</span>
                    </div>
                  </template>
                  <template v-else-if="item.action.actionType === 1">
                    <div class="card-title no-wrap" :id="`comment-${item.id}`">
                      <span>{{item.account}}</span>
                      <span>评论于</span>
                      <span>{{readableTimeComparingNow(item.created)}}</span>
                      <span
                        class="reply-btn"
                        @click="selectReply(item)"
                        v-if="prStore.prStatus === 1"
                      >回复</span>
                      <span
                        class="del-btn"
                        @click="deleteComment(item.id)"
                        v-if="user.account === item.account"
                      >删除</span>
                    </div>
                    <div class="card-content">
                      <div class="comment-text">{{item.action.comment.comment}}</div>
                    </div>
                  </template>
                  <template v-else-if="item.action.actionType === 2">
                    <div class="card-title no-wrap" :id="`comment-${item.id}`">
                      <span>{{item.account}}</span>
                      <span>评论于</span>
                      <span>{{readableTimeComparingNow(item.created)}}</span>
                      <span
                        class="reply-btn"
                        @click="selectReply(item)"
                        v-if="prStore.prStatus === 1"
                      >回复</span>
                      <span
                        class="del-btn"
                        @click="deleteComment(item.id)"
                        v-if="user.account === item.account"
                      >删除</span>
                    </div>
                    <div class="card-content">
                      <div
                        class="comment-reply no-wrap"
                        @click="scrollToComment(item.action.reply.fromId)"
                      >{{item.action.reply.fromAccount}}:{{item.action.reply.fromComment}}</div>
                      <div class="comment-text">{{item.action.reply.replyComment}}</div>
                    </div>
                  </template>
                </div>
              </a-timeline-item>
              <a-timeline-item v-if="prStore.prStatus === 1">
                <div class="message-card">
                  <div class="card-title">编写评论</div>
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
                      <a-button danger @click="closePr" ghost>关闭合并请求</a-button>
                      <a-button type="primary" style="margin-left:10px" @click="addComment">提交</a-button>
                    </div>
                  </div>
                </div>
              </a-timeline-item>
            </a-timeline>
          </div>
        </a-tab-pane>
        <a-tab-pane key="diff" tab="代码差异" v-if="prStore.prStatus === 1 || prStore.prStatus === 3">
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
        </a-tab-pane>
        <a-tab-pane
          key="reviewer"
          tab="评审人"
          v-if="prStore.prStatus === 1 || prStore.prStatus === 3"
        ></a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>
<script setup>
import ZNaviBack from "@/components/common/ZNaviBack";
import FileDiffDetail from "@/components/git/FileDiffDetail";
import CommitList from "@/components/git/CommitList";
import ConflictFiles from "@/components/git/ConflictFiles";
import { ref, reactive, createVNode, nextTick, inject } from "vue";
import { usePrStore } from "@/pinia/prStore";
import { useUserStore } from "@/pinia/userStore";
import {
  getPullRequestRequest,
  closePullRequestRequest,
  listTimelineRequest,
  addCommentRequest,
  deleteCommentRequest,
  canMergeRequest,
  mergePullRequestRequest
} from "@/api/git/prApi";
import { useRoute } from "vue-router";
import PrStatusTag from "@/components/git/PrStatusTag";
import { diffRefsRequest } from "@/api/git/repoApi";
import { message, Modal } from "ant-design-vue";
import {
  ExclamationCircleOutlined,
  LoadingOutlined,
  CheckCircleOutlined
} from "@ant-design/icons-vue";
import { readableTimeComparingNow } from "@/utils/time";
import { CloseCircleFilled, WarningOutlined } from "@ant-design/icons-vue";
import { prCommentRegexp } from "@/utils/regexp";
const user = useUserStore();
const reload = inject("gitRepoLayoutReload");
const scrollToElem = inject("gitRepoLayoutScrollToElem");
const scrollToBottom = inject("gitRepoLayoutScrollToBottom");
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
const diffRefsLoaded = ref(false);
const selectTab = key => {
  if (key === "diff" && !diffRefsLoaded.value) {
    diffRefs();
  }
};
const merging = ref(false);
const canMerge = ref(false);
const canMergeDetectLoaded = ref(false);
const canMergeDetect = ref({});
const route = useRoute();
const prStore = usePrStore();
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
  listTimelineRequest(prStore.id).then(res => {
    timelines.value = res.data;
  });
};
const detectCanMerge = () => {
  canMergeRequest(prStore.id).then(res => {
    canMergeDetectLoaded.value = true;
    canMergeDetect.value = res.data;
  });
};
if (prStore.id === 0 || prStore.id !== parseInt(route.params.prId)) {
  getPullRequestRequest(route.params.prId).then(res => {
    prStore.commentCount = res.data.commentCount;
    prStore.createBy = res.data.createBy;
    prStore.created = res.data.created;
    prStore.head = res.data.head;
    prStore.headCommitId = res.data.headCommitId;
    prStore.headType = res.data.headType;
    prStore.id = res.data.id;
    prStore.prComment = res.data.prComment;
    prStore.prStatus = res.data.prStatus;
    prStore.prTitle = res.data.prTitle;
    prStore.repoId = res.data.repoId;
    prStore.target = res.data.target;
    prStore.targetCommitId = res.data.targetCommitId;
    prStore.targetType = res.data.targetType;
    listTimeline();
    if (prStore.prStatus === 1) {
      detectCanMerge();
    }
  });
} else {
  listTimeline();
  if (prStore.prStatus === 1) {
    detectCanMerge();
  }
}
const diffRefs = () => {
  if (prStore.id === 0) {
    return;
  }
  diffRefsLoaded.value = true;
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
    title: "你确定要关闭?",
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      closePullRequestRequest(route.params.prId).then(() => {
        message.success("关闭成功");
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
const addComment = () => {
  if (!prCommentRegexp.test(replyItem.replyComment)) {
    message.warn("评论格式不合法");
    return;
  }
  addCommentRequest({
    prId: prStore.id,
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
      mergePullRequestRequest(prStore.id)
        .then(() => {
          merging.value = false;
          message.success("合并成功");
          getPullRequestRequest(route.params.prId).then(res => {
            prStore.commentCount = res.data.commentCount;
            prStore.createBy = res.data.createBy;
            prStore.created = res.data.created;
            prStore.head = res.data.head;
            prStore.headCommitId = res.data.headCommitId;
            prStore.headType = res.data.headType;
            prStore.id = res.data.id;
            prStore.prComment = res.data.prComment;
            prStore.prStatus = res.data.prStatus;
            prStore.prTitle = res.data.prTitle;
            prStore.repoId = res.data.repoId;
            prStore.target = res.data.target;
            prStore.targetCommitId = res.data.targetCommitId;
            prStore.targetType = res.data.targetType;
            reload();
          });
        })
        .catch(() => {
          merging.value = false;
        });
    },
    onCancel() {}
  });
};
</script>
<style scoped>
.pr-id {
  color: gray;
}
.title {
  font-size: 18px;
  margin-bottom: 10px;
}
.title > span + span {
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
.message-card:before {
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
.card-content {
  border-top: 1px solid #d9d9d9;
  padding: 14px;
}
.comment-reply {
  font-size: 14px;
  color: gray;
  position: relative;
  border-left: 4px solid gray;
  padding: 6px 10px;
  cursor: pointer;
}
.comment-text {
  font-size: 14px;
  padding: 10px;
}
.timeline {
  width: 100%;
  display: flex;
}
.reply-btn,
.del-btn {
  float: right;
}
.del-btn {
  color: darkred;
}
.del-btn:hover {
  cursor: pointer;
  color: red;
}
.reply-btn:hover {
  color: #1677ff;
  cursor: pointer;
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
</style>