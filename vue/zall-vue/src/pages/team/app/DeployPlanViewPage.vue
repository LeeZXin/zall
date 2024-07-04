<template>
  <div style="padding:14px">
    <ul class="info-list" v-if="planDetail">
      <li>
        <div class="info-name">名称</div>
        <div class="info-value">
          <span>{{planDetail.name}}</span>
          <a-tag color="orange" style="font-weight:bold;margin-left: 8px">{{planDetail.env}}</a-tag>
        </div>
      </li>
      <li>
        <div class="info-name">状态</div>
        <div class="info-value">
          <PlanStatusTag :status="planDetail.planStatus" />
        </div>
      </li>
      <li>
        <div class="info-name">创建人</div>
        <div class="info-value">{{planDetail.creator}}</div>
      </li>
    </ul>
    <ul class="info-list" v-if="planDetail">
      <li>
        <div class="info-name">创建时间</div>
        <div class="info-value">{{planDetail.created}}</div>
      </li>
      <li>
        <div class="info-name">制品号</div>
        <div class="info-value">{{planDetail.productVersion}}</div>
      </li>
      <li>
        <div class="info-name">配置</div>
        <div class="info-value">
          <span class="check-config-btn" @click="yamlModalOpen = true">查看配置</span>
        </div>
      </li>
    </ul>
    <div>
      <a-button
        type="primary"
        :icon="h(PlayCircleOutlined)"
        v-if="planDetail?.planStatus === 1"
        @click="startPlan"
      >执行发布计划</a-button>
    </div>
    <div class="service-name no-wrap">{{planDetail?planDetail.pipelineName:''}}</div>
    <div class="flex" style="padding-left:calc(50% - 420px)">
      <ul class="stages">
        <li v-for="(item, index) in stageList" v-bind:key="index">
          <div class="stage-card">
            <div class="stage-info">
              <div
                class="stage-name flex-between"
                :style="{'background-color': getStageNameBackgroundColor(item)}"
                @click="selectStage(item, index)"
              >
                <span>{{item.isAutomatic?'自动阶段':'交互阶段'}}-{{item.name}}</span>
                <span v-if="item.isRunning || item.hasError">
                  <LoadingOutlined v-if="item.isRunning" />
                  <span style="margin-left:8px">
                    <span v-if="item.isRunning">执行中</span>
                    <span v-else-if="item.hasError">执行失败</span>
                  </span>
                </span>
              </div>
              <div class="stage-progress">
                <a-progress
                  :percent="item.percent"
                  size="small"
                  style="padding:9px"
                  :format="()=>`${item.done}/${item.total}`"
                />
                <div
                  class="wait-interact-btn"
                  v-if="item.waitInteract"
                  @click="showConfirmModal(item, index)"
                >
                  <LoadingOutlined />
                  <span style="margin-left:8px">等待交互</span>
                </div>
                <div class="kill-btn" v-if="item.isRunning" @click="killStage(index)">
                  <CloseOutlined />
                  <span style="margin-left:8px">中止执行</span>
                </div>
              </div>
            </div>
            <div class="arrow-down" v-if="index < stageList.length - 1">
              <ArrowDownOutlined />
            </div>
          </div>
          <div class="stage-line"></div>
        </li>
      </ul>
      <div class="stage-detail" v-if="selectedStage">
        <div class="stage-agent-name flex-between">
          <span>{{selectedStage.name}}</span>
          <a-popover placement="top" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="forceRedoUnSuccessAgentStages">
                  <span style="margin-left:8px">强制重新执行未成功任务</span>
                </li>
              </ul>
            </template>
            <EllipsisOutlined
              v-if="selectedStage.canForceRedoUnSuccessAgentStages"
              style="font-size:22px"
            />
          </a-popover>
        </div>
        <ul class="detail-list">
          <li
            v-for="(item, index) in selectedStage.subStages"
            v-bind:key="item.id +'_' + index"
            @click="selectSubStage(item)"
          >
            <div class="flex-between">
              <span class="no-wrap">{{item.agent}}</span>
              <StageStatusTag :status="item.stageStatus" />
            </div>
          </li>
        </ul>
      </div>
    </div>
    <a-drawer
      :title="selectedStage ? selectedStage.name:''"
      v-model:open="drawerVisible"
      :bodyStyle="{padding: '10px 20px'}"
      :closable="false"
      :width="500"
    >
      <ul class="agent-info-ul">
        <li>
          <div>代理:</div>
          <div>{{selectedSubStage.agent}}</div>
        </li>
        <li>
          <div>Host:</div>
          <div>{{selectedSubStage.agentHost}}</div>
        </li>
        <li>
          <div>状态:</div>
          <div>
            <StageStatusTag :status="selectedSubStage.stageStatus" />
          </div>
        </li>
        <li v-show="selectedSubStage.executeLog?.length > 0">
          <div>执行日志:</div>
          <div>
            <Codemirror
              v-model="selectedSubStage.executeLog"
              style="height:200px;width:100%"
              :extensions="extensions"
              :disabled="true"
            />
          </div>
        </li>
        <li v-show="selectedSubStage.rollbackLog?.length > 0">
          <div>回滚日志:</div>
          <div>
            <Codemirror
              v-model="selectedSubStage.rollbackLog"
              style="height:200px;width:100%"
              :extensions="extensions"
              :disabled="true"
            />
          </div>
        </li>
      </ul>
      <div
        class="agent-info-btn"
        v-if="planDetail.planStatus === 2 && selectedSubStage.stageStatus === 4"
      >
        <a-button type="primary" style="width:100%" @click="redoAgentStage">重新执行</a-button>
      </div>
    </a-drawer>
    <a-modal v-model:open="yamlModalOpen" title="服务配置" :footer="null">
      <Codemirror
        v-model="yamlConfig"
        style="height:380px;width:100%"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
    <a-modal v-model:open="confirmModal.open" :title="confirmModal.title" @ok="confirmOk">
      <div style="padding:10px;">
        <div style="font-size:14px;margin-bottom:10px;font-weight:bold">{{confirmModal.message}}</div>
        <ul class="confirm-form-ul">
          <li v-for="item in confirmModal.formItems" v-bind:key="item.bindKey">
            <div>{{item.label || item.key}}</div>
            <div v-if="item.options?.length > 0">
              <a-select style="width: 100%" v-model:value="item.value" :options="item.options" />
            </div>
            <div v-else>
              <a-input v-model:value="item.value" style="width:100%" />
            </div>
          </li>
        </ul>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import PlanStatusTag from "@/components/app/PlanStatusTag";
