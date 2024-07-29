<template>
  <div style="padding:10px">
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0" v-if="totalCount > 0">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'pullRequest'">
          <template v-if="dataItem[dataIndex]">
            <PrIdTag :repoId="route.params.repoId" :prId="dataItem[dataIndex].id" />
            <PrStatusTag :status="dataItem[dataIndex].prStatus" />
          </template>
        </template>
        <template v-else-if="dataIndex === 'isProtectedBranch'">
          <span>{{dataItem[dataIndex]?'是':'否'}}</span>
        </template>
        <template v-else-if="dataIndex === 'lastCommitTime'">
          <span>{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div class="op-icon" v-if="!dataItem['isProtectedBranch']">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Branch</span>
              </template>
              <delete-outlined @click="deleteBranch(dataItem['name'])" />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="goToHistoryCommits(dataItem['name'])">
                  <control-outlined />
                  <span style="margin-left:4px">查看所有的活动</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </template>
      </template>
    </ZTable>
    <ZNoData v-else>
      <template #desc>
        <div style="font-size:14px;text-align:center">
          <span>无分支数据, 尝试去</span>
          <span class="suggest-text" @click="gotoIndex">提交代码</span>
        </div>
      </template>
    </ZNoData>
    <a-pagination
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listBranch()"
    />
  </div>
</template>
<script setup>
import PrIdTag from "@/components/git/PrIdTag";
import ZNoData from "@/components/common/ZNoData";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  DeleteOutlined,
  ControlOutlined,
  ExclamationCircleOutlined,
  EllipsisOutlined
} from "@ant-design/icons-vue";
import {
  pageBranchCommitsRequest,
  deleteBranchRequest
} from "@/api/git/repoApi";
import { readableTimeComparingNow } from "@/utils/time";
import PrStatusTag from "@/components/git/PrStatusTag";
import { Modal, message } from "ant-design-vue";
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const columns = [
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
    title: "是否是保护分支",
    dataIndex: "isProtectedBranch",
    key: "isProtectedBranch"
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
];
const listBranch = () => {
  pageBranchCommitsRequest({
    repoId: route.params.repoId,
    pageNum: currPage.value
  }).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        name: item.name,
        lastCommitTime: item.lastCommit.committedTime,
        pullRequest: item.lastPullRequest,
        isProtectedBranch: item.isProtectedBranch
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
        if (totalCount.value - 1 <= (currPage.value - 1) * pageSize) {
          currPage.value -= 1;
        }
        message.success("删除成功");
        listBranch();
      });
    },
    onCancel() {}
  });
};
const goToHistoryCommits = branch => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/commit/list/${branch}`
  );
};
const gotoIndex = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/index`
  );
};
listBranch();
</script>
<style scoped>
</style>