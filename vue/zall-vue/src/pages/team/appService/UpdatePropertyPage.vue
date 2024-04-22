<template>
  <div style="padding:14px">
    <ZNaviBack url="/appService/property/list" name="配置列表" />
    <div class="container">
      <div class="form">
        <div class="title">编辑配置</div>
        <div class="form-item">
          <div class="label">配置名称</div>
          <div class="form-item-text">sentinel-flow</div>
        </div>
        <div class="form-item">
          <div class="label">格式</div>
          <div class="form-item-text">{{format}}</div>
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
          <a-button type="primary">保存</a-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import ZNaviBack from "@/components/common/ZNaviBack";
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
  yamlContent: ""
});
const format = ref("json");
const propertiesLang = StreamLanguage.define(properties);
const changeExtension = format => {
  switch (format) {
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
changeExtension(format.value);
</script>
<style scoped>
</style>