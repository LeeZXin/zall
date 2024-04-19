<template>
  <a-layout>
    <a-layout-header style="font-size:22px;color:white">
      <span>{{t("app")}}</span>
      <AvatarName style="float:right;" />
      <I18nSelect style="float:right;margin-right: 20px" />
    </a-layout-header>
    <a-layout-content style="height: calc(100vh - 64px); overflow: scroll;padding:0 24px">
      <div class="org-img">
        <img src="http://www.dd.com" />
      </div>
      <div class="org-name">业务开发中心</div>
      <div class="create-team">
        <a-button type="primary" @click="toCreateTeam">{{t("createTeamText")}}</a-button>
      </div>
      <a-list bordered :data-source="data" style="margin-top:20px">
        <template #renderItem="{ item }">
          <a-list-item class="select-item">{{ item }}</a-list-item>
        </template>
        <template #header>
          <div class="list-header">{{t("myTeam")}}</div>
        </template>
      </a-list>
    </a-layout-content>
  </a-layout>
</template>
<script setup>
import I18nSelect from "../components/i18n/I18nSelect";
import AvatarName from "../components/user/AvatarName";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
const { t } = useI18n();
const data = [
  "地盘开发组-前端",
  "地盘开发组-后端",
  "用户运营组-前端",
  "用户运营组-后端"
];
const router = useRouter();
const toCreateTeam = () => {
  router.push("/createTeam");
};
</script>
<style scoped>
.org-img {
  padding: 20px 0;
}
.org-img > img {
  width: 80px;
  height: 80px;
  margin-left: 50%;
  transform: translateX(-50%);
}
.org-name {
  text-align: center;
  font-size: 32px;
  font-weight: bold;
  line-height: 60px;
  height: 60px;
}
.list-header {
  font-weight: bold;
  font-size: 18px;
}
.create-team {
  margin-top: 10px;
}
.select-item {
  cursor: pointer;
}
.select-item:hover {
  background-color: #f0f0f0;
}
</style>