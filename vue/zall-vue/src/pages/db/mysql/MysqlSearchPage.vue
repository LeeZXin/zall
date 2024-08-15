<template>
  <div style="padding:10px" class="flex">
    <div class="left">
      <div class="base-select">
        <a-select
          style="width:100%"
          placeholder="请选择数据库"
          v-model:value="selectedDbId"
          :options="dbList"
        />
      </div>
      <ul class="access-ul">
        <li v-for="item in accessBases" v-bind:key="item">
          <div class="flex-center sticky-top">
            <div class="base-arrow" @click="showOrHideTables(item)">
              <CaretDownOutlined v-if="item.open" />
              <CaretRightOutlined v-else />
            </div>
            <div class="base-name" @click="showSearchTab(item.base)">{{item.base}}</div>
          </div>
          <ul class="tables-ul" v-show="item.open">
            <li
              v-for="table in item.tables"
              v-bind:key="table"
              @click="showCreateTableSql(item.base, table)"
            >{{table}}</li>
          </ul>
        </li>
      </ul>
    </div>
    <div class="right" v-show="tableInfoDiv === 'search'">
      <div class="right-body">
        <div class="right-body-title">{{selectedBase}}</div>
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
            <span>耗时: {{duration}}</span>
            <div>
              <span style="padding-right:6px">Limit:</span>
              <a-select v-model:value="searchLimit" style="width:80px" :options="limitList" />
              <a-button
                type="primary"
                :icon="h(SearchOutlined)"
                @click="doSearchSql(false)"
                style="margin-left:6px"
              >搜索</a-button>
              <a-button
                type="primary"
                style="margin-left:6px"
                :icon="h(TableOutlined)"
                @click="doExplainSql"
              >执行计划</a-button>
            </div>
          </div>
        </div>
        <div v-show="showSearchResult">
          <ZTable :columns="searchColumns" :dataSource="searchDataSource" />
          <div class="flex-between" style="margin-top:10px">
            <div>
              <a-pagination
                v-model:current="searchCurrPage"
                :total="searchTotalCount"
                show-less-items
                :pageSize="10"
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
            >导出结果</a-button>
          </div>
        </div>
      </div>
    </div>
    <div class="right" v-show="tableInfoDiv === 'table'">
      <div class="right-body">
        <div class="right-body-title">{{selectedTable}}</div>
        <a-tabs
          type="card"
          size="small"
          @change="tableInfoTabsChange"
          v-model:activeKey="tableInfoTabActiveKey"
        >
          <a-tab-pane key="table" tab="建表语句">
            <Codemirror
              v-model="createTableSql"
              style="height:380px;width:100%"
              :extensions="extensions"
              :disabled="true"
            />
          </a-tab-pane>
          <a-tab-pane key="index" tab="索引">
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
import ZTable from "@/components/common/ZTable";
import { ref, h, watch } from "vue";
import {
  SearchOutlined,
  ExportOutlined,
  TableOutlined,
  CaretRightOutlined,
  CaretDownOutlined
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
const extensions = [sql()];
const createTableSql = ref("");
const searchSql = ref("");
const searchCurrPage = ref(1);
const searchTotalCount = ref(0);
const searchColumns = ref([]);
const searchDataSource = ref([]);
const searchDataPartition = ref([]);
const searchAllData = ref([]);
const selectedDbId = ref(null);
const selectedBase = ref("");
const dbList = ref([]);
const accessBases = ref([]);
const selectedTable = ref("");
const tableInfoDiv = ref("");
const indexColumns = ref([]);
const indexDataSource = ref([]);
const tableInfoTabActiveKey = ref("table");
const searchCodeMirror = ref(null);
const showSearchResult = ref(false);
const searchLimit = ref(1);
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
const searchCodeMirrorReady = e => {
  searchCodeMirror.value = e.view;
};

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

const showSearchTab = base => {
  tableInfoDiv.value = "search";
  selectedBase.value = base;
};

const showCreateTableSql = (base, table) => {
  indexColumns.value = [];
  indexDataSource.value = [];
  createTableSql.value = "";
  selectedTable.value = base + "/" + table;
  selectedBase.value = base;
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
};

const tableInfoTabsChange = key => {
  if (key === "index" && indexColumns.value.length === 0) {
    let t = selectedTable.value.split("/");
    showTableIndex(t[0], t[1]);
  }
};

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
      searchTotalCount.value = res.data.data.length;
      searchCurrPage.value = 1;
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

const doExplainSql = () => {
  doSearchSql(true);
};

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

const changeSearchPage = () => {
  searchDataSource.value = searchDataPartition.value[searchCurrPage.value - 1];
};

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

watch(
  () => selectedDbId.value,
  () => selectDbIdChange()
);

listDb();
</script>

<style scoped>
.left {
  width: 300px;
  margin-right: 10px;
  background-color: white;
}

.right {
  width: calc(100% - 310px);
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
  padding: 0 5px;
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
  font-weight: bold;
}

.access-ul {
  width: 100%;
  max-height: 500px;
  overflow: scroll;
  border-left: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  border-right: 1px solid #d9d9d9;
  border-bottom-right-radius: 4px;
  border-bottom-left-radius: 4px;
}

.base-select {
  padding: 10px;
  border: 1px solid #d9d9d9;
  border-top-right-radius: 4px;
  border-top-left-radius: 4px;
}
</style>