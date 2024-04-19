<template>
  <div style="padding:10px;">
    <div>
      <BranchTagSelect @select="onBranchTagSelect" />
      <a-popover v-model:open="cloneDownloadVisible" trigger="click" placement="bottomRight">
        <template #content>
          <a-tabs style="width: 300px;padding-bottom:12px" size="small">
            <a-tab-pane key="1" tab="HTTPS">
              <div style="display:flex;align-items:center;margin-top:10px">
                <a-input style="width:90%" :disabled="true" />
                <div class="copy-icon">
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
                <a-input style="width:90%" :disabled="true" />
                <div class="copy-icon">
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
          <a-button type="default" style="width:100%;margin-bottom:12px">下载ZIP</a-button>
        </template>
        <a-button type="primary" style="float:right">
          <span>克隆/下载</span>
          <caret-down-outlined style="font-size:12px" />
        </a-button>
      </a-popover>
    </div>
    <div v-show="showDir">
      <div class="dir-table">
        <div class="first-line">
          <span>Lee init commit</span>
          <span style="float:right">
            <span>huhdfug99</span>
            <span>4个月前</span>
          </span>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>
        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
        </div>

        <div class="dir-line">
          <div
            class="dir-line-item dir-line-file"
          >modules/apifukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div
            class="dir-line-item"
          >feat:fukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkkfukckkkkkkkkkkk</div>
          <div class="dir-line-item" style="text-align:right">4个月前</div>
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
</template>
<script setup>
import { ref } from "vue";
import {
  CaretDownOutlined,
  FileOutlined,
  CopyOutlined
} from "@ant-design/icons-vue";
import VMdEditor from "@kangc/v-md-editor";
import "@kangc/v-md-editor/lib/style/base-editor.css";
import githubTheme from "@kangc/v-md-editor/lib/theme/github.js";
import "@kangc/v-md-editor/lib/theme/style/github.css";
import BranchTagSelect from "@/components/git/BranchTagSelect";
VMdEditor.use(githubTheme);
const readmeContent = ref("fuck you");
const cloneDownloadVisible = ref(false);
const showAddReadmeContent = ref(false);
const showDir = ref(true);
const onBranchTagSelect = event => {
  console.log(event);
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
  padding: 0 10px;
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
</style>