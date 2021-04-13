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
export function _getUserPlaylist() {
  return axios.get("/playlist/user");
}

// 用户登录
export function _userLogin(obj) {
  return axios.post("/user/login", obj);
}

export function _userInfo() {
  return axios.get("/user/me");
}
