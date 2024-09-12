<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('gitWorkflow.createWorkflow')}}</span>
        <span v-else-if="mode === 'update'">{{t('gitWorkflow.updateWorkflow')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWorkflow.name')}}</div>
        <div class="section-body">
          <a-input type="input" v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWorkflow.desc')}}</div>
        <div class="section-body">
          <a-textarea :auto-size="{ minRows: 3, maxRows: 5 }" v-model:value="formState.desc" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWorkflow.zalletNode')}}</div>
        <div class="section-body">
          <a-select
            v-model:value="formState.agentId"
            style="width:100%;"
            :options="zalletNodeList"
            show-search
            :filter-option="filterZalletNodeListOption"
          />
          <div class="input-desc">{{t('gitWorkflow.zalletNodeDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('gitWorkflow.triggerType')}}</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.source">
            <a-radio :style="radioStyle" :value="1">{{t('gitWorkflow.branchTrigger')}}</a-radio>
            <div class="radio-option-desc">{{t('gitWorkflow.branchTriggerDesc')}}</div>
            <a-radio :style="radioStyle" :value="2">{{t('gitWorkflow.pullRequestTrigger')}}</a-radio>
            <div class="radio-option-desc">{{t('gitWorkflow.pullRequestTriggerDesc')}}</div>
          </a-radio-group>
          <a-input type="input" v-model:value="formState.wildBranches" />
          <div class="input-desc">{{t('gitWorkflow.triggerBranchDesc')}}</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title flex-between">
          <span>{{t('gitWorkflow.workflowCfg')}}</span>
          <span @click="formatYaml" class="format-yaml-btn">{{t('gitWorkflow.formatYaml')}}</span>
        </div>
        <Codemirror
          v-model="formState.yamlContent"
          :style="codemirrorStyle"
          :extensions="extensions"
        />
        <div class="section-body">
          <div class="input-desc" style="margin:0">
            <span>{{t('gitWorkflow.yamlDesc')}}</span>
            <span class="insert-template" @click="insertTemplate">{{t('gitWorkflow.insertTpl')}}</span>
          </div>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateWorkflow">{{t('gitWorkflow.save')}}</a-button>
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
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
// 模式
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
// 表单
const formState = reactive({
  id: 0,
  name: "",
  yamlContent: "",
  wildBranches: "",
  desc: "",
  source: 1,
  agentId: undefined
});
// 格式化yaml
const formatYaml = () => {
  if (formState.yamlContent) {
    const parsedYaml = jsyaml.load(formState.yamlContent);
    formState.yamlContent = jsyaml.dump(parsedYaml);
  }
};
// 插入模板
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
      title: `${t("gitWorkflow.overrideWarn")}?`,
      icon: createVNode(ExclamationCircleOutlined),
      onOk() {
        formState.yamlContent = content;
      }
    });
  } else {
    formState.yamlContent = content;
  }
};
// 新增或编辑工作流
const createOrUpdateWorkflow = () => {
  if (!workflowNameRegexp.test(formState.name)) {
    message.warn(t("gitWorkflow.nameFormatErr"));
    return;
  }
  if (!workflowDescRegexp.test(formState.desc)) {
    message.warn(t("gitWorkflow.descFormatErr"));
    return;
  }
  if (!formState.wildBranches) {
    message.warn(t("gitWorkflow.pleaseFillWildBranches"));
    return;
  }
  if (!formState.agentId) {
    message.warn(t("gitWorkflow.pleaseSelectAgent"));
    return;
  }
  let branches = formState.wildBranches.split(";");
  for (let index in branches) {
    let branch = branches[index];
    if (!workflowWildBranchRegexp.test(branch)) {
      message.warn(t("gitWorkflow.wildBranchesFormatErr"));
      return;
    }
  }
  if (!formState.yamlContent) {
    message.warn(t("gitWorkflow.pleaseFillYamlContent"));
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
      message.success(t("operationSuccess"));
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
      message.success(t("operationSuccess"));
      router.push(
        `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/workflow/list`
      );
    });
  }
};
// zallet下拉框过滤
const filterZalletNodeListOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 获取zallet节点
const listAllZalletNode = () => {
  listAllZalletNodeRequest().then(res => {
    zalletNodeList.value = res.data.map(item => {
      return {
        value: item.nodeId,
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
  font-weight: normal;
}
.format-yaml-btn:hover {
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