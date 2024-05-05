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
    prCommentRegexp
}