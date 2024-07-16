<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button type="primary" :icon="h(PlusOutlined)" @click="gotoCreatePage">创建发布计划</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <PLanStatusTag v-if="dataIndex === 'planStatus'" :status="dataItem[dataIndex]" />
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li
                  @click="closePlan(dataItem)"
                  v-if="dataItem['planStatus'] === 1 || dataItem['planStatus'] === 2"
                >
                  <close-outlined />
                  <span style="margin-left:8px">关闭发布计划</span>
                </li>
                <li @click="viewPlan(dataItem)">
                  <eye-outlined />
                  <span style="margin-left:8px">查看发布计划</span>
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
      v-model:current="currPage"
      :total="totalCount"
      show-less-items
      :pageSize="pageSize"
      style="margin-top:10px"
      :hideOnSinglePage="true"
      :showSizeChanger="false"
      @change="()=>listDeployPlan()"
    />
  </div>
</template>
<script setup>
import {
  CloseOutlined,
  PlusOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined,
  EyeOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listDeployPlanRequest,
  closeDeployPlanRequest
} from "@/api/app/deployPlanApi";
import { Modal, message } from "ant-design-vue";
import PLanStatusTag from "@/components/app/PlanStatusTag";
import { useDeloyPlanStore } from "@/pinia/deployPlanStore";
const planStore = useDeloyPlanStore();
const route = useRoute();
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const columns = [
  {
    title: "名称",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "制品号",
    dataIndex: "productVersion",
    key: "productVersion"
  },
  {
    title: "部署流水线",
    dataIndex: "pipelineName",
    key: "pipelineName"
  },
  {
    title: "状态",
    dataIndex: "planStatus",
    key: "planStatus"
  },
  {
    title: "创建时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "创建人",
    dataIndex: "creator",
    key: "creator"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
];

const closePlan = item => {
  Modal.confirm({
    title: `你确定要关闭${item.name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      closeDeployPlanRequest(item.id).then(() => {
        message.success("关闭成功");
        listDeployPlan();
      });
    },
    onCancel() {}
  });
};

const listDeployPlan = () => {
  listDeployPlanRequest({
    appId: route.params.appId,
    pageNum: currPage.value,
    env: selectedEnv.value
  }).then(res => {
    totalCount.value = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};

const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/create?env=${selectedEnv.value}`
  );
};

const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  currPage.value = 1;
  listDeployPlan();
};

const viewPlan = item => {
  planStore.id = item.id;
  planStore.env = item.env;
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/${item.id}/view`
  );
};
</script>
<style scoped>
</style>