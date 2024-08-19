<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建通知模板</span>
        <span v-else-if="mode === 'update'">编辑通知模板</span>
      </div>
      <div class="section">
        <div class="section-title">名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">通知模板名称, 长度为1-32</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">Webhook Url</div>
        <div class="section-body">
          <a-input v-model:value="formState.url" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">类型</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.notifyType">
            <a-radio value="wework">企业微信</a-radio>
            <a-radio value="feishu">飞书</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.notifyType === 'feishu'">
        <div class="section-title">飞书签名密钥</div>
        <div class="section-body">
          <a-input v-model:value="formState.feishuSignKey" />
          <div class="input-desc">非必填, 飞书签名验证专用</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">模板</div>
        <Codemirror
          v-model="formState.template"
          style="height:380px;width:100%"
          :extensions="extensions"
        />
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateTpl">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { Codemirror } from "vue-codemirror";
import { oneDark } from "@codemirror/theme-one-dark";
import { reactive, ref } from "vue";
import { notifyTplNameRegexp, notifyTplUrlRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import {
  createNotifyTplRequest,
  updateNotifyTplRequest
} from "@/api/notify/notifyApi";
import { useRoute, useRouter } from "vue-router";
import { useNotifyTplStore } from "@/pinia/notifyTplStore";
import { json } from "@codemirror/lang-json";
// code mirror扩展项
const extensions = ref([json(), oneDark]);
const notifyTplStore = useNotifyTplStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
// 表单数据
const formState = reactive({
  name: "",
  template: "",
  url: "",
  notifyType: "wework",
  feishuSignKey: ""
});
// 点击“立即保存”
const saveOrUpdateTpl = () => {
  if (!notifyTplNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!notifyTplUrlRegexp.test(formState.url)) {
    message.warn("webhook url格式错误");
    return;
  }
  if (!formState.template) {
    message.warn("模板为空");
    return;
  }
  if (formState.template !== "feishu") {
    formState.feishuSignKey = "";
  }
  if (mode === "create") {
    // 创建模板
    createNotifyTplRequest({
      name: formState.name,
      teamId: parseInt(route.params.teamId),
      cfg: {
        notifyType: formState.notifyType,
        url: formState.url,
        template: formState.template,
        feishuSignKey: formState.feishuSignKey
      }
    }).then(() => {
      message.success("创建成功");
      router.push(`/team/${route.params.teamId}/notifyTpl/list`);
    });
  } else if (mode === "update") {
    // 编辑模板
    updateNotifyTplRequest({
      id: notifyTplStore.id,
      name: formState.name,
      cfg: {
        notifyType: formState.notifyType,
        url: formState.url,
        template: formState.template,
        feishuSignKey: formState.feishuSignKey
      }
    }).then(() => {
      message.success("保存成功");
      router.push(`/team/${route.params.teamId}/notifyTpl/list`);
    });
  }
};

if (mode === "update") {
  // store没有跳转list
  if (notifyTplStore.id === 0) {
    router.push(`/team/${route.params.teamId}/notifyTpl/list`);
  } else {
    formState.name = notifyTplStore.name;
    formState.notifyType = notifyTplStore.notifyType;
    formState.url = notifyTplStore.url;
    formState.feishuSignKey = notifyTplStore.feishuSignKey;
    formState.template = notifyTplStore.template;
  }
}
</script>
<style scoped>
</style>