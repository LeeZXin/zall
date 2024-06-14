<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建部署配置</a-button>
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
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <div class="op-icon" @click="deleteConfig(dataItem)">
              <a-tooltip placement="top">
                <template #title>
                  <span>Delete File</span>
                </template>
                <delete-outlined/>
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
import { listDeployConfigRequest, deleteDeployConfigRequest} from "@/api/app/deployApi";
import { useDeloyConfigStore } from "@/pinia/deployConfigStore";
import { Modal, message } from "ant-design-vue";
const deployConfigStore = useDeloyConfigStore();
const selectedEnv = ref("");
const envList = ref([]);
const route = useRoute();
const router = useRouter();
const dataSource = ref([]);

const columns = ref([
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

const listDeployConfig = () => {
  listDeployConfigRequest(
    {
      appId: route.params.appId,
      env: selectedEnv.value
    },
    selectedEnv.value
  ).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

const deleteConfig = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteDeployConfigRequest(item.id, item.env).then(() => {
        message.success("删除成功");
        listDeployConfig();
      });
    },
    onCancel() {}
  });
}

const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployConfig/create?env=${selectedEnv.value}`
  );
};

const gotoUpdatePage = item => {
  deployConfigStore.id = item.id;
  deployConfigStore.name = item.name;
  deployConfigStore.env = item.env;
  deployConfigStore.content = item.content;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployConfig/${item.id}/update`
  );
};

getEnvList();

watch(
  () => selectedEnv.value,
  newVal => {
    router.replace(
      `/team/${route.params.teamId}/app/${route.params.appId}/deployConfig/list/${newVal}`
    );
    listDeployConfig();
  }
);
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>