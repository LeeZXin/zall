<template>
  <div style="padding:14px">
    <div style="margin-bottom:20px">
      <span class="header" @click="backToAppList">
        <arrow-left-outlined />
        <span style="margin-left:8px">工作流列表</span>
      </span>
    </div>
    <div class="container">
      <div class="body">
        <div class="title">编辑工作流</div>
        <div style="display:flex">
          <div class="form-item">
            <div class="label">
              <span>Id</span>
            </div>
            <div class="form-item-text">{{teamName}}</div>
          </div>
          <div class="form-item">
            <div class="label">
              <span>所属团队</span>
            </div>
            <div class="form-item-text">{{teamName}}</div>
          </div>
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
          <a-button type="primary" style="margin-top:20px;">编辑</a-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import { ArrowLeftOutlined } from "@ant-design/icons-vue";
import { useTeamStore } from "@/pinia/teamStore";
import { useRouter } from "vue-router";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
import jsyaml from "js-yaml";
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
.body > .title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
}
.format-yaml-text {
  float: right;
  cursor: pointer;
}
.format-yaml-text:hover {
  color: #1677ff;
}
</style>