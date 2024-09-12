<template>
  <div style="padding:10px">
    <div style="margin-bottom:16px">
      <div style="font-size:16px;font-weight:bold;margin-bottom:8px">{{t('gitRepo.recycle')}}</div>
    </div>
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchRepo"
        :placeholder="t('gitRepo.searchName')"
        style="width:240px;margin-right:10px"
        @change="searchChange"
      />
    </div>
    <ZTable :columns="columns" :dataSource="repoList" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex === 'operation'">
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="deleteRepoPermanently(dataItem)">
                  <DeleteOutlined />
                  <span style="margin-left:4px">{{t('gitRepo.deletePermanently')}}</span>
                </li>
                <li @click="recoverRepo(dataItem)">
                  <RollbackOutlined />
                  <span style="margin-left:4px">{{t('gitRepo.recoverFromRecycle')}}</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </span>
        <span v-else-if="dataIndex === 'gitSize'">{{readableVolumeSize(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'lfsSize'">{{readableVolumeSize(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex === 'deleted'">{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        <span v-else>{{dataItem[dataIndex]}}</span>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  DeleteOutlined,
  EllipsisOutlined,
  RollbackOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import { ref, createVNode } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";
import {
  getDeletedRepoListRequest,
  deleteRepoPermanentlyRequest,
  reoverFromRecycleRequest
} from "@/api/git/repoApi";
import { readableVolumeSize } from "@/utils/size";
import { readableTimeComparingNow } from "@/utils/time";
import { Modal, message } from "ant-design-vue";
const { t } = useI18n();
const route = useRoute();
// 搜索框key
const searchRepo = ref("");
// 所有仓库列表
const wholeRepoList = ref([]);
// 搜索框检索后仓库列表
const repoList = ref(wholeRepoList.value);
// searchChange 搜索框触发搜索
const searchChange = () => {
  let searchKey = searchRepo.value;
  if (!searchKey || searchKey === "") {
    repoList.value = wholeRepoList.value;
    return;
  }
  repoList.value = wholeRepoList.value.filter(item => {
    return item.name.indexOf(searchKey) >= 0;
  });
};
const columns = [
  {
    i18nTitle: "gitRepo.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "gitRepo.repoDesc",
    dataIndex: "repoDesc",
    key: "repoDesc"
  },
  {
    i18nTitle: "gitRepo.gitSize",
    dataIndex: "gitSize",
    key: "gitSize"
  },
  {
    i18nTitle: "gitRepo.lfsSize",
    dataIndex: "lfsSize",
    key: "lfsSize"
  },
  {
    i18nTitle: "gitRepo.deleteTime",
    dataIndex: "deleted",
    key: "deleted"
  },
  {
    i18nTitle: "gitRepo.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
const getDeletedRepo = () => {
  // 获取仓库列表
  getDeletedRepoListRequest(route.params.teamId).then(res => {
    const ret = res.data.map(item => {
      return {
        key: item.repoId,
        ...item
      };
    });
    wholeRepoList.value = ret;
    repoList.value = ret;
  });
};
const deleteRepoPermanently = item => {
  Modal.confirm({
    title: `你确定要永久删除${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteRepoPermanentlyRequest(item.repoId).then(() => {
        message.success(t("operationSuccess"));
        getDeletedRepo();
      });
    },
    onCancel() {}
  });
};
const recoverRepo = item => {
  Modal.confirm({
    title: `你确定要恢复${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      reoverFromRecycleRequest(item.repoId).then(() => {
        message.success(t("operationSuccess"));
        getDeletedRepo();
      });
    },
    onCancel() {}
  });
};
getDeletedRepo();
</script>
<style scoped>
.no-data {
  font-size: 16px;
  text-align: center;
}
</style>