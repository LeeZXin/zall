<template>
  <div style="padding:14px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建配置文件</span>
        <span v-else-if="mode === 'new'">新增版本</span>
      </div>
      <div class="section" v-if="mode==='create'">
        <div class="section-title">选择环境</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            placeholder="选择环境"
            v-model:value="formState.selectedEnv"
            :options="envList"
          />
          <div class="input-desc">多环境选择, 选择其中一个环境</div>
        </div>
      </div>
      <div class="section" v-if="mode==='new'">
        <div class="section-title">跟随版本号</div>
        <div class="section-body">{{route.query.from}}</div>
      </div>
      <div class="section" v-if="mode==='new'">
        <div class="section-title">已选环境</div>
        <div class="section-body">{{formState.selectedEnv}}</div>
      </div>
      <div class="section" v-if="mode === 'create'">
        <div class="section-title">文件名称</div>
        <div class="section-body">
          <a-input v-model:value="formState.name" />
          <div
            class="input-desc"
          >配置名称,文件的唯一标识,不能包含特殊字符,最后文件会拼接下面格式作为后缀,例如 test-xxx 选择格式为json, 则最后保存的名称为test-xxx.json</div>
        </div>
      </div>
      <div class="section" v-if="mode === 'new'">
        <div class="section-title">文件名称</div>
        <div class="section-body">{{formState.name}}</div>
      </div>
      <div class="section" v-if="mode === 'create'">
        <div class="section-title">格式</div>
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
          <span>配置内容</span>
          <span v-if="mode==='new'" class="diff-btn" @click="showDiffModal">对比</span>
        </div>
        <Codemirror
          v-model="formState.content"
          style="height:380px;width:100%"
          :extensions="extensions"
        />
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="saveOrUpdateFile">立即保存</a-button>
      </div>
    </div>
    <a-modal title="新旧对比" :footer="null" v-model:open="diffModalOpen" :width="800">
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
const diffModalOpen = ref(false);
const propertyFileStore = usePropertyFileStore();
const route = useRoute();
const router = useRouter();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const extensions = ref([json(), oneDark]);
const formState = reactive({
  name: "",
  format: "json",
  content: "",
  selectedEnv: "",
  oldContent: ""
});
const envList = ref([]);
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
const showDiffModal = () => {
  diffModalOpen.value = true;
};
const propertiesLang = StreamLanguage.define(properties);
const onFormatChange = event => {
  switch (event.target.value) {
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

const saveOrUpdateFile = () => {
  if (mode === "create") {
    if (!propertyFileNameRegexp.test(formState.name)) {
      message.warn("文件名称格式错误");
      return;
    }
    createPropertyFileRequest(
      {
        env: formState.selectedEnv,
        appId: route.params.appId,
        content: formState.content,
        name: formState.name + "." + formState.format
      },
      formState.selectedEnv
    ).then(() => {
      message.success("创建成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/property/list/${formState.selectedEnv}`
      );
    });
  } else if (mode === "new") {
    newVersionRequest(
      {
        fileId: propertyFileStore.id,
        content: formState.content,
        lastVersion: route.query.from
      },
      formState.selectedEnv
    ).then(() => {
      message.success("保存成功");
      router.push(
        `/team/${route.params.teamId}/app/${route.params.appId}/property/${route.params.fileId}/history/list`
      );
    });
  }
};

if (mode === "create") {
  getEnvCfg();
} else if (mode === "new") {
  if (propertyFileStore.id === 0) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/property/list`
    );
  } else if (!route.query.from) {
    router.push(
      `/team/${route.params.teamId}/app/${route.params.appId}/property/${route.params.fileId}/history/list`
    );
  } else {
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
          `/team/${route.params.teamId}/app/${route.params.appId}/property/${route.params.fileId}/history/list`
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