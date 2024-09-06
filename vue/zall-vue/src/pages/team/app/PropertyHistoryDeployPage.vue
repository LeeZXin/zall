<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span>{{t('propertyFile.deployVersion')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertyFile.lastVersion')}}</div>
        <div class="section-body">{{formState.lastVersion}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertyFile.currentVersion')}}</div>
        <div class="section-body">{{formState.version}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertyFile.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertyFile.name')}}</div>
        <div class="section-body">{{formState.name}}</div>
      </div>
      <div class="section">
        <div class="section-title">{{t('propertyFile.deployNode')}}</div>
        <div class="section-body">
          <a-checkbox-group v-model:value="formState.nodeList" style="width: 100%">
            <ul class="node-ul">
              <li v-for="(item, index) in nodeList" v-bind:key="index">
                <a-checkbox :value="item.id">{{item.name}}</a-checkbox>
              </li>
            </ul>
          </a-checkbox-group>
        </div>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>{{t('propertyFile.content')}}</span>
          <span class="diff-btn" @click="showDiffModal">{{t('propertyFile.compare')}}</span>
        </div>
        <Codemirror
          v-model="formState.content"
          style="height:380px;width:100%"
          :extensions="extensions"
          :disabled="true"
        />
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="deployFile">{{t('propertyFile.deploy')}}</a-button>
      </div>
    </div>
    <a-modal
      :title="t('propertyFile.compare')"
      :footer="null"
      v-model:open="diffModalOpen"
      :width="800"
    >
      <code-diff
        :old-string="formState.oldContent"
        :new-string="formState.content"
        :context="10"
        outputFormat="side-by-side"
        :hideStat="true"
        :filename="formState.lastVersion"
        :newFilename="formState.version"
        style="max-height:400px;overflow:scroll"
      />
    </a-modal>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { Codemirror } from "vue-codemirror";
import { oneDark } from "@codemirror/theme-one-dark";
import { message } from "ant-design-vue";
import {
  getHistoryByVersionRequest,
  listPropertySourceByFileIdRequest,
  deployHistoryRequest
} from "@/api/app/propertyApi";
import { CodeDiff } from "v-code-diff";
import { useRoute, useRouter } from "vue-router";
import { usePropertyHistoryStore } from "@/pinia/propertyHistoryStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
// 透传的配置版本历史
const propertyHistoryStore = usePropertyHistoryStore();
// 对比modal是否展开
const diffModalOpen = ref(false);
const route = useRoute();
const router = useRouter();
const extensions = ref([oneDark]);
// 表单数据
const formState = reactive({
  name: "",
  content: "",
  selectedEnv: "",
  oldContent: "",
  lastVersion: "",
  version: "",
  nodeList: []
});
// 发布节点列表
const nodeList = ref([]);
// 展开对比modal
const showDiffModal = () => {
  diffModalOpen.value = true;
};
// 发布配置
const deployFile = () => {
  if (formState.nodeList.length === 0) {
    message.warn(t('propertyFile.pleaseSelectNodes'));
    return;
  }
  deployHistoryRequest({
    historyId: propertyHistoryStore.id,
    sourceIdList: formState.nodeList
  }).then(() => {
    message.success(t('operationSuccess'));
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${route.params.fileId}/history/list`
    );
  });
};
if (propertyHistoryStore.id === 0) {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/list`
  );
} else {
  formState.name = propertyHistoryStore.fileName;
  formState.selectedEnv = propertyHistoryStore.env;
  formState.content = propertyHistoryStore.content;
  formState.lastVersion = propertyHistoryStore.lastVersion;
  formState.version = propertyHistoryStore.version;
  if (formState.lastVersion) {
    getHistoryByVersionRequest({
      fileId: propertyHistoryStore.fileId,
      version: formState.lastVersion
    }).then(res => {
      formState.oldContent = res.data.value.content;
    });
  }
  listPropertySourceByFileIdRequest(propertyHistoryStore.fileId).then(res => {
    nodeList.value = res.data;
  });
}
</script>
<style scoped>
.diff-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
.node-ul > li + li {
  margin-top: 10px;
  font-size: 14px;
}
.auditor-ul > li {
  width: 100%;
  display: flex;
  align-items: center;
  word-break: break-all;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.auditor-ul > li + li {
  margin-top: 8px;
}
</style>