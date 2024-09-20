<template>
  <div>
    <div class="body" :style="props.style">
      <div class="left">
        <div class="title">
          <ordered-list-outlined />
          <span style="margin-left:10px">{{t('repoIndex.fileList')}}</span>
        </div>
        <a-tree
          style="margin-top:10px"
          :load-data="lsTree"
          :tree-data="treeData"
          @select="selectNode"
          :showLine="true"
          v-model:expandedKeys="expandedKeys"
        />
      </div>
      <div class="right">
        <div class="file-path">
          <template v-for="(item, index) in files" v-bind:key="item">
            <div class="file-path-item" @click="clickFileItem(index)">{{item}}</div>
            <div class="file-path-split" v-if="index < files.length - 1">/</div>
          </template>
        </div>
        <template v-if="showRight">
          <template v-if="showIndex">
            <div class="dir-table">
              <div class="first-line">
                <div class="dir-commit-text flex-center">
                  <ZAvatar
                    :url="dirLatestCommit.avatarUrl"
                    :name="dirLatestCommit.name"
                    :account="dirLatestCommit.committer"
                    :showName="true"
                  />
                  <span>{{dirLatestCommit.commitMsg}}</span>
                </div>
                <div class="dir-commit-text">
                  <span style="margin-right:4px">{{dirLatestCommit.shortCommitId}}</span>
                  <span>{{readableTimeComparingNow(dirLatestCommit.committedTime)}}</span>
                </div>
              </div>
              <div class="dir-line" v-for="item in dirFiles" v-bind:key="item.commit.commitId">
                <div class="dir-line-item dir-line-file" @click="getAndCatFile(item.rawPath)">
                  <folder-outlined v-if="item.mode === 'directory'" style="margin-right:8px" />
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
          </template>
          <template v-else>
            <div class="commit-info">
              <div class="left flex-center">
                <ZAvatar
                  :url="latestCommit.avatarUrl"
                  :name="latestCommit.name"
                  :account="latestCommit.committer"
                  :showName="true"
                />
                <span>{{latestCommit.commitMsg}}</span>
              </div>
              <div class="right">
                <span>{{latestCommit.shortCommitId}}</span>
                <span>{{readableTimeComparingNow(latestCommit.committedTime)}}</span>
              </div>
            </div>
            <div class="code-body">
              <div class="code-top">
                <div class="code-info">
                  <span>{{t('repoIndex.fileSize')}}:</span>
                  <span>{{fileSize}}</span>
                </div>
              </div>
              <div style="max-height:calc(100vh - 240px);overflow:scroll;width:100%">
                <div class="code-code">
                  <Codemirror
                    v-model="fileContent"
                    style="width:100%;"
                    :extensions="extensions"
                    :disabled="true"
                  />
                </div>
              </div>
            </div>
          </template>
        </template>
      </div>
    </div>
  </div>
