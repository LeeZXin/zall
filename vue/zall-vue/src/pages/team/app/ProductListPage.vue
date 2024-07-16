<template>
  <div style="padding:10px;height:100%" v-show="!showStatusList">
    <div style="margin-bottom:10px;" class="flex-end">
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div class="op-icon" @click="deleteProduct(dataItem)" v-else>
          <a-tooltip placement="top">
            <template #title>
              <span>Delete File</span>
            </template>
            <delete-outlined />
          </a-tooltip>
        </div>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import {
  ExclamationCircleOutlined,
  DeleteOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import { listProductRequest, deleteProductRequest } from "@/api/app/productApi";
import { message, Modal } from "ant-design-vue";
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
const columns = [
  {
    title: "制品号",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "创建人",
    dataIndex: "creator",
    key: "creator"
  },
  {
    title: "创建时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];
const listProduct = () => {
  listProductRequest({
    appId: route.params.appId,
    env: selectedEnv.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

const deleteProduct = item => {
  Modal.confirm({
    title: `你确定要删除${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteProductRequest(item.id).then(() => {
        message.success("删除成功");
        listProduct();
      });
    },
    onCancel() {}
  });
};

const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/product/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  listProduct();
};
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