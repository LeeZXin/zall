<template>
  <div style="padding:10px;">
    <div v-if="branches.length > 0">
      <BranchTagSelect @select="onBranchTagSelect" :branches="branches" :tags="tags" />
      <a-popover v-model:open="cloneDownloadVisible" trigger="click" placement="bottomRight">
        <template #content>
          <a-tabs style="width: 300px;padding-bottom:12px" size="small">
            <a-tab-pane key="1" tab="HTTP">
              <div class="clone-input">
                <a-input v-model:value="gitHttpUrl" readonly />
                <div class="copy-icon" @click="copy(0)">
                  <CopyOutlined />
                </div>
              </div>
            </a-tab-pane>
            <a-tab-pane key="2" tab="SSH">
              <div class="clone-input">
                <a-input v-model:value="gitSshUrl" readonly />
                <div class="copy-icon" @click="copy(1)">
                  <CopyOutlined />
                </div>
              </div>
            </a-tab-pane>
          </a-tabs>
          <div>
            <a-button
              type="primary"
              ghost
              style="width:100%"
              @click="downloadZip"
              v-if="branches.length > 0 && selectedRef.refType === 'branch'"
              :icon="h(DownloadOutlined)"
            >{{t('repoIndex.downloadZip')}}</a-button>
          </div>
        </template>
        <a-button type="primary" style="float:right;font-size:14px">
          <span>{{t('repoIndex.clone')}}</span>
          <CaretDownOutlined />
        </a-button>
      </a-popover>
    </div>
    <div v-if="branches.length > 0">
      <div v-show="showDir">
        <div class="dir-table">
          <div class="first-line">
            <div class="commit-text flex-center">
              <ZAvatar :url="latestCommit.avatarUrl" :name="latestCommit.name" :showName="true" />
              <span>{{latestCommit.commitMsg}}</span>
            </div>
            <div class="commit-text">
              <span style="margin-right:4px">{{latestCommit.shortCommitId}}</span>
              <span>{{readableTimeComparingNow(latestCommit.committedTime)}}</span>
            </div>
          </div>
          <div class="dir-line" v-for="item in files" v-bind:key="item.commit.commitId">
            <div class="dir-line-item dir-line-file" @click="toRepoTree(item.rawPath)">
              <folder-outlined v-if="item.mode === 'directory'" style="margin-right:4px" />
              <span>{{item.path}}</span>
            </div>
            <div class="dir-line-item">{{item.commit.commitMsg}}</div>
            <div
              class="dir-line-item"
              style="text-align:right"
            >{{readableTimeComparingNow(item.commit.committedTime)}}</div>
          </div>
        </div>
        <div class="dir-table" v-if="showAddReadmeContent">
          <div class="first-line">
            <file-outlined />
            <span style="padding-left:6px">README.md</span>
          </div>
          <div class="readme-content">
            <v-md-editor v-model="readmeContent" mode="preview" />
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <div class="help-section">
        <div class="help-title">{{t('repoIndex.createRepoFromCmd')}}</div>
        <div class="help-text">
          <div>touch README.md</div>
          <div>git init</div>
          <div>git checkout -b main</div>
          <div>git add README.md</div>
          <div>git commit -m "first commit"</div>
          <div>git remote add origin {{gitSshUrl}}</div>
          <div>git push -u origin {{defaultBranch}}</div>
        </div>
      </div>
      <div class="help-section">
        <div class="help-title">{{t('repoIndex.pushCreatedRepoFromCmd')}}</div>
        <div class="help-text">
          <div>git remote add origin {{gitSshUrl}}</div>
          <div>git push -u origin {{defaultBranch}}</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import { ref, reactive, h } from "vue";
