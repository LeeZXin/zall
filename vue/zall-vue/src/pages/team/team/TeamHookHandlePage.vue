<template>
  <div style="padding:10px">
    <div class="container">
      <div class="header">
        <span v-if="mode === 'create'">添加Team Hook</span>
        <span v-else-if="mode === 'update'">更新Team Hook</span>
      </div>
      <div class="section">
        <div class="section-title">
          <span>名称</span>
        </div>
        <div class="section-body">
          <a-input style="width:100%" v-model:value="formState.name" />
          <div class="input-desc">标识team hook</div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>类型</span>
          <span style="color:darkred">*</span>
        </div>
        <div class="section-body">
          <a-radio-group v-model:value="formState.hookType">
            <a-radio :value="1">Webhook</a-radio>
            <a-radio :value="2">外部通知</a-radio>
          </a-radio-group>
        </div>
      </div>
      <div class="section" v-if="formState.hookType === 1">
        <div class="section-title">
          <span>Webhook</span>
        </div>
        <div class="section-body">
          <div>
            <div style="font-size: 12px;margin-bottom: 6px">hook url</div>
            <a-input style="width:100%" v-model:value="formState.hookUrl" placeholder="请填写" />
          </div>
          <div style="margin-top: 10px">
            <div style="font-size: 12px;margin-bottom: 6px">签名密钥</div>
            <a-input-password
              style="width:100%"
              v-model:value="formState.secret"
              placeholder="请填写"
            />
          </div>
        </div>
      </div>
      <div class="section" v-else-if="formState.hookType === 2">
        <div class="section-title">
          <span>外部通知模板</span>
        </div>
        <div class="section-body">
          <a-select
            style="width: 100%"
            v-model:value="formState.tplId"
            :options="tplList"
            show-search
            :filter-option="filterTplListOption"
            placeholder="请选择"
          />
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>代码仓库事件</span>
        </div>
        <div class="section-body">
          <ul class="event-list">
            <li v-for="(item, index) in gitCheckboxes" v-bind:key="index">
              <div style="font-size:12px;margin-bottom:8px">{{item.title}}</div>
              <ul class="action-list">
                <li v-for="action in item.actions" v-bind:key="`${index}-${action.key}`">
                  <a-checkbox v-model:checked="action.value">{{action.title}}</a-checkbox>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="section">
        <div class="section-title">
          <span>团队事件</span>
        </div>
        <div class="section-body">
          <ul class="event-list">
            <li v-for="(item, index) in teamCheckboxes" v-bind:key="index">
              <div style="font-size:12px;margin-bottom:8px">{{item.title}}</div>
              <ul class="action-list">
                <li v-for="action in item.actions" v-bind:key="`${index}-${action.key}`">
                  <a-checkbox v-model:checked="action.value">{{action.title}}</a-checkbox>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="section" v-if="envList.length > 0">
        <div class="section-title">
          <span>环境相关事件</span>
        </div>
        <div class="section-body">
          <div style="margin-bottom: 18px">
            <div style="font-size:12px;margin-bottom:8px">环境</div>
            <a-select style="width: 100%" v-model:value="selectedEnv" :options="envList" />
          </div>
          <ul class="event-list">
            <li v-for="(item, index) in envRelatedEvents" v-bind:key="index">
              <div style="font-size:12px;margin-bottom:8px">{{item.title}}</div>
              <ul class="action-list">
                <li v-for="action in item.actions" v-bind:key="`${index}-${action.key}`">
                  <a-checkbox
                    v-model:checked="envRelatedCheckboxes[selectedEnv][item.key][action.key]"
                  >{{action.title}}</a-checkbox>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
      <div class="save-btn-line">
        <a-button type="primary" @click="createOrUpdateTeamHook">立即保存</a-button>
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
const selectedEnv = ref(null);
const route = useRoute();
const getMode = () => {
  let s = route.path.split("/");
  return s[s.length - 1];
};
const envRelatedCheckboxes = ref({});
const envList = ref([]);
const tplList = ref([]);
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
    deregister: false,
    reRegister: false,
    deleteDownService: false
  },
  appProduct: {
    delete: false
  },
  appPromScrape: {
    create: false,
    update: false,
    delete: false
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
const envRelatedEvents = [
  {
    key: "appSource",
    title: "应用服务资源",
    actions: [
      {
        key: "managePropertySource",
        title: "管理配置中心来源"
      },
      {
        key: "manageDiscoverySource",
        title: "管理注册中心来源"
      },
      {
        key: "manageServiceSource",
        title: "管理服务状态来源"
      }
    ]
  },
  {
    key: "appPropertyFile",
    title: "配置文件",
    actions: [
      {
        key: "create",
        title: "新增"
      },
      {
        key: "delete",
        title: "删除"
      }
    ]
  },
  {
    key: "appPropertyVersion",
    title: "配置版本",
    actions: [
      {
        key: "new",
        title: "新增"
      },
      {
        key: "deploy",
        title: "发布"
      }
    ]
  },
  {
    key: "appDeployPipeline",
    title: "部署流水线",
    actions: [
      {
        key: "create",
        title: "新增"
      },
      {
        key: "update",
        title: "编辑"
      },
      {
        key: "delete",
        title: "删除"
      }
    ]
  },
  {
    key: "appDeployPipelineVars",
    title: "部署流水线变量",
    actions: [
      {
        key: "create",
        title: "新增"
      },
      {
        key: "update",
        title: "编辑"
      },
      {
        key: "delete",
        title: "删除"
      }
    ]
  },
  {
    key: "appDeployPlan",
    title: "发布计划",
    actions: [
      {
        key: "create",
        title: "新增"
      },
      {
        key: "close",
        title: "关闭"
      },
      {
        key: "start",
        title: "开始"
      }
    ]
  },
  {
    key: "appDeployService",
    title: "部署服务",
    actions: [
      {
        key: "triggerAction",
        title: "触发指令"
      }
    ]
  },
  {
    key: "appDiscovery",
    title: "注册中心",
    actions: [
      {
        key: "deregister",
        title: "下线服务"
      },
      {
        key: "reRegister",
        title: "上线服务"
      },
      {
        key: "deleteDownService",
        title: "删除下线服务"
      }
    ]
  },
  {
    key: "appProduct",
    title: "应用制品",
    actions: [
      {
        key: "delete",
        title: "删除"
      }
    ]
  },
  {
    key: "appPromScrape",
    title: "Prometheus抓取任务",
    actions: [
      {
        key: "create",
        title: "新增"
      },
      {
        key: "update",
        title: "编辑"
      },
      {
        key: "delete",
        title: "删除"
      }
    ]
  },
  {
    key: "timer",
    title: "定时任务",
    actions: [
      {
        key: "create",
        title: "新增"
      },
      {
        key: "update",
        title: "编辑"
      },
      {
        key: "delete",
        title: "删除"
      },
      {
        key: "enable",
        title: "启用"
      },
      {
        key: "disable",
        title: "禁用"
      },
      {
        key: "manuallyTrigger",
        title: "手动触发"
      }
    ]
  },
  {
    key: "timerTask",
    title: "定时任务执行",
    actions: [
      {
        key: "fail",
        title: "任务失败"
      }
    ]
  }
];
const teamCheckboxes = reactive([
  {
    key: "team",
    title: "团队",
    actions: [
      {
        key: "create",
        title: "新增",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      }
    ]
  },
  {
    key: "teamRole",
    title: "团队角色",
    actions: [
      {
        key: "create",
        title: "新增",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      }
    ]
  },
  {
    key: "teamUser",
    title: "团队成员",
    actions: [
      {
        key: "create",
        title: "新增",
        value: false
      },
      {
        key: "changeRole",
        title: "变更角色",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      }
    ]
  },
  {
    key: "app",
    title: "应用服务",
    actions: [
      {
        key: "create",
        title: "新增",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      },
      {
        key: "transfer",
        title: "迁移",
        value: false
      }
    ]
  }
]);
const gitCheckboxes = reactive([
  {
    key: "protectedBranch",
    title: "保护分支",
    actions: [
      {
        key: "create",
        title: "新增",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      }
    ]
  },
  {
    key: "gitPush",
    title: "代码提交",
    actions: [
      {
        key: "commit",
        title: "提交",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      }
    ]
  },
  {
    key: "pullRequest",
    title: "合并请求",
    actions: [
      {
        key: "submit",
        title: "提交",
        value: false
      },
      {
        key: "close",
        title: "关闭",
        value: false
      },
      {
        key: "merge",
        title: "合并",
        value: false
      },
      {
        key: "review",
        title: "评审",
        value: false
      },
      {
        key: "addComment",
        title: "添加评论",
        value: false
      },
      {
        key: "deleteComment",
        title: "删除评论",
        value: false
      }
    ]
  },
  {
    key: "gitRepo",
    title: "代码仓库",
    actions: [
      {
        key: "create",
        title: "新增",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "deleteTemporarily",
        title: "临时删除",
        value: false
      },
      {
        key: "deletePermanently",
        title: "永久删除",
        value: false
      },
      {
        key: "archived",
        title: "归档",
        value: false
      },
      {
        key: "unArchived",
        title: "取消归档",
        value: false
      },
      {
        key: "recoverFromRecycle",
        title: "恢复删除",
        value: false
      }
    ]
  },
  {
    key: "gitWorkflow",
    title: "git工作流",
    actions: [
      {
        key: "create",
        title: "创建",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      },
      {
        key: "trigger",
        title: "触发",
        value: false
      },
      {
        key: "kill",
        title: "停止任务",
        value: false
      }
    ]
  },
  {
    key: "gitWorkflowVars",
    title: "git工作流变量",
    actions: [
      {
        key: "create",
        title: "创建",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      }
    ]
  },
  {
    key: "gitWebhook",
    title: "git webhook",
    actions: [
      {
        key: "create",
        title: "创建",
        value: false
      },
      {
        key: "update",
        title: "编辑",
        value: false
      },
      {
        key: "delete",
        title: "删除",
        value: false
      }
    ]
  }
]);
const teamHookStore = useTeamHookStore();
const router = useRouter();
const mode = getMode();
const formState = reactive({
  hookUrl: "",
  secret: "",
  name: "",
  hookType: 1,
  tplId: null
});
const createOrUpdateTeamHook = () => {
  if (!teamHookNameRegexp.test(formState.name)) {
    message.warn("名称格式错误");
    return;
  }
  if (formState.hookType === 1) {
    if (!teamHookUrlRegexp.test(formState.hookUrl)) {
      message.warn("url格式错误");
      return;
    }
    if (!teamHookSecretRegexp.test(formState.secret)) {
      message.warn("密钥格式错误");
      return;
    }
    formState.tplId = null;
  } else if (formState.hookType === 2) {
    if (!formState.tplId) {
      message.warn("请选择通知模板");
      return;
    }
    formState.hookUrl = "";
    formState.secret = "";
  } else {
    message.warn("未选择类型");
    return;
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
      message.success("添加成功");
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
      message.success("更新成功");
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
.header {
  font-size: 18px;
  margin-bottom: 10px;
  font-weight: bold;
}
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