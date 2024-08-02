<template>
  <div style="padding:10px;" v-show="!selectedSource.showServices">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-button
          type="primary"
          :icon="h(SettingOutlined)"
          @click="showBindDiscoverySourceModal"
          v-if="appStore.perm?.canManageDiscoverySource"
        >管理注册中心来源绑定</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <div>
      <ZTable :columns="sourceColumns" :dataSource="sourceDataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <span v-else @click="listAndShowService(dataItem)" class="check-btn">查看</span>
        </template>
      </ZTable>
    </div>
  </div>
  <div style="padding:10px;" v-show="selectedSource.showServices">
    <div style="margin-bottom:10px">
      <span class="check-btn" @click="selectedSource.showServices = false">返回集群选择</span>
    </div>
    <div style="margin-bottom:10px;">
      <span style="font-weight:bold">{{selectedSource.name}}</span>
      <a-tag color="orange" style="margin-left:10px;font-weight:bold">{{selectedSource.env}}</a-tag>
      <span class="check-btn" style="margin-left:10px" @click="listService">
        <ReloadOutlined />
        <span style="margin-left: 6px">刷新</span>
      </span>
    </div>
    <div>
      <ZTable :columns="serviceColumns" :dataSource="serviceDataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <template v-if="dataIndex === 'up'">
            <CheckCircleFilled style="color:green" v-if="dataItem[dataIndex]" />
            <CloseCircleFilled style="color:red" v-else />
          </template>
          <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <li @click="deregisterService(dataItem)" v-if="dataItem['up']">
                    <CloseOutlined />
                    <span style="margin-left:4px">下线服务</span>
                  </li>
                  <template v-else>
                    <li @click="reRegisterService(dataItem)">
                      <UploadOutlined />
                      <span style="margin-left:4px">上线服务</span>
                    </li>
                    <li @click="deleteDownService(dataItem)">
                      <CloseOutlined />
                      <span style="margin-left:4px">删除服务</span>
                    </li>
                  </template>
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
  EllipsisOutlined,
  ReloadOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  CloseOutlined,
  UploadOutlined,
  ExclamationCircleOutlined,
  SettingOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, h, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listBindDiscoverySourceRequest,
  listDiscoveryServiceRequest,
  deregisterServiceRequest,
  reRegisterServiceRequest,
  deleteDownServiceRequest,
  listAllDiscoverySourceRequest,
  bindAppAndDiscoverySourceRequest
} from "@/api/app/discoveryApi";
import EnvSelector from "@/components/app/EnvSelector";
import { message, Modal } from "ant-design-vue";
import { useAppStore } from "@/pinia/appStore";
const appStore = useAppStore();
const selectedEnv = ref("");
const route = useRoute();
const router = useRouter();
const bindModal = reactive({
  open: false,
  selectIdList: [],
  sourceList: []
});
// 来源列表
const sourceDataSource = ref([]);
// 服务列表
const serviceDataSource = ref([]);
// 选择的注册中心
const selectedSource = reactive({
  showServices: false,
  id: 0,
  name: "",
  bindId: 0,
  env: ""
});
// 注册中心来源数据项
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
// 注册服务数据项
const serviceColumns = [
  {
    title: "协议",
    dataIndex: "protocol",
    key: "protocol"
  },
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "host",
    dataIndex: "host",
    key: "host"
  },
  {
    title: "端口",
    dataIndex: "port",
    key: "port"
  },
  {
    title: "权重",
    dataIndex: "weight",
    key: "weight"
  },
  {
    title: "版本",
    dataIndex: "version",
    key: "version"
  },
  {
    title: "地域",
    dataIndex: "region",
    key: "region"
  },
  {
    title: "地区",
    dataIndex: "zone",
    key: "zone"
  },
  {
    title: "是否在线",
    dataIndex: "up",
    key: "up"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];
// 获取注册中心列表
const listDiscoverySource = () => {
  listBindDiscoverySourceRequest({
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
// 获取注册服务
const listService = () => {
  listDiscoveryServiceRequest(selectedSource.bindId).then(res => {
    serviceDataSource.value = res.data.map((item, index) => {
      return {
        key: index,
        ...item
      };
    });
  });
};
// 获取注册服务并切换界面
const listAndShowService = item => {
  selectedSource.name = item.name;
  selectedSource.showServices = true;
  selectedSource.id = item.id;
  selectedSource.bindId = item.bindId;
  selectedSource.env = item.env;
  listService();
};
// 选择环境变化
const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/discoveryService/list/${e.newVal}`
  );
  listDiscoverySource();
};
// 下线服务
const deregisterService = item => {
  Modal.confirm({
    title: `你确定要下线${item.host}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deregisterServiceRequest({
        bindId: selectedSource.bindId,
        instanceId: item.instanceId
      }).then(() => {
        message.success("下线成功");
        listService();
      });
    },
    onCancel() {}
  });
};
// 重新上线服务
const reRegisterService = item => {
  Modal.confirm({
    title: `你确定要上线${item.host}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      reRegisterServiceRequest({
        bindId: selectedSource.bindId,
        instanceId: item.instanceId
      }).then(() => {
        message.success("上线成功");
        listService();
      });
    },
    onCancel() {}
  });
};
// 删除下线服务
const deleteDownService = item => {
  Modal.confirm({
    title: `你确定要删除${item.host}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteDownServiceRequest({
        bindId: selectedSource.bindId,
        instanceId: item.instanceId
      }).then(() => {
        message.success("删除成功");
        listService();
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
const showBindDiscoverySourceModal = () => {
  if (!selectedEnv.value) {
    return;
  }
  listAllDiscoverySourceRequest(selectedEnv.value).then(res => {
    bindModal.sourceList = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    listBindDiscoverySourceRequest({
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
  bindAppAndDiscoverySourceRequest({
    appId: route.params.appId,
    sourceIdList: bindModal.selectIdList,
    env: selectedEnv.value
  }).then(() => {
    message.success("操作成功");
    bindModal.open = false;
    listDiscoverySource();
  });
};
</script>
<style scoped>
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>