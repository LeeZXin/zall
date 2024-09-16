<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">{{t('teamHook.createHook')}}</span>
        <span v-else-if="mode === 'update'">{{t('teamHook.updateHook')}}</span>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamHook.name')}}</div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamHook.hookType')}}</div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.hookType">
            <a-radio :value="1">{{t('teamHook.webhook')}}</a-radio>
            <a-radio :value="2">{{t('teamHook.notification')}}</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.hookType === 1">
        <div class="section-title">{{t('teamHook.webhook')}}</div>
        <div class="section-body">
          <div>
            <div style="font-size: 12px;margin-bottom: 6px">{{t('teamHook.webhookUrl')}}</div>
            <a-input style="width:100%" v-model:value="formState.hookUrl" />
          </div>
          <div style="margin-top: 10px">
            <div style="font-size: 12px;margin-bottom: 6px">{{t('teamHook.webhookSecret')}}</div>
            <a-input-password style="width:100%" v-model:value="formState.secret" />
          </div>
        </div>
      </div>
      <div class="section" v-else-if="formState.hookType === 2">
        <div class="section-title">{{t('teamHook.notification')}}</div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.tplId"
            :options="tplList"
            show-search
            :filter-option="filterTplListOption"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamHook.gitEvents')}}</div>
        <div class="section-body">
          <ul class="event-list">
            <li v-for="(item, index) in gitCheckboxes" v-bind:key="index">
              <div style="font-size:12px;margin-bottom:8px">{{t(item.title)}}</div>
              <ul class="action-list">
                <li v-for="action in item.actions" v-bind:key="`${index}-${action.key}`">
                  <a-checkbox v-model:checked="action.value">{{t(action.title)}}</a-checkbox>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="section">
        <div class="section-title">{{t('teamHook.teamEvents')}}</div>
        <div class="section-body">
          <ul class="event-list">
            <li v-for="(item, index) in teamCheckboxes" v-bind:key="index">
              <div style="font-size:12px;margin-bottom:8px">{{t(item.title)}}</div>
              <ul class="action-list">
                <li v-for="action in item.actions" v-bind:key="`${index}-${action.key}`">
                  <a-checkbox v-model:checked="action.value">{{t(action.title)}}</a-checkbox>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="envList.length > 0">
        <div class="section-title">{{t('teamHook.envRelatedEvents')}}</div>
        <div class="section-body">
          <div style="margin-bottom: 18px">
            <div style="font-size:12px;margin-bottom:8px">{{t('teamHook.env')}}</div>
            <a-select style="width: 100%" v-model:value="selectedEnv" :options="envList" />
          </div>
          <ul class="event-list">
            <li v-for="(item, index) in envRelatedEvents" v-bind:key="index">
              <div style="font-size:12px;margin-bottom:8px">{{t(item.title)}}</div>
              <ul class="action-list">
                <li v-for="action in item.actions" v-bind:key="`${index}-${action.key}`">
                  <a-checkbox
                    v-model:checked="envRelatedCheckboxes[selectedEnv][item.key][action.key]"
                  >{{t(action.title)}}</a-checkbox>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateTeamHook">{{t('teamHook.save')}}</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from "vue";
