<template>
  <div style="padding:10px">
    <div style="margin-bottom:16px">
      <div style="font-size:16px;font-weight:bold;margin-bottom:8px">仓库回收站</div>
    </div>
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchRepo"
        :placeholder="t('gitRepo.searchText')"
        style="width:240px;margin-right:10px"
        @change="searchChange"
      />
    </div>
    <ZTable :columns="columns" :dataSource="repoList">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex === 'operation'">
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="deleteRepoPermanently(dataItem)">
                  <DeleteOutlined />
                  <span style="margin-left:4px">永久删除</span>
                </li>
                <li @click="recoverRepo(dataItem)">
                  <RollbackOutlined />
                  <span style="margin-left:4px">恢复</span>
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
    title: "仓库名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "描述",
    dataIndex: "repoDesc",
    key: "repoDesc"
  },
  {
    title: "仓库大小",
    dataIndex: "gitSize",
    key: "gitSize"
  },
  {
    title: "lfs大小",
    dataIndex: "lfsSize",
    key: "lfsSize"
  },
  {
    title: "删除时间",
    dataIndex: "deleted",
    key: "deleted"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
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
    title: `你确定要永久删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteRepoPermanentlyRequest(item.repoId).then(() => {
        message.success("删除成功");
        getDeletedRepo();
      });
    },
    onCancel() {}
  });
};
const recoverRepo = item => {
  Modal.confirm({
    title: `你确定要恢复${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      reoverFromRecycleRequest(item.repoId).then(() => {
        message.success("恢复成功");
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