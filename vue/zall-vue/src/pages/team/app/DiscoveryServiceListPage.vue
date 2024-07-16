<template>
  <div style="padding:10px;" v-show="!showServices">
    <div style="margin-bottom:10px" class="flex-right">
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <div>
      <ZTable :columns="sourceColumns" :dataSource="sourceDataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <span v-else @click="getServices(dataItem)" class="check-btn">查看</span>
        </template>
      </ZTable>
    </div>
  </div>
  <div style="padding:10px;" v-show="showServices">
    <div style="margin-bottom:10px">
      <span class="check-btn" @click="showServices = false">返回集群选择</span>
    </div>
    <div style="margin-bottom:10px;">
      <span style="font-weight:bold">{{selectedSourceName}}</span>
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
</template>
<script setup>
import {
  EllipsisOutlined,
  ReloadOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  CloseOutlined,
  UploadOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listSimpleDiscoverySourceRequest,
  listDiscoveryServiceRequest,
  deregisterServiceRequest,
  reRegisterServiceRequest,
  deleteDownServiceRequest
} from "@/api/app/discoveryApi";
import EnvSelector from "@/components/app/EnvSelector";
import { message, Modal } from "ant-design-vue";
const selectedEnv = ref("");
const route = useRoute();
const router = useRouter();
const sourceDataSource = ref([]);
const serviceDataSource = ref([]);
const selectedSourceId = ref(0);
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
const selectedSourceName = ref("");
const showServices = ref(false);
const listDiscoverySource = () => {
  listSimpleDiscoverySourceRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    sourceDataSource.value = res.data.map(item => {
      return {
        key: item.name,
        ...item
      };
    });
  });
};

const listService = () => {
  listDiscoveryServiceRequest(selectedSourceId.value).then(res => {
    serviceDataSource.value = res.data.map((item, index) => {
      return {
        key: index,
        ...item
      };
    });
  });
};

const getServices = item => {
  selectedSourceName.value = item.name;
  showServices.value = true;
  selectedSourceId.value = item.id;
  listService();
};

const onEnvChange = e => {
  selectedEnv.value = e.newVal;
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/discoveryService/list/${e.newVal}`
  );
  listDiscoverySource();
};

const deregisterService = item => {
  Modal.confirm({
    title: `你确定要下线${item.host}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deregisterServiceRequest({
        sourceId: selectedSourceId.value,
        instanceId: item.instanceId
      }).then(() => {
        message.success("下线成功");
        listService();
      });
    },
    onCancel() {}
  });
};

const reRegisterService = item => {
  Modal.confirm({
    title: `你确定要上线${item.host}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      reRegisterServiceRequest({
        sourceId: selectedSourceId.value,
        instanceId: item.instanceId
      }).then(() => {
        message.success("上线成功");
        listService();
      });
    },
    onCancel() {}
  });
};

const deleteDownService = item => {
  Modal.confirm({
    title: `你确定要删除${item.host}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteDownServiceRequest({
        sourceId: selectedSourceId.value,
        instanceId: item.instanceId
      }).then(() => {
        message.success("删除成功");
        listService();
      });
    },
    onCancel() {}
  });
};
</script>
<style scoped>
.check-btn:hover {
  color: #1677ff;
  cursor: pointer;
}
</style>