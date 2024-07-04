<template>
  <div style="padding:10px;height:100%" v-show="!showStatusList">
    <div style="margin-bottom:10px;" class="flex-end">
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="sourceColumns" :dataSource="sourceDataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <span v-else class="check-btn" @click="getStatusList(dataItem)">查看</span>
      </template>
    </ZTable>
  </div>
  <div style="padding:10px;height:100%" v-show="showStatusList">
    <div style="margin-bottom:20px;">
      <span class="check-btn" @click="backToSource">返回集群选择</span>
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
</template>
<script setup>
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listServiceSourceRequest,
  listServiceStatusRequest,
  listStatusActionsRequest,
  doStatusActionRequest
} from "@/api/app/serviceStatusApi";
import { message } from "ant-design-vue";
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const sourceDataSource = ref([]);
const showStatusList = ref(false);
const sourceColumns = ref([
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
]);
const statusColumns = ref([
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
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);
const statusDataSource = ref([]);
const statusInterval = ref(null);
const selectedSource = ref(null);
const listServiceSource = () => {
  listServiceSourceRequest(
    {
      appId: route.params.appId,
      env: selectedEnv.value
    },
    selectedEnv.value
  ).then(res => {
    sourceDataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
const actionList = ref([]);
const clearStatusInterval = () => {
  if (statusInterval.value) {
    clearInterval(statusInterval.value);
  }
};

const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/serviceStatus/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  listServiceSource();
};

const backToSource = () => {
  showStatusList.value = false;
  clearStatusInterval();
};

const getStatusList = data => {
  selectedSource.value = data;
  if (actionList.value.length === 0) {
    listStatusActionsRequest(data.id, data.env).then(res => {
      actionList.value = res.data;
    });
  }
  listServiceStatusRequest(data.id, data.env).then(res => {
    showStatusList.value = true;
    statusDataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
    refreshStatus();
  });
};

const listStatus = () => {
  listServiceStatusRequest(
    selectedSource.value.id,
    selectedSource.value.env
  ).then(res => {
    statusDataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

const refreshStatus = () => {
  clearStatusInterval();
  statusInterval.value = setInterval(listStatus, 5000);
};

const doAction = (item, action) => {
  doStatusActionRequest(
    {
      serviceId: item.id,
      sourceId: selectedSource.value.id,
      action: action
    },
    selectedSource.value.env
  ).then(() => {
    message.success("操作成功");
    listStatus();
  });
};

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