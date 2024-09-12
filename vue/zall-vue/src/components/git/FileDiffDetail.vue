<template>
  <div class="diff-table">
    <div class="header">
      <span class="arrow" @click="clickArrow">
        <right-outlined v-show="!showCode" />
        <down-outlined v-show="showCode" />
      </span>
      <span style="color:green;margin-right:12px">
        <span>{{props.stat.insertNums}}</span>
        <span>++</span>
      </span>
      <span style="color:red;margin-right:12px">
        <span>{{props.stat.deleteNums}}</span>
        <span>--</span>
      </span>
      <span>{{props.stat.rawPath}}</span>
    </div>
    <div class="is-binary-text" v-show="showCode && loadData && isBinary">{{t('binary')}}</div>
    <table v-show="showCode && loadData && !isBinary">
      <colgroup>
        <col width="44" />
        <col />
        <col width="44" />
        <col />
      </colgroup>
      <tr v-for="(item, index) in diffLines" v-bind:key="index">
        <template v-if="item.prefix === '-'">
          <td class="blob-num blob-num-deletion">{{item.leftNo}}</td>
          <td class="blob-code blob-code-deletion">
            <span class="blob-code-inner blob-code-marker" data-code-marker="-">
              <span class="x">{{item.text}}</span>
            </span>
          </td>
          <td class="blob-num blob-num-deletion"></td>
          <td class="blob-code blob-code-deletion">
            <span class="blob-code-inner blob-code-marker">
              <span class="x"></span>
            </span>
          </td>
        </template>
        <template v-else-if="item.prefix === '*'">
          <td class="blob-num blob-num-deletion">{{item.leftNo}}</td>
          <td class="blob-code blob-code-deletion">
            <span class="blob-code-inner blob-code-marker" data-code-marker="-">
              <span class="x">{{item.text}}</span>
            </span>
          </td>
          <td class="blob-num blob-num-addition">{{item.rightNo}}</td>
          <td class="blob-code blob-code-addition">
            <span class="blob-code-inner blob-code-marker" data-code-marker="+">
              <span class="x">{{item.updateText}}</span>
            </span>
          </td>
        </template>
        <template v-else-if="item.prefix === '+'">
          <td class="blob-num blob-num-addition"></td>
          <td class="blob-code blob-code-addition">
            <span class="blob-code-inner">
              <span class="x"></span>
            </span>
          </td>
          <td class="blob-num blob-num-addition">{{item.rightNo}}</td>
          <td class="blob-code blob-code-addition">
            <span class="blob-code-inner blob-code-marker" data-code-marker="+">
              <span class="x">{{item.text}}</span>
            </span>
          </td>
        </template>
        <template v-else-if="item.prefix === ' '">
          <td class="blob-num blob-num-context">{{item.leftNo}}</td>
          <td class="blob-code blob-code-context split-side-left">
            <span class="blob-code-inner blob-code-marker" data-code-marker>{{item.text}}</span>
          </td>
          <td class="blob-num blob-num-context">{{item.rightNo}}</td>
          <td class="blob-code blob-code-context split-side-right">
            <span class="blob-code-inner blob-code-marker" data-code-marker>{{item.text}}</span>
          </td>
        </template>
        <template v-else>
          <td class="blob-num blob-num-hunk" colspan="1">
            <column-height-outlined />
          </td>
          <td
            class="blob-code blob-code-inner blob-code-hunk"
            colspan="3"
            align="left"
          >{{item.text}}</td>
        </template>
      </tr>
    </table>
  </div>
