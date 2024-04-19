<template>
  <div style="padding:20px 0">
    <div class="container">
      <div class="form">
        <div class="title">创建配置</div>
        <div class="form-item">
          <div class="label">配置名称</div>
          <a-input type="input" placeholder="请输入" v-model:value="formState.name" />
        </div>
        <div class="form-item">
          <div class="label">格式</div>
          <div>
            <a-radio-group v-model:value="formState.format" @change="onFormatChange">
              <a-radio value="json">json</a-radio>
              <a-radio value="yaml">yaml</a-radio>
              <a-radio value="xml">xml</a-radio>
              <a-radio value="properties">properties</a-radio>
              <a-radio value="text">text</a-radio>
            </a-radio-group>
          </div>
        </div>
        <div class="form-item">
          <div class="label">配置内容</div>
          <Codemirror
            v-model="formState.yamlContent"
            style="height:380px;width:100%"
            :extensions="extensions"
          />
        </div>
        <div class="form-item">
          <a-button type="primary">创建</a-button>
        </div>
      </div>
    </div>
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
const extensions = ref([json(), oneDark]);
const formState = reactive({
  name: "",
  format: "json",
  yamlContent: ""
});
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
</script>
<style scoped>
</style>