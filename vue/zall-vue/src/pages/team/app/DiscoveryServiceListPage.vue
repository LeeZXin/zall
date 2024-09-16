<template>
  <div style="padding:10px;" v-show="!selectedSource.showServices">
    <div style="margin-bottom:10px" class="flex-between">
      <div>
        <a-button
          type="primary"
          :icon="h(SettingOutlined)"
          @click="showBindDiscoverySourceModal"
          v-if="appStore.perm?.canManageDiscoverySource"
        >{{t('discoveryService.manageSource')}}</a-button>
      </div>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <div>
      <ZTable :columns="sourceColumns" :dataSource="sourceDataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <span
            v-else
            @click="listAndShowService(dataItem)"
            class="check-btn"
          >{{t('discoveryService.view')}}</span>
        </template>
      </ZTable>
    </div>
  </div>
  <div style="padding:10px;" v-show="selectedSource.showServices">
    <div style="margin-bottom:10px">
      <span class="check-btn" @click="selectedSource.showServices = false">
        <LeftOutlined />
        <span style="margin-left:4px">{{t('discoveryService.backToSelectSource')}}</span>
      </span>
    </div>
    <div style="margin-bottom:10px;">
      <span style="font-weight:bold">{{selectedSource.name}}</span>
      <a-tag color="orange" style="margin-left:10px;">{{selectedSource.env}}</a-tag>
      <span class="check-btn" style="margin-left:10px" @click="listService">
        <ReloadOutlined />
        <span style="margin-left: 6px">{{t('discoveryService.refresh')}}</span>
      </span>
    </div>
    <div>
      <ZTable :columns="serviceColumns" :dataSource="serviceDataSource" :scroll="{x:1300}">
        <template #bodyCell="{dataIndex, dataItem}">
          <template v-if="dataIndex === 'isDown'">
            <LoadingOutlined v-if="dataItem.loading" />
            <template v-else>
              <CloseCircleFilled style="color:red" v-if="dataItem[dataIndex]" />
              <CheckCircleFilled style="color:green" v-else />
            </template>
          </template>
          <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <a-popover placement="bottomRight" trigger="hover">
              <template #content>
                <ul class="op-list">
                  <template v-if="dataItem.isDown">
                    <li @click="markAsUp(dataItem)">
                      <UploadOutlined />
                      <span style="margin-left:4px">{{t('discoveryService.markAsUp')}}</span>
                    </li>
                  </template>
                  <li @click="markAsDown(dataItem)" v-else>
                    <CloseOutlined />
                    <span style="margin-left:4px">{{t('discoveryService.markAsDown')}}</span>
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
  <a-modal
    v-model:open="bindModal.open"
    :title="t('discoveryService.bindSource')"
    @ok="handleBindModalOk"
  >
    <div>
      <div style="font-size:12px;margin-bottom:3px">{{t('discoveryService.selectedEnv')}}</div>
      <div>{{selectedEnv}}</div>
    </div>
    <div style="margin-top: 10px">
      <div style="font-size:12px;margin-bottom:3px">{{t('discoveryService.source')}}</div>
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
  EllipsisOutlined,
  ReloadOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  UploadOutlined,
  ExclamationCircleOutlined,
  SettingOutlined,
  LeftOutlined,
  CloseOutlined,
  LoadingOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, h, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listBindDiscoverySourceRequest,
  listDiscoveryServiceRequest,
  markAsDownServiceRequest,
  markAsUpServiceRequest,
  listAllDiscoverySourceRequest,
  bindAppAndDiscoverySourceRequest
} from "@/api/app/discoveryApi";
import EnvSelector from "@/components/app/EnvSelector";
import { message, Modal } from "ant-design-vue";
import { useAppStore } from "@/pinia/appStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
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
    i18nTitle: "discoveryService.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "discoveryService.operation",
    dataIndex: "operation",
    key: "operation"
  }
];
// 注册服务数据项
const serviceColumns = [
  {
    i18nTitle: "discoveryService.serviceProtocol",
    dataIndex: "protocol",
    key: "protocol"
  },
  {
    i18nTitle: "discoveryService.serviceName",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "discoveryService.serviceHost",
    dataIndex: "host",
    key: "host"
  },
  {
    i18nTitle: "discoveryService.servicePort",
    dataIndex: "port",
    key: "port"
  },
  {
    i18nTitle: "discoveryService.serviceWeight",
    dataIndex: "weight",
    key: "weight"
  },
  {
    i18nTitle: "discoveryService.serviceVersion",
    dataIndex: "version",
    key: "version"
  },
  {
    i18nTitle: "discoveryService.serviceRegion",
    dataIndex: "region",
    key: "region"
  },
  {
    i18nTitle: "discoveryService.serviceZone",
    dataIndex: "zone",
    key: "zone"
  },
  {
    i18nTitle: "discoveryService.serviceUp",
    dataIndex: "isDown",
    key: "isDown"
  },
  {
    i18nTitle: "discoveryService.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
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
        loading: false,
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
const markAsDown = item => {
  Modal.confirm({
    title: `${t("discoveryService.confirmMarkAsDown")} ${item.host}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      item.loading = true;
      markAsDownServiceRequest({
        bindId: selectedSource.bindId,
        instanceId: item.instanceId
      })
        .then(() => {
          item.loading = false;
          item.isDown = true;
          message.success(t("operationSuccess"));
        })
        .catch(() => {
          item.loading = false;
        });
    },
    onCancel() {}
  });
};
// 重新上线服务
const markAsUp = item => {
  Modal.confirm({
    title: `${t("discoveryService.confirmMarkAsUp")} ${item.host}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      item.loading = true;
      markAsUpServiceRequest({
        bindId: selectedSource.bindId,
        instanceId: item.instanceId
      })
        .then(() => {
          message.success(t("operationSuccess"));
          item.loading = false;
          item.isDown = false;
        })
        .catch(() => {
          item.loading = false;
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
    message.success(t("operationSuccess"));
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