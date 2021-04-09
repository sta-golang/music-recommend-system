export function _Swiper (curUl, curOl) {
  var timer = null
  if (!curOl) return
  if (!curUl) return
  const btns = [...curOl.querySelectorAll('li')]
  const sz = [...curUl.querySelectorAll('li')]
  const next = curUl.querySelector('.next')
  const pre = curUl.querySelector('.pre')
  for (let i = 1; i <= sz.length; i++) {
    sz[i - 1].id = i
    btns[i - 1].name = i
  }
  // 最后一位的索引
  const len = sz.length - 1
  reset()
  syncBtns()
  // 下一张其实就是把第一张放到最后、
  function getNext () {

    const giveUp = sz.shift()
    sz.push(giveUp)
    // 重置一下
    for (let i = 0; i < sz.length; i++) {
      sz[i].style.zIndex = i
      sz[i].style.transform = 'scale(1)'
      sz[i].style.filter = 'contrast(50%)'
    }
    reset()
    syncBtns()

  }
  function getPre () {

    const giveUp = sz.pop()
    sz.unshift(giveUp)
    // 重置一下
    for (let i = 0; i < sz.length; i++) {
      sz[i].style.zIndex = i
      sz[i].style.transform = 'scale(1)'
      sz[i].style.filter = 'contrast(50%)'
    }
    reset()
    syncBtns()


  }
  function reset () {
    // 通过操作数组来控制
    // 就是把最后的设置left400
    // 倒数第二的设置成zIndex = 100 left = '257px transform = 'scale(1.3)'
    // 倒数第三的设置left 200
    // 先把所有放到257px  藏起来
    for (let i = 0; i < sz.length; i++) {
      sz[i].style.left = '257px'
    }
    sz[len].style.left = '0px'
    sz[0].style.zIndex = 100
    sz[0].style.left = '257px'
    sz[0].style.transform = 'scale(1.3)'
    sz[0].style.opacity = 1
    sz[0].style.filter = 'contrast(100%)'
    sz[1].style.left = '515px'
  }
  next.addEventListener('click', function () {
    clearInterval(timer)
    getNext()
    timer = setInterval(() => {
      getNext()
    }, 3000)
  })
  pre.addEventListener('click', function () {
    clearInterval(timer)
    getPre()
    timer = setInterval(() => {
      getPre()
    }, 3000)
  })
  curUl.onmouseover = function () {
    clearInterval(timer)
  }
  curUl.onmouseout = function () {
    clearInterval(timer)
    timer = setInterval(() => {
      getNext()
    }, 3000)
  }
  function syncBtns () {
    for (let i = 0; i < btns.length; i++) {
      if (btns[i].name == sz[0].id) {
        btns[i].style.backgroundColor = 'red'
      } else {
        btns[i].style.backgroundColor = '#999'
      }
    }
  }

  // 经过滑块的时候
  for (let i = 0; i < btns.length; i++) {
    btns[i].addEventListener('mouseenter', function () {
      clearInterval(timer)
      const len1 = sz[0].id
      const len2 = btns[i].name
      // 比较之间的差值
      let dis = Math.max(len1, len2) - Math.min(len1, len2)
      // 根据结果判断是往前还是往后
      if (len1 > len2) {
        while (dis--) {
          // 循环直到差值为0就退出循环
          getPre()
        }
      } else {
        while (dis--) {
          getNext()
        }
      }
    })
  }
  // 自动播放
  timer = setInterval(() => {
    getNext()
  }, 3000)
}
