<template>
  <div style="display:inline-block">
    <img v-if="props.url" :src="props.url" :class="['avatar', avatarSizeMap[props.size]]" />
    <div
      v-else
      :class="['text', textSizeMap[props.size]]"
    >{{props.name?.length > 0?props.name.substring(0,1):''}}</div>
  </div>
</template>

<script setup>
import { defineProps } from "vue";

const props = defineProps(["url", "name", "size"]);

const avatarSizeMap = {
  small: "small-avatar",
  medium: "medium-avatar"
};
const textSizeMap = {
  small: "small-text",
  medium: "medium-text"
};
</script>

<style scoped>
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
  font-size: 18px;
}
.small-text {
  width: 20px;
  height: 20px;
  line-height: 20px;
}
.medium-text {
  width: 40px;
  height: 40px;
  line-height: 38px;
}
</style>