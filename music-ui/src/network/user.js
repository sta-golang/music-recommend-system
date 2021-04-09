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
  return axios.get("/user/playlist", {
    params: {
      uid
    }
  });
}

// 用户登录
export function _userLogin(obj) {
  return axios.post("/user/login", obj);
}

export function _userInfo() {
  return axios.get("/user/me");
}
