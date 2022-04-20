# plugin_random_setu

## [ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin)的插件

- 该插件维护一个默认长度为100的图片FIFO队列，每隔20s随机拉取pixiv的一张图
- imgFIFOLimit可修改FIFO队列长度，updateInterval可修改图片更新速度
- 输入"随机涩图"即可获取当前队列中的随机一张图

## 插件安装方式
1. 在ZeroBot-Plugin目录下输入 go get github.com/tlogn/plugin_random_setu@main
2. 在ZeroBot-Plugin的main.go中，在import里面加一句 _ "github.com/tlogn/plugin_random_setu/randsetu"
3. 在ZeroBot-Plugin目录下输入 go mod tidy
4. 在ZeroBot-Plugin目录下输入 go run main.go config.go 直接运行，或者输入 sh run.sh 编译后运行
