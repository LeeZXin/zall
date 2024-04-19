<template>
  <div>
    <div class="container" :style="props.style">
      <div class="left">
        <div class="title">
          <ordered-list-outlined />
          <span style="margin-left:10px">文件列表</span>
        </div>
        <a-tree
          style="margin-top:10px"
          v-model:expandedKeys="expandedKeys"
          v-model:selectedKeys="selectedKeys"
          :load-data="onLoadData"
          :tree-data="treeData"
        />
      </div>
      <div class="right">
        <div class="file-path">
          <div class="file-path-item">zall</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
          <div class="file-path-split">/</div>
          <div class="file-path-item">api.go</div>
        </div>
        <div class="commit-info">
          <div class="left">
            <span>LeeZXin</span>
            <span>feat: some</span>
          </div>
          <div class="right">
            <span>ffff</span>
            <span>-</span>
            <span>四天前</span>
            <div class="history-btn">
              <HistoryOutlined />
              <span>历史</span>
            </div>
          </div>
        </div>
        <div class="code-body" v-if="showCodeBody">
          <div class="code-top">
            <a-radio-group @change="codeOrBlameChange" v-model:value="codeOrBlame">
              <a-radio-button value="code">代码</a-radio-button>
              <a-radio-button value="blame">Blame</a-radio-button>
            </a-radio-group>
            <div class="code-info">
              <span>200</span>
              <span>行</span>
              <span>30</span>
              <span>KB</span>
            </div>
          </div>
          <div class="code-code">
            <ul class="blame-info" v-if="showBlame">
              <li>
                <a-popover>
                  <template #content>
                    <div class="commit-content">
                      <div class="title">
                        <link-outlined />
                        <span>fffff</span>
                      </div>
                      <div class="bottom">
                        <span class="author-name">LeeZXin</span>
                        <span class="gray-text">提交于</span>
                        <span class="gray-text">三个月前</span>
                      </div>
                    </div>
                  </template>
                  <div class="commit-text">ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff</div>
                </a-popover>
                <span>四个月前</span>
              </li>
            </ul>
            <Codemirror
              v-model="code"
              :style="codemirroStyle"
              :extensions="extensions"
              :disabled="true"
            />
          </div>
        </div>
        <div class="code-dir" v-if="!showCodeBody">
          <ul class="code-line">
            <li>名称</li>
            <li>上次提交信息</li>
            <li class="commit-time">提交时间</li>
          </ul>
          <ul class="code-line">
            <li>
              <span class="dir-name">...</span>
            </li>
            <li></li>
            <li class="commit-time"></li>
          </ul>
          <ul class="code-line">
            <li>
              <span class="dir-name">apppp/ddd</span>
            </li>
            <li>nmn</li>
            <li class="commit-time">四天前</li>
          </ul>
          <ul class="code-line">
            <li>
              <span class="dir-name">fff/dddd</span>
            </li>
            <li>nmn</li>
            <li class="commit-time">四天前</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import {
  OrderedListOutlined,
  HistoryOutlined,
  LinkOutlined
} from "@ant-design/icons-vue";
import { ref, defineProps } from "vue";
import { Codemirror } from "vue-codemirror";
import { javascript } from "@codemirror/lang-javascript";
import { oneDark } from "@codemirror/theme-one-dark";
const props = defineProps(["style"]);
const code = ref(
  `cnm`);
const extensions = [javascript(), oneDark];
const expandedKeys = ref([]);
const selectedKeys = ref([]);
const codemirroStyle = ref({ height: "100%", width: "100%" });
const codeOrBlame = ref("code");
const showBlame = ref(false);
const treeData = ref([
  {
    title: "Expand to load",
    key: "0"
  },
  {
    title: "Expand to load",
    key: "1"
  },
  {
    title: "Tree Node",
    key: "2",
    isLeaf: true
  }
]);
const showCodeBody = ref(true);
const onLoadData = treeNode => {
  return new Promise(resolve => {
    if (treeNode.dataRef.children) {
      resolve();
      return;
    }
    setTimeout(() => {
      treeNode.dataRef.children = [
        {
          title: "Child Node",
          key: `${treeNode.eventKey}-0`
        },
        {
          title: "Child Node",
          key: `${treeNode.eventKey}-1`
        }
      ];
      treeData.value = [...treeData.value];
      resolve();
    }, 1000);
  });
};
const codeOrBlameChange = () => {
    if (codeOrBlame.value === "code") {
        showBlame.value = false;
    } else if (codeOrBlame.value === "blame") {
        showBlame.value = true;
    }
};
</script>
<style scoped>
.container {
  width: 100%;
  display: flex;
}
.container > .left {
  width: 20%;
  height: 100%;
  overflow: scroll;
  padding: 10px;
}
.container > .left > .title {
  height: 32px;
  line-height: 32px;
  font-size: 14px;
}
.container > .right {
  width: 80%;
  overflow: scroll;
  padding: 10px;
  border-left: 1px solid #d9d9d9;
  min-height: calc(100vh - 64px);
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
.file-path-item:hover {
  background-color: #f0f0f0;
  cursor: pointer;
}
.code-body {
  border-radius: 4px;
  width: 100%;
  border: 1px solid #d9d9d9;
  height: calc(100% - 104px);
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
  padding: 5px;
  display: flex;
  justify-content: space-between;
  border-radius: 4px;
  border: 1px solid #d9d9d9;
  margin-bottom: 10px;
}
.commit-info > .left {
  line-height: 32px;
  font-size: 14px;
}
.commit-info > .left > span {
  padding-left: 8px;
}
.commit-info > .right {
  line-height: 32px;
  font-size: 14px;
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
.code-line {
  display: flex;
  align-items: center;
  width: 100%;
}
.code-line > li {
  width: 33.33%;
  height: 32px;
  line-height: 32px;
  font-size: 14px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0px 12px;
}
.code-line + .code-line {
  border-top: 1px solid #d9d9d9;
}
.commit-time {
  text-align: right;
}
.dir-name {
  cursor: pointer;
}
.dir-name:hover {
  text-decoration: underline;
  color: #1677ff;
}
.code-code {
  display: flex;
  align-items: center;
  border-top: 1px solid #d9d9d9;
  height: calc(100% - 52px);
}
.blame-info {
  width: 25%;
  height: 100%;
  overflow: scroll;
  padding: 4px;
  min-width: 200px;
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
  font-size: 14px;
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
</style>