</template>
<script setup>
import {
  ColumnHeightOutlined,
  RightOutlined,
  DownOutlined
} from "@ant-design/icons-vue";
import { ref, defineProps } from "vue";
import { diffFileRequest } from "@/api/git/repoApi";
import { useI18n } from "vue-i18n";
const { t } = useI18n();
const showCode = ref(false);
const loadData = ref(false);
const isBinary = ref(false);
const diffLines = ref([]);
const clickArrow = () => {
  let show = !showCode.value;
  if (show && !loadData.value) {
    loadData.value = true;
    diffFileRequest({
      repoId: props.repoId,
      target: props.target,
      head: props.head,
      filePath: props.stat.rawPath
    }).then(res => {
      isBinary.value = res.data.isBinary;
      let lines = res.data.lines;
      let addMap = {};
      let delMap = {};
      lines.forEach(item => {
        if (item.prefix === "+") {
          addMap[item.rightNo] = item;
        } else if (item.prefix === "-") {
          delMap[item.leftNo] = item;
        }
      });
      let ret = [];
      lines.forEach(item => {
        let line = { ...item };
        if (item.prefix === "-") {
          let add = addMap[item.rightNo];
          if (add) {
            line.updateText = add.text;
            line.prefix = "*";
          }
          ret.push(line);
        } else if (item.prefix === "+") {
          let del = delMap[item.leftNo];
          if (!del) {
            ret.push(line);
          }
        } else {
          ret.push(line);
        }
      });
      diffLines.value = ret;
    });
  }
  showCode.value = show;
};
const props = defineProps(["stat", "head", "target", "repoId"]);
</script>
<style scoped>
.diff-table {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  margin-top: 10px;
}
.diff-table > .header {
  padding: 12px;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  word-break: break-all;
  white-space: nowrap;
}
.diff-table * {
  position: static;
}
.diff-table > table,
.diff-table {
  width: 100%;
}
.diff-table > table {
  border-spacing: 0;
  table-layout: fixed;
  border-top: 1px solid #d9d9d9;
  margin-bottom: 4px;
}
.blob-code-deletion,
.blob-code-addition,
.blob-code-context {
  padding-left: 22px !important;
}
.blob-code {
  position: relative;
  padding-right: 10px;
  padding-left: 10px;
  line-height: 20px;
  vertical-align: top;
}
.blob-code-deletion {
  background-color: #ffebe9;
  outline: 1px dashed transparent;
}
.blob-code-inner {
  display: table-cell;
  overflow: visible;
  font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas,
    Liberation Mono, monospace;
  font-size: 12px;
  color: #1f2328;
  word-wrap: anywhere;
  white-space: pre-wrap;
}
.blob-code + .blob-num {
  border-left: 1px solid #d9d9d9;
}
.blob-code-addition {
  background-color: #e6ffec;
  outline: 1px dotted transparent;
}
.blob-code-addition .x {
  color: #1f2328;
  background-color: #abf2bc;
}
.blob-num {
  position: relative;
  width: 1%;
  min-width: 50px;
  padding-right: 10px;
  padding-left: 10px;
  font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas,
    Liberation Mono, monospace;
  font-size: 12px;
  line-height: 20px;
  color: var(--color-fg-subtle);
  text-align: right;
  white-space: nowrap;
  vertical-align: top;
  cursor: pointer;
  -webkit-user-select: none;
  user-select: none;
}
.blob-num-addition {
  color: #1f2328;
  background-color: #ccffd8;
  border-color: #1f883d;
}
.blob-num-deletion {
  color: #1f2328;
  background-color: #ffd7d5;
  border-color: #cf222e;
}
.blob-code-marker::before {
  position: absolute;
  top: 1px;
  left: 8px;
  padding-right: 8px;
  content: attr(data-code-marker);
}
.blob-num-hunk {
  background-color: rgba(84, 174, 255, 0.4);
}
.blob-code-hunk {
  background-color: #ddf4ff;
}
.blob-code-deletion .x {
  color: #1f2328;
  background-color: rgba(255, 129, 130, 0.4);
}
.arrow {
  cursor: pointer;
  margin-right: 8px;
  font-size: 12px;
}
.is-binary-text {
  border-top: 1px solid #d9d9d9;
  line-height: 80px;
  font-size: 14px;
  text-align: center;
  width: 100%;
}
</style>