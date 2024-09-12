<template>
  <div style="padding:10px">
    <div class="header">
      <a-button
        type="primary"
        @click="gotoCreatePage"
        :icon="h(PlusOutlined)"
      >{{t('protectedBranch.createBranch')}}</a-button>
    </div>
    <ul class="branch-list" v-if="branches.length > 0">
      <li v-for="item in branches" v-bind:key="item.id">
        <div class="branch-pattern no-wrap">{{item.pattern}}</div>
        <ul class="op-btns">
          <li
            class="update-btn"
            @click="handleProtectedBranch(item)"
          >{{t('protectedBranch.update')}}</li>
          <li class="del-btn" @click="deleteProtectedBranch(item)">{{t('protectedBranch.delete')}}</li>
        </ul>
      </li>
    </ul>
    <ZNoData v-else />
  </div>
</template>
<script setup>
import { ref, createVNode, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import ZNoData from "@/components/common/ZNoData";
import {
  listProtectedBranchRequest,
  deleteProtectedBranchRequest
} from "@/api/git/branchApi";
import { ExclamationCircleOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { message, Modal } from "ant-design-vue";
import { useProtectedBranchStore } from "@/pinia/protectedBranchStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const route = useRoute();
const branches = ref([]);
const protectedBranchStore = useProtectedBranchStore();
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/protectedBranch/create`
  );
};
// 删除保护分支
const deleteProtectedBranch = item => {
  Modal.confirm({
    title: `${t('protectedBranch.confirmDelete')} ${item.pattern}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteProtectedBranchRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listProtectedBranch();
      });
    },
    onCancel() {}
  });
};
// 获取列表
const listProtectedBranch = () => {
  listProtectedBranchRequest(route.params.repoId).then(res => {
    branches.value = res.data;
  });
};
// 跳转编辑页面
const handleProtectedBranch = item => {
  protectedBranchStore.id = item.id;
  protectedBranchStore.pattern = item.pattern;
  protectedBranchStore.cfg = item.cfg;
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/protectedBranch/${item.id}/update`
  );
};
listProtectedBranch();
</script>
<style scoped>
.branch-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.branch-list > li {
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.branch-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.branch-pattern {
  font-size: 14px;
  line-height: 32px;
  width: 60%;
  padding-left: 10px;
}
.op-btns {
  display: flex;
  align-items: center;
}
.op-btns > li {
  line-height: 32px;
  font-size: 14px;
  padding: 0 10px;
  cursor: pointer;
}
.op-btns > li:first-child {
  border-top: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}
.op-btns > li:not(:first-child, :last-child) {
  border-top: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
}
.op-btns > li:last-child {
  border-top: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
}
.op-btns > li + li {
  border-left: 1px solid #d9d9d9;
}
.header {
  margin-bottom: 10px;
}
.header > span {
  font-size: 18px;
  font-weight: bold;
  line-height: 32px;
  padding-left: 8px;
}
.del-btn {
  color: darkred;
}
.del-btn:hover {
  color: white;
  background-color: darkred;
}
.update-btn:hover,
.view-btn:hover {
  background-color: #f0f0f0;
}
.no-data-text {
  font-size: 14px;
  line-height: 18px;
  padding: 10px;
  text-align: center;
  word-break: break-all;
}
</style>