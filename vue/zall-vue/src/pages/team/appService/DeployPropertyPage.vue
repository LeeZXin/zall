<template>
  <div style="padding:14px">
    <ZNaviBack url="/appService/property/list" name="返回配置列表" />
    <ZTable :columns="columns" :dataSource="dataSource">
      <template #bodyCell="{dataIndex, dataItem}">
        <span v-if="dataIndex !== 'operation'">{{dataItem[dataIndex]}}</span>
        <div v-else>
          <div class="op-icon" @click="deleteApp">
            <a-tooltip placement="top">
              <template #title>
                <span>发布配置</span>
              </template>
              <upload-outlined />
            </a-tooltip>
          </div>
        </div>
      </template>
    </ZTable>
    <div style="margin-top:10px">
      <div class="selected-node-text">
        <span>已选择配置节点</span>
        <span>:</span>
        <span>v00001</span>
      </div>
      <div class="flex-center">
        <span class="select-version-text">选择新版本:</span>
        <a-select style="width: 120px" size="small">
          <a-select-option value="jack">Jack</a-select-option>
          <a-select-option value="lucy">Lucy</a-select-option>
          <a-select-option value="disabled" disabled>Disabled</a-select-option>
          <a-select-option value="Yiminghe">yiminghe</a-select-option>
        </a-select>
      </div>
      <div style="width:100%">
        <code-diff
          :old-string="oldContent"
          :new-string="newContent"
          :context="10"
          outputFormat="side-by-side"
          :hideStat="true"
          filename="fuck"
          newFilename="nmb"
        />
      </div>
      <div style="margin-top:10px">
          <a-button type="primary">发布</a-button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { UploadOutlined } from "@ant-design/icons-vue";
import ZNaviBack from "@/components/common/ZNaviBack";
import ZTable from "@/components/common/ZTable";
import { CodeDiff } from "v-code-diff";
import { ref } from "vue";
const oldContent = ref("草泥马草泥马才能出门");
const newContent = ref("");
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
    title: "配置节点",
    dataIndex: "name",
    key: "name"
  },
  {
    title: "当前版本",
    dataIndex: "age",
    key: "age"
  },
  {
    title: "上次发布时间",
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
.selected-node-text {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
}
.select-version-text {
  line-height: 32px;
  height: 32px;
  font-size: 14px;
  margin-right: 10px;
}
</style>