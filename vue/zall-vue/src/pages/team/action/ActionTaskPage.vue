<template>
  <div style="padding:14px">
    <div style="margin-bottom:20px">
      <span class="header" @click="backToActionList">
        <arrow-left-outlined />
        <span style="margin-left:8px">工作流列表</span>
      </span>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" label="任务列表">
      <template #bodyCell="{dataIndex, dataItem}">
        <a-tag color="purple" v-if="dataIndex === 'pullRequest'">{{dataItem[dataIndex]}}</a-tag>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="click">
            <template #content>
              <ul class="op-list">
                <li>
                  <eye-outlined />
                  <span style="margin-left:4px">查看</span>
                </li>
                <li>
                  <close-outlined />
                  <span style="margin-left:4px">停止</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">...</div>
          </a-popover>
        </div>
      </template>
    </ZTable>
    <ZPagination :disableLastPage="true" @change="paginationChange"/>
  </div>
</template>
<script setup>
import ZPagination from "@/components/common/ZPagination";
import ZTable from "@/components/common/ZTable";
import { EyeOutlined, ArrowLeftOutlined, CloseOutlined } from "@ant-design/icons-vue";
import { ref } from "vue";
import { useRouter } from "vue-router";
const router = useRouter();
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
    title: "任务id",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "工作流名称",
    dataIndex: "age",
    key: "age"
  },
  {
    title: "触发方式",
    dataIndex: "pullRequest",
    key: "pullRequest"
  },
  {
    title: "任务状态",
    dataIndex: "pullRequest",
    key: "pullRequest"
  },
  {
    title: "操作人",
    dataIndex: "pullRequest",
    key: "pullRequest"
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

const backToActionList = () => {
  router.push("/team/action/list");
};

const paginationChange = (key) => {
    console.log(key);
}
</script>
<style scoped>
.header {
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 10px;
}
.header:hover {
  color: #1677ff;
}
</style>