import {
  createTeamHookRequest,
  updateTeamHookRequest
} from "@/api/team/teamHookApi";
import { listAllTplByTeamIdRequest } from "@/api/team/notifyApi";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import { useRoute, useRouter } from "vue-router";
import {
  teamHookUrlRegexp,
  teamHookSecretRegexp,
  teamHookNameRegexp
} from "@/utils/regexp";
import { message } from "ant-design-vue";
import { useTeamHookStore } from "@/pinia/teamHookStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const selectedEnv = ref(null);
const route = useRoute();
// 模式
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
// 环境相关checkbox
const envRelatedCheckboxes = ref({});
// 环境列表
const envList = ref([]);
// 消息通知列表
const tplList = ref([]);
// 默认数据
const defaultEnvRelated = {
  appSource: {
    managePropertySource: false,
    manageDiscoverySource: false,
    manageServiceSource: false
  },
  appPropertyFile: {
    create: false,
    delete: false
  },
  appPropertyVersion: {
    new: false,
    deploy: false
  },
  appDeployPipeline: {
    create: false,
    delete: false,
    update: false
  },
  appDeployPipelineVars: {
    create: false,
    delete: false,
    update: false
  },
  appDeployPlan: {
    create: false,
    close: false,
    start: false
  },
  appDeployService: {
    triggerAction: false
  },
  appDiscovery: {
    markAsDown: false,
    markAsUp: false
  },
  appArtifact: {
    delete: false
  },
  appPromScrape: {
    create: false,
    update: false,
    delete: false
  },
  appAlertConfig: {
    create: false,
    update: false,
    delete: false,
    enable: false,
    disable: false
  },
  timer: {
    create: false,
    delete: false,
    update: false,
    enable: false,
    disable: false,
    manuallyTrigger: false
  },
  timerTask: {
    fail: false
  }
};
// events列表
const envRelatedEvents = [
  {
    key: "appSource",
    title: "teamHook.appSource.title",
    actions: [
      {
        key: "managePropertySource",
        title: "teamHook.appSource.managePropertySource"
      },
      {
        key: "manageDiscoverySource",
        title: "teamHook.appSource.manageDiscoverySource"
      },
      {
        key: "manageServiceSource",
        title: "teamHook.appSource.manageServiceSource"
      }
    ]
  },
  {
    key: "appPropertyFile",
    title: "teamHook.appPropertyFile.title",
    actions: [
      {
        key: "create",
        title: "teamHook.appPropertyFile.create"
      },
      {
        key: "delete",
        title: "teamHook.appPropertyFile.delete"
      }
    ]
  },
  {
    key: "appPropertyVersion",
    title: "teamHook.appPropertyVersion.title",
    actions: [
      {
        key: "new",
        title: "teamHook.appPropertyVersion.new"
      },
      {
        key: "deploy",
        title: "teamHook.appPropertyVersion.deploy"
      }
    ]
  },
  {
    key: "appDeployPipeline",
    title: "teamHook.appDeployPipeline.title",
    actions: [
      {
        key: "create",
        title: "teamHook.appDeployPipeline.create"
      },
      {
        key: "update",
        title: "teamHook.appDeployPipeline.update"
      },
      {
        key: "delete",
        title: "teamHook.appDeployPipeline.delete"
      }
    ]
  },
  {
    key: "appDeployPipelineVars",
    title: "teamHook.appDeployPipelineVars.title",
    actions: [
      {
        key: "create",
        title: "teamHook.appDeployPipelineVars.create"
      },
      {
        key: "update",
        title: "teamHook.appDeployPipelineVars.update"
      },
      {
        key: "delete",
        title: "teamHook.appDeployPipelineVars.delete"
      }
    ]
  },
  {
    key: "appDeployPlan",
    title: "teamHook.appDeployPlan.title",
    actions: [
      {
        key: "create",
        title: "teamHook.appDeployPlan.create"
      },
      {
        key: "close",
        title: "teamHook.appDeployPlan.close"
      },
      {
        key: "start",
        title: "teamHook.appDeployPlan.start"
      }
    ]
  },
  {
    key: "appDeployService",
    title: "teamHook.appDeployService.title",
    actions: [
      {
        key: "triggerAction",
        title: "teamHook.appDeployService.triggerAction"
      }
    ]
  },
  {
    key: "appDiscovery",
    title: "teamHook.appDiscovery.title",
    actions: [
      {
        key: "markAsDown",
        title: "teamHook.appDiscovery.markAsDown"
      },
      {
        key: "markAsUp",
        title: "teamHook.appDiscovery.markAsUp"
      }
    ]
  },
  {
    key: "appArtifact",
    title: "teamHook.appArtifact.title",
    actions: [
      {
        key: "upload",
        title: "teamHook.appArtifact.upload"
      },
      {
        key: "delete",
        title: "teamHook.appArtifact.delete"
      }
    ]
  },
  {
    key: "appPromScrape",
    title: "teamHook.appPromScrape.title",
    actions: [
      {
        key: "create",
        title: "teamHook.appPromScrape.create"
      },
      {
        key: "update",
        title: "teamHook.appPromScrape.update"
      },
      {
        key: "delete",
        title: "teamHook.appPromScrape.delete"
      }
    ]
  },
  {
    key: "appAlertConfig",
    title: "teamHook.appAlertConfig.title",
    actions: [
      {
        key: "create",
        title: "teamHook.appAlertConfig.create"
      },
      {
        key: "update",
        title: "teamHook.appAlertConfig.update"
      },
      {
        key: "delete",
        title: "teamHook.appAlertConfig.delete"
      },
      {
        key: "enable",
        title: "teamHook.appAlertConfig.enable"
      },
      {
        key: "disable",
        title: "teamHook.appAlertConfig.disable"
      }
    ]
  },
  {
    key: "timer",
    title: "teamHook.timer.title",
    actions: [
      {
        key: "create",
        title: "teamHook.timer.create"
      },
      {
        key: "update",
        title: "teamHook.timer.update"
      },
      {
        key: "delete",
        title: "teamHook.timer.delete"
      },
      {
        key: "enable",
        title: "teamHook.timer.enable"
      },
      {
        key: "disable",
        title: "teamHook.timer.disable"
      },
      {
        key: "manuallyTrigger",
        title: "teamHook.timer.manuallyTrigger"
      }
    ]
  },
  {
    key: "timerTask",
    title: "teamHook.timerTask.title",
    actions: [
      {
        key: "fail",
        title: "teamHook.timerTask.fail"
      }
    ]
  }
];
// 团队events列表
const teamCheckboxes = reactive([
  {
    key: "team",
    title: "teamHook.team.title",
    actions: [
      {
        key: "create",
        title: "teamHook.team.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.team.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.team.delete",
        value: false
      }
    ]
  },
  {
    key: "teamRole",
    title: "teamHook.teamRole.title",
    actions: [
      {
        key: "create",
        title: "teamHook.teamRole.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.teamRole.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.teamRole.delete",
        value: false
      }
    ]
  },
  {
    key: "teamUser",
    title: "teamHook.teamUser.title",
    actions: [
      {
        key: "create",
        title: "teamHook.teamUser.create",
        value: false
      },
      {
        key: "changeRole",
        title: "teamHook.teamUser.changeRole",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.teamUser.delete",
        value: false
      }
    ]
  },
  {
    key: "app",
    title: "teamHook.app.title",
    actions: [
      {
        key: "create",
        title: "teamHook.app.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.app.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.app.delete",
        value: false
      },
      {
        key: "transfer",
        title: "teamHook.app.transfer",
        value: false
      }
    ]
  },
  {
    key: "notifyTpl",
    title: "teamHook.notifyTpl.title",
    actions: [
      {
        key: "create",
        title: "teamHook.notifyTpl.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.notifyTpl.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.notifyTpl.delete",
        value: false
      },
      {
        key: "changeApiKey",
        title: "teamHook.notifyTpl.changeApiKey",
        value: false
      }
    ]
  },
  {
    key: "weworkAccessToken",
    title: "teamHook.weworkAccessToken.title",
    actions: [
      {
        key: "create",
        title: "teamHook.weworkAccessToken.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.weworkAccessToken.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.weworkAccessToken.delete",
        value: false
      },
      {
        key: "changeApiKey",
        title: "teamHook.weworkAccessToken.changeApiKey",
        value: false
      },
      {
        key: "refresh",
        title: "teamHook.weworkAccessToken.refresh",
        value: false
      }
    ]
  },
  {
    key: "feishuAccessToken",
    title: "teamHook.feishuAccessToken.title",
    actions: [
      {
        key: "create",
        title: "teamHook.feishuAccessToken.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.feishuAccessToken.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.feishuAccessToken.delete",
        value: false
      },
      {
        key: "changeApiKey",
        title: "teamHook.feishuAccessToken.changeApiKey",
        value: false
      },
      {
        key: "refresh",
        title: "teamHook.feishuAccessToken.refresh",
        value: false
      }
    ]
  }
]);
// git events列表
const gitCheckboxes = reactive([
  {
    key: "protectedBranch",
    title: "teamHook.protectedBranch.title",
    actions: [
      {
        key: "create",
        title: "teamHook.protectedBranch.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.protectedBranch.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.protectedBranch.delete",
        value: false
      }
    ]
  },
  {
    key: "gitPush",
    title: "teamHook.gitPush.title",
    actions: [
      {
        key: "commit",
        title: "teamHook.gitPush.commit",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.gitPush.delete",
        value: false
      }
    ]
  },
  {
    key: "pullRequest",
    title: "teamHook.pullRequest.title",
    actions: [
      {
        key: "submit",
        title: "teamHook.pullRequest.submit",
        value: false
      },
      {
        key: "close",
        title: "teamHook.pullRequest.close",
        value: false
      },
      {
        key: "merge",
        title: "teamHook.pullRequest.merge",
        value: false
      },
      {
        key: "review",
        title: "teamHook.pullRequest.review",
        value: false
      },
      {
        key: "addComment",
        title: "teamHook.pullRequest.addComment",
        value: false
      },
      {
        key: "deleteComment",
        title: "teamHook.pullRequest.deleteComment",
        value: false
      }
    ]
  },
  {
    key: "gitRepo",
    title: "teamHook.gitRepo.title",
    actions: [
      {
        key: "create",
        title: "teamHook.gitRepo.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.gitRepo.update",
        value: false
      },
      {
        key: "deleteTemporarily",
        title: "teamHook.gitRepo.deleteTemporarily",
        value: false
      },
      {
        key: "deletePermanently",
        title: "teamHook.gitRepo.deletePermanently",
        value: false
      },
      {
        key: "archived",
        title: "teamHook.gitRepo.archived",
        value: false
      },
      {
        key: "unArchived",
        title: "teamHook.gitRepo.unArchived",
        value: false
      },
      {
        key: "recoverFromRecycle",
        title: "teamHook.gitRepo.recoverFromRecycle",
        value: false
      }
    ]
  },
  {
    key: "gitWorkflow",
    title: "teamHook.gitWorkflow.title",
    actions: [
      {
        key: "create",
        title: "teamHook.gitWorkflow.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.gitWorkflow.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.gitWorkflow.delete",
        value: false
      },
      {
        key: "trigger",
        title: "teamHook.gitWorkflow.trigger",
        value: false
      },
      {
        key: "kill",
        title: "teamHook.gitWorkflow.kill",
        value: false
      }
    ]
  },
  {
    key: "gitWorkflowVars",
    title: "teamHook.gitWorkflowVars.title",
    actions: [
      {
        key: "create",
        title: "teamHook.gitWorkflowVars.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.gitWorkflowVars.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.gitWorkflowVars.delete",
        value: false
      }
    ]
  },
  {
    key: "gitWebhook",
    title: "teamHook.gitWebhook.title",
    actions: [
      {
        key: "create",
        title: "teamHook.gitWebhook.create",
        value: false
      },
      {
        key: "update",
        title: "teamHook.gitWebhook.update",
        value: false
      },
      {
        key: "delete",
        title: "teamHook.gitWebhook.delete",
        value: false
      }
    ]
  }
]);
const teamHookStore = useTeamHookStore();
const router = useRouter();
const mode = getMode();
// 表单数据
const formState = reactive({
  hookUrl: "",
  secret: "",
  name: "",
  hookType: 1,
  tplId: null
});
// 新增或编辑teamhook
const createOrUpdateTeamHook = () => {
  if (!teamHookNameRegexp.test(formState.name)) {
    message.warn(t("teamHook.nameFormatErr"));
    return;
  }
  if (formState.hookType === 1) {
    if (!teamHookUrlRegexp.test(formState.hookUrl)) {
      message.warn(t("teamHook.webhookUrlFormatErr"));
      return;
    }
    if (!teamHookSecretRegexp.test(formState.secret)) {
      message.warn(t("teamHook.webhookSecretFormatErr"));
      return;
    }
    formState.tplId = null;
  } else if (formState.hookType === 2) {
    if (!formState.tplId) {
      message.warn(t("teamHook.pleaseSelectNotificationTpl"));
      return;
    }
    formState.hookUrl = "";
    formState.secret = "";
  }
  let events = {};
  gitCheckboxes.forEach(git => {
    let actions = {};
    git.actions.forEach(action => {
      actions[action.key] = action.value;
    });
    events[git.key] = actions;
  });
  teamCheckboxes.forEach(team => {
    let actions = {};
    team.actions.forEach(action => {
      actions[action.key] = action.value;
    });
    events[team.key] = actions;
  });
  events["envRelated"] = JSON.parse(JSON.stringify(envRelatedCheckboxes.value));
  if (mode === "create") {
    createTeamHookRequest({
      name: formState.name,
      teamId: parseInt(route.params.teamId),
      hookType: formState.hookType,
      events: events,
      hookCfg: {
        hookUrl: formState.hookUrl,
        secret: formState.secret,
        notifyTplId: formState.tplId
      }
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/team/${route.params.teamId}/teamHook/list`);
    });
  } else if (mode === "update") {
    updateTeamHookRequest({
      id: teamHookStore.id,
      name: formState.name,
      teamId: parseInt(route.params.teamId),
      hookType: formState.hookType,
      events: events,
      hookCfg: {
        hookUrl: formState.hookUrl,
        secret: formState.secret,
        notifyTplId: formState.tplId
      }
    }).then(() => {
      message.success(t("operationSuccess"));
      router.push(`/team/${route.params.teamId}/teamHook/list`);
    });
  }
};
// 下拉框过滤
const filterTplListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 获取环境列表
const getEnvList = callback => {
  getEnvCfgRequest().then(res => {
    if (callback) {
      callback([...res.data]);
    }
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (res.data.length > 0) {
      selectedEnv.value = res.data[0];
    }
  });
};
// 获取消息模板列表
const getTplList = () => {
  listAllTplByTeamIdRequest(route.params.teamId).then(res => {
    tplList.value = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
  });
};
if (mode === "create") {
  getEnvList(envList => {
    envList.forEach(item => {
      let v = envRelatedCheckboxes.value;
      // 深拷贝
      v[item] = JSON.parse(JSON.stringify(defaultEnvRelated));
    });
  });
} else if (mode === "update") {
  if (teamHookStore.id === 0) {
    router.push(`/team/${route.params.teamId}/teamHook/list`);
  } else {
    formState.name = teamHookStore.name;
    formState.hookType = teamHookStore.hookType;
    formState.hookUrl = teamHookStore.hookCfg?.hookUrl;
    formState.secret = teamHookStore.hookCfg?.secret;
    if (teamHookStore.hookCfg?.notifyTplId === 0) {
      formState.tplId = null;
    } else {
      formState.tplId = teamHookStore.hookCfg?.notifyTplId;
    }
    gitCheckboxes.forEach(git => {
      git.actions.forEach(action => {
        if (teamHookStore.events[git.key]) {
          action.value = teamHookStore.events[git.key][action.key];
        }
      });
    });
    teamCheckboxes.forEach(team => {
      team.actions.forEach(action => {
        if (teamHookStore.events[team.key]) {
          action.value = teamHookStore.events[team.key][action.key];
        }
      });
    });
    getEnvList(envList => {
      if (teamHookStore.events["envRelated"]) {
        envList.forEach(item => {
          let v = envRelatedCheckboxes.value;
          if (teamHookStore.events["envRelated"][item]) {
            v[item] = teamHookStore.events["envRelated"][item];
          } else {
            // 深拷贝
            v[item] = JSON.parse(JSON.stringify(defaultEnvRelated));
          }
        });
      } else {
        envList.forEach(item => {
          let v = envRelatedCheckboxes.value;
          // 深拷贝
          v[item] = JSON.parse(JSON.stringify(defaultEnvRelated));
        });
      }
    });
  }
}
getTplList();
</script>
<style scoped>
.action-list {
  font-size: 14px;
  display: flex;
  flex-wrap: wrap;
}
.action-list > li {
  width: 33.33%;
  margin-bottom: 8px;
}
.event-list > li + li {
  margin-top: 16px;
}
</style>