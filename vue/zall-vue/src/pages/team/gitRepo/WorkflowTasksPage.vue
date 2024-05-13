<template>
  <div style="padding:14px">
    <ZTable :columns="columns" :dataSource="taskList" label="任务列表" style="margin-top:10px" v-if="taskList.length > 0">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="click">
            <template #content>
              <ul class="op-list">
                <li>
                  <eye-outlined />
                  <span style="margin-left:4px" @click="gotoTaskDetail">查看</span>
                </li>
                <li>
                  <close-outlined />
                  <span style="margin-left:4px">停止</span>
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
    <ZNoData v-else/>
    <a-pagination
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      v-show="totalCount > pageSize"
      @change="()=>listBranch()"
    />
    
  </div>
</template>
<script setup>
import ZNoData from "@/components/common/ZNoData";
import ZTable from "@/components/common/ZTable";
import {
  EyeOutlined,
  CloseOutlined,
  EllipsisOutlined
} from "@ant-design/icons-vue";
import { ref } from "vue";
import { useRouter } from "vue-router";
const totalCount = ref(0);
const pageSize = 10;
const currPage = ref(1);
const router = useRouter();
const taskList = ref([]);

const columns = ref([
  {
    title: "工作流名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "触发方式",
    dataIndex: "triggerType",
    key: "triggerType"
  },
  {
    title: "任务状态",
    dataIndex: "taskStatus",
    key: "taskStatus"
  },
  {
    title: "操作人",
    dataIndex: "operator",
    key: "operator"
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
]);
const gotoTaskDetail = () => {
  router.push("/");
};
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