<template>
  <div style="padding:14px">
    <a-date-picker v-model:value="dateVal" style="margin-bottom:10px" @change="pageLog" />
    <ZTable
      :columns="columns"
      :dataSource="dataSource"
      style="margin-top:0px"
      :label="timerTaskStore.name"
    >
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'isSuccess'">
          <CheckCircleFilled style="color:green" v-if="dataItem[dataIndex]" />
          <CloseCircleFilled style="color:red" v-else />
        </template>
        <template v-else-if="dataIndex === 'triggerType'">
          <span>{{t(triggerTypeMap[dataItem[dataIndex]])}}</span>
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="showErrLog(dataItem['errLog'])" v-if="!dataItem['isSuccess']">
                  <ControlOutlined />
                  <span style="margin-left:4px">查看错误日志</span>
                </li>
                <li @click="showTask(dataItem)">
                  <SettingOutlined />
                  <span style="margin-left:4px">查看任务配置</span>
                </li>
              </ul>
            </template>
            <div class="op-icon">
              <EllipsisOutlined />
            </div>
          </a-popover>
        </template>
      </template>
    </ZTable>
    <a-pagination
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>pageLog()"
    />
    <a-modal v-model:open="errLogModalOpen" title="错误日志" :footer="null">
      <Codemirror
        v-model="errLog"
        :style="codemirrorStyle"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
    <a-modal v-model:open="taskModalOpen" title="任务配置" :footer="null">
      <ul class="task-ul">
        <li>
          <div class="item-name">任务类型</div>
          <div class="item-value">{{task.taskType}}</div>
        </li>
        <template v-if="task.taskType === 'http'">
          <li>
            <div class="item-name">Url</div>
            <div class="item-value">{{task.httpTask.url}}</div>
          </li>
          <li>
            <div class="item-name">Method</div>
            <div class="item-value">{{task.httpTask.method}}</div>
          </li>
          <li>
            <div class="item-name">Headers</div>
            <div class="item-value">{{task.httpTask.headers}}</div>
          </li>
          <li>
            <div class="item-name">Body</div>
            <div class="item-value">{{task.httpTask.bodyStr}}</div>
          </li>
          <li>
            <div class="item-name">Content-Type</div>
            <div class="item-value">{{task.httpTask.contentType}}</div>
          </li>
          <li>
            <div class="item-name">Zones</div>
            <div class="item-value">{{task.httpTask.zones}}</div>
          </li>
        </template>
      </ul>
    </a-modal>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  EllipsisOutlined,
  ControlOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  SettingOutlined
} from "@ant-design/icons-vue";
import { pageTimerTaskLogRequest } from "@/api/team/timerApi";
import dayjs from "dayjs";
import { Codemirror } from "vue-codemirror";
import { json } from "@codemirror/lang-json";
import { oneDark } from "@codemirror/theme-one-dark";
import { useI18n } from "vue-i18n";
import { useTimerTaskStore } from "@/pinia/timerTaskStore";
const router = useRouter();
const timerTaskStore = useTimerTaskStore();
const { t } = useI18n();
const extensions = [json(), oneDark];
const codemirrorStyle = { height: "380px", width: "100%" };
const dateVal = ref(dayjs());
const route = useRoute();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const errLog = ref("");
const errLogModalOpen = ref(false);
const taskModalOpen = ref(false);
const task = ref();
const columns = ref([
  {
    title: "操作人",
    dataIndex: "triggerBy",
    key: "triggerBy"
  },
  {
    title: "操作方式",
    dataIndex: "triggerType",
    key: "triggerType"
  },
  {
    title: "是否成功",
    dataIndex: "isSuccess",
    key: "isSuccess"
  },
  {
    title: "触发时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);
const triggerTypeMap = {
  1: "timerTask.autoTriggerType",
  2: "timerTask.manualTriggerType"
};
const showErrLog = log => {
  errLogModalOpen.value = true;
  errLog.value = log;
};
const showTask = item => {
  taskModalOpen.value = true;
  task.value = item.task;
};
const pageLog = () => {
  pageTimerTaskLogRequest(
    {
      taskId: route.params.taskId,
      pageNum: currPage.value,
      dateStr: dateVal.value.format("YYYY-MM-DD")
    },
    timerTaskStore.env
  ).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map((item, index) => {
      return {
        key: index,
        ...item
      };
    });
  });
};
if (timerTaskStore.id === 0) {
  router.push(`/team/${route.params.teamId}/timerTask/list`);
} else {
  pageLog();
}
</script>
<style scoped>
.task-ul {
  width: 100%;
  padding-bottom: 20px;
}
.task-ul > li {
  width: 100%;
}
.task-ul > li + li {
  margin-top: 12px;
}
.item-name {
  font-size: 12px;
  margin-bottom: 4px;
}
.item-value {
  font-size: 14px;
  line-height: 18px;
  padding-left: 20px;
  min-height: 18px;
}
</style>