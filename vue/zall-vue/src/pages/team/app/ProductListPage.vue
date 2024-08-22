<template>
  <div style="padding:10px;height:100%">
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
    <a-pagination
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listProduct()"
    />
  </div>
</template>
<script setup>
import {
  ExclamationCircleOutlined,
  DeleteOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import { listProductRequest, deleteProductRequest } from "@/api/app/productApi";
import { message, Modal } from "ant-design-vue";
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
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
  listProductRequest(
    {
      appId: route.params.appId,
      env: selectedEnv.value,
      pageNum: dataPage.current
    },
    selectedEnv.value
  ).then(res => {
    dataPage.totalCount = res.totalCount;
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
      deleteProductRequest(item.id, item.env).then(() => {
        message.success("删除成功");
        dataPage.current = 1;
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