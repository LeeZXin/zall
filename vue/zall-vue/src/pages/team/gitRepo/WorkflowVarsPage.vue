<template>
  <div style="padding:10px">
    <div class="header">
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">添加变量</a-button>
    </div>
    <ul class="vars-list" v-if="varsList.length > 0">
      <li v-for="item in varsList" v-bind:key="item.id">
        <div class="vars-pattern no-wrap">{{item.name}}</div>
        <ul class="op-btns">
          <li class="update-btn" @click="handleVars(item)">编辑</li>
          <li class="del-btn" @click="deleteVars(item)">删除</li>
        </ul>
      </li>
    </ul>
    <ZNoData v-else>
      <template #desc>
        <div
          class="no-data-text"
        >Variables are encrypted and are used for sensitive or long data</div>
      </template>
    </ZNoData>
  </div>
</template>
<script setup>
import { ref, createVNode, h } from "vue";
import { useRouter, useRoute } from "vue-router";
import ZNoData from "@/components/common/ZNoData";
import { ExclamationCircleOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { message, Modal } from "ant-design-vue";
import { listVarsRequest, deleteVarsRequest } from "@/api/git/workflowApi";
const router = useRouter();
const route = useRoute();
const varsList = ref([]);
const gotoCreatePage = () => {
  router.push(`/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/vars/create`);
};
const deleteVars = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteVarsRequest(item.id).then(() => {
        message.success("删除成功");
        listVars();
      });
    },
    onCancel() {}
  });
};
const listVars = () => {
  listVarsRequest(route.params.repoId).then(res => {
    varsList.value = res.data;
  });
};
const handleVars = item => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/vars/${item.id}/update`
  );
};
listVars();
</script>
<style scoped>
.vars-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.vars-list > li {
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.vars-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.vars-pattern {
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
.ping-btn:hover {
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