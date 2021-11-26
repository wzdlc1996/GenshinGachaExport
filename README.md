# README

## Build on Windows

执行命令

```sh
env GOOS=windows GOARCH=amd64 go build
```

测试在 `windows 11` 上通过

## 使用

确保在原神游戏中打开一次祈愿记录. 然后双击运行 `GenshinGachaExport.exe`. 它将会在当前目录下生成若干文件:
1.  `GachaLog*.json` 它们以json格式存储所有祈愿记录, 包含其完整信息, 涉及到用户自己的uid请谨慎发布
2.  `data.xlsx` 它以xlsx格式保存抽卡数据, 不包含敏感信息, 没有编辑格式. 但可以直接使用于 [原神抽卡记录分析工具](https://genshin-gacha-analyzer.vercel.app/)

## 鸣谢

2.  [github.com/sunfkny/genshin-gacha-export-js](https://github.com/sunfkny/genshin-gacha-export-js/blob/main/index.js)
1.  [github.com/voderl/genshin-gacha-analyzer](https://github.com/voderl/genshin-gacha-analyzer)
3.  [github.com/qax-os/excelize](https://github.com/qax-os/excelize)