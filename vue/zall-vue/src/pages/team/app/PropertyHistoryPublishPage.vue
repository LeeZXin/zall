<template>
  <div style="padding:14px">
    <div class="container">
      <div class="title">
        <span>发布配置</span>
      </div>
      <div class="section">
        <div class="section-title">跟随版本号</div>
        <div class="section-body">{{formState.lastVersion}}</div>
      </div>
      <div class="section">
        <div class="section-title">当前版本号</div>
        <div class="section-body">{{formState.version}}</div>
      </div>
      <div class="section">
        <div class="section-title">已选环境</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section">
        <div class="section-title">文件名称</div>
        <div class="section-body">{{formState.name}}</div>
      </div>
      <div class="section">
        <div class="section-title">发布节点</div>
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
          <span>配置内容</span>
          <span class="diff-btn" @click="showDiffModal">对比</span>
        </div>
        <Codemirror
          v-model="formState.content"
          style="height:380px;width:100%"
          :extensions="extensions"
          :disabled="true"
        />
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="publishFile">立即发布</a-button>
      </div>
    </div>
    <a-modal title="新旧对比" :footer="null" v-model:open="diffModalOpen" :width="800">
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
import { StreamLanguage } from "@codemirror/language";
import { properties } from "@codemirror/legacy-modes/mode/properties";
import { message } from "ant-design-vue";
import {
  getHistoryByVersionRequest,
  listPropertySourceByFileIdRequest,
  deployHistoryRequest
} from "@/api/app/propertyApi";
import { CodeDiff } from "v-code-diff";
import { useRoute, useRouter } from "vue-router";
import { usePropertyHistoryStore } from "@/pinia/propertyHistoryStore";
const propertyHistoryStore = usePropertyHistoryStore();
const diffModalOpen = ref(false);
const route = useRoute();
const router = useRouter();
const propertiesLang = StreamLanguage.define(properties);
const extensions = ref([propertiesLang, oneDark]);
const formState = reactive({
  name: "",
  content: "",
  selectedEnv: "",
  oldContent: "",
  lastVersion: "",
  version: "",
  nodeList: []
});
const nodeList = ref([]);
const showDiffModal = () => {
  diffModalOpen.value = true;
};

const publishFile = () => {
  if (formState.nodeList.length === 0) {
    message.warn("请选择节点");
    return;
  }
  deployHistoryRequest(
    {
      historyId: propertyHistoryStore.id,
      sourceIdList: formState.nodeList
    },
    propertyHistoryStore.env
  ).then(() => {
    message.success("发布成功");
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
    getHistoryByVersionRequest(
      {
        fileId: propertyHistoryStore.fileId,
        version: formState.lastVersion
      },
      propertyHistoryStore.env
    ).then(res => {
      formState.oldContent = res.data.value.content;
    });
  }
  listPropertySourceByFileIdRequest(
    propertyHistoryStore.fileId,
    propertyHistoryStore.env
  ).then(res => {
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
</style>