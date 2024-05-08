<template>
  <div>
    <div class="body" :style="props.style">
      <div class="left">
        <div class="title">
          <ordered-list-outlined />
          <span style="margin-left:10px">文件列表</span>
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
            <div class="file-path-item">{{item}}</div>
            <div class="file-path-split" v-if="index < files.length - 1">/</div>
          </template>
        </div>
        <template v-if="showRight">
          <div class="commit-info">
            <div class="left">
              <span>{{latestCommit.committer}}</span>
              <span>{{latestCommit.commitMsg}}</span>
            </div>
            <div class="right">
              <span>{{latestCommit.shortCommitId}}</span>
              <span>{{readableTimeComparingNow(latestCommit.committedTime)}}</span>
            </div>
          </div>
          <div class="code-body">
            <div class="code-top">
              <a-radio-group @change="codeOrBlameChange" v-model:value="codeOrBlame">
                <a-radio-button value="code">代码</a-radio-button>
                <a-radio-button value="blame">Blame</a-radio-button>
              </a-radio-group>
              <div class="code-info">
                <span>文件大小:</span>
                <span>{{fileSize}}</span>
              </div>
            </div>
            <div style="max-height:calc(100vh - 234px);overflow:scroll;width:100%">
              <div class="code-code">
                <ul class="blame-info" v-if="showBlame">
                  <li
                    v-for="item in blameList"
                    v-bind:key="item.commit.commitId"
                    :style="blameLineStyle"
                  >
                    <a-popover>
                      <template #content>
                        <div class="commit-content">
                          <div class="title">
                            <div class="commit-link-icon">
                              <link-outlined />
                            </div>
                            <div class="commit-msg">{{item.commit.commitMsg}}</div>
                          </div>
                          <div class="bottom">
                            <span class="author-name">{{item.commit.committer.account}}</span>
                            <span class="gray-text">提交于</span>
                            <span
                              class="gray-text"
                            >{{readableTimeComparingNow(item.commit.committedTime)}}</span>
                          </div>
                        </div>
                      </template>
                      <div class="commit-text">{{item.commit.committer.account}}</div>
                    </a-popover>
                    <span>{{readableTimeComparingNow(item.commit.committedTime)}}</span>
                  </li>
                </ul>
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
      </div>
    </div>
  </div>
</template>
<script setup>
import { OrderedListOutlined, LinkOutlined } from "@ant-design/icons-vue";
import { ref, defineProps, reactive } from "vue";
import { Codemirror } from "vue-codemirror";
import { javascript } from "@codemirror/lang-javascript";
import { oneDark } from "@codemirror/theme-one-dark";
import { useRoute } from "vue-router";
import { readableTimeComparingNow } from "@/utils/time";
import {
  entriesRepoRequest,
  catFileRequest,
  blameRequest
} from "@/api/git/repoApi";
const route = useRoute();
const props = defineProps(["style"]);
const fileContent = ref("");
const fileSize = ref("");
const extensions = [oneDark, javascript()];
const codeOrBlame = ref("code");
const showBlame = ref(false);
const files = ref(route.params.files ? route.params.files : []);
const treeData = ref([]);
const showRight = ref(false);
// 展开节点
const expandedKeys = ref([]);
// 最后一次提交
const latestCommit = reactive({
  committer: "",
  commitMsg: "",
  shortCommitId: "",
  committedTime: ""
});
const blameLineStyle = ref({});
// 每行提交信息
const blameList = ref([]);
let hasBlame = false;
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
              isLeaf: item.mode !== "directory"
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
const codeOrBlameChange = () => {
  if (codeOrBlame.value === "code") {
    showBlame.value = false;
  } else if (codeOrBlame.value === "blame") {
    showBlame.value = true;
    if (!hasBlame) {
      blameRequest({
        repoId: parseInt(route.params.repoId),
        ref: route.params.ref,
        refType: route.params.refType,
        filePath: files.value.join("/")
      }).then(res => {
        hasBlame = true;
        blameList.value = res.data;
        if (res.data && res.data.length > 0) {
          let nodes = document.getElementsByClassName("cm-gutterElement");
          if (nodes && nodes.length > 0) {
            for (let index in nodes) {
              let h = nodes[index].style.height;
              if (h && h !== "0px" && h !== "0") {
                blameLineStyle.value = {
                  lineHeight: h
                };
                break;
              }
            }
          }
        }
      });
    }
  }
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
  if (e.selected && e.node.dataRef.isLeaf) {
    let filePath = node[0];
    history.replaceState(
      {},
      "",
      `/gitRepo/${route.params.repoId}/tree/${route.params.refType}/${route.params.ref}/${filePath}`
    );
    files.value = filePath.split("/");
    getAndCatFile(filePath);
    showBlame.value = false;
    hasBlame = false;
    codeOrBlame.value = "code";
  }
};
// 获取文件详细内容
const getAndCatFile = (filePath, init) => {
  getFiles(filePath).then(res => {
    if (res.data && res.data.length === 1) {
      if (res.data[0].mode !== "directory") {
        catFile(filePath).then(res => {
          fileContent.value = res.data.content;
          fileSize.value = res.data.size;
          let commit = res.data.commit;
          if (commit) {
            latestCommit.committer = commit.committer.account;
            latestCommit.commitMsg = commit.commitMsg;
            latestCommit.shortCommitId = commit.shortId;
            latestCommit.committedTime = commit.committedTime;
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
        showRight.value = true;
      } else {
        expandedKeys.value = [filePath];
      }
    } else if (!res.data || res.data.length === 0) {
      files.value = [];
      showRight.value = false;
    }
  });
};
const catFile = filePath => {
  return catFileRequest({
    repoId: parseInt(route.params.repoId),
    ref: route.params.ref,
    refType: route.params.refType,
    filePath
  });
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
  padding: 10px;
  display: flex;
  align-items: center;
}
.code-info {
  margin-left: 8px;
}
.code-info > span {
  line-height: 32px;
  font-size: 14px;
  padding-left: 3px;
}
.commit-info {
  padding: 5px 10px;
  display: flex;
  justify-content: space-between;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  margin-bottom: 10px;
}
.commit-info > .left {
  line-height: 32px;
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
  line-height: 32px;
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
.blame-info {
  width: 25%;
  min-height: calc(100vh - 236px);
  padding: 4px;
  min-width: 200px;
  background-color: #282c34;
  color: #abb2bf;
}
.blame-info > li {
  font-size: 12px;
  line-height: 1.4;
  width: 100%;
  display: flex;
  justify-content: space-between;
  padding: 0 6px;
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
</style>