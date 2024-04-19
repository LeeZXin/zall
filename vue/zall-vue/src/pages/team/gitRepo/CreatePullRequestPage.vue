<template>
  <div style="padding:14px">
    <div class="header">
      <div class="title">Compare changes</div>
      <div class="desc">Compare changes across branches, commits, tags, and more below</div>
    </div>
    <div class="merge-select">
      <div class="merge-select-left">
        <BranchTagSelect style="margin-right:6px" />
        <arrow-left-outlined />
        <BranchTagSelect style="margin-left:6px" />
        <div class="merge-warn" v-show="!canMerge">
          <close-outlined style="color:red" />
          <span style="padding-left:4px">Can’t merge</span>
        </div>
        <div class="merge-warn" v-show="canMerge">
          <check-outlined style="color:green" />
          <span style="padding-left:4px">Can merge</span>
        </div>
      </div>
      <a-button
        type="primary"
        :disabled="!canMerge"
        v-show="!submitFormVisible"
        @click="showSubmitForm"
      >创建合并请求</a-button>
    </div>
    <div class="submit-form" v-if="submitFormVisible">
      <a-input placeholder="标题" />
      <a-textarea
        style="margin-top:10px"
        placeholder="评论"
        :auto-size="{ minRows: 5, maxRows: 10 }"
      />
      <div style="margin-top:10px;text-align:right">
        <a-button type="primary">创建合并请求</a-button>
      </div>
    </div>
    <div class="commit-list">
      <div class="title">
        <span style="font-weight:bold">3</span>
        <span style="color:orange">个提交</span>
        <span style="font-weight:bold">182</span>
        <span style="color:green">次新增</span>
        <span style="font-weight:bold">3</span>
        <span style="color:red">次删除</span>
      </div>
      <ul class="commit-item-list">
        <li>
          <div class="commit-item">
            <div class="title">feat: fuufufufufuufufufufu</div>
            <div class="desc">
              <span>LeeZXin</span>
              <span>提交于</span>
              <span>2024/02/12</span>
            </div>
          </div>
          <CommitSha>cmn,b</CommitSha>
        </li>
        <li>
          <div class="commit-item">
            <div class="title">feat: fuufufufufuufufufufu</div>
            <div class="desc">
              <span>LeeZXin</span>
              <span>提交于</span>
              <span>2024/02/12</span>
            </div>
          </div>
          <CommitSha>cmn,b</CommitSha>
        </li>
        <li>
          <div class="commit-item">
            <div class="title">feat: fuufufufufuufufufufu</div>
            <div class="desc">
              <span>LeeZXin</span>
              <span>提交于</span>
              <span>2024/02/12</span>
            </div>
          </div>
          <CommitSha>hh</CommitSha>
        </li>
      </ul>
    </div>

    <FileDiffDetail />
    <FileDiffDetail />
    <FileDiffDetail />
    <FileDiffDetail />
  </div>
</template>
<script setup>
import {
  ArrowLeftOutlined,
  CloseOutlined,
  CheckOutlined
} from "@ant-design/icons-vue";
import CommitSha from "@/components/git/CommitSha";
import FileDiffDetail from "@/components/git/FileDiffDetail";
import BranchTagSelect from "@/components/git/BranchTagSelect";
import { ref } from "vue";
const canMerge = ref(true);
const submitFormVisible = ref(false);
const showSubmitForm = () => {
  submitFormVisible.value = true;
};
</script>
<style scoped>
.header {
  margin-bottom: 10px;
}
.header > .title {
  font-size: 16px;
  padding-bottom: 4px;
}
.header > .desc {
  font-size: 12px;
  color: gray;
}
.merge-select {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  display: flex;
  align-items: center;
  padding: 10px;
  justify-content: space-between;
}
.commit-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  margin-top: 10px;
}
.merge-warn {
  font-size: 14px;
  margin-left: 6px;
}
.merge-select-left {
  display: flex;
  align-items: center;
}
.commit-list > .title {
  font-size: 14px;
  line-height: 32px;
  padding: 0 10px;
}
.commit-list > .title > span {
  padding-left: 6px;
}
.commit-item-list > li {
  padding: 10px;
  border-top: 1px solid #d9d9d9;
  display: flex;
  width: 100%;
  justify-content: space-between;
  align-items: baseline;
}
.commit-item > .title {
  font-size: 16px;
  padding-bottom: 8px;
  font-weight: bold;
}
.commit-item > .desc {
  font-size: 12px;
  color: gray;
}
.submit-form {
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  margin-top: 10px;
  padding: 10px;
}
</style>