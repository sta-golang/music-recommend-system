module github.com/sta-golang/ml-music-data

go 1.14

require (
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/sta-golang/go-lib-utils v0.3.3
	github.com/sta-golang/music-recommend v0.0.0
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/tencentyun/cos-go-sdk-v5 v0.7.24
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
)

replace github.com/sta-golang/music-recommend => ../music-recommend

replace github.com/sta-golang/music-algorithm => ../music-algorithm