import StageStatusTag from "@/components/app/StageStatusTag";
import {
  ArrowDownOutlined,
  PlayCircleOutlined,
  LoadingOutlined,
  ExclamationCircleOutlined,
  CloseOutlined,
  EllipsisOutlined
} from "@ant-design/icons-vue";
import { ref, h, createVNode, onUnmounted, reactive } from "vue";
import {
  getDeployPlanDetailRequest,
  startDeployPlanRequest,
  listDeployPlanStagesRequest,
  redoDeployAgentStageRequest,
  killDeployStageRequest,
  confirmInteractStageRequest,
  forceRedoStageRequest
} from "@/api/app/deployPlanApi";
import { Codemirror } from "vue-codemirror";
import { yaml } from "@codemirror/lang-yaml";
import { oneDark } from "@codemirror/theme-one-dark";
import { useDeloyPlanStore } from "@/pinia/deployPlanStore";
import { useRouter, useRoute } from "vue-router";
import { Modal, message } from "ant-design-vue";
const yamlConfig = ref("");
const extensions = [yaml(), oneDark];
const yamlModalOpen = ref(false);
const router = useRouter();
const route = useRoute();
const planStore = useDeloyPlanStore();
const planDetail = ref(null);
const stageList = ref([]);
const selectedStage = ref(null);
const selectedIndex = ref(-1);
const selectedSubStage = ref(null);
const refreshInterval = ref(null);
const drawerVisible = ref(false);
const confirmModal = reactive({
  open: false,
  title: "",
  formItems: [],
  stageIndex: -1,
  type: ""
});
const selectStage = (item, index) => {
  selectedStage.value = item;
  selectedIndex.value = index;
};

