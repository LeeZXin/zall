<template>
  <div style="padding:10px;height:100%">
    <div style="margin-bottom:10px;" class="flex-end">
      <EnvSelector @change="onEnvChange" :defaultEnv="route.params.env" />
    </div>
    <ZTable :columns="columns" :dataSource="dataSource" :scroll="{x:1300}">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div class="op-icon" @click="deleteArtifact(dataItem)" v-else>
          <DeleteOutlined />
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
      @change="()=>listArtifact()"
    />
  </div>
</template>
<script setup>
import {
  ExclamationCircleOutlined,
  DeleteOutlined
} from "@ant-design/icons-vue";
import EnvSelector from "@/components/app/EnvSelector";
import ZTable from "@/components/common/ZTable";
import { ref, createVNode, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  listArtifactRequest,
  deleteArtifactRequest
} from "@/api/app/artifactApi";
import { message, Modal } from "ant-design-vue";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
// 选择环境
const selectedEnv = ref("");
const router = useRouter();
const dataSource = ref([]);
// 分页
const dataPage = reactive({
  current: 1,
  pageSize: 10,
  totalCount: 0
});
// 数据项
const columns = [
  {
    i18nTitle: "artifacts.name",
    dataIndex: "name",
    key: "name"
  },
  {
    i18nTitle: "artifacts.creator",
    dataIndex: "creator",
    key: "creator"
  },
  {
    i18nTitle: "artifacts.created",
    dataIndex: "created",
    key: "created"
  },
  {
    i18nTitle: "artifacts.operation",
    dataIndex: "operation",
    key: "operation",
    width: 130,
    fixed: "right"
  }
];
// 制品列表
const listArtifact = () => {
  listArtifactRequest(
    {
      appId: route.params.appId,
      env: selectedEnv.value,
      pageNum: dataPage.current
    },
    selectedEnv.value
  ).then(res => {
    dataPage.totalCount = res.totalCount;
    dataSource.value = res.data.map(item => {
      return {
        key: item.id,
        ...item
      };
    });
  });
};
// 删除制品
const deleteArtifact = item => {
  Modal.confirm({
    title: `${t("artifacts.confirmDelete")} ${item.name}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteArtifactRequest(item.id, item.env).then(() => {
        message.success(t("operationSuccess"));
        dataPage.current = 1;
        listArtifact();
      });
    },
    onCancel() {}
  });
};
// 环境变化
const onEnvChange = e => {
  router.replace(
    `/team/${route.params.teamId}/app/${route.params.appId}/artifact/list/${e.newVal}`
  );
  selectedEnv.value = e.newVal;
  listArtifact();
};
</script>
<style scoped>
</style>