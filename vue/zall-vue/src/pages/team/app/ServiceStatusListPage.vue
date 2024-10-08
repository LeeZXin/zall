<template>
  <div style="padding:10px;height:100%" v-if="!showStatusList">
    <div style="margin-bottom:10px;" class="flex-between">
      <div>
        <a-button
          type="primary"
          :icon="h(SettingOutlined)"
          @click="showBindServiceSourceModal"
          v-if="appStore.perm?.canManageServiceSource"
        >{{t('serviceStatus.manageSource')}}</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="sourceColumns" :dataSource="sourceDataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <span
          v-else
          class="check-btn"
          @click="listAndShowStatusList(dataItem)"
        >{{t('serviceStatus.view')}}</span>
      </template>
    </ZTable>
  </div>
  <div style="padding:10px;height:100%" v-if="showStatusList">
    <div style="margin-bottom:10px;">
      <span class="check-btn" @click="backToSource">
        <LeftOutlined />
        <span style="margin-left:4px">{{t('serviceStatus.backToSelectSource')}}</span>
      </span>
    </div>
    <div style="margin-bottom:10px;">
      <span style="font-weight:bold">{{selectedSource?selectedSource.name:""}}</span>
      <a-tag color="orange" style="margin-left:10px">{{selectedSource.env}}</a-tag>
    </div>
    <ZTable :columns="statusColumns" :dataSource="statusDataSource" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex === 'created'">{{readableTimeComparingNow(dataItem[dataIndex])}}</span>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="killService(dataItem)">
                  <CloseOutlined />
                  <span style="margin-left:4px">{{t('serviceStatus.killService')}}</span>
                </li>
                <li @click="restartService(dataItem)">
                  <ReloadOutlined />
                  <span style="margin-left:4px">{{t('serviceStatus.restartService')}}</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </div>
      </template>
    </ZTable>
  </div>
  <a-modal
    v-model:open="bindModal.open"
    :title="t('serviceStatus.bindSource')"
    @ok="handleBindModalOk"
  >
    <div>
      <div style="font-size:12px;margin-bottom:3px">{{t('serviceStatus.selectedEnv')}}</div>
      <div>{{selectedEnv}}</div>
    </div>
    <div style="margin-top: 10px">
      <div style="font-size:12px;margin-bottom:3px">{{t('serviceStatus.source')}}</div>
      <a-select
        style="width: 100%"
        v-model:value="bindModal.selectIdList"
        :options="bindModal.sourceList"
        show-search
        mode="multiple"
        :filter-option="filterSourceListOption"
      />
    </div>
  </a-modal>
</template>
<script setup>
import {
  ExclamationCircleOutlined,
  SettingOutlined,
  LeftOutlined,
  EllipsisOutlined,
  CloseOutlined,
  ReloadOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, onUnmounted, createVNode, reactive, h } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listBindServiceSourceRequest,
  listServiceStatusRequest,
  listAllServiceSourceRequest,
  bindAppAndServiceSourceRequest,
  killServiceRequest,
  restartServiceRequest
} from "@/api/app/serviceApi";
import { message, Modal } from "ant-design-vue";
import { useAppStore } from "@/pinia/appStore";
import { useI18n } from "vue-i18n";
import { readableTimeComparingNow } from "@/utils/time";
const { t } = useI18n();
const appStore = useAppStore();
const route = useRoute();
const bindModal = reactive({
  open: false,
  selectIdList: [],
  sourceList: []
});
// 当前环境
const selectedEnv = ref("");
const router = useRouter();
// 服务来源数据
const sourceDataSource = ref([]);
// 状态 <-> 来源切换flag
const showStatusList = ref(false);
// 来源数据项
const sourceColumns = [
  {
    i18nTitle: "serviceStatus.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "serviceStatus.operation",
    dataIndex: "operation",
    key: "operation"
  }
];
// 状态数据项
const statusColumns = [
  {
    title: "id",
    dataIndex: "id",
    key: "id"
  },
  {
    i18nTitle: "serviceStatus.host",
    dataIndex: "host",
    key: "host"
  },
  {
    i18nTitle: "serviceStatus.status",
    dataIndex: "status",
    key: "status"
  },
  {
    i18nTitle: "serviceStatus.cpuPercent",
    dataIndex: "cpuPercent",
    key: "cpuPercent"
  },
  {
    i18nTitle: "serviceStatus.memPercent",
    dataIndex: "memPercent",
    key: "memPercent"
  },
  {
    i18nTitle: "serviceStatus.createTime",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "serviceStatus.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 状态数据
const statusDataSource = ref([]);
// 间隔获取状态interval
const statusInterval = ref(null);
// 选择的来源
const selectedSource = ref(null);
// 获取服务来源状态
const listServiceSource = () => {
  listBindServiceSourceRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    sourceDataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
// 清除interval
const clearStatusInterval = () => {
  if (statusInterval.value) {
    clearInterval(statusInterval.value);
  }
};
// 环境变化
const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/serviceStatus/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  listServiceSource();
};
// 返回来源界面
const backToSource = () => {
  showStatusList.value = false;
  clearStatusInterval();
};
// 获取状态数据并切换界面
const listAndShowStatusList = data => {
  selectedSource.value = data;
  listServiceStatusRequest(data.bindId).then(res => {
    showStatusList.value = true;
    statusDataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item,
        cpuPercent: item.cpuPercent + "%",
        memPercent: item.memPercent + "%"
      };
    });
    refreshStatus();
  });
};
// 获取状态数据
const listStatus = () => {
  listServiceStatusRequest(selectedSource.value.bindId).then(res => {
    statusDataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item,
        cpuPercent: item.cpuPercent + "%",
        memPercent: item.memPercent + "%"
      };
    });
  });
};
// 刷新状态数据
const refreshStatus = () => {
  clearStatusInterval();
  statusInterval.value = setInterval(listStatus, 5000);
};
// 杀死服务
const killService = item => {
  Modal.confirm({
    title: `${t("serviceStatus.confirmKillService")} ${item.id}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      killServiceRequest({
        bindId: selectedSource.value.bindId,
        serviceId: item.id
      }).then(() => {
        message.success(t("operationSuccess"));
      });
    },
    onCancel() {}
  });
};
// 重启服务
const restartService = item => {
  Modal.confirm({
    title: `${t("serviceStatus.confirmRestartService")} ${item.id}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      restartServiceRequest({
        bindId: selectedSource.value.bindId,
        serviceId: item.id
      }).then(() => {
        message.success(t("operationSuccess"));
      });
    },
    onCancel() {}
  });
};
// 下拉框搜索过滤
const filterSourceListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 展示绑定来源modal
const showBindServiceSourceModal = () => {
  if (!selectedEnv.value) {
    return;
  }
  listAllServiceSourceRequest(selectedEnv.value).then(res => {
    bindModal.sourceList = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    listBindServiceSourceRequest({
      appId: route.params.appId,
      env: selectedEnv.value
    }).then(res => {
      bindModal.selectIdList = res.data.map(item => item.id);
      bindModal.open = true;
    });
  });
};
// 绑定modal点击“确定”按钮
const handleBindModalOk = () => {
  bindAppAndServiceSourceRequest({
    appId: route.params.appId,
    sourceIdList: bindModal.selectIdList,
    env: selectedEnv.value
  }).then(() => {
    message.success(t("operationSuccess"));
    bindModal.open = false;
    listServiceSource();
  });
};
// 解除interval
onUnmounted(() => clearStatusInterval());
</script>
<style scoped>
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>