const selectSubStage = item => {
  selectedSubStage.value = item;
  drawerVisible.value = true;
};

const getStageNameBackgroundColor = data => {
  if (data.hasError) {
    return "#ff4d4f";
  }
  if (data.waitInteract) {
    return "#1677ff";
  }
  if (data.isAllDone) {
    return "#52c41a";
  }
  if (data.isRunning) {
    return "#1677ff";
  }
  return "#6b6b6b";
};

const forceRedoUnSuccessAgentStages = () => {
  if (selectedStage.value.isAutomatic) {
    Modal.confirm({
      title: `你确定要强制重新执行吗?`,
      icon: createVNode(ExclamationCircleOutlined),
      onOk() {
        forceRedoStageRequest(
          {
            planId: planStore.id,
            stageIndex: selectedIndex.value,
          },
          planStore.env
        ).then(() => {
          message.success("开始执行");
          getDetail();
        });
      },
      onCancel() {}
    });
  } else {
    handleConfirmModal("force", selectedStage.value, selectedIndex.value);
  }
};

const getDetail = () => {
  getDeployPlanDetailRequest(planStore.id, planStore.env).then(res => {
    planDetail.value = res.data;
    yamlConfig.value = res.data.pipelineConfig;
    listStages();
    if (res.data?.planStatus === 3 || res.data?.planStatus === 4) {
      clearRefreshInterval();
    }
  });
};

const startPlan = () => {
  Modal.confirm({
    title: `你确定要执行发布计划吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      startDeployPlanRequest(planStore.id, planStore.env).then(() => {
        message.success("开始执行");
        getDetail();
      });
    },
    onCancel() {}
  });
};

const listStages = () => {
  listDeployPlanStagesRequest(planStore.id, planStore.env).then(res => {
    stageList.value = res.data.map((item, index) => {
      return {
        id: index + 1,
        ...item
      };
    });
    if (selectedStage.value) {
      let stage = stageList.value.find(
        item => item.id === selectedStage.value.id
      );
      if (stage) {
        selectedStage.value = stage;
      }
      if (drawerVisible.value && stage) {
        // 如果详情展开 修改详情界面
        let sub = stage.subStages.find(
          item => item.id === selectedSubStage.value.id
        );
        if (sub) {
          selectedSubStage.value = sub;
        }
      }
    }
  });
};

const clearRefreshInterval = () => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value);
  }
};

const redoAgentStage = () => {
  Modal.confirm({
    title: `你确定要重新执行吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      redoDeployAgentStageRequest(
        selectedSubStage.value.id,
        planStore.env
      ).then(() => {
        message.success("重新执行成功");
        getDetail();
      });
    },
    onCancel() {}
  });
};

