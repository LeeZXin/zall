<template>
  <div :style="props.style" class="body">
    <div class="left">
      <ul class="job-list">
        <li
          v-for="job in jobList"
          :class="{
            'job-item-selected': job === selectedJob
        }"
          @click="selectJob(job)"
          v-bind:key="job"
        >{{job}}</li>
      </ul>
    </div>
    <div class="right">
      <div class="run-status">
        <span>成功</span>
        <span style="margin-left:8px">运行1m12s</span>
      </div>
      <ul class="step-list">
        <li class="step-item-selected">
          <div class="flex-center step-item">
            <span class="left-down-icon">
              <right-outlined />
            </span>
            <check-outlined style="color:green" />
            <span class="step-name">编译</span>
            <span class="step-duration">1m12s</span>
          </div>
          <div class="step-log">
            <table>
              <colgroup>
                <col width="44" />
                <col />
                <col width="44" />
                <col />
              </colgroup>
              <tr>
                <td>1</td>
                <td>cd wihte</td>
              </tr>
              <tr>
                <td>1</td>
                <td>cd wihte</td>
              </tr>
              <tr>
                <td>1</td>
                <td>cd wihte</td>
              </tr>
              <tr>
                <td>1</td>
                <td>cd wihte</td>
              </tr>
              <tr>
                <td>1</td>
                <td>cd wihte</td>
              </tr>
              <tr>
                <td>1</td>
                <td>cd wihte</td>
              </tr>
              <tr>
                <td>1</td>
                <td>cd wihte</td>
              </tr>
            </table>
          </div>
        </li>
        <li>
          <div class="flex-center step-item">
            <span class="left-down-icon">
              <right-outlined />
            </span>
            <check-outlined style="color:green" />
            <span class="step-name">编译</span>
            <span class="step-duration">1m12s</span>
          </div>
        </li>
        <li>
          <div class="flex-center step-item">
            <span class="left-down-icon">
              <right-outlined />
            </span>
            <close-outlined style="color:red" />
            <span class="step-name">编译</span>
            <span class="step-duration">1m12s</span>
          </div>
        </li>
        <li>
          <div class="flex-center step-item">
            <span class="left-down-icon">
              <down-outlined />
            </span>
            <minus-outlined />
            <span class="step-name">编译</span>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>
<script setup>
import { defineProps, ref } from "vue";
import {
  CheckOutlined,
  CloseOutlined,
  MinusOutlined,
  RightOutlined,
  DownOutlined
} from "@ant-design/icons-vue";
const props = defineProps(["style"]);
const jobList = ["job1", "job2", "job3"];
const selectedJob = ref("");

const selectJob = job => {
  selectedJob.value = job;
};
</script>
<style scoped>
.body {
  display: flex;
}
.body > .left {
  width: 20%;
  background-color: white;
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
  border-left: 1px solid #d9d9d9;
  border-top: 1px solid #d9d9d9;
  border-bottom: 1px solid #d9d9d9;
  padding-bottom: 10px;
}
.body > .right {
  width: 80%;
  background-color: #001529;
  color: white;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
  padding: 10px;
}
.job-list > li {
  font-size: 14px;
  white-space: pre-wrap;
  padding: 12px;
  cursor: pointer;
}
.job-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.job-list > li:last-child {
  border-bottom: 1px solid #d9d9d9;
}
.job-item-selected {
  color: white;
  background-color: #1677ff;
}
.run-status {
  margin-bottom: 10px;
  font-size: 12px;
}
.step-list > li {
  font-size: 14px;
}
.step-name {
  margin-left: 20px;
  font-weight: bold;
}
.step-duration {
  margin-left: 12px;
  color: gray;
}
.left-down-icon {
  margin-right: 20px;
  font-size: 12px;
}
.step-item-selected {
  border: 1px solid white;
  border-radius: 4px;
}
.step-log {
  background-color: #2c323a;
  border-bottom-right-radius: 4px;
  border-bottom-left-radius: 4px;
  color: #768391;
}
.step-log > table > tbody > tr > td {
    line-height: 22px;
    font-size: 14px;
}
.step-log > table {
    width: 100%;
}
.step-log > table > tbody > tr > td:first-child {
  text-align: center;
}
.step-item {
  line-height: 38px;
  padding: 0 10px;
  cursor: pointer;
}

</style>