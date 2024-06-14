<template>
  <div style="padding:14px">
    <a-date-picker v-model:value="dateVal" style="margin-bottom:10px" @change="pageLog" />
    <ZTable
      :columns="columns"
      :dataSource="dataSource"
      style="margin-top:0px"
    >
      <template #bodyCell="{dataIndex, dataItem}">
        <template v-if="dataIndex !== 'operation'">
          <span>{{dataItem[dataIndex]}}</span>
        </template>
        <template v-else>
          <a-popover placement="bottomRight" trigger="hover">
            <template #content>
              <ul class="op-list">
                <li @click="showReqBody(dataItem['reqBody'])">
                  <ControlOutlined />
                  <span style="margin-left:4px">查看请求体</span>
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
    <a-modal v-model:open="reqBodyModalOpen" title="请求体" :footer="null">
      <Codemirror
        v-model="reqBody"
        :style="codemirrorStyle"
        :extensions="extensions"
        :disabled="true"
      />
    </a-modal>
  </div>
</template>
<script setup>
import ZTable from "@/components/common/ZTable";
import { ref } from "vue";
import { useRoute } from "vue-router";
import { EllipsisOutlined, ControlOutlined } from "@ant-design/icons-vue";
import { pageLogRequest } from "@/api/git/oplogApi";
import dayjs from "dayjs";
import { Codemirror } from "vue-codemirror";
import { json } from "@codemirror/lang-json";
import { oneDark } from "@codemirror/theme-one-dark";
const extensions = [json(), oneDark];
const codemirrorStyle = { height: "380px", width: "100%" };
const dateVal = ref(dayjs());
const route = useRoute();
const dataSource = ref([]);
const currPage = ref(1);
const pageSize = 10;
const totalCount = ref(0);
const reqBody = ref("");
const reqBodyModalOpen = ref(false);
const columns = ref([
  {
    title: "操作人",
    dataIndex: "account",
    key: "account"
  },
  {
    title: "操作时间",
    dataIndex: "created",
    key: "created"
  },
  {
    title: "内容",
    dataIndex: "content",
    key: "content"
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation"
  }
]);
const showReqBody = (req) => {
  reqBodyModalOpen.value = true;
  let jsonReq = JSON.parse(req);
  reqBody.value = JSON.stringify(jsonReq, null, 4);
}
const pageLog = () => {
  pageLogRequest({
    repoId: route.params.repoId,
    pageNum: currPage.value,
    dateStr: dateVal.value.format("YYYY-MM-DD")
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
pageLog();
</script>
<style scoped>
</style>