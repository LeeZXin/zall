<template>
  <a-popover v-model:open="branchTagVisible" trigger="click" placement="bottomLeft">
    <template #content>
      <div style="width:280px">
        <a-input
          placeholder="搜索分支"
          v-show="showSearchBranch"
          v-model:value="searchBranch"
          @change="branchInputChange"
        >
          <template #prefix>
            <search-outlined />
          </template>
        </a-input>
        <a-input
          placeholder="搜索标签"
          v-show="!showSearchBranch"
          v-model:value="searchTag"
          @change="tagInputChange"
        >
          <template #prefix>
            <search-outlined />
          </template>
        </a-input>
      </div>
      <a-tabs style="width: 280px;" size="small" @change="onSearchBranchTagTabChange">
        <a-tab-pane key="branch" tab="分支">
          <ul class="branch-tag-list">
            <li @click="select(item)" v-for="item in branches" v-bind:key="item">
              <div class="branch-tag-name">{{item}}</div>
            </li>
          </ul>
        </a-tab-pane>
        <a-tab-pane key="tag" tab="标签">
          <ul class="branch-tag-list">
            <li @click="select(item)" v-for="item in tags" v-bind:key="item">
              <div class="branch-tag-name">{{item}}</div>
            </li>
          </ul>
        </a-tab-pane>
      </a-tabs>
    </template>
    <div class="branch-tag-select" :style="props.style">
      <branches-outlined v-show="selectedBranchOrTagType === 'branch'" />
      <tag-outlined v-show="selectedBranchOrTagType === 'tag'" />
      <span class="branch-tag-select-text">{{selectedBranchOrTagName}}</span>
      <caret-down-outlined />
    </div>
  </a-popover>
</template>
<script setup>
import { ref, defineEmits, defineProps } from "vue";
import {
  SearchOutlined,
  BranchesOutlined,
  TagOutlined,
  CaretDownOutlined
} from "@ant-design/icons-vue";
const props = defineProps(["style", "branches", "tags"]);
const selectedBranchOrTagName = ref("");
const selectedBranchOrTagType = ref("");
const searchBranch = ref("");
const searchTag = ref("");
const branchTagVisible = ref(false);
const showSearchBranch = ref(true);
const branches = ref(props.branches);
const tags = ref(props.tags);
const emit = defineEmits(["select"]);
const onSearchBranchTagTabChange = activeKey => {
  if (activeKey === "branch") {
    showSearchBranch.value = true;
  } else if (activeKey === "tag") {
    showSearchBranch.value = false;
  }
};
// 选择分支或标签
const select = name => {
  let selectedType = showSearchBranch.value ? "branch" : "tag";
  emit("select", {
    key: selectedType,
    value: name
  });
  selectedBranchOrTagType.value = selectedType;
  selectedBranchOrTagName.value = name;
  branchTagVisible.value = false;
};
if (props.branches && props.branches.length > 0) {
  select(props.branches[0]);
}
// 搜索分支
const branchInputChange = () => {
  const val = searchBranch.value;
  if (val === "") {
    branches.value = props.branches;
  } else {
    branches.value = props.branches.filter(item => {
      return item.indexOf(val) >= 0;
    });
  }
};
// 搜索标签
const tagInputChange = () => {
  const val = searchTag.value;
  if (val === "") {
    tags.value = props.tags;
  } else {
    tags.value = props.tags.filter(item => {
      return item.indexOf(val) >= 0;
    });
  }
};
</script>
<style scoped>
.branch-tag-select {
  display: inline-block;
  height: 32px;
  line-height: 32px;
  border: 1px solid #dadee3;
  border-radius: 4px;
  padding: 0 10px;
  cursor: pointer;
}
.branch-tag-select:hover {
  background-color: #f0f0f0;
}
.branch-tag-select-text {
  padding: 0px 4px;
}
.branch-tag-list > li {
  padding: 6px 0px;
}

.branch-tag-list > li + li {
  border-top: 1px solid #f0f0f0;
}
.branch-tag-name {
  height: 32px;
  line-height: 32px;
  font-size: 14px;
  cursor: pointer;
  border-radius: 4px;
  padding-left: 4px;
  width: 100%;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.branch-tag-name:hover {
  background-color: #f0f0f0;
}
</style>