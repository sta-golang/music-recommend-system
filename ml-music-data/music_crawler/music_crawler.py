from urllib import parse
from lxml import etree
from urllib3 import disable_warnings
import requests
import json
import sys

all_link = ['/playlist?id=19723756','/playlist?id=3779629','/playlist?id=2884035','/playlist?id=3778678','/playlist?id=991319590',]

class Wangyiyun(object):

    def __init__(self, **kwargs):
        # 歌单的歌曲风格
        self.types = kwargs['types']
        # 歌单的发布类型
        # self.years = kwargs['years']
        # 这是当前爬取的页数
        self.pages = pages
        # 这是请求的url参数（页数）
        self.limit = 35
        self.offset = 35 * self.pages - self.limit
        # 这是请求的url
        self.url = "https://music.163.com/discover/playlist/?"



    # 设置请求头部信息(可扩展：不同的User - Agent)
    def set_header(self):
        self.header = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
            "Referer": "https://music.163.com/",
            "Upgrade-Insecure-Requests": '1',
        }
        return self.header

    # 设置请求表格信息
    def set_froms(self):
        self.key = parse.quote(self.types)
        self.froms = {
            "cat": self.key,
            #"order": self.years,
            "limit": self.limit,
            "offset": self.offset,
        }
        return self.froms

    # 解析代码，获取有用的数据
    def parsing_codes(self):
        page = etree.HTML(self.code)
        # 标题
        self.title = page.xpath('//div[@class="u-cover u-cover-1"]/a[@title]/@title')
        # 作者
        self.author = page.xpath('//p/a[@class="nm nm-icn f-thide s-fc3"]/text()')
        # 阅读量
        self.listen = page.xpath('//span[@class="nb"]/text()')
        # 歌单链接
        self.link = page.xpath('//div[@class="u-cover u-cover-1"]/a[@href]/@href')
        if len(self.link) <= 0:
            return False
        all_link.extend([x for x in self.link])
        return True

    # 获取网页源代码
    def get_code(self):
        disable_warnings()
        self.froms['cat'] = self.types
        disable_warnings()
        self.new_url = self.url + parse.urlencode(self.froms)
        self.code = requests.get(
            url=self.new_url,
            headers=self.header,
            data=self.froms,
            verify=False,
        ).text

    # 爬取多页时刷新offset
    def multi(self, page):
        self.offset = self.limit * page - self.limit


if __name__ == '__main__':
    types = ["华语","欧美","日语","韩语","粤语","流行","摇滚","民谣","电子","舞曲","说唱","经典"," 90后","00后","古风","80后"]
    pages = 50
    for type in types:
        music = Wangyiyun(
            types=type,
        )
        for i in range(pages):
            page = i + 1  # 因为没有第0页
            music.multi(page)  # 爬取多页时指定，传入当前页数，刷新offset
            music.set_header()  # 调用头部方法，构造请求头信息
            music.set_froms()  # 调用froms方法，构造froms信息
            music.get_code()  # 获取当前页面的源码
            flag = music.parsing_codes()  # 处理源码，获取指定数据
            if flag == False:
                break
    link_set = set(all_link)
    bys = json.dumps(all_link)
    print(bys)
    sys.stdout.flush()