import {
  CaretDownOutlined,
  FileOutlined,
  CopyOutlined,
  FolderOutlined,
  DownloadOutlined
} from "@ant-design/icons-vue";
import VMdEditor from "@kangc/v-md-editor";
import "@kangc/v-md-editor/lib/style/base-editor.css";
import todoList from "@kangc/v-md-editor/lib/plugins/todo-list/index";
import "@kangc/v-md-editor/lib/plugins/todo-list/todo-list.css";
import githubTheme from "@kangc/v-md-editor/lib/theme/github.js";
import "@kangc/v-md-editor/lib/theme/style/github.css";
import BranchTagSelect from "@/components/git/BranchTagSelect";
import { treeRepoRequest, getBaseInfoRequest } from "@/api/git/repoApi";
import { useRoute, useRouter } from "vue-router";
import { message } from "ant-design-vue";
import { readableTimeComparingNow } from "@/utils/time";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const route = useRoute();
const router = useRouter();
VMdEditor.use(githubTheme);
VMdEditor.use(todoList());
const readmeContent = ref("");
const cloneDownloadVisible = ref(false);
const showAddReadmeContent = ref(false);
const showDir = ref(true);
// 仓库id
const repoId = parseInt(route.params.repoId);
// 分支列表
const branches = ref([]);
// 标签列表
const tags = ref([]);
// 文件列表
const files = ref([]);
// git clone http://xxxxx
const gitHttpUrl = ref("");
// git clone ssh://xxxx
const gitSshUrl = ref("");
// 默认分支
const defaultBranch = ref("");
// 选择的分支或标签
const selectedRef = reactive({
  ref: "",
  refType: ""
});
//
const onBranchTagSelect = event => {
  selectedRef.refType = event.key;
  selectedRef.ref = event.value;
  getTreeRepo(event.value, event.key);
};
// 最后一次提交
const latestCommit = reactive({
  committer: "",
  commitMsg: "",
  shortCommitId: "",
  committedTime: "",
  avatarUrl: "",
  name: ""
});
// 获取基本信息
const getBaseInfo = () => {
  getBaseInfoRequest(repoId).then(res => {
    branches.value = res.data.branches;
    tags.value = res.data.tags;
    gitHttpUrl.value = res.data.cloneHttpUrl;
    gitSshUrl.value = res.data.cloneSshUrl;
    defaultBranch.value = res.data.defaultBranch;
  });
};
// getTreeRepo 获取代码信息
const getTreeRepo = (ref, refType) => {
  treeRepoRequest({
    repoId,
    ref,
    refType
  }).then(res => {
    if (res.latestCommit) {
      latestCommit.committer = res.latestCommit.committer.account;
      latestCommit.avatarUrl = res.latestCommit.committer.avatarUrl;
      latestCommit.name = res.latestCommit.committer.name;
      latestCommit.commitMsg = res.latestCommit.commitMsg;
      latestCommit.shortCommitId = res.latestCommit.shortId;
      latestCommit.committedTime = res.latestCommit.committedTime;
    }
    if (res.tree && res.tree.files) {
      files.value = res.tree.files;
    }
    readmeContent.value = res.readmeText;
    showAddReadmeContent.value = res.hasReadme ? true : false;
  });
};
// 复制
const copy = type => {
  if (type === 0) {
    window.navigator.clipboard.writeText(gitHttpUrl.value);
  } else {
    window.navigator.clipboard.writeText(gitSshUrl.value);
  }
  message.success(t("copySuccess"));
};
// 跳转代码详情页
const toRepoTree = path => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/tree/${selectedRef.refType}/${selectedRef.ref}/${path}`
  );
};
const downloadZip = () => {
  window.open(
    `/api/gitRepo/archive?repoId=${repoId}&fileName=${selectedRef.ref}.zip`
  );
};
getBaseInfo();
</script>
<style scoped>
.dir-table {
  margin-top: 12px;
  border-radius: 4px;
  border: 1px solid #dadee3;
}
.dir-table > .first-line {
  line-height: 48px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.dir-line {
  width: 100%;
  display: flex;
  align-items: center;
  border-top: 1px solid #dadee3;
}
.dir-line-item {
  line-height: 48px;
  padding: 0 20px;
  overflow: hidden;
  width: 33.33%;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.dir-line-file:hover {
  cursor: pointer;
  color: #1677ff;
}
.copy-icon {
  text-align: center;
  width: 10%;
  cursor: pointer;
  color: gray;
}
.copy-icon:hover {
  color: black;
}
.add-readme,
.readme-content {
  border-top: 1px solid #dadee3;
}
.commit-text {
  max-width: 50%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.help-section {
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}
.help-section + .help-section {
  margin-top: 10px;
}
.help-title {
  padding: 14px;
  border-bottom: 1px solid #d9d9d9;
  font-size: 16px;
  font-weight: bold;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}
.help-text {
  padding: 10px 14px;
  font-size: 14px;
  line-height: 32px;
}
.clone-input {
  display: flex;
  align-items: center;
  margin-top: 10px;
}
</style>