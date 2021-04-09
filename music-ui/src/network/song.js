import axios from './axios'
// 歌曲 

// 发送评论
export function _SendComments (obj) {
  return axios({
    url: '/comment',
    params: {
      ...obj
    }
  })
}

// 获取歌曲的mp3
export function _getSongUrl (id) {
  return axios({
    url: '/song/url',
    params: {
      id
    }
  })
}


// 获取歌词信息
export function _getLyric (id) {
  return axios({
    url: '/lyric',
    params: {
      id
    }
  })
}

// 获取歌曲评论
export function _getSongComment (obj) {
  return axios({
    url: '/comment/music',
    params: {
      ...obj
    }
  })
}

