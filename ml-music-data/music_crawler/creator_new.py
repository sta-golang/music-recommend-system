from urllib import parse
from lxml import etree
from urllib3 import disable_warnings
import requests
import sys, getopt
from lxml import html
import json
import traceback

class CreatorCrawler(object):

    def __init__(self, id):

        self.url = "https://music.163.com/artist?"
        self.id = id
        self.superstar = False
        self.fans_num = 0

    # 设置请求头部信息(可扩展：不同的User - Agent)
    def set_header(self):
        self. header={'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8',
             'Accept-Encoding': 'gzip, deflate',
             'Accept-Language': 'zh-CN,zh;q=0.9',
             'Connection': 'keep-alive',
             'Cookie': '_iuqxldmzr_=32; _ntes_nnid=0e6e1606eb78758c48c3fc823c6c57dd,1527314455632; '
                       '_ntes_nuid=0e6e1606eb78758c48c3fc823c6c57dd; __utmc=94650624; __utmz=94650624.1527314456.1.1.'
                       'utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); WM_TID=blBrSVohtue8%2B6VgDkxOkJ2G0VyAgyOY;'
                       ' JSESSIONID-WYYY=Du06y%5Csx0ddxxx8n6G6Dwk97Dhy2vuMzYDhQY8D%2BmW3vlbshKsMRxS%2BJYEnvCCh%5CKY'
                       'x2hJ5xhmAy8W%5CT%2BKqwjWnTDaOzhlQj19AuJwMttOIh5T%5C05uByqO%2FWM%2F1ZS9sqjslE2AC8YD7h7Tt0Shufi'
                       '2d077U9tlBepCx048eEImRkXDkr%3A1527321477141; __utma=94650624.1687343966.1527314456.1527314456'
                       '.1527319890.2; __utmb=94650624.3.10.1527319890',
             'Host': 'music.163.com',
             'Referer': 'http://music.163.com/',
             'Upgrade-Insecure-Requests': '1',
             'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) '
                           'Chrome/66.0.3359.181 Safari/537.36'}
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

        self.tmp = page.xpath('//p[@class="note s-fc3"]/text()')

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

    def to_json(self, id, name):
        superType = 1
        if self.superstar:
            superType = 9

        if len(self.tmp) > 1:
            print(id, name, self.tmp)
            info = {
                    "image_url" : 'http://p2.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg',
                    "description" : '未知歌手,可能为用户上传',
                    "similar_creator" : '[]',
                    "superstar": 0,
                    "fans_num": 0,
                    "unknow": True,
                    "id":int(id),
                    "name":name
                }
            return json.dumps(info)

        info = {"image_url":self.image[0],
                "description": self.description[0],
                "similar_creator":self.similar_creator,
                "superstar":superType,
                "fans_num":self.fans_num,
                "id":int(id),
                "name":name}
        return json.dumps(info)





def main(argv):
    print(argv[0])
    music = CreatorCrawler(argv[0])
    music.set_header()  # 调用头部方法，构造请求头信息
    music.set_froms()  # 调用froms方法，构造froms信息
    music.get_code()  # 获取当前页面的源码
    music.parsing_codes()  # 处理源码，获取指定数据
    bys = music.to_json()
    with open("result11.txt", "w") as fp:
        fp.write(bys)
        fp.close()
    sys.stdout.flush()

if __name__ == '__main__':
    with open('../creator.txt', encoding='utf-8') as f:
        fileText = f.read()
        f.close()
    jsonText = json.loads(fileText)
    for i, val in enumerate(jsonText):
        x = val.split('-')
        music = CreatorCrawler(x[0])
        music.set_header()  # 调用头部方法，构造请求头信息
        music.set_froms()  # 调用froms方法，构造froms信息
        music.get_code()  # 获取当前页面的源码
        music.parsing_codes()  # 处理源码，获取指定数据
        bys = music.to_json(x[0],x[1])
        print(bys)
