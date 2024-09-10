<template>
  <div style="padding:10px" class="flex">
    <div class="left">
      <div class="base-select">
        <a-select
          style="width:100%"
          v-model:value="selectedDbId"
          :options="dbList"
          show-search
          :filter-option="filterDbListOption"
        />
      </div>
      <div class="access-bordered">
        <ul class="access-ul" v-if="accessBases.length > 0">
          <li v-for="item in accessBases" v-bind:key="item">
            <div class="flex-center sticky-top">
              <div class="base-arrow" @click="showOrHideTables(item)">
                <CaretDownOutlined v-if="item.open" />
                <CaretRightOutlined v-else />
              </div>
              <div class="base-name" @click="showSearchTab(item.base)">
                <DatabaseOutlined />
                <span style="margin-left:4px">{{item.base}}</span>
              </div>
            </div>
            <ul class="tables-ul" v-show="item.open">
              <li
                v-for="table in item.tables"
                v-bind:key="table"
                @click="showCreateTableSql(item.base, table.table, table.size)"
              >
                <TableOutlined />
                <span style="margin-left:4px">
                  <span>{{table.table}}</span>
                  <span style="font-size:10px;color:gray">{{table.size}}</span>
                </span>
              </li>
            </ul>
          </li>
        </ul>
        <ZNoData v-else :unbordered="true" />
      </div>
    </div>
    <div class="right" v-show="tableInfoDiv === 'search'">
      <div class="right-body">
        <div class="right-body-title">
          <span style="font-weight:bold">{{selectedBase}}</span>
        </div>
        <div>
          <div style="margin-bottom:10px;">
            <Codemirror
              v-model="searchSql"
              style="height:280px;width:100%"
              :extensions="extensions"
              @ready="searchCodeMirrorReady"
            />
          </div>
          <div class="flex-between">
            <span>{{t('mysqlSearch.duration')}}: {{duration}}</span>
            <div>
              <span style="padding-right:6px">Limit:</span>
              <a-select v-model:value="searchLimit" style="width:80px" :options="limitList" />
              <a-button
                type="primary"
                :icon="h(PlayCircleOutlined)"
                @click="doSearchSql(false)"
                style="margin-left:6px"
              >{{t('mysqlSearch.search')}}</a-button>
              <a-button
                type="primary"
                style="margin-left:6px"
                :icon="h(TableOutlined)"
                @click="doExplainSql"
              >{{t('mysqlSearch.explain')}}</a-button>
            </div>
          </div>
        </div>
        <div v-show="showSearchResult">
          <ZTable :columns="searchColumns" :dataSource="searchDataSource" />
          <div class="flex-between" style="margin-top:10px">
            <div>
              <a-pagination
                v-model:current="searchDataPage.current"
                :total="searchDataPage.totalCount"
                show-less-items
                :pageSize="searchDataPage.pageSize"
                :hideOnSinglePage="true"
                :showSizeChanger="false"
                @change="()=>changeSearchPage()"
              />
            </div>
            <a-button
              type="primary"
              :icon="h(ExportOutlined)"
              v-show="searchAllData.length > 0"
              @click="exportToExcel"
            >{{t('mysqlSearch.export')}}</a-button>
          </div>
        </div>
      </div>
    </div>
    <div class="right" v-show="tableInfoDiv === 'table'">
      <div class="right-body">
        <div class="right-body-title">
          <span style="font-weight:bold">{{selectedTable.table}}</span>
          <span style="font-size:10px;color:gray">{{selectedTable.size}}</span>
        </div>
        <a-tabs
          type="card"
          size="small"
          @change="tableInfoTabsChange"
          v-model:activeKey="tableInfoTabActiveKey"
        >
          <a-tab-pane key="table" :tab="t('mysqlSearch.createTableSql')">
            <Codemirror
              v-model="createTableSql"
              style="height:380px;width:100%"
              :extensions="extensions"
              :disabled="true"
            />
          </a-tab-pane>
          <a-tab-pane key="index" :tab="t('mysqlSearch.index')">
            <div style="padding: 10px 0">
              <ZTable
                :columns="indexColumns"
                :dataSource="indexDataSource"
                style="margin-top:unset"
              />
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </div>
  </div>
</template>

