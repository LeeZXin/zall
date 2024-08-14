<template>
  <div style="padding:10px">
    <div class="container">
      <div class="title">
        <span v-if="mode === 'create'">创建工作流</span>
        <span v-if="mode === 'update'">编辑工作流</span>
      </div>
      <div class="section">
        <div class="section-title">工作流名称</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">工作流描述</div>
        <div class="section-body">
          <a-textarea :auto-size="{ minRows: 3, maxRows: 5 }" v-model:value="formState.desc" />
          <div class="input-desc">简单的话来描述工作流的作用</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">Zallet代理</div>
        <div class="section-body">
          <a-select
            v-model:value="formState.agentId"
            style="width:100%;"
            :options="zalletNodeList"
            show-search
            :filter-option="filterZalletNodeListOption"
          />
          <div class="input-desc">选择一个zallet代理来执行</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">触发方式</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.source">
            <a-radio :style="radioStyle" :value="1">分支</a-radio>
            <div class="radio-option-desc">当分支push操作时, 将触发工作流</div>
            <a-radio :style="radioStyle" :value="2">合并请求</a-radio>
            <div class="radio-option-desc">当提交分支的合并请求时, 将触发工作流</div>
          </a-radio-group>
          <a-input type="input" v-model:value="formState.wildBranches" />
          <div class="input-desc">以glob方式保存, 以分号隔开, 例如dev_*</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>工作流配置</span>
          <span @click="formatYaml" class="format-yaml-btn">格式化yaml</span>
        </div>
        <Codemirror
          v-model="formState.yamlContent"
          :style="codemirrorStyle"
          :extensions="extensions"
        />
        <div class="section-body">
          <div class="input-desc" style="margin:0">
            <span>以yaml配置工作流脚本配置, 模板可查看</span>
            <span class="insert-template" @click="insertTemplate">插入模版</span>
          </div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateWorkflow">立即保存</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, createVNode, ref } from "vue";
import { ExclamationCircleOutlined } from "@ant-design/icons-vue";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
import jsyaml from "js-yaml";
import { Modal, message } from "ant-design-vue";
import {
  createWorkflowRequest,
  getWorkflowDetailRequest,
  updateWorkflowRequest
} from "@/api/git/workflowApi";
import { listAllZalletNodeRequest } from "@/api/zallet/zalletApi";
import { useRouter, useRoute } from "vue-router";
import {
  workflowNameRegexp,
  workflowWildBranchRegexp,
  workflowDescRegexp
} from "@/utils/regexp";
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const mode = getMode();
const router = useRouter();
const extensions = [yaml(), oneDark];
const codemirrorStyle = { height: "380px", width: "100%" };
const radioStyle = reactive({
  display: "flex",
  alignItems: "flex-start"
});
const zalletNodeList = ref([]);
const formState = reactive({
  id: 0,
  name: "",
  yamlContent: "",
  wildBranches: "",
  desc: "",
  source: 1,
  agentId: undefined
});
const formatYaml = () => {
  if (formState.yamlContent) {
    const parsedYaml = jsyaml.load(formState.yamlContent);
    formState.yamlContent = jsyaml.dump(parsedYaml);
  }
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
const createOrUpdateWorkflow = () => {
  if (!workflowNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (!workflowDescRegexp.test(formState.desc)) {
    message.warn("描述格式错误");
    return;
  }
  if (!formState.wildBranches) {
    message.warn("分支不能为空");
    return;
  }
  if (!formState.agentId) {
    message.warn("请选择代理");
    return;
  }
  let branches = formState.wildBranches.split(";");
  for (let index in branches) {
    let branch = branches[index];
    if (!workflowWildBranchRegexp.test(branch)) {
      message.warn("分支格式错误");
      return;
    }
  }
  if (!formState.yamlContent) {
    message.warn("yaml配置错误");
    return;
  }
  let source = {
    sourceType: formState.source
  };
  switch (formState.source) {
    case 1:
      source.branchSource = branches;
      break;
    case 2:
      source.pullRequestSource = {
        branches
      };
      break;
  }
  if (mode === "update") {
    updateWorkflowRequest({
      name: formState.name,
      workflowId: formState.id,
      source: source,
      agentId: formState.agentId,
      yamlContent: formState.yamlContent,
      desc: formState.desc
    }).then(() => {
      message.success("编辑成功");
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/list`
      );
    });
  } else if (mode === "create") {
    createWorkflowRequest({
      name: formState.name,
      repoId: parseInt(route.params.repoId),
      source: source,
      agentId: formState.agentId,
      yamlContent: formState.yamlContent,
      desc: formState.desc
    }).then(() => {
      message.success("创建成功");
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/list`
      );
    });
  }
};
const filterZalletNodeListOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
const listAllZalletNode = () => {
  listAllZalletNodeRequest().then(res => {
    zalletNodeList.value = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
  });
};
listAllZalletNode();
if (mode === "update") {
  getWorkflowDetailRequest(route.params.workflowId).then(res => {
    let wf = res.data;
    formState.id = wf.id;
    formState.name = wf.name;
    formState.agentId = wf.agentId;
    formState.yamlContent = wf.yamlContent;
    if (wf.source.sourceType === 1) {
      formState.wildBranches = wf.source.branchSource.join(";");
    } else if (wf.source.sourceType === 2) {
      formState.wildBranches = wf.source.pullRequestSource.branches.join(";");
    }
    formState.desc = wf.desc;
    formState.source = wf.source.sourceType;
  });
}
</script>
<style scoped>
.format-yaml-btn {
  cursor: pointer;
}
.format-yaml-text:hover {
  color: #1677ff;
}
.insert-template {
  cursor: pointer;
  padding: 0 8px;
  color: black;
}
.insert-template:hover {
  color: #1677ff;
}
</style>