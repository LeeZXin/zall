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
          <span v-if="dataIndex === 'isEnabled'">{{dataItem[dataIndex]?'已启用':'已关闭'}}</span>
          <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <div class="op-icon" @click="deleteApp">
              <a-tooltip placement="top">
                <template #title>
                  <span>Delete File</span>
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
                  <li>
                    <eye-outlined />
                    <span style="margin-left:4px">发布历史</span>
                  </li>
                  <li
                    v-if="dataItem['isEnabled'] === false"
                    @click="enableOrDisableConfig(dataItem, true)"
                  >
                    <check-outlined />
                    <span style="margin-left:4px">启用配置</span>
                  </li>
                  <li
                    v-else-if="dataItem['isEnabled'] === true"
                    @click="enableOrDisableConfig(dataItem, false)"
                  >
                    <close-outlined />
                    <span style="margin-left:4px">关闭配置</span>
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
  EyeOutlined,
  PlusOutlined,
  EllipsisOutlined,
  CheckOutlined,
  CloseOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, h, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getEnvCfgRequest } from "@/api/cfg/cfgApi";
import {
  listDeployConfigRequest,
  enableDeployConfigRequest,
  disableDeployConfigRequest
} from "@/api/app/deployApi";
import { useDeloyConfigStore } from "@/pinia/deployConfigStore";
import { message } from "ant-design-vue";
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
        key: item.name,
        ...item
      };
    });
  });
};

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

const enableOrDisableConfig = (item, result) => {
  let request = result ? enableDeployConfigRequest : disableDeployConfigRequest;
  request(item.id, item.env).then(() => {
    message.success("操作成功");
    listDeployConfig();
  });
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