<script setup>
import ZNoData from "@/components/common/ZNoData";
import ZTable from "@/components/common/ZTable";
import { ref, h, watch, reactive } from "vue";
import {
  PlayCircleOutlined,
  ExportOutlined,
  TableOutlined,
  CaretRightOutlined,
  CaretDownOutlined,
  DatabaseOutlined
} from "@ant-design/icons-vue";
import {
  listAuthorizedDbRequest,
  listAuthorizedBaseRequest,
  listAuthorizedTableRequest,
  getCreateTableSqlRequest,
  showTableIndexRequest,
  executeSelectSqlRequest
} from "@/api/db/mysqlApi";
import { Codemirror } from "vue-codemirror";
import { sql } from "@codemirror/lang-sql";
import { saveAs } from "file-saver";
import * as XLSX from "xlsx";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const extensions = [sql()];
// 建表语句
const createTableSql = ref("");
// 搜索sql
const searchSql = ref("");
// 搜索数据分页
const searchDataPage = reactive({
  current: 1,
  totalCount: 0,
  pageSize: 10
});
// 搜索结果数据项
const searchColumns = ref([]);
// 搜索结果列表
const searchDataSource = ref([]);
// 搜索结果分页
const searchDataPartition = ref([]);
// 搜索所有数据
const searchAllData = ref([]);
// 选择的数据库
const selectedDbId = ref(null);
// 选择的库
const selectedBase = ref("");
// 数据库列表
const dbList = ref([]);
// 数据库表
const accessBases = ref([]);
// 选择的表
const selectedTable = reactive({
  table: "",
  size: ""
});
const tableInfoDiv = ref("");
// 索引数据项
const indexColumns = ref([]);
// 索引数据
const indexDataSource = ref([]);
const tableInfoTabActiveKey = ref("table");
const searchCodeMirror = ref(null);
const showSearchResult = ref(false);
const searchLimit = ref(1);
// 执行耗时
const duration = ref("");
const limitList = [
  {
    value: 1,
    label: 1
  },
  {
    value: 100,
    label: 100
  },
  {
    value: 500,
    label: 500
  },
  {
    value: 1000,
    label: 1000
  }
];
// codemirror加载完成触发 为了能获取鼠标选中语句结果
const searchCodeMirrorReady = e => {
  searchCodeMirror.value = e.view;
};
// 数据库列表
const listDb = () => {
  listAuthorizedDbRequest().then(res => {
    dbList.value = res.data.map(item => {
      return {
        value: item.id,
        label: item.name
      };
    });
    if (res.data.length > 0) {
      selectedDbId.value = res.data[0].id;
    }
  });
};
// 选择数据库
const selectDbIdChange = () => {
  listAuthorizedBaseRequest(selectedDbId.value).then(res => {
    accessBases.value = res.data.map(item => {
      return {
        base: item,
        open: false,
        tables: []
      };
    });
  });
};
// 是否展示数据表列表
const showOrHideTables = item => {
  if (item.open) {
    item.open = false;
  } else if (item.tables.length > 0) {
    item.open = true;
  } else {
    listAuthorizedTableRequest({
      dbId: selectedDbId.value,
      accessBase: item.base
    }).then(res => {
      item.tables = res.data;
      item.open = true;
    });
  }
};
// 展示搜索的tab
const showSearchTab = base => {
  tableInfoDiv.value = "search";
  selectedBase.value = base;
  selectedTable.table = "";
};
// 展示建表语句
const showCreateTableSql = (base, table, size) => {
  let baseAndTable = base + "/" + table;
  if (selectedTable.table !== baseAndTable) {
    indexColumns.value = [];
    indexDataSource.value = [];
    createTableSql.value = "";
    selectedBase.value = base;
    selectedTable.table = baseAndTable;
    selectedTable.size = size;
    tableInfoDiv.value = "table";
    getCreateTableSqlRequest({
      dbId: selectedDbId.value,
      accessBase: base,
      accessTable: table
    }).then(res => {
      createTableSql.value = res.data;
      if (tableInfoTabActiveKey.value === "index") {
        showTableIndex(base, table);
      }
    });
  }
};
// 索引和建表语句切换
const tableInfoTabsChange = key => {
  if (key === "index" && indexColumns.value.length === 0) {
    let t = selectedTable.table.split("/");
    showTableIndex(t[0], t[1]);
  }
};
// 展示索引
const showTableIndex = (base, table) => {
  showTableIndexRequest({
    dbId: selectedDbId.value,
    accessBase: base,
    accessTable: table
  }).then(res => {
    indexColumns.value = res.data.columns.map(item => {
      return {
        title: item,
        dataIndex: item,
        key: item
      };
    });
    indexDataSource.value = res.data.data.map((item, index) => {
      return {
        key: index,
        ...item
      };
    });
  });
};
// 搜索sql
const doSearchSql = isExplain => {
  const ranges = searchCodeMirror.value.state.selection.ranges;
  let sql = searchSql.value;
  const subList = ranges.map(range => {
    return sql.substring(range.from, range.to);
  });
  if (subList.length > 0 && subList[0]) {
    sql = subList[0];
  }
  if (sql) {
    duration.value = "";
    executeSelectSqlRequest({
      dbId: selectedDbId.value,
      accessBase: selectedBase.value,
      cmd: isExplain ? "explain " + sql : sql,
      limit: searchLimit.value
    }).then(res => {
      duration.value = res.data.duration;
      searchDataPage.totalCount = res.data.data.length;
      searchDataPage.current = 1;
      showSearchResult.value = true;
      searchColumns.value = res.data.columns.map(item => {
        return {
          title: item,
          dataIndex: item,
          key: item
        };
      });
      let allData = res.data.data.map((item, index) => {
        return {
          key: index + 1,
          ...item
        };
      });
      searchAllData.value = allData;
      searchDataPartition.value = partition(allData, 10);
      if (searchDataPartition.value.length > 0) {
        searchDataSource.value = searchDataPartition.value[0];
      } else {
        searchDataSource.value = [];
      }
    });
  }
};
// 执行计划
const doExplainSql = () => {
  doSearchSql(true);
};
// 分页
const partition = (arr, length) => {
  let result = [];
  for (let i = 0, j = arr.length; i < j; i++) {
    if (i % length === 0) {
      result.push([]);
    }
    result[result.length - 1].push(arr[i]);
  }
  return result;
};
// 搜索结果分页触发
const changeSearchPage = () => {
  searchDataSource.value =
    searchDataPartition.value[searchDataPage.current - 1];
};
// 导出数据表
const exportToExcel = () => {
  const ws = XLSX.utils.json_to_sheet(searchAllData.value, {
    header: searchColumns.value.map(item => item.dataIndex)
  });
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, "Sheet1");

  const wbout = XLSX.write(wb, { bookType: "xlsx", type: "binary" });
  function s2ab(s) {
    const buf = new ArrayBuffer(s.length);
    const view = new Uint8Array(buf);
    for (let i = 0; i < s.length; ++i) {
      view[i] = s.charCodeAt(i) & 0xff;
    }
    return buf;
  }

  saveAs(
    new Blob([s2ab(wbout)], { type: "application/octet-stream" }),
    `${new Date().getTime()}.xlsx`
  );
};

