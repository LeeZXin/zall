<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('propertyFile.createFile')}}</span>
        <span v-else-if="mode === 'new'">{{t('propertyFile.newVersion')}}</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">{{t('propertyFile.selectEnv')}}</div>
        <div class="section-body">
          <a-select style="width: 100%" v-model:value="formState.selectedEnv" :options="envList" />
        </div>
      </div>
      <div class="section" v-if="mode==='new'">
        <div class="section-title">{{t('propertyFile.lastVersion')}}</div>
        <div class="section-body">{{route.query.from}}</div>
      </div>
      <div class="section" v-if="mode==='new'">
        <div class="section-title">{{t('propertyFile.selectedEnv')}}</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section" v-if="mode === 'create'">
        <div class="section-title">{{t('propertyFile.name')}}</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div class="input-desc">{{t('propertyFile.nameDesc')}}</div>
        </div>
      </div>
      <div class="section" v-if="mode === 'new'">
        <div class="section-title">{{t('propertyFile.name')}}</div>
        <div class="section-body">{{formState.name}}</div>
      </div>
      <div class="section" v-if="mode === 'create'">
        <div class="section-title">{{t('propertyFile.format')}}</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.format" @change="onFormatChange">
            <a-radio value="json">json</a-radio>
            <a-radio value="yaml">yaml</a-radio>
            <a-radio value="xml">xml</a-radio>
            <a-radio value="properties">properties</a-radio>
            <a-radio value="text">text</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>{{t('propertyFile.content')}}</span>
          <span
            v-if="mode==='new'"
            class="diff-btn"
            @click="showDiffModal"
          >{{t('propertyFile.compare')}}</span>
        </div>
        <Codemirror
          v-model="formState.content"
          style="height:380px;width:100%"
          :extensions="extensions"
        />
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateFile">{{t('propertyFile.save')}}</a-button>
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
        filename="old"
        newFilename="new"
        style="max-height:400px;overflow:scroll"
      />
    </a-modal>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import { Codemirror } from "vue-codemirror";
import { xml } from "@codemirror/lang-xml";
import { json } from "@codemirror/lang-json";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
import { StreamLanguage } from "@codemirror/language";
import { properties } from "@codemirror/legacy-modes/mode/properties";
import { propertyFileNameRegexp } from "@/utils/regexp";
import { message } from "ant-design-vue";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  createPropertyFileRequest,
  getHistoryByVersionRequest,
  newVersionRequest
} from "@/api/app/propertyApi";
import { CodeDiff } from "v-code-diff";
import { useRoute, useRouter } from "vue-router";
import { usePropertyFileStore } from "@/pinia/propertyFileStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const diffModalOpen = ref(false);
const propertyFileStore = usePropertyFileStore();
const route = useRoute();
const router = useRouter();
// 模式
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const extensions = ref([json(), oneDark]);
// 表单数据
const formState = reactive({
  name: "",
  format: "json",
  content: "",
  selectedEnv: "",
  oldContent: ""
});
// 环境列表
const envList = ref([]);
// 获取环境列表
const getEnvCfg = () => {
  getEnvCfgRequest().then(res => {
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (route.query.env && res.data?.includes(route.query.env)) {
      formState.selectedEnv = route.query.env;
    } else if (res.data.length > 0) {
      formState.selectedEnv = res.data[0];
    }
  });
};
// 展示差异modal
const showDiffModal = () => {
  diffModalOpen.value = true;
};
// codemirror对不同格式文件支持
const propertiesLang = StreamLanguage.define(properties);
const onFormatChange = event => {
  formatChange(event.target.value);
};
const formatChange = ext => {
  switch (ext) {
    case "json":
      extensions.value = [json(), oneDark];
      break;
    case "yaml":
      extensions.value = [yaml(), oneDark];
      break;
    case "xml":
      extensions.value = [xml(), oneDark];
      break;
    case "properties":
      extensions.value = [propertiesLang, oneDark];
      break;
    default:
      extensions.value = [oneDark];
      break;
  }
};
// 新增或编辑配置文件
const saveOrUpdateFile = () => {
  if (mode === "create") {
    if (!propertyFileNameRegexp.test(formState.name)) {
      message.warn(t("propertyFile.nameFormatErr"));
      return;
    }
    createPropertyFileRequest({
      env: formState.selectedEnv,
      appId: route.params.appId,
      content: formState.content,
      name: formState.name + "." + formState.format
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "new") {
    newVersionRequest({
      fileId: propertyFileStore.id,
      content: formState.content,
      lastVersion: route.query.from
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${route.params.fileId}/history/list`
      );
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "new") {
  if (propertyFileStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/list`
    );
  } else if (!route.query.from) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${route.params.fileId}/history/list`
    );
  } else {
    let dot = propertyFileStore.name.lastIndexOf(".");
    formatChange(propertyFileStore.name.substr(dot + 1));
    getHistoryByVersionRequest({
      fileId: propertyFileStore.id,
      version: route.query.from
    }).then(res => {
      if (res.data.exist) {
        formState.name = propertyFileStore.name;
        formState.selectedEnv = propertyFileStore.env;
        formState.content = res.data.value.content;
        formState.oldContent = res.data.value.content;
      } else {
        message.warn("跟随版本数据不存在");
        router.push(
          `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${route.params.fileId}/history/list`
        );
      }
    });
  }
}
</script>
<style scoped>
.diff-btn:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>