<template>
  <div style="padding:14px">
    <div style="margin-bottom:20px">
      <span class="header" @click="backToAppList">
        <arrow-left-outlined />
        <span style="margin-left:8px">工作流列表</span>
      </span>
    </div>
    <div class="container">
      <div class="title">创建工作流</div>
      <div class="form-item">
        <div class="label">
          <span>所属团队</span>
        </div>
        <div class="form-item-text">{{teamName}}</div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>工作流名称</span>
        </div>
        <div>
          <a-input type="input" v-model:value="formState.actionName" />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>代理地址</span>
        </div>
        <div>
          <a-input type="input" v-model:value="formState.agentHost" />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>代理token</span>
        </div>
        <div>
          <a-input type="input" v-model:value="formState.agentToken" />
        </div>
      </div>
      <div class="form-item">
        <div class="label">
          <span>yaml配置</span>
          <span class="insert-template" @click="insertTemplate">插入模版</span>
          <span @click="formatYaml" class="format-yaml-text">格式化yaml</span>
        </div>
        <div>
          <Codemirror
            v-model="formState.yamlContent"
            :style="codemirroStyle"
            :extensions="extensions"
          />
        </div>
      </div>
      <div class="form-item">
        <a-button type="primary" style="margin-top:20px;">立即创建</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, createVNode } from "vue";
import {
  ArrowLeftOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import { useTeamStore } from "@/pinia/teamStore";
import { useRouter } from "vue-router";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
import jsyaml from "js-yaml";
import { Modal } from "ant-design-vue";
const extensions = [yaml(), oneDark];
const codemirroStyle = { height: "380px", width: "100%" };
const team = useTeamStore();
const teamName = team.teamName;
const router = useRouter();
const formState = reactive({
  actionName: "",
  agentHost: "",
  agentToken: "",
  yamlContent: ""
});
const formatYaml = () => {
  if (formState.yamlContent) {
    const parsedYaml = jsyaml.load(formState.yamlContent);
    formState.yamlContent = jsyaml.dump(parsedYaml);
  }
};
const backToAppList = () => {
  router.push("/team/action/list");
};
const insertTemplate = () => {
  let content = `jobs:
  job1:
    needs: []
    steps:
      - name: job1 step 1
        with:
          a: '1'
        script: |
          echo $a
  job2:
    needs: []
    steps:
      - name: job2 step 1
        with:
          b: '2'
        script: |
          echo $b
`;
  if (formState.yamlContent) {
    Modal.confirm({
      title: "确认",
      icon: createVNode(ExclamationCircleOutlined),
      content: "插入模版会覆盖已写入的内容, 确定要插入吗?",
      okText: "确定",
      cancelText: "取消",
      onOk() {
        formState.yamlContent = content;
      }
    });
  } else {
    formState.yamlContent = content;
  }
};
</script>
<style scoped>
.header {
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 10px;
}
.header:hover {
  color: #1677ff;
}
.format-yaml-text {
  float: right;
  cursor: pointer;
}
.format-yaml-text:hover {
  color: #1677ff;
}
.insert-template {
  margin-left: 12px;
  cursor: pointer;
}
.insert-template:hover {
  color: #1677ff;
}
</style>