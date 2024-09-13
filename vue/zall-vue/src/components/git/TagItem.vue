<template>
  <div style="padding: 10px 16px">
    <div class="tag-title">
      <div class="flex-center no-wrap" style="max-width: 60%">
        <span class="tag-name">{{props.data.name}}</span>
        <div class="toggle-dot">
          <EllipsisOutlined @click="toggleShowCommitMsg" />
        </div>
      </div>
      <div class="flex-center">
        <a-popover placement="bottomRight" trigger="hover">
          <template #content>
            <ul class="op-list">
              <li @click="deleteTag(props.data.name)">
                <DeleteOutlined />
                <span style="margin-left:4px">{{t('tagList.deleteTag')}}</span>
              </li>
            </ul>
          </template>
          <div class="op-icon">
            <EllipsisOutlined />
          </div>
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
      <li @click="gotoCommitDiff(props.data.longCommitId)" class="btn">
        <LinkOutlined />
        <span style="margin-left:4px">{{props.data.commitId}}</span>
      </li>
      <li @click="download(`${props.data.name}.tar.gz`)" class="btn">
        <FileZipOutlined />
        <span style="margin-left:4px">tar.gz</span>
      </li>
      <li @click="download(`${props.data.name}.zip`)" class="btn">
        <FileZipOutlined />
        <span style="margin-left:4px">zip</span>
      </li>
      <li v-if="props.data.verified">
        <a-popover placement="bottom">
          <template #content>
            <div style="width: 300px;font-size:14px;padding:6px">
              <div style="margin-bottom: 12px;" class="flex-center no-wrap">
                <CheckCircleFilled style="color:green;margin-right:10px" />
                <span>{{t('tagList.thisTagIsVerified')}}</span>
              </div>
              <div class="flex-center" style="margin-bottom: 12px;">
                <ZAvatar
                  :url="props.data.signer?.avatarUrl"
                  :name="props.data.signer?.name"
                  :disablePopover="true"
                  size="medium"
                />
                <div style="margin-left:8px">
                  <div style="margin-bottom: 3px" class="no-wrap">{{props.data.signer?.account}}</div>
                  <div class="no-wrap">{{props.data.signer?.name}}</div>
                </div>
              </div>
              <div class="no-wrap">{{props.data.signer?.type}} KEY</div>
              <div style="color:gray;word-break:break-all">{{props.data.signer?.key}}</div>
            </div>
          </template>
          <span style="cursor:pointer">
            <a-tag color="green">{{t('tagList.verified')}}</a-tag>
          </span>
        </a-popover>
      </li>
    </ul>
  </div>
</template>
<script setup>
import {
  LinkOutlined,
  FileZipOutlined,
  ClockCircleOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  EllipsisOutlined,
  CheckCircleFilled
} from "@ant-design/icons-vue";
import ZAvatar from "@/components/user/ZAvatar";
import { ref, defineProps, createVNode, defineEmits } from "vue";
import { readableTimeComparingNow } from "@/utils/time";
import { Modal, message } from "ant-design-vue";
import { deleteTagRequest } from "@/api/git/repoApi";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const router = useRouter();
const emit = defineEmits(["delete"]);
const showCommitMsg = ref(false);
// 展示提交信息
const toggleShowCommitMsg = () => {
  showCommitMsg.value = !showCommitMsg.value;
};
const props = defineProps(["data", "repoId", "teamId"]);
// 下载文件
const download = path => {
  window.open(`/api/gitRepo/archive?repoId=${props.repoId}&fileName=${path}`);
};
// 跳转commit页面
const gotoCommitDiff = commitId => {
  router.push(
    `/team/${props.teamId}/gitRepo/${props.repoId}/commit/diff/${commitId}`
  );
};
// 删除标签
const deleteTag = tag => {
  Modal.confirm({
    title: `${t("tagList.deleteTag")} ${tag}?`,
    icon: createVNode(ExclamationCircleOutlined),
    onOk() {
      deleteTagRequest({
        repoId: props.repoId,
        tag
      }).then(() => {
        message.success(t("operationSuccess"));
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
.tag-content {
  margin-top: 10px;
  display: flex;
  align-items: center;
}
.tag-content > li {
  font-size: 14px;
}
.btn:hover {
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
.toggle-dot {
  padding: 0 4px;
  border-radius: 4px;
  display: inline-block;
  font-size: 16px;
  cursor: pointer;
}
</style>