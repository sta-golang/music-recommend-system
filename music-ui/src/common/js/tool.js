/**时间格式化函数 */
export function formatDate (date, fmt) {
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(2 - RegExp.$1.length));
  }
  let o = {
    'M+': date.getMonth() + 1,
    'd+': date.getDate(),
    'h+': date.getHours(),
    'm+': date.getMinutes(),
    's+': date.getSeconds()
  };
  for (let k in o) {
    if (new RegExp(`(${k})`).test(fmt)) {
      let str = o[k] + '';
      fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? str : padLeftZero(str));
    }
  }
  return fmt;
};
function padLeftZero (str) {
  return ('00' + str).substr(str.length);
};

// 十万以上截取
export function bignumSlice (num) {
  if ((num - 0) > 100000) {
    let sum = num + ''
    return sum.slice(0, sum.length - 4) + '万'
  }
  return num
}
// 实现深拷贝
export function deepClone (newObj, oldObj) {
  for (var k in oldObj) {
    var item = oldObj[k];
    if (item instanceof Array) {
      newObj[k] = [];
      deepClone(newObj[k], item)
    } else if (item instanceof Object) {
      newObj[k] = {};
      deepClone(newObj[k], item)
    } else {
      newObj[k] = item
    }
  }
}