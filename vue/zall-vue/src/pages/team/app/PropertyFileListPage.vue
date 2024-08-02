<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建配置文件</a-button>
        <a-button
          type="primary"
          :icon="h(SettingOutlined)"
          @click="showBindPropertySourceModal"
          style="margin-left:6px"
          v-if="appStore.perm?.canManagePropertySource"
        >管理配置来源绑定</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteFile(dataItem)">
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
                <li @click="gotoHistoryListPage(dataItem)">
                  <file-text-outlined />
                  <span style="margin-left:4px">版本列表</span>
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
  DeleteOutlined,
  FileTextOutlined,
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  SettingOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listPropertyFileRequest,
  deletePropertyFileRequest,
  listAllPropertySourceRequest,
  listBindPropertySourceRequest,
  bindAppAndPropertySourceRequest
} from "@/api/app/propertyApi";
import { usePropertyFileStore } from "@/pinia/propertyFileStore";
import EnvSelector from "@/components/app/EnvSelector";
import { Modal, message } from "ant-design-vue";
import { useAppStore } from "@/pinia/appStore";
const appStore = useAppStore();
const propertyFileStore = usePropertyFileStore();
const bindModal = reactive({
  open: false,
  selectIdList: [],
  sourceList: []
});
// 当前环境
const selectedEnv = ref("");
const route = useRoute();
const router = useRouter();
// 表格数据
const dataSource = ref([]);
// 数据项
const columns = [
  {
    title: "配置文件",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

// 获取配置文件列表
const listPropertyFile = () => {
  listPropertyFileRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.name,
        ...item
      };
    });
  });
};
// 跳转创建配置文件页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/create?env=${selectedEnv.value}`
  );
};
// 跳转历史版本页面
const gotoHistoryListPage = item => {
  propertyFileStore.id = item.id;
  propertyFileStore.name = item.name;
  propertyFileStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/${item.id}/history/list`
  );
};
// 删除配置文件
const deleteFile = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deletePropertyFileRequest(item.id).then(() => {
        message.success("删除成功");
        listPropertyFile();
      });
    },
    onCancel() {}
  });
};
// 选择环境变动
const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/propertyFile/list/${e.newVal}`
  );
  listPropertyFile();
};
// 配置下拉框搜索过滤
const filterSourceListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};
// 展示绑定配置来源modal
const showBindPropertySourceModal = () => {
  if (!selectedEnv.value) {
    return;
  }
  listAllPropertySourceRequest(selectedEnv.value).then(res => {
    bindModal.sourceList = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    listBindPropertySourceRequest({
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
  bindAppAndPropertySourceRequest({
    appId: route.params.appId,
    sourceIdList: bindModal.selectIdList,
    env: selectedEnv.value
  }).then(() => {
    message.success("操作成功");
    bindModal.open = false;
  });
};
</script>
<style scoped>
</style>