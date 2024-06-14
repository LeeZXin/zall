<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建探针配置</a-button>
      <div>
        <span style="margin-right:6px">环境:</span>
        <a-select
          style="width: 200px"
          placeholder="选择环境"
          v-model:value="selectedEnv"
          :options="envList"
        />
      </div>
    </div>
    <div class="body">
      <ZTable :columns="columns" :dataSource="dataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex === 'isEnabled'">
            <a-switch :checked="dataItem[dataIndex]" @click="enableOrDisableProbe(dataItem)" />
          </span>
          <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <div class="op-icon" @click="deleteProbe(dataItem)">
              <a-tooltip placement="top">
                <template #title>
                  <span>Delete Probe</span>
                </template>
                <delete-outlined />
              </a-tooltip>
            </div>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="gotoUpdatePage(dataItem)">
                    <edit-outlined />
                    <span style="margin-left:4px">编辑配置</span>
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
      <a-pagination
        v-model:current="currPage"
        :total="totalCount"
        show-less-items
        :pageSize="pageSize"
        style="margin-top:10px"
        :hideOnSinglePage="true"
        :showSizeChanger="false"
        @change="()=>listProbe()"
      />
    </div>
  </div>
</template>
<script setup>
import {
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, h, watch, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  listProbeRequest,
  deleteProbeRequest,
  enableProbeRequest,
  disableProbeRequest
} from "@/api/app/probeApi";
import { useProbeStore } from "@/pinia/probeStore";
import { Modal, message } from "ant-design-vue";
const probeStore = useProbeStore();
const selectedEnv = ref("");
const envList = ref([]);
const route = useRoute();
const router = useRouter();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);

const columns = ref([
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "是否启用",
    dataIndex: "isEnabled",
    key: "isEnabled"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const getEnvList = () => {
  getEnvCfgRequest().then(res => {
    envList.value = res.data.map(item => {
      return {
        value: item,
        label: item
      };
    });
    if (route.params.env && res.data?.includes(route.params.env)) {
      selectedEnv.value = route.params.env;
    } else if (res.data.length > 0) {
      selectedEnv.value = res.data[0];
    }
  });
};

const deleteProbe = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteProbeRequest(item.id, item.env).then(() => {
        if (totalCount.value - 1 <= (currPage.value - 1) * pageSize) {
          currPage.value -= 1;
        }
        message.success("删除成功");
        listProbe();
      });
    },
    onCancel() {}
  });
};

const listProbe = () => {
  listProbeRequest(
    {
      appId: route.params.appId,
      pageNum: currPage.value,
      env: selectedEnv.value
    },
    selectedEnv.value
  ).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/probe/create?env=${selectedEnv.value}`
  );
};

const gotoUpdatePage = item => {
  probeStore.id = item.id;
  probeStore.name = item.name;
  probeStore.env = item.env;
  probeStore.config = item.config;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/probe/${item.id}/update`
  );
};

const enableOrDisableProbe = item => {
  if (item.isEnabled === true) {
    disableProbeRequest(item.id, item.env).then(() => {
      message.success("关闭成功");
      item.isEnabled = false;
    });
  } else if (item.isEnabled === false) {
    enableProbeRequest(item.id, item.env).then(() => {
      message.success("启动成功");
      item.isEnabled = true;
    });
  }
};
getEnvList();

watch(
  () => selectedEnv.value,
  newVal => {
    router.replace(
      `/team/${route.params.teamId}/app/${route.params.appId}/probe/list/${newVal}`
    );
    listProbe();
  }
);
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>