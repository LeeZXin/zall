<template>
  <div class="diff-table">
    <div class="header">
      <span class="arrow" @click="clickArrow">
        <right-outlined v-show="!showCode" />
        <down-outlined v-show="showCode"/>
      </span>
      <span style="color:green;margin-right:8px">
        <span>7</span>
        <span>++</span>
      </span>
      <span style="color:red;margin-right:8px">
        <span>6</span>
        <span>--</span>
      </span>
      <span>dfdsf.json</span>
    </div>
    <table v-show="showCode">
      <colgroup>
        <col width="44" />
        <col />
        <col width="44" />
        <col />
      </colgroup>
      <tr>
        <td class="blob-num blob-num-deletion">1</td>
        <td class="blob-code blob-code-deletion">
          <span class="blob-code-inner blob-code-marker" data-code-marker="-">
            <span
              class="x"
            >fucj youfucj youfucj youfucj youfucj youfucj youfucj youfucj youfucj youfucj youfucj youfucj you</span>
          </span>
        </td>
        <td class="blob-num blob-num-addition">1</td>
        <td class="blob-code blob-code-addition">
          <span class="blob-code-inner blob-code-marker" data-code-marker="+">
            <span class="x">fuc</span>
          </span>
        </td>
      </tr>
      <tr>
        <td class="blob-num blob-num-context">2</td>
        <td class="blob-code blob-code-context split-side-left">
          <span class="blob-code-inner blob-code-marker" data-code-marker>yes</span>
        </td>
        <td class="blob-num blob-num-context">2</td>
        <td class="blob-code blob-code-context split-side-right">
          <span class="blob-code-inner blob-code-marker" data-code-marker>yes</span>
        </td>
      </tr>
      <tr>
        <td class="blob-num blob-num-context">3</td>
        <td class="blob-code blob-code-context split-side-left">
          <span class="blob-code-inner blob-code-marker" data-code-marker></span>
        </td>
        <td class="blob-num blob-num-context">3</td>
        <td class="blob-code blob-code-context split-side-right">
          <span class="blob-code-inner blob-code-marker" data-code-marker></span>
        </td>
      </tr>
      <tr>
        <td class="blob-num blob-num-hunk" colspan="1">
          <column-height-outlined />
        </td>
        <td
          class="blob-code blob-code-inner blob-code-hunk"
          colspan="3"
          align="left"
        >@@ 1,23 2,34 @@</td>
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
import { ref } from "vue";
const showCode = ref(false);
const clickArrow = () => {
    showCode.value = !showCode.value;
}
</script>
<style scoped>
.diff-table {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  margin-top: 10px;
}
.diff-table > .header {
  padding: 0 10px;
  line-height: 32px;
  font-size: 14px;
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
.blob-code-marker:before {
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
  font-size: 10px;
}
</style>