// 数据库下拉框搜索
const filterDbListOption = (input, option) => {
  return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};

watch(
  () => selectedDbId.value,
  () => selectDbIdChange()
);

listDb();
</script>

<style scoped>
.left {
  width: 280px;
  margin-right: 10px;
  background-color: white;
}

.right {
  width: calc(100% - 290px);
  overflow: scroll;
}

.right-body {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 10px;
  background-color: white;
}

.base-name {
  width: calc(100% - 32px);
  line-height: 32px;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  padding-right: 5px;
  background-color: white;
  cursor: pointer;
}

.base-arrow {
  width: 28px;
  height: 32px;
  line-height: 32px;
  text-align: center;
  font-size: 12px;
  cursor: pointer;
  background-color: white;
}

.tables-ul > li:hover {
  background-color: #f0f0f0;
  cursor: pointer;
}

.tables-ul > li {
  line-height: 26px;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  padding: 0 40px;
}

.right-body-title {
  padding: 8px;
  margin-bottom: 10px;
  font-size: 14px;
}

.access-bordered {
  padding: 1px;
  border-left: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
  border-bottom-right-radius: 4px;
  border-bottom-left-radius: 4px;
}

.access-ul {
  width: 100%;
  max-height: 500px;
  overflow: scroll;
}

.base-select {
  padding: 10px;
  border: 1px solid #d9d9d9;
  border-top-right-radius: 4px;
  border-top-left-radius: 4px;
}
</style>