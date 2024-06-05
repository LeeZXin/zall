<template>
  <div style="padding:10px;">
    <div style="height:32px">
      <BranchTagSelect
        @select="onBranchTagSelect"
        :branches="branches"
        :tags="tags"
        v-if="branches.length > 0"
      />
      <a-popover v-model:open="cloneDownloadVisible" trigger="click" placement="bottomRight">
        <template #content>
          <a-tabs style="width: 300px;padding-bottom:12px" size="small">
            <a-tab-pane key="1" tab="HTTP">
              <div style="display:flex;align-items:center;margin-top:10px">
                <a-input type="input" v-model:value="gitHttpUrl" readonly />
                <div class="copy-icon" @click="copy(0)">
                  <a-tooltip placement="top">
                    <template #title>
                      <span>Copy</span>
                    </template>
                    <copy-outlined />
                  </a-tooltip>
                </div>
              </div>
            </a-tab-pane>
            <a-tab-pane key="2" tab="SSH">
              <div style="display:flex;align-items:center;margin-top:10px">
                <a-input type="input" v-model:value="gitSshUrl" readonly />
                <div class="copy-icon" @click="copy(1)">
                  <a-tooltip placement="top">
                    <template #title>
                      <span>Copy</span>
                    </template>
                    <copy-outlined />
                  </a-tooltip>
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
            >下载zip</a-button>
          </div>
        </template>
        <a-button type="primary" style="float:right">
          <span>克隆</span>
          <caret-down-outlined style="font-size:12px" />
        </a-button>
      </a-popover>
    </div>
    <div v-if="branches.length > 0">
      <div v-show="showDir">
        <div class="dir-table">
          <div class="first-line">
            <div class="commit-text">
              <span style="margin-right:4px">{{latestCommit.committer}}</span>
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
        <div class="dir-table">
          <div class="first-line">
            <file-outlined />
            <span style="padding-left:6px">README.md</span>
          </div>
          <div class="readme-content" v-if="showAddReadmeContent">
            <v-md-editor v-model="readmeContent" mode="preview" />
          </div>
          <div class="add-readme" v-if="!showAddReadmeContent">
            <div style="text-align:center;font-size:24px;line-height:60px;">
              <file-outlined />
            </div>
            <div
              style="font-weight:bold;text-align:center;font-size:22px;line-height:30px;padding:24px;"
            >Try to Add a README to let everyone interested in this repository understands yours project</div>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <div class="help-section">
        <div class="help-title">从命令行创建一个新的仓库</div>
        <div class="help-text">
          <div>touch README.md</div>
          <div>git init</div>
          <div>git checkout -b main</div>
          <div>git add README.md</div>
          <div>git commit -m "first commit"</div>
          <div>git remote add origin {{gitSshUrl}}</div>
          <div>git push -u origin main</div>
        </div>
      </div>
      <div class="help-section">
        <div class="help-title">从命令行推送已经创建的仓库</div>
        <div class="help-text">
          <div>git remote add origin {{gitSshUrl}}</div>
          <div>git push -u origin main</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, reactive } from "vue";
import {
  CaretDownOutlined,
  FileOutlined,
  CopyOutlined,
  FolderOutlined
} from "@ant-design/icons-vue";
import VMdEditor from "@kangc/v-md-editor";
import "@kangc/v-md-editor/lib/style/base-editor.css";
import todoList from "@kangc/v-md-editor/lib/plugins/todo-list/index";
import "@kangc/v-md-editor/lib/plugins/todo-list/todo-list.css";
import githubTheme from "@kangc/v-md-editor/lib/theme/github.js";
import "@kangc/v-md-editor/lib/theme/style/github.css";
import BranchTagSelect from "@/components/git/BranchTagSelect";
import { treeRepoRequest, simpleInfoRequest } from "@/api/git/repoApi";
import { useRoute, useRouter } from "vue-router";
import { message } from "ant-design-vue";
import { readableTimeComparingNow } from "@/utils/time";
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
  committedTime: ""
});
simpleInfoRequest(repoId).then(res => {
  branches.value = res.data.branches;
  tags.value = res.data.tags;
  gitHttpUrl.value = res.data.cloneHttpUrl;
  gitSshUrl.value = res.data.cloneSshUrl;
});
// getTreeRepo 获取代码信息
const getTreeRepo = (ref, refType) => {
  treeRepoRequest({
    repoId,
    ref,
    refType
  }).then(res => {
    if (res.latestCommit) {
      latestCommit.committer = res.latestCommit.committer.account;
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
  message.success("复制成功");
};
// 跳转代码详情页
const toRepoTree = path => {
  router.push(
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/tree/${selectedRef.refType}/${selectedRef.ref}/` +
      path
  );
};
const downloadZip = () => {
  window.open(
    `/api/gitRepo/archive?repoId=${repoId}&fileName=${selectedRef.ref}.zip`
  );
};
</script>
<style scoped>
.dir-table {
  margin-top: 12px;
  border-radius: 4px;
  border: 1px solid #dadee3;
}
.dir-table > .first-line {
  height: 42px;
  line-height: 42px;
  padding: 0 16px;
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
  height: 42px;
  line-height: 42px;
  padding: 0 16px;
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
  padding: 20px 16px;
  font-size: 14px;
  line-height: 22px;
}
</style>