<template>
  <div style="padding:14px">
    <div style="margin-bottom:10px">
      <a-radio-group v-model:value="value1">
        <a-radio-button value="all">
          <span>所有</span>
          <span>(0)</span>
        </a-radio-button>
        <a-radio-button value="open">
          <span>已打开</span>
          <span>(0)</span>
        </a-radio-button>
        <a-radio-button value="merged">
          <span>已合并</span>
          <span>(0)</span>
        </a-radio-button>
        <a-radio-button value="a">
          <span>已关闭</span>
          <span>(0)</span>
        </a-radio-button>
      </a-radio-group>
    </div>
    <div class="header">
      <a-input placeholder="搜索合并请求">
        <template #prefix>
          <search-outlined />
        </template>
      </a-input>
      <a-button type="primary" style="margin-left:10px">创建合并请求</a-button>
    </div>
    <ul class="pr-list">
      <li>
        <div class="pr-title">
          <close-circle-outlined />
          <span class="pr-title-span">sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
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
          <span class="pr-title-span">sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
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
          <span class="pr-title-span">sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
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
          <span class="pr-title-span">sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
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
          <span class="pr-title-span">sgojwoiegjoiwejgioejroigjeriogneroigoierhgiorehjgrioejiofew</span>
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
</template>
<script setup>
import {
  SearchOutlined,
  CloseCircleOutlined,
  MessageOutlined
} from "@ant-design/icons-vue";
</script>
<style scoped>
.header {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
}
.pr-list {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
}
.pr-list > li {
  padding: 12px;
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