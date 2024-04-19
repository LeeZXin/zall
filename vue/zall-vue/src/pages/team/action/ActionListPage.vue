<template>
  <div style="padding:14px">
    <div style="margin-bottom:10px">
      <a-input
        v-model:value="searchApp"
        placeholder="搜索工作流"
        style="width:240px;margin-right:10px"
        @change="searchChange"
      />
      <a-button type="primary" @click="gotoCreatePage">创建工作流</a-button>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <a-tag color="purple" v-if="dataIndex === 'pullRequest'">{{dataItem[dataIndex]}}</a-tag>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteApp">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Action</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="click">
            <template #content>
              <ul class="op-list">
                <li>
                  <play-circle-outlined />
                  <span style="margin-left:4px">执行</span>
                </li>
                <li>
                  <file-text-outlined />
                  <span style="margin-left:4px">日志</span>
                </li>
                <li>
                  <eye-outlined />
                  <span style="margin-left:4px">查看</span>
                </li>
                <li>
                  <edit-outlined />
                  <span style="margin-left:4px">编辑</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">...</div>
          </a-popover>
        </div>
      </template>
    </ZTable>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import {
  DeleteOutlined,
  ExclamationCircleOutlined,
  PlayCircleOutlined,
  EditOutlined,
  EyeOutlined,
  FileTextOutlined
} from "@ant-design/icons-vue";
import { ref, createVNode } from "vue";
import { useRouter } from "vue-router";
import { Modal } from "ant-design-vue";
const router = useRouter();
const searchApp = ref("");
const dataSource = ref([
  {
    key: "1",
    name: "胡彦斌",
    age: 32,
    pullRequest: "西湖区湖底公园1号"
  },
  {
    key: "2",
    name: "胡彦祖",
    age: 42,
    pullRequest: "西湖区湖底公园1号"
  }
]);

const columns = ref([
  {
    title: "AppId",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "名称",
    dataIndex: "age",
    key: "age"
  },
  {
    title: "创建时间",
    dataIndex: "pullRequest",
    key: "pullRequest"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);

const gotoCreatePage = () => {
  router.push("/team/action/create");
};

const deleteApp = () => {
  Modal.confirm({
    title: "你确定要删除?",
    icon: createVNode(ExclamationCircleOutlined),
    content:
      "When clicked the OK button, this dialog will be closed after 1 second",
    okText: "fuc",
    cancelText: "nimba",
    async onOk() {
      try {
        return await new Promise((resolve, reject) => {
          setTimeout(Math.random() > 0.5 ? resolve : reject, 1000);
        });
      } catch {
        return console.log("Oops errors!");
      }
    },
    onCancel() {}
  });
};
</script>
<style scoped>
</style>