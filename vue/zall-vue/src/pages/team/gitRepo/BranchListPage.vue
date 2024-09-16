<template>
  <div style="padding:10px">
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'pullRequest'">
          <template v-if="dataItem[dataIndex]">
            <PrIndexTag
              :repoId="route.params.repoId"
              :prIndex="dataItem[dataIndex].prIndex"
              :teamId="route.params.teamId"
            />
            <PrStatusTag :status="dataItem[dataIndex].prStatus" />
          </template>
        </template>
        <template v-else-if="dataIndex === 'isProtectedBranch'">
          <span>{{dataItem[dataIndex]?t('branchList.yes'):t('branchList.no')}}</span>
        </template>
        <template v-else-if="dataIndex === 'lastCommitTime'">
          <span>{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <div
            class="op-icon"
            v-if="!dataItem['isProtectedBranch']"
            @click="deleteBranch(dataItem['name'])"
          >
            <DeleteOutlined />
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="goToHistoryCommits(dataItem['name'])">
                  <ControlOutlined />
                  <span style="margin-left:4px">{{t('branchList.viewCommits')}}</span>
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
    <a-pagination
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listBranch()"
    />
  </div>
</template>
<script setup>
import PrIndexTag from "@/components/git/PrIndexTag";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  DeleteOutlined,
  ControlOutlined,
  ExclamationCircleOutlined,
  EllipsisOutlined
} from "@ant-design/icons-vue";
import {
  listBranchCommitsRequest,
  deleteBranchRequest
} from "@/api/git/repoApi";
import { readableTimeComparingNow } from "@/utils/time";
import PrStatusTag from "@/components/git/PrStatusTag";
import { Modal, message } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const dataSource = ref([]);
// 分页
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 数据项
const columns = [
  {
    i18nTitle: "branchList.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "branchList.lastCommitTime",
    dataIndex: "lastCommitTime",
    key: "lastCommitTime"
  },
  {
    i18nTitle: "branchList.isProtectedBranch",
    dataIndex: "isProtectedBranch",
    key: "isProtectedBranch"
  },
  {
    i18nTitle: "branchList.pullRequest",
    dataIndex: "pullRequest",
    key: "pullRequest"
  },
  {
    i18nTitle: "branchList.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 获取列表
const listBranch = () => {
  listBranchCommitsRequest({
    repoId: route.params.repoId,
    pageNum: dataPage.current
  }).then(res => {
    dataPage.totalCount = res.totalCount;
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
// 删除分支
const deleteBranch = branch => {
  Modal.confirm({
    title: `${t("branchList.confirmDelete")} ${branch}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteBranchRequest({
        repoId: parseInt(route.params.repoId),
        branch
      }).then(() => {
        message.success(t("operationSuccess"));
        dataPage.current = 1;
        listBranch();
      });
    },
    onCancel() {}
  });
};
// 跳转提交列表
const goToHistoryCommits = branch => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/commit/list/${branch}`
  );
};
listBranch();
</script>
<style scoped>
</style>