<template>
  <div style="padding:14px">
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'pullRequest'">
          <template v-if="dataItem[dataIndex]">
            <a-button
              type="link"
              @click="toPrDetail(dataItem[dataIndex].id)"
            >#{{dataItem[dataIndex].id}}</a-button>
            <PrStatusTag :status="dataItem[dataIndex].prStatus" />
          </template>
        </template>
        <template v-else-if="dataIndex === 'lastCommitTime'">
          <span>{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Branch</span>
              </template>
              <delete-outlined @click="deleteBranch(dataItem['name'])" />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="click">
            <template #content>
              <ul class="op-list">
                <li @click="goToHistoryCommits(dataItem['name'])">
                  <control-outlined />
                  <span style="margin-left:4px">查看所有的活动</span>
                </li>
                <li></li>
              </ul>
            </template>
            <div class="op-icon">...</div>
          </a-popover>
        </template>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  DeleteOutlined,
  ControlOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import { allBranchCommitRequest, deleteBranchRequest } from "@/api/git/repoApi";
import { readableTimeComparingNow } from "@/utils/time";
import PrStatusTag from "@/components/git/PrStatusTag";
import { Modal, message } from "ant-design-vue";
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
const columns = ref([
  {
    title: "分支",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "最后更新时间",
    dataIndex: "lastCommitTime",
    key: "lastCommitTime"
  },
  {
    title: "合并请求",
    dataIndex: "pullRequest",
    key: "pullRequest"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);
const listBranch = () => {
  allBranchCommitRequest(route.params.repoId).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        name: item.name,
        lastCommitTime: item.lastCommit.committedTime,
        pullRequest: item.lastPullRequest
      };
    });
  });
};
const deleteBranch = branch => {
  Modal.confirm({
    title: `你确定要删除${branch}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteBranchRequest({
        repoId: parseInt(route.params.repoId),
        branch
      }).then(() => {
        message.success("删除成功");
        listBranch();
      });
    },
    onCancel() {}
  });
};
const toPrDetail = id => {
  router.push(`/gitRepo/${route.params.repoId}/pullRequest/${id}/detail`);
};
const goToHistoryCommits = branch => {
  router.push(`/gitRepo/${route.params.repoId}/commit/list/${branch}`);
};
listBranch();
</script>
<style scoped>
.header {
  line-height: 32px;
  font-size: 18px;
}
</style>