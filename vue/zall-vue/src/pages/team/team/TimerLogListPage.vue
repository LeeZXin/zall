<template>
  <div style="padding:10px">
    <div style="margin-bottom:10px;font-size: 14px;" class="flex-between">
      <div>
        <span style="font-weight:bold;font-size:18px">{{timerTaskStore.name}}</span>
        <a-tag style="margin-left:4px" color="orange">{{timerTaskStore.env}}</a-tag>
        <a-button
          type="primary"
          :icon="h(ReloadOutlined)"
          @click="reloadLog"
        >{{t('timerTask.reloadLog')}}</a-button>
      </div>
      <div>
        <span>{{t('timerTask.searchMonthly')}}:</span>
        <a-date-picker
          v-model:value="monthVal"
          style="margin-left:6px"
          @change="listLog"
          picker="month"
          :allowClear="false"
        />
      </div>
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" style="margin-top:0px" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex === 'isSuccess'">
          <span v-if="dataItem[dataIndex]">
            <CheckCircleFilled style="color:green" />
            <span style="margin-left:4px">{{t('timerTask.successful')}}</span>
          </span>
          <span v-else>
            <CloseCircleFilled style="color:red" />
            <span style="margin-left:4px">{{t('timerTask.failed')}}</span>
          </span>
        </template>
        <template v-else-if="dataIndex === 'triggerType'">
          <a-tag
            :color="triggerTypeMap[dataItem[dataIndex]].color"
          >{{t(triggerTypeMap[dataItem[dataIndex]].type)}}</a-tag>
        </template>
        <template v-else-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="showErrLog(dataItem)" v-if="!dataItem['isSuccess']">
                  <ControlOutlined />
                  <span style="margin-left:4px">{{t('timerTask.viewErrLog')}}</span>
                </li>
                <li @click="showTask(dataItem)">
                  <SettingOutlined />
                  <span style="margin-left:4px">{{t('timerTask.viewTaskCfg')}}</span>
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
      v-model:current="dataPage.current"
      :total="dataPage.totalCount"
      show-less-items
      :pageSize="dataPage.pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listLog()"
    />
    <a-modal v-model:open="errLogModal.open" :title="t('timerTask.errLog')" :footer="null">
      <Codemirror
        v-model="errLogModal.log"
        :style="codemirrorStyle"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
    <a-modal v-model:open="taskModal.open" :title="t('timerTask.taskCfg')" :footer="null">
      <ul class="task-ul">
        <li>
          <div class="item-name">{{t('timerTask.taskType')}}</div>
          <div class="item-value">{{taskModal.task?.taskType}}</div>
        </li>
        <template v-if="taskModal.task?.taskType === 'http'">
          <li>
            <div class="item-name">Url</div>
            <div class="item-value">{{taskModal.task?.httpTask.url}}</div>
          </li>
          <li>
            <div class="item-name">Method</div>
            <div class="item-value">{{taskModal.task?.httpTask.method}}</div>
          </li>
          <li>
            <div class="item-name">Headers</div>
            <div class="item-value">{{taskModal.task?.httpTask.headers}}</div>
          </li>
          <li>
            <div class="item-name">Body</div>
            <div class="item-value">{{taskModal.task?.httpTask.bodyStr}}</div>
          </li>
          <li>
            <div class="item-name">Content-Type</div>
            <div class="item-value">{{taskModal.task?.httpTask.contentType}}</div>
          </li>
        </template>
      </ul>
    </a-modal>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref, reactive, h } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  EllipsisOutlined,
  ControlOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  SettingOutlined,
  ReloadOutlined
} from "@ant-design/icons-vue";
import { listTimerTaskLogRequest } from "@/api/team/timerApi";
import dayjs from "dayjs";
import { Codemirror } from "vue-codemirror";
import { json } from "@codemirror/lang-json";
import { oneDark } from "@codemirror/theme-one-dark";
import { useI18n } from "vue-i18n";
import { useTimerTaskStore } from "@/pinia/timerTaskStore";
import "moment/locale/zh-cn";
const router = useRouter();
const timerTaskStore = useTimerTaskStore();
const { t } = useI18n();
const extensions = [json(), oneDark];
const codemirrorStyle = { height: "380px", width: "100%" };
// 选择月份
const monthVal = ref(dayjs());
const route = useRoute();
// 列表数据
const dataSource = ref([]);
// 分页数据
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
// 错误日志modal
const errLogModal = reactive({
  log: "",
  open: false
});
// 任务配置modal
const taskModal = reactive({
  open: false,
  task: null
});
// 表项
const columns = [
  {
    i18nTitle: "timerTask.logColumns.triggerBy",
    dataIndex: "triggerBy",
    key: "triggerBy"
  },
  {
    i18nTitle: "timerTask.logColumns.triggerType",
    dataIndex: "triggerType",
    key: "triggerType"
  },
  {
    i18nTitle: "timerTask.logColumns.isSuccess",
    dataIndex: "isSuccess",
    key: "isSuccess"
  },
  {
    i18nTitle: "timerTask.logColumns.created",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "timerTask.logColumns.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 触发类型i18n map
const triggerTypeMap = {
  1: {
    type: "timerTask.autoTriggerType",
    color: "orange"
  },
  2: {
    type: "timerTask.manualTriggerType",
    color: "purple"
  }
};
// 展示错误日志
const showErrLog = item => {
  errLogModal.open = true;
  errLogModal.log = item.errLog;
};
// 展示任务配置
const showTask = item => {
  taskModal.open = true;
  taskModal.task = item.task;
};
// 日志列表
const listLog = () => {
  listTimerTaskLogRequest({
    id: route.params.timerId,
    pageNum: dataPage.current,
    month: monthVal.value.format("YYYY-MM")
  }).then(res => {
    dataPage.totalCount = res.totalCount;
    dataSource.value = res.data.map((item, index) => {
      return {
        key: index,
        ...item
      };
    });
  });
};
// 重新加载日志
const reloadLog = () => {
  dataPage.current = 1;
  listLog();
};
if (timerTaskStore.id === 0) {
  router.push(`/team/${route.params.teamId}/timer/list`);
} else {
  listLog();
}
</script>
<style scoped>
.task-ul {
  width: 100%;
  padding: 10px 0;
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