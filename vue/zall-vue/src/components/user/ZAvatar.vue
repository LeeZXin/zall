<template>
  <template v-if="props.disablePopover">
    <div class="inline-flex-center" style="margin-right: 6px;">
      <img
        v-if="props.url"
        :src="props.url"
        :class="['avatar', avatarSizeMap[props.size]?avatarSizeMap[props.size]:'small-avatar']"
      />
      <div
        v-else
        :class="['text', textSizeMap[props.size]?textSizeMap[props.size]:'small-text']"
      >{{props.name?.length > 0?props.name.substring(0,1):''}}</div>
      <div
        v-if="props.showName"
        style="font-size:14px;margin-left:6px"
        :class="[nameSizeMap[props.size]?nameSizeMap[props.size]:'small-name']"
      >{{props.name}}</div>
    </div>
  </template>
  <a-popover placement="bottom" v-else>
    <template #content>
      <div class="flex-center" style="width: 200px;font-size:14px;padding:6px">
        <img v-if="props.url" :src="props.url" class="avatar medium-avatar" />
        <div class="text medium-text" v-else>{{props.name?.length > 0?props.name.substring(0,1):''}}</div>
        <div style="margin-left:8px">
          <div style="margin-bottom: 3px" class="no-wrap">{{props.account}}</div>
          <div class="no-wrap">{{props.name}}</div>
        </div>
      </div>
    </template>
    <div class="inline-flex-center" style="margin-right: 6px;cursor:pointer">
      <img
        v-if="props.url"
        :src="props.url"
        :class="['avatar', avatarSizeMap[props.size]?avatarSizeMap[props.size]:'small-avatar']"
      />
      <div
        v-else
        :class="['text', textSizeMap[props.size]?textSizeMap[props.size]:'small-text']"
      >{{props.name?.length > 0?props.name.substring(0,1):''}}</div>
      <div
        v-if="props.showName"
        style="font-size:14px;margin-left:6px"
        :class="[nameSizeMap[props.size]?nameSizeMap[props.size]:'small-name']"
      >{{props.name}}</div>
    </div>
  </a-popover>
</template>

<script setup>
import { defineProps } from "vue";

const props = defineProps([
  "url",
  "name",
  "size",
  "showName",
  "account",
  "disablePopover"
]);

const avatarSizeMap = {
  small: "small-avatar",
  medium: "medium-avatar"
};
const textSizeMap = {
  small: "small-text",
  medium: "medium-text"
};
const nameSizeMap = {
  small: "small-name",
  medium: "medium-name"
};
</script>

<style scoped>
.img-name {
  display: inline-flex;
  align-items: center;
}
.avatar {
  border-radius: 50%;
}
.small-avatar {
  width: 20px;
  height: 20px;
}
.medium-avatar {
  width: 40px;
  height: 40px;
}
.text {
  border-radius: 50%;
  color: white;
  background-color: purple;
  text-align: center;
  font-size: 14px;
}
.small-text {
  width: 20px;
  height: 20px;
  line-height: 18px;
}
.medium-text {
  width: 40px;
  height: 40px;
  line-height: 38px;
}
.small-name {
  line-height: 20px;
}
.medium-name {
  line-height: 40px;
}
</style>