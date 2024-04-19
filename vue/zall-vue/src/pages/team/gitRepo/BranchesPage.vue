<template>
  <div style="padding:14px">
    <a-radio-group v-model:value="branchType">
      <a-radio-button value="all">所有</a-radio-button>
      <a-radio-button value="active">活跃</a-radio-button>
      <a-radio-button value="inactive">非活跃</a-radio-button>
    </a-radio-group>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <a-tag color="purple" v-if="dataIndex === 'pullRequest'">{{dataItem[dataIndex]}}</a-tag>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon">
            <a-tooltip placement="top">
              <template #title>
                <span>Delete Branch</span>
              </template>
              <delete-outlined />
            </a-tooltip>
          </div>
          <a-popover placement="bottomRight" trigger="click">
            <template #content>
              <ul class="op-list">
                <li>
                  <control-outlined />
                  <span style="margin-left:4px">查看所有的活动</span>
                </li>
                <li></li>
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
import { ref } from "vue";
import { DeleteOutlined, ControlOutlined } from "@ant-design/icons-vue";
const branchType = ref("all");
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
    title: "分支",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "最后更新时间",
    dataIndex: "age",
    key: "age"
  },
  {
    title: "合并请求",
    dataIndex: "pullRequest",
    key: "pullRequest"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);
</script>
<style scoped>
.header {
  line-height: 32px;
  font-size: 18px;
}
</style>