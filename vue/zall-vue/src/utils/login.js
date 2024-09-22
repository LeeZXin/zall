const tag = "loginUser";

// 从localstorage获取登录用户
export const getLoginUser = () => {
    try {
        return JSON.parse(localStorage.getItem(tag) || "");
    } catch {
        return {};
    }
}

// 保存登录用户信息
export const setLoginUser = user => {
    localStorage.setItem(tag, JSON.stringify(user));
}