</template>
<script setup>
import ZAvatar from "@/components/user/ZAvatar";
import {
  OrderedListOutlined,
  FileOutlined,
  FolderOutlined
} from "@ant-design/icons-vue";
import { ref, defineProps, reactive } from "vue";
import { Codemirror } from "vue-codemirror";
import { oneDark } from "@codemirror/theme-one-dark";
import { useRoute } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
import {
  entriesRepoRequest,
  catFileRequest,
  treeRepoRequest
} from "@/api/git/repoApi";
import { useI18n } from "vue-i18n";
import VMdEditor from "@kangc/v-md-editor";
import "@kangc/v-md-editor/lib/style/base-editor.css";
import todoList from "@kangc/v-md-editor/lib/plugins/todo-list/index";
import "@kangc/v-md-editor/lib/plugins/todo-list/todo-list.css";
import githubTheme from "@kangc/v-md-editor/lib/theme/github.js";
import "@kangc/v-md-editor/lib/theme/style/github.css";
VMdEditor.use(githubTheme);
VMdEditor.use(todoList());
const { t } = useI18n();
const route = useRoute();
const props = defineProps(["style"]);
const fileContent = ref("");
const fileSize = ref("");
const extensions = [oneDark];
const files = ref(route.params.files ? route.params.files : []);
const treeData = ref([]);
const showRight = ref(false);
const showIndex = ref(false);
const readmeContent = ref("");
const showAddReadmeContent = ref(false);
// 文件列表
const dirFiles = ref([]);
// 选择文件夹最后一次提交
const dirLatestCommit = reactive({
  committer: "",
  commitMsg: "",
  shortCommitId: "",
  committedTime: "",
  avatarUrl: "",
  name: ""
});
// 展开节点
const expandedKeys = ref([]);
// 最后一次提交
const latestCommit = reactive({
  committer: "",
  commitMsg: "",
  shortCommitId: "",
  committedTime: "",
  avatarUrl: "",
  name: ""
});
const lsTree = treeNode => {
  return new Promise(resolve => {
    if (treeNode.dataRef.children) {
      return;
    }
    getFiles(treeNode.dataRef.key + "/")
      .then(res => {
        if (res.data) {
          treeNode.dataRef.children = res.data.map(item => {
            return {
              title: item.path,
              key: item.rawPath,
              isLeaf: item.mode !== "directory" && item.mode !== "subModule"
            };
          });
          treeData.value = [...treeData.value];
        }
        resolve();
      })
      .catch(() => {
        resolve();
      });
  });
};
// 获取文件列表
const getFiles = dir => {
  return entriesRepoRequest({
    repoId: parseInt(route.params.repoId),
    ref: route.params.ref,
    refType: route.params.refType,
    dir
  });
};
// 树形控件里数据
getFiles("").then(res => {
  if (res.data) {
    treeData.value = res.data.map(item => {
      return {
        title: item.rawPath,
        key: item.rawPath,
        isLeaf: item.mode !== "directory"
      };
    });
  }
});
// 点击叶子节点时触发
const selectNode = (node, e) => {
  if (e.selected) {
    if (e.node.dataRef.isLeaf) {
      getAndCatFile(node[0]);
    } else {
      treeRepo(e.node.key + "/");
    }
  }
};
// 获取文件夹内容
const treeRepo = dir => {
  let path = dir.substring(0, dir.length - 1);
  history.replaceState(
    {},
    "",
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/tree/${route.params.refType}/${route.params.ref}/${path}`
  );
  files.value = path.split("/");
  treeRepoRequest({
    repoId: parseInt(route.params.repoId),
    ref: route.params.ref,
    refType: route.params.refType,
    dir: dir
  }).then(res => {
    showRight.value = true;
    showIndex.value = true;
    if (res.latestCommit) {
      dirLatestCommit.committer = res.latestCommit.committer.account;
      dirLatestCommit.avatarUrl = res.latestCommit.committer.avatarUrl;
      dirLatestCommit.name = res.latestCommit.committer.name;
      dirLatestCommit.commitMsg = res.latestCommit.commitMsg;
      dirLatestCommit.shortCommitId = res.latestCommit.shortId;
      dirLatestCommit.committedTime = res.latestCommit.committedTime;
    }
    if (res.tree && res.tree.files) {
      dirFiles.value = res.tree.files;
    }
    readmeContent.value = res.readmeText;
    showAddReadmeContent.value = res.hasReadme ? true : false;
  });
};
// 获取文件详细内容
const getAndCatFile = (filePath, init) => {
  let path = filePath;
  if (filePath.endsWith("/")) {
    path = path.substring(0, path.length - 1);
  }
  history.replaceState(
    {},
    "",
    `/team/${route.params.teamId}/gitRepo/${route.params.repoId}/tree/${route.params.refType}/${route.params.ref}/${path}`
  );
  files.value = path.split("/");
  getFiles(filePath).then(res => {
    if (!res.data) {
      return;
    }
    if (res.data.length === 1 && res.data[0].mode !== "directory") {
      catFile(filePath).then(res => {
        fileContent.value = res.data.content;
        fileSize.value = res.data.size;
        let commit = res.data.commit;
        if (commit) {
          latestCommit.committer = commit.committer.account;
          latestCommit.commitMsg = commit.commitMsg;
          latestCommit.shortCommitId = commit.shortId;
          latestCommit.committedTime = commit.committedTime;
          latestCommit.avatarUrl = commit.committer.avatarUrl;
          latestCommit.name = commit.committer.name;
        }
      });
      if (init) {
        let keys = [];
        let filesVal = files.value;
        for (let i = 1; i < filesVal.length; i++) {
          keys.push(filesVal.slice(0, i).join("/"));
        }
        expandedKeys.value = keys;
      }
      showIndex.value = false;
      showRight.value = true;
    } else if (res.data.length === 1 && res.data[0].mode === "directory") {
      treeRepo(res.data[0].rawPath + "/");
      showIndex.value = true;
      showRight.value = true;
    } else if (res.data.length > 0) {
      let split = filePath.split("/");
      let path = "";
      let keys = [];
      split.forEach(str => {
        if (path === "") {
          path = str;
        } else {
          path = path + "/" + str;
        }
        keys.push(path);
      });
      expandedKeys.value = keys;
    } else {
      files.value = [];
      showIndex.value = false;
      showRight.value = false;
    }
  });
};
// 获取文件内容
const catFile = filePath => {
  return catFileRequest({
    repoId: parseInt(route.params.repoId),
    ref: route.params.ref,
    refType: route.params.refType,
    filePath
  });
};
// 文件路径导航点击单个文件夹
const clickFileItem = index => {
  if (index < files.value.length - 1) {
    treeRepo(files.value.slice(0, index + 1).join("/") + "/");
  }
};
getAndCatFile(files.value.join("/"), true);
</script>
<style scoped>
.body {
  width: 100%;
  display: flex;
}
.body > .left {
  width: 20%;
  height: calc(100vh - 64px);
  overflow: scroll;
  padding: 10px;
  border-right: 1px solid #d9d9d9;
}
.body > .left > .title {
  height: 32px;
  line-height: 32px;
  font-size: 14px;
}
.body > .right {
  width: 80%;
  overflow: scroll;
  padding: 10px;
  height: calc(100vh - 64px);
}
.file-path {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  width: 100%;
  white-space: nowrap;
  overflow-x: scroll;
}
.file-path-item,
.file-path-split {
  height: 32px;
  line-height: 32px;
  font-size: 14px;
}
.file-path-item {
  padding: 0 8px;
  border-radius: 4px;
}
.file-path-split {
  padding: 0 2px;
}
.code-body {
  border-radius: 4px;
  width: 100%;
  border: 1px solid #d9d9d9;
}
.code-dir {
  border-radius: 4px;
  width: 100%;
  border: 1px solid #d9d9d9;
  max-height: calc(100% - 104px);
  overflow: scroll;
}
.code-top {
  padding: 10px 16px;
  display: flex;
  align-items: center;
}
.code-info > span {
  line-height: 32px;
  font-size: 14px;
  padding-left: 3px;
}
.commit-info {
  line-height: 48px;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  margin-bottom: 10px;
}
.commit-info > .left {
  font-size: 14px;
  max-width: 50%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.commit-info > .left > span + span {
  padding-left: 8px;
}
.commit-info > .right {
  font-size: 14px;
  max-width: 50%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.commit-info > .right > span + span {
  padding-left: 8px;
}
.history-btn {
  display: inline-block;
  line-height: 32px;
  padding: 0 8px;
  cursor: pointer;
  border-radius: 4px;
  margin-left: 4px;
}
.history-btn > span {
  padding-left: 4px;
}
.history-btn:hover {
  background-color: #f0f0f0;
}
.code-code {
  display: flex;
  border-top: 1px solid #d9d9d9;
  min-height: calc(100vh - 236px);
}
.commit-text {
  display: inline-block;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.commit-text:hover {
  text-decoration: underline;
  color: #1677ff;
  cursor: pointer;
}
.dir-commit-text {
  max-width: 50%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.commit-content {
  width: 240px;
}
.commit-content > .title {
  width: 100%;
  font-size: 14px;
  display: flex;
}
.commit-content > .title > span {
  padding-left: 8px;
}
.commit-content > .bottom {
  font-size: 14px;
  line-height: 32px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.commit-content > .bottom > .gray-text {
  font-size: 12px;
  color: gray;
  padding-left: 6px;
}
.commit-content > .bottom > .author-name:hover {
  text-decoration: underline;
  color: #1677ff;
  cursor: pointer;
}
.commit-link-icon {
  display: inline-block;
  width: 20px;
  height: 20px;
  line-height: 20px;
  text-align: center;
}
.commit-msg {
  display: inline-block;
  width: calc(100% - 20px);
  white-space: pre-wrap;
  word-break: break-all;
}
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
.file-path-item:hover {
  cursor: pointer;
  color: #1677ff;
}
</style>