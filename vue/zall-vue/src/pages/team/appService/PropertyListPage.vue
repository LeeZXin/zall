<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px">
      <a-button type="primary">创建配置</a-button>
    </div>
    <div class="body">
      <ZTable :columns="columns" :dataSource="dataSource">
        <template #bodyCell="{dataIndex, dataItem}">
          <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
          <div v-else>
            <div class="op-icon" @click="deleteApp">
              <a-tooltip placement="top">
                <template #title>
                  <span>Delete File</span>
                </template>
                <delete-outlined />
              </a-tooltip>
            </div>
            <a-popover placement="bottomRight" trigger="click">
              <template #content>
                <ul class="op-list">
                  <li>
                    <play-circle-outlined />
                    <span style="margin-left:4px">编辑发布</span>
                  </li>
                  <li>
                    <file-text-outlined />
                    <span style="margin-left:4px">版本对比</span>
                  </li>
                  <li>
                    <eye-outlined />
                    <span style="margin-left:4px">发布历史</span>
                  </li>
                </ul>
              </template>
              <div class="op-icon">...</div>
            </a-popover>
          </div>
        </template>
      </ZTable>
    </div>
  </div>
</template>
<script setup>
import { DeleteOutlined } from "@ant-design/icons-vue";
import ZTable from "@/components/common/ZTable";
import { ref } from "vue";
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
    title: "配置文件",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "当前版本",
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
</script>
<style scoped>
.body {
  width: 100%;
  height: calc(100% - 42px);
}
</style>