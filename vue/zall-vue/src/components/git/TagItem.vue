<template>
  <div style="padding: 10px 16px">
    <div class="tag-title">
      <div class="flex-center no-wrap" style="max-width: 60%">
        <span class="tag-name">{{props.data.name}}</span>
        <div class="toggle-arrow">
          <DownOutlined v-show="showCommitMsg" @click="setShowCommitMsg(false)" />
          <RightOutlined v-show="!showCommitMsg" @click="setShowCommitMsg(true)" />
        </div>
      </div>
      <div class="flex-center">
        <a-popover placement="bottomRight" trigger="hover">
          <template #content>
            <ul class="op-list">
              <li @click="deleteTag(props.data.name)">
                <DeleteOutlined />
                <span style="margin-left:4px">删除tag</span>
              </li>
            </ul>
          </template>
          <div class="op-icon">...</div>
        </a-popover>
      </div>
    </div>
    <div class="tag-commit-msg" v-show="showCommitMsg">
      <pre>{{props.data.tagCommitMsg}}</pre>
    </div>
    <ul class="tag-content">
      <li>
        <ClockCircleOutlined />
        <span style="margin-left:4px">{{readableTimeComparingNow(props.data.taggerTime)}}</span>
      </li>
      <li @click="gotoCommitDiff(props.data.longCommitId)">
        <LinkOutlined />
        <span style="margin-left:4px">{{props.data.commitId}}</span>
      </li>
      <li @click="download(`${props.data.name}.tar.gz`)">
        <FileZipOutlined />
        <span style="margin-left:4px">tar.gz</span>
      </li>
      <li @click="download(`${props.data.name}.zip`)">
        <FileZipOutlined />
        <span style="margin-left:4px">zip</span>
      </li>
      <li v-if="props.data.verified">
        <div class="verify-tag">verified</div>
      </li>
    </ul>
  </div>
</template>
<script setup>
import {
  LinkOutlined,
  FileZipOutlined,
  ClockCircleOutlined,
  DownOutlined,
  RightOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined
} from "@ant-design/icons-vue";
import { ref, defineProps, createVNode, defineEmits } from "vue";
import { readableTimeComparingNow } from "@/utils/time";
import { Modal, message } from "ant-design-vue";
import { deleteTagRequest } from "@/api/git/repoApi";
import { useRouter } from "vue-router";
const router = useRouter();
const emit = defineEmits(["delete"]);
const showCommitMsg = ref(false);
const setShowCommitMsg = show => {
  showCommitMsg.value = show;
};
const props = defineProps(["data", "repoId"]);
const download = path => {
  window.open(`/api/gitRepo/archive?repoId=${props.repoId}&fileName=${path}`);
};
const gotoCommitDiff = commitId => {
  router.push(`/gitRepo/${props.repoId}/commit/diff/${commitId}`);
};
const deleteTag = tag => {
  Modal.confirm({
    title: `你确定要删除${tag}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    okText: "ok",
    cancelText: "cancel",
    onOk() {
      deleteTagRequest({
        repoId: props.repoId,
        tag
      }).then(() => {
        message.success("删除成功");
        emit("delete", tag);
      });
    },
    onCancel() {}
  });
};
</script>
<style scoped>
.tag-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}
.verify-tag {
  border-radius: 8px;
  border: 1px solid green;
  padding: 6px;
  font-size: 12px;
  color: green;
}
.tag-content {
  margin-top: 10px;
  display: flex;
  align-items: center;
}
.tag-content > li {
  font-size: 14px;
}
.tag-content > li:hover {
  cursor: pointer;
  color: #1677ff;
}
.tag-content > li + li {
  margin-left: 16px;
}
.tag-name {
  line-height: 32px;
  font-size: 16px;
  font-weight: bold;
}
.tag-commit-msg > pre {
  font-size: 14px;
}
.toggle-arrow {
  padding: 4px;
  border-radius: 4px;
  display: inline-block;
  font-size: 12px;
  cursor: pointer;
}
</style>