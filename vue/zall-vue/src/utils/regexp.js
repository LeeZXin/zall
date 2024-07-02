// 用户校验
const accountRegexp = /^\w{4,32}$/;
const passwordRegexp = /^.{6,}$/;
const usernameRegexp = /^.{1,32}$/;
const emailRegexp = /^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$/;
// 团队校验
const teamNameRegexp = /^.{1,32}$/;
// 仓库校验
const repoNameRegexp = /^[\w-]{1,32}$/;
const defaultBranchRegexp = /^[\w.-]{0,128}$/;
const repoDescRegexp = /^.{0,255}$/;
// 合并请求校验
const prTitleRegexp = /^.{1,255}$/;
const prCommentRegexp = /^.{1,1024}$/;
// 保护分支
const protectedBranchPatternRegexp = /^\S{1,32}$/;
// webhook
const webhookUrlRegexp = /^https?:\/\/.+$/;
const webhookSecretRegexp = /^.{1,1024}$/;
// 工作流
const workflowNameRegexp = /^.{1,32}$/;
const workflowWildBranchRegexp = /^.{1,32}$/;
const workflowBranchRegexp = /^.{1,1024}$/;
const workflowAgentHostRegexp = /^(\d{1,3}\.){3}\d{1,3}:\d+/;
const workflowAgentTokenRegexp = /^.{1,1024}$/;
const workflowDescRegexp = /^.{1,1024}$/;
const workflowVarsNameRegexp = /^\w{1,32}$/;
const workflowVarsContentRegexp = /^[\s\S]{1,10240}$/;
// 角色
const teamRoleNameRegexp = /^.{1,32}$/;
// 定时任务
const timerTaskNameRegexp = /^.{1,64}$/;
// 应用服务
const appIdRegexp = /^[\w-]{1,32}$/;
const appNameRegexp = /^.{1,32}$/;
// 配置中心
const propertyFileNameRegexp = /^[\w-]{1,32}$/;
// 流水线配置
const pipelineNameRegexp = /^.{1,32}$/;
// 发布计划
const deployPlanNameRegexp = /^.{1,32}$/;
const deployProductVersionRegexp = /^.{1,128}$/;
// 服务来源
const serviceSourceNameRegexp = /^.{1,32}$/;
const serviceSourceHostRegexp = /^(\d{1,3}\.){3}\d{1,3}:\d+/;
const serviceSourceApiKeyRegexp = /^.{0,128}$/;
export {
    accountRegexp,
    passwordRegexp,
    usernameRegexp,
    emailRegexp,
    teamNameRegexp,
    repoNameRegexp,
    defaultBranchRegexp,
    repoDescRegexp,
    prTitleRegexp,
    prCommentRegexp,
    protectedBranchPatternRegexp,
    webhookUrlRegexp,
    webhookSecretRegexp,
    workflowNameRegexp,
    workflowWildBranchRegexp,
    workflowBranchRegexp,
    workflowAgentHostRegexp,
    workflowAgentTokenRegexp,
    workflowDescRegexp,
    workflowVarsNameRegexp,
    workflowVarsContentRegexp,
    teamRoleNameRegexp,
    timerTaskNameRegexp,
    appIdRegexp,
    appNameRegexp,
    propertyFileNameRegexp,
    pipelineNameRegexp,
    deployPlanNameRegexp,
    deployProductVersionRegexp,
    serviceSourceNameRegexp,
    serviceSourceHostRegexp,
    serviceSourceApiKeyRegexp
}