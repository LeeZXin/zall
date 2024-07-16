<template>
  <div style="padding:14px">
    <div class="container">
      <div class="section">
        <div class="section-title">
          <span>应用名称</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-input v-model:value="appName" />
            <div class="input-desc">简单描述应用的作用</div>
          </div>
          <div class="input-item">
            <a-button type="primary" @click="updateApp">保存名称</a-button>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>危险操作</span>
        </div>
        <div class="section-body">
          <div class="input-item">
            <a-button type="primary" danger @click="deleteApp">删除应用</a-button>
            <div class="input-desc">删除应用后, 将删除跟应用相关的配置、部署流水线、相关记录等信息, 不可逆</div>
          </div>
          <div class="input-item">
            <a-button type="primary" danger @click="deleteRepo">迁移至其他团队</a-button>
            <div class="input-desc">将应用迁移至其他团队, 该团队成员无法再看到此应用</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, createVNode } from "vue";
import {
  getAppRequest,
  updateAppRequest,
  deleteAppRequest
} from "@/api/app/appApi";
import { useRoute, useRouter } from "vue-router";
import { appNameRegexp } from "@/utils/regexp";
import { message, Modal } from "ant-design-vue";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
const appName = ref("");
const route = useRoute();
const router = useRouter();
const getApp = () => {
  getAppRequest(route.params.appId).then(res => {
    appName.value = res.data.name;
  });
};

const updateApp = () => {
  if (!appNameRegexp.test(appName.value)) {
    message.warn("名称格式错误");
    return;
  }
  updateAppRequest({
    appId: route.params.appId,
    name: appName.value
  }).then(() => {
    message.success("编辑成功");
  });
};

const deleteApp = () => {
  Modal.confirm({
    title: `你确定要删除该应用吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteAppRequest(route.params.appId).then(() => {
        message.success("删除成功");
        router.push(`/team/${route.params.teamId}/app/list`);
      });
    },
    onCancel() {}
  });
};

getApp();
</script>
<style>
</style>