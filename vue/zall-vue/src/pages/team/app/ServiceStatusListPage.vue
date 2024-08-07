<template>
  <div style="padding:10px;height:100%" v-if="!showStatusList">
    <div style="margin-bottom:10px;" class="flex-between">
      <div>
        <a-button
          type="primary"
          :icon="h(SettingOutlined)"
          @click="showBindServiceSourceModal"
          v-if="appStore.perm?.canManageServiceSource"
        >管理服务状态来源绑定</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="sourceColumns" :dataSource="sourceDataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <span v-else class="check-btn" @click="listAndShowStatusList(dataItem)">查看</span>
      </template>
    </ZTable>
  </div>
  <div style="padding:10px;height:100%" v-if="showStatusList">
    <div style="margin-bottom:10px;">
      <span class="check-btn" @click="backToSource">返回集群选择</span>
    </div>
    <div style="margin-bottom:10px;font-weight:bold">
      <span>{{selectedSource?selectedSource.name:""}}</span>
      <a-tag color="orange" style="margin-left:10px">{{selectedSource.env}}</a-tag>
    </div>
    <ZTable :columns="statusColumns" :dataSource="statusDataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <ul class="action-ul" v-else>
          <li
            v-for="item in actionList"
            v-bind:key="item"
            @click="doAction(dataItem, item)"
          >{{item}}</li>
        </ul>
      </template>
    </ZTable>
  </div>
  <a-modal v-model:open="bindModal.open" title="绑定配置来源" @ok="handleBindModalOk">
    <a-select
      style="width: 100%"
      placeholder="请选择"
      v-model:value="bindModal.selectIdList"
      :options="bindModal.sourceList"
      show-search
      mode="multiple"
      :filter-option="filterSourceListOption"
    />
  </a-modal>
</template>
<script setup>
import {
  ExclamationCircleOutlined,
  SettingOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, onUnmounted, createVNode, reactive, h } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listBindServiceSourceRequest,
  listServiceStatusRequest,
  listStatusActionsRequest,
  doStatusActionRequest,
  listAllServiceSourceRequest,
  bindAppAndServiceSourceRequest
} from "@/api/app/serviceApi";
import { message, Modal } from "ant-design-vue";
import { useAppStore } from "@/pinia/appStore";
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
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "操作",
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
    title: "主机",
    dataIndex: "host",
    key: "host"
  },
  {
    title: "状态",
    dataIndex: "status",
    key: "status"
  },
  {
    title: "CPU百分比",
    dataIndex: "cpuPercent",
    key: "cpuPercent"
  },
  {
    title: "内存百分比",
    dataIndex: "memPercent",
    key: "memPercent"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
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
// 操作列表
const actionList = ref([]);
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
  if (actionList.value.length === 0) {
    listStatusActionsRequest(data.bindId).then(res => {
      actionList.value = res.data;
    });
  }
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
// 执行自定义操作
const doAction = (item, action) => {
  Modal.confirm({
    title: `你确定要${action}${item.id}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      doStatusActionRequest({
        serviceId: item.id,
        bindId: selectedSource.value.bindId,
        action: action
      }).then(() => {
        message.success("操作成功");
        listStatus();
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
    message.success("操作成功");
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
.action-ul > li {
  display: inline;
}
.action-ul > li:hover {
  color: #1677ff;
  cursor: pointer;
}
.action-ul > li + li {
  margin-left: 6px;
}
</style>