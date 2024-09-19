// 用户校验
const accountRegexp = /^\w{4,32}$/;
const passwordRegexp = /^.{6,255}$/;
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
const pipelineVarsNameRegexp = /^\w{1,32}$/;
const pipelineVarsContentRegexp = /^[\s\S]{1,10240}$/;
// 发布计划
const deployPlanNameRegexp = /^.{1,32}$/;
// 服务来源
const serviceSourceNameRegexp = /^.{1,32}$/;
const serviceSourceDatasourceRegexp = /^.{1,255}$/;
// 配置来源
const propertySourceNameRegexp = /^.{1,32}$/;
// 注册中心来源
const discoverySourceNameRegexp = /^.{1,32}$/;
// prometheus抓取
const promScrapeEndpointRegexp = /^[\w-]{1,32}$/;
const promScrapeTargetRegexp = /^.{1,}$/;
// db
const dbHostRegexp = /^(\d{1,3}\.){3}\d{1,3}:\d+/;
const dbNameRegexp = /^.{1,32}$/;
const dbUsernameRegexp = /^.+$/;
const dbAccessBaseRegexp = /^.+$/;
const dbAccessTablesRegexp = /^.+$/;
const dbApplyReasonRegexp = /^.{1,255}$/;
const dbDataUpdateCmdRegexp = /^.{1,10240}$/;
// git repo server
const gitRepoServerHostRegexp = /^(\d{1,3}\.){3}\d{1,3}:\d+/;
// env
const envRegexp = /^[a-zA-Z]{1,16}$/;
// ssh key
const sshKeyNameRegexp = /^.{1,128}$/;
// gpg key
const gpgKeyNameRegexp = /^.{1,128}$/;
// zallet node
const zalletNodeIdRegexp = /^[\w-]{1,32}$/;
const zalletNameRegexp = /^.{1,32}$/;
const zalletAgentHostRegexp = /^(\d{1,3}\.){3}\d{1,3}:\d+/;
const zalletAgentTokenRegexp = /^.{0,1024}$/;
// notify tpl
const notifyTplNameRegexp = /^.{1,32}$/;
const notifyTplUrlRegexp = /^https?:\/\/.+$/;
// teamHook
const teamHookUrlRegexp = /^https?:\/\/.+$/;
const teamHookSecretRegexp = /^.{1,1024}$/;
const teamHookNameRegexp = /^.{1,32}$/;
// alert config
const alertConfigNameRegexp = /^.{1,32}$/;
const alertConfigHookUrlRegexp = /^https?:\/\/.+$/;
const alertConfigSecretRegexp = /^.{1,1024}$/;
const alertMysqlHostRegexp = /^(\d{1,3}\.){3}\d{1,3}:\d+/;
const alertHttpHostRegexp = /^https?:\/\/.+$/;
const alertIpPortHostRegexp = /^(\d{1,3}\.){3}\d{1,3}:\d+/;
// wework access token
const weworkAccessTokenNameRegexp = /^.{1,32}$/;
const weworkAccessTokenCorpIdRegexp = /^\S{1,32}$/
const weworkAccessTokenSecretRegexp = /^\S{1,64}$/
    // feishu access token
const feishuAccessTokenNameRegexp = /^.{1,32}$/;
const feishuAccessTokenAppIdRegexp = /^\S{1,32}$/
const feishuAccessTokenSecretRegexp = /^\S{1,64}$/
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
    workflowDescRegexp,
    workflowVarsNameRegexp,
    workflowVarsContentRegexp,
    teamRoleNameRegexp,
    timerTaskNameRegexp,
    appIdRegexp,
    appNameRegexp,
    propertyFileNameRegexp,
    pipelineNameRegexp,
    pipelineVarsNameRegexp,
    pipelineVarsContentRegexp,
    deployPlanNameRegexp,
    serviceSourceNameRegexp,
    serviceSourceDatasourceRegexp,
    propertySourceNameRegexp,
    discoverySourceNameRegexp,
    promScrapeEndpointRegexp,
    promScrapeTargetRegexp,
    dbHostRegexp,
    dbNameRegexp,
    dbUsernameRegexp,
    dbAccessBaseRegexp,
    dbAccessTablesRegexp,
    dbApplyReasonRegexp,
    dbDataUpdateCmdRegexp,
    gitRepoServerHostRegexp,
    envRegexp,
    sshKeyNameRegexp,
    gpgKeyNameRegexp,
    zalletNodeIdRegexp,
    zalletNameRegexp,
    zalletAgentHostRegexp,
    zalletAgentTokenRegexp,
    notifyTplNameRegexp,
    notifyTplUrlRegexp,
    teamHookUrlRegexp,
    teamHookSecretRegexp,
    teamHookNameRegexp,
    alertConfigNameRegexp,
    alertConfigHookUrlRegexp,
    alertConfigSecretRegexp,
    alertMysqlHostRegexp,
    alertHttpHostRegexp,
    alertIpPortHostRegexp,
    weworkAccessTokenNameRegexp,
    weworkAccessTokenCorpIdRegexp,
    weworkAccessTokenSecretRegexp,
    feishuAccessTokenNameRegexp,
    feishuAccessTokenAppIdRegexp,
    feishuAccessTokenSecretRegexp,
}