const killStage = index => {
  Modal.confirm({
    title: `你确定要中止执行吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      killDeployStageRequest(planStore.id, index, planStore.env).then(() => {
        message.success("中止成功");
        getDetail();
      });
    },
    onCancel() {}
  });
};

const showConfirmModal = (item, index) => {
  handleConfirmModal("normal", item, index);
};

const handleConfirmModal = (type, item, index) => {
  console.log(type, item, index);
  if (item.confirm) {
    confirmModal.type = type;
    confirmModal.title = item.name;
    confirmModal.message = item.confirm.message;
    if (item.confirm.form) {
      confirmModal.formItems = item.confirm.form.map((f, i) => {
        return {
          bindKey: `${index}_${i}`,
          value: undefined,
          ...f
        };
      });
    } else {
      confirmModal.formItems = [];
    }
    confirmModal.stageIndex = index;
    confirmModal.open = true;
  }
};

const confirmOk = () => {
  let args = {};
  if (confirmModal.formItems?.length > 0) {
    for (let index in confirmModal.formItems) {
      let fi = confirmModal.formItems[index];
      if (!fi.options || fi.options.length === 0) {
        let re = new RegExp(fi.regexp);
        if (!re || !re.test(fi.value)) {
          message.warn(`${fi.label || fi.key} 格式错误`);
          return;
        }
      } else {
        if (!fi.options.find(item => item.value === fi.value)) {
          message.warn(`${fi.label || fi.key} 格式错误`);
          return;
        }
      }
      args[fi.key] = fi.value;
    }
  }
  let request;
  switch (confirmModal.type) {
    case "force":
      request = forceRedoStageRequest;
      break;
    case "normal":
      request = confirmInteractStageRequest;
      break;
  }

  request(
    {
      planId: planStore.id,
      stageIndex: confirmModal.stageIndex,
      args
    },
    planStore.env
  ).then(() => {
    message.success("操作成功");
    confirmModal.open = false;
    getDetail();
  });
};

if (planStore.id && planStore.id > 0) {
  getDetail();
  onUnmounted(clearRefreshInterval);
  refreshInterval.value = setInterval(getDetail, 10000);
} else {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/list`
  );
}
</script>

<style scoped>
.check-config-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
.service-name {
  font-size: 16px;
  font-weight: bold;
  line-height: 42px;
  padding-left: 10px;
  margin-bottom: 20px;
}

.stages {
  margin-right: 40px;
}

.stages > li {
  display: flex;
}

.stage-card {
  width: 400px;
}

.stage-info {
  width: 100%;
}

.stage-progress {
  border-bottom: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
}

.stage-name {
  width: 100%;
  font-size: 14px;
  word-break: break-all;
  line-height: 40px;
  padding: 0 9px;
  line-height: 36px;
  color: white;
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
  background-color: #6b6b6b;
  cursor: pointer;
}

.stage-agent-name {
  width: 100%;
  font-size: 14px;
  word-break: break-all;
  line-height: 40px;
  padding: 0 9px;
  line-height: 36px;
  color: white;
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
  background-color: #1677ff;
}

.stage-name-error {
  background-color: #ff4d4f;
}
.arrow-down {
  width: 100%;
  text-align: center;
  font-size: 22px;
  color: gray;
  padding: 6px 0px;
}

.stage-detail {
  width: 400px;
  overflow: scroll;
  max-height: 600px;
}

.detail-list {
  border-bottom-left-radius: 4px;
  border-bottom-right-radius: 4px;
  border-bottom: 1px solid #d9d9d9;
  border-left: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
}

.detail-list > li {
  padding: 10px;
  font-size: 14px;
}

.detail-list > li + li {
  border-top: 1px solid #d9d9d9;
}

.detail-list > li:hover {
  cursor: pointer;
  background-color: #f0f0f0;
}

.info-list > li {
  width: 33.33%;
}

.agent-info-ul > li {
  line-height: 32px;
  width: 100%;
  display: flex;
  font-size: 14px;
}

.agent-info-ul > li + li {
  margin-top: 10px;
}
.agent-info-ul > li > div:nth-child(1) {
  width: 20%;
  word-break: break-all;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.agent-info-ul > li > div:nth-child(2) {
  width: 80%;
  word-break: break-all;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.agent-info-btn {
  margin-top: 20px;
}
.kill-btn {
  text-align: center;
  line-height: 32px;
  color: #ff4d4f;
}
.kill-btn:hover {
  cursor: pointer;
  background-color: #ffe7e7;
}
.wait-interact-btn {
  text-align: center;
  line-height: 32px;
  color: #1677ff;
}
.wait-interact-btn:hover {
  cursor: pointer;
  background-color: #afcdf8;
}
.confirm-form-ul > li > div:nth-child(1) {
  font-size: 14px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  word-break: break-all;
  margin-bottom: 6px;
}
.confirm-form-ul > li + li {
  margin-top: 8px;
}
</style>