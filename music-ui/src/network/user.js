import axios from "./axios";
// 获取用户详情信息
export function _getUserDetail(uid, cookie) {
  return axios.get("/user/detail", {
    params: {
      cookie,
      uid
    }
  });
}

// 获取用户歌单
export function _getUserSongList(uid) {
  return axios.get("/playlist/user", {
    params: {
      uid
    }
  });
}
export function _getUserPlaylist(username) {
  return axios.get("/playlist/user", {
    params: {
      username
    }
  });
}

// 用户登录
export function _userLogin(obj) {
  return axios.post("/user/login", obj);
}

// 用户注册
export function _userRegister(obj) {
  return axios.post("/user/register", obj);
}

// 验证码
export function _userVerification(obj) {
  return axios.post("/user/code", "username=" + obj);
}

export function _userInfo() {
  return axios.get("/user/me");
}
