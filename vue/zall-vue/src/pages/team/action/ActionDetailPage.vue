<template>
  <div style="padding:14px">
    <div style="margin-bottom:20px">
      <span class="header" @click="backToTaskList">
        <arrow-left-outlined />
        <span style="margin-left:8px">任务列表</span>
      </span>
    </div>
    <div class="container">
      <div class="body">
        <div class="title">工作流详情</div>
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
          <div class="form-item-text">{{teamName}}</div>
        </div>
        <div class="form-item">
          <div class="label">
            <span>代理地址</span>
          </div>
          <div class="form-item-text">{{teamName}}</div>
        </div>
        <div class="form-item">
          <div class="label">
            <span>代理token</span>
          </div>
          <div class="form-item-text">{{teamName}}</div>
        </div>
        <div class="form-item">
          <div class="label">
            <span>yaml配置</span>
          </div>
          <div>
            <Codemirror
              v-model="formState.yamlContent"
              :style="codemirroStyle"
              :extensions="extensions"
              :disabled="true"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive } from "vue";
import { ArrowLeftOutlined } from "@ant-design/icons-vue";
import { useTeamStore } from "@/pinia/TeamStore";
import { useRouter } from "vue-router";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
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
const backToTaskList = () => {
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