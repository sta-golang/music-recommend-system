import axios from './axios'
// 获取搜索建议
export function _getSearchSuggest (keywords) {
  return axios({
    url: '/search/suggest',
    params: {
      keywords
    }
  })
}

// 获取搜索内容
export function _getSearchList (keywords, type) {
  return axios({
    url: '/search',
    params: {
      keywords,
      type
    }
  })
}


// 获取热搜列表
export function _getSearchHot () {
  const result = axios.get('/search/hot/detail')
  return result
}