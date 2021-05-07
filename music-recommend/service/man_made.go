package service

import "github.com/sta-golang/go-lib-utils/algorithm/data_structure/set"

var manMadeArr = []string{"1463165983", "32507038", "32507038", "27955653", "35528482", "27808044", "412902689",
	"26305527", "26305527", "25727803", "108138", "108138", "26305534", "26305534", "26305550", "26305531", "66282", "28563317", "31877628", "65919",
	"64625", "1293886117", "407450223", "27731176", "27678655", "29710981", "31134170", "400074326", "209643", "209758", "5254811", "209539",
	"209801", "36861903", "36861903", "210049", "191254", "191248", "191060", "191179", "191278", "191268", "28996036", "28427772", "191195", "191134",
	"27789126", "190995", "400161259", "1444738318", "167799", "168039", "168016", "28427707", "1833633769", "486194122", "287063", "287575"}
var manMadeSet *set.StringSet

func init() {
	manMadeSet = set.NewStringSet(len(manMadeArr) << 1)
	manMadeSet.Add(manMadeArr...)
}

func GetManMadeList() []string {
	return manMadeArr
}

func ContainsManMade(id string) bool {
	return manMadeSet.Contains(id)
}
