<template>
  <div>
    <div v-show="!showList">
      <CreatePr />
    </div>
    <div v-show="showList">
      <div class="header">
        <a-input placeholder="搜索合并请求">
          <template #prefix>
            <search-outlined />
          </template>
        </a-input>
        <a-button type="primary" style="margin-left:10px">创建合并请求</a-button>
      </div>
      <div class="body">
        <div class="title">
          <div class="open-num" :style="openClosedSelected ? {color: 'black'}:{}">
            <apartment-outlined />
            <span>177324</span>
            <span>Open</span>
          </div>
          <div class="closed-num" :style="!openClosedSelected ? {color: 'black'}:{}">
            <check-outlined />
            <span>177324</span>
            <span>Closed</span>
          </div>
        </div>
        <ul class="pr-list">
          <li>
            <div class="pr-title">
              <close-circle-outlined />
              <span
                class="pr-title-span"
              >sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
              <span class="pr-title-right">
                <message-outlined />
                <span class="message-num">3</span>
              </span>
            </div>
            <div class="pr-desc">
              <a-tag color="purple">Closed</a-tag>
              <span>#13923 by ywwg was merged 2 days ago • Approved</span>
            </div>
          </li>
          <li>
            <div class="pr-title">
              <close-circle-outlined />
              <span
                class="pr-title-span"
              >sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
              <span class="pr-title-right">
                <message-outlined />
                <span class="message-num">3</span>
              </span>
            </div>
            <div class="pr-desc">
              <a-tag color="purple">Closed</a-tag>
              <span>#13923 by ywwg was merged 2 days ago • Approved</span>
            </div>
          </li>
          <li>
            <div class="pr-title">
              <close-circle-outlined />
              <span
                class="pr-title-span"
              >sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
              <span class="pr-title-right">
                <message-outlined />
                <span class="message-num">3</span>
              </span>
            </div>
            <div class="pr-desc">
              <a-tag color="purple">Closed</a-tag>
              <span>#13923 by ywwg was merged 2 days ago • Approved</span>
            </div>
          </li>
          <li>
            <div class="pr-title">
              <close-circle-outlined />
              <span
                class="pr-title-span"
              >sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
              <span class="pr-title-right">
                <message-outlined />
                <span class="message-num">3</span>
              </span>
            </div>
            <div class="pr-desc">
              <a-tag color="purple">Closed</a-tag>
              <span>#13923 by ywwg was merged 2 days ago • Approved</span>
            </div>
          </li>
          <li>
            <div class="pr-title">
              <close-circle-outlined />
              <span
                class="pr-title-span"
              >sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
              <span class="pr-title-right">
                <message-outlined />
                <span class="message-num">3</span>
              </span>
            </div>
            <div class="pr-desc">
              <a-tag color="purple">Closed</a-tag>
              <span>#13923 by ywwg was merged 2 days ago • Approved</span>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
<script setup>
import {
  SearchOutlined,
  ApartmentOutlined,
  CheckOutlined,
  CloseCircleOutlined,
  MessageOutlined
} from "@ant-design/icons-vue";
import CreatePr from "@/components/git/CreatePullRequest";
import { ref } from "vue";
const showList = ref(false);
const openClosedSelected = ref(false);
</script>
<style scoped>
.header {
  margin: 10px 0px;
  display: flex;
  align-items: center;
}
.body {
  border-radius: 4px;
  border: 1px solid #d9d9d9;
}
.body > .title {
  display: flex;
  align-items: center;
  height: 32px;
  line-height: 32px;
  padding-left: 10px;
  font-size: 14px;
}
.closed-num {
  margin-left: 18px;
}
.closed-num > span + span,
.open-num > span + span {
  margin-left: 4px;
}
.open-num,
.closed-num {
  color: gray;
  cursor: pointer;
}
.open-num:hover,
.closed-num:hover {
  color: #1677ff;
}
.pr-list {
  border-top: 1px solid #d9d9d9;
}
.pr-list > li {
  padding: 6px 12px;
}
.pr-list > li + li {
  border-top: 1px solid #d9d9d9;
}
.pr-list > li:hover {
  cursor: pointer;
  background-color: #f0f0f0;
}
.pr-title {
  font-size: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.pr-title-span {
  margin-left: 6px;
  font-weight: bold;
}
.pr-title-right {
  float: right;
  font-size: 14px;
}
.pr-desc {
  margin-top: 6px;
  color: gray;
  font-size: 12px;
}
.pr-title-right > .message-num {
  padding-left: 4px;
}
</style>