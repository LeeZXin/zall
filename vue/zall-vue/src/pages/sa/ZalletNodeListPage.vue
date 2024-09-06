<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchKey"
        placeholder="搜索名称"
        style="width:240px;margin-right:6px"
        @pressEnter="searchZalletNode"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
      <a-button type="primary" @click="gotoCreatePage" :icon="h(PlusOutlined)">创建节点</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteZalletNode(dataItem)">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Node</span>
              </template>
              <DeleteOutlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="gotoUpdatePage(dataItem)">
                  <EditOutlined />
                  <span style="margin-left:4px">编辑节点</span>
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
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listZalletNode()"
    />
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  DeleteOutlined
} from "@ant-design/icons-vue";
import {
  listZalletNodeRequest,
  deleteZalletNodeRequest
} from "@/api/zallet/zalletApi";
import { ref, h, reactive, createVNode } from "vue";
import { useRouter } from "vue-router";
import { Modal, message } from "ant-design-vue";
import { useZalletNodeStore } from "@/pinia/zalletNodeStore";
const zalletNodeStore = useZalletNodeStore();
// 分页数据
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 搜索关键词
const searchKey = ref("");
const router = useRouter();
// 数据
const dataSource = ref([]);
// 数据项
const columns = [{
    title: "nodeId",
    dataIndex: "nodeId",
    key: "nodeId"
  },
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "代理host",
    dataIndex: "agentHost",
    key: "agentHost"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];
// 跳转创建节点界面
const gotoCreatePage = () => {
  router.push(`/sa/zalletNode/create`);
};
// 跳转编辑节点界面
const gotoUpdatePage = item => {
  zalletNodeStore.id = item.id;
  zalletNodeStore.nodeId = item.nodeId;
  zalletNodeStore.name = item.name;
  zalletNodeStore.agentHost = item.agentHost;
  zalletNodeStore.agentToken = item.agentToken;
  router.push(`/sa/zalletNode/${item.id}/update`);
};
// 获取节点列表
const listZalletNode = () => {
  listZalletNodeRequest({
    pageNum: dataPage.current,
    name: searchKey.value
  }).then(res => {
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
    dataPage.totalCount = res.totalCount;
  });
};
// 删除节点
const deleteZalletNode = item => {
  Modal.confirm({
    title: `是否要删除${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteZalletNodeRequest(item.id).then(() => {
        message.success("操作成功");
        searchZalletNode();
      });
    }
  });
};
const searchZalletNode = () => {
    dataPage.current = 1;
    listZalletNode();
}
listZalletNode();
</script>
<style scoped>
</style>