<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px" class="flex-between">
      <a-button
        type="primary"
        :icon="h(PlusOutlined)"
        @click="gotoCreatePage"
      >{{t('deployPlan.createPlan')}}</a-button>
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <PLanStatusTag v-if="dataIndex === 'planStatus'" :status="dataItem[dataIndex]" />
        <div v-else-if="dataIndex === 'creator'" class="flex-center">
          <ZAvatar
            :url="dataItem.creator?.avatarUrl"
            :name="dataItem.creator?.name"
            :account="dataItem.creator?.account"
            :showName="true"
          />
        </div>
        <span v-else-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li
                  @click="closePlan(dataItem)"
                  v-if="dataItem['planStatus'] === 1 || dataItem['planStatus'] === 2"
                >
                  <CloseOutlined />
                  <span style="margin-left:8px">{{t('deployPlan.closePlan')}}</span>
                </li>
                <li @click="viewPlan(dataItem)">
                  <EyeOutlined />
                  <span style="margin-left:8px">{{t('deployPlan.viewPlan')}}</span>
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
import ZAvatar from "@/components/user/ZAvatar";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, h, createVNode, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listDeployPlanRequest,
  closeDeployPlanRequest
} from "@/api/app/deployPlanApi";
import { Modal, message } from "ant-design-vue";
import PLanStatusTag from "@/components/app/PlanStatusTag";
import { useDeloyPlanStore } from "@/pinia/deployPlanStore";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const planStore = useDeloyPlanStore();
const route = useRoute();
// 选择的环境
const selectedEnv = ref("");
const router = useRouter();
// 列表
const dataSource = ref([]);
// 分页
const dataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 数据项
const columns = [
  {
    i18nTitle: "deployPlan.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "deployPlan.artifactVersion",
    dataIndex: "artifactVersion",
    key: "artifactVersion"
  },
  {
    i18nTitle: "deployPlan.pipelineName",
    dataIndex: "pipelineName",
    key: "pipelineName"
  },
  {
    i18nTitle: "deployPlan.planStatus",
    dataIndex: "planStatus",
    key: "planStatus"
  },
  {
    i18nTitle: "deployPlan.createTime",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "deployPlan.creator",
    dataIndex: "creator",
    key: "creator"
  },
  {
    i18nTitle: "deployPlan.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 关闭发布计划
const closePlan = item => {
  Modal.confirm({
    title: `${t("deployPlan.confirmClose")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      closeDeployPlanRequest(item.id).then(() => {
        message.success(t("operationSuccess"));
        listDeployPlan();
      });
    },
    onCancel() {}
  });
};
// 计划列表
const listDeployPlan = () => {
  listDeployPlanRequest({
    appId: route.params.appId,
    pageNum: dataPage.current,
    env: selectedEnv.value
  }).then(res => {
    dataPage.totalCount = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
// 跳转创建页面
const gotoCreatePage = () => {
  router.push(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/create?env=${selectedEnv.value}`
  );
};
// 环境变化
const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/deployPlan/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  dataPage.current = 1;
  listDeployPlan();
};
// 查看计划
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