<template>
  <a-popover v-model:open="branchTagVisible" trigger="click" placement="bottomLeft">
    <template #content>
      <div style="width:280px">
        <a-input placeholder="搜索分支" v-show="showSearchBranch">
          <template #prefix>
            <search-outlined />
          </template>
        </a-input>
        <a-input placeholder="搜索标签" v-show="!showSearchBranch">
          <template #prefix>
            <search-outlined />
          </template>
        </a-input>
      </div>
      <a-tabs style="width: 280px;" size="small" @change="onSearchBranchTagTabChange">
        <a-tab-pane key="branch" tab="分支">
          <ul class="branch-tag-list">
            <li @click="select('master')">
              <div class="branch-tag-name">master</div>
            </li>
            <li @click="select('sit')">
              <div class="branch-tag-name">sit</div>
            </li>
          </ul>
        </a-tab-pane>
        <a-tab-pane key="tag" tab="标签">
          <ul class="branch-tag-list">
            <li @click="select('master')">
              <div class="branch-tag-name">master</div>
            </li>
            <li @click="select('sit')">
              <div class="branch-tag-name">sit</div>
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
const props = defineProps(["style"]);
const selectedBranchOrTagName = ref("master");
const selectedBranchOrTagType = ref("branch");
const branchTagVisible = ref(false);
const showSearchBranch = ref(true);
const emit = defineEmits(["select"]);
const onSearchBranchTagTabChange = activeKey => {
  if (activeKey === "branch") {
    showSearchBranch.value = true;
  } else if (activeKey === "tag") {
    showSearchBranch.value = false;
  }
};
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