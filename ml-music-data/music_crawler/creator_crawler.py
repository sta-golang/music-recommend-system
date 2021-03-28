from urllib import parse
from lxml import etree
from urllib3 import disable_warnings
import requests
import sys, getopt
from lxml import html
import json

class CreatorCrawler(object):

    def __init__(self, id):

        self.url = "https://music.163.com/artist?"
        self.id = id
        self.superstar = False
        self.fans_num = 0

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
        self.froms = {
            "id":self.id,
        }
        return self.froms

    # 解析代码，获取有用的数据
    def parsing_codes(self):
        page = etree.HTML(self.code)

        tree3 = html.tostring(page[0], encoding='utf-8').decode('utf-8')
        # # 标题
        self.image = page.xpath('//meta[@property="og:image"]/@content')

        self.description = page.xpath('//meta[@name="description"]/@content')

        self.similar_creator = page.xpath('//a[@class="nm nm-icn f-ib f-thide"]/@href')

        self.home = page.xpath('//a[@class="btn-rz f-tid"]/@href')

        for i, str in enumerate(self.similar_creator):
            arr = str.split("=", 1)
            self.similar_creator[i] = arr[1]

        if len(self.home) == 0:
            self.superstar = True
        else:
            self.parsing_home()

    def parsing_home(self):
        disable_warnings()
        disable_warnings()
        home_url = "https://music.163.com" + self.home[0]
        self.code = requests.get(
            url=home_url,
            headers=self.header,
            data=self.froms,
            verify=False,
        ).text

        page = etree.HTML(self.code)

        temp = page.xpath('//i[@class="tag u-icn2 u-icn2-pfv"]/@class')
        if len(temp) > 0:
            self.superstar = True
        else:
            temp = page.xpath('//strong[@id="fan_count"]/text()')

            self.fans_num = int(temp[0])

    # 获取网页源代码
    def get_code(self):
        disable_warnings()
        disable_warnings()
        self.new_url = self.url + parse.urlencode(self.froms)
        self.code = requests.get(
            url=self.new_url,
            headers=self.header,
            data=self.froms,
            verify=False,
        ).text

    def to_json(self):
        info = {"image_url":self.image[0],
                "description": self.description[0],
                "similar_creator":self.similar_creator,
                "superstar":self.superstar,
                "fans_num":self.fans_num
                }
        bys = json.dumps(info)
        return bys




def main(argv):
    print("进来了")
    music = CreatorCrawler(argv[0])
    music.set_header()  # 调用头部方法，构造请求头信息
    music.set_froms()  # 调用froms方法，构造froms信息
    music.get_code()  # 获取当前页面的源码
    music.parsing_codes()  # 处理源码，获取指定数据
    bys = music.to_json()
    print(bys)

if __name__ == '__main__':
    main(sys.argv[1:])
    sys.stdout.flush()