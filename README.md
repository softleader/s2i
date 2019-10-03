# s2i (Source-to-Image)

The [slctl](https://github.com/softleader/slctl) plugin to build source to image to SoftLeader docker swarm ecosystem

> s2i (Source-to-Image) 是一個非常針對性, 只設計給符合松凌科技開發環境 docker swarm 使用的 command, 請注意: 將來可能會因為全面轉 k8s 而廢棄使用

## Install

```sh
$ slctl plugin install github.com/softleader/s2i
```

## Usage

### prerelease, pre

![](./docs/command-prerelease.svg)

`slctl s2i prerelease` 或 `slctl s2i pre` 的目的是快速的將當前修改的異動版更到開發環境 docker swarm 中, 並且自動的在 github 上將當前的 branch 保留一個版本 (pre-release),
由於這段是在 local 進行, 因此使用前請務必確保當前專案 *local 已經跟 remote 同步過程式了*

請執行 `slctl s2i pre -h` 取得更多說明

> 此 command 僅適用於已經依照[此篇](https://github.com/softleader/softleader-microservice-wiki/wiki/Using-JIB-to-build-image)步驟調整成 jib 包版的專案

### release

![](./docs/command-release.svg)

`slctl s2i release` 的目的是快速化標準的定版程序, 如將手動去 github 下 tag 等多的步驟結合為單一 command

請執行 `slctl s2i release -h` 取得更多說明

> 2-3. 的 update service 需要專案的 Jenkinsfile 配合做些調整, 請參考 [Jenkins Hook to Update Service on Deployer](https://github.com/softleader/softleader-microservice-wiki/wiki/Jenkins-Hook-to-Update-Service-on-Deployer)

### tag

`slctl s2i tag` 的目的是快速的整理某個 repo 下得 tags 及其 releases, 可控制的 sub command 有:

#### tag delete 

`slctl s2i tag delete <TAG..>` 可刪除一至多個 tag 及其 release, 範例:

```sh
# 以互動的問答方式, 詢問所有可控制的問題
slctl s2i tag delete -i

# 在當前目錄的專案中, 刪除名稱 1.0.0 及 1.1.0 的 tag 及 release (完整比對)
slctl s2i tag delete 1.0.0 1.1.0

# 在當前目錄的專案中, 刪除所有名稱為 1 開頭或 2 開頭的 tag 及其 release (以 regex 比對)
slctl s2i tag delete ^1 ^2 -r

# 在當前目錄的專案中, "模擬" 刪除所有 1 開頭的 tag 及其 release (以 regex 比對)
# "模擬" 通常可用於檢視 regex 正確性, 不會真的作用到 GitHub 上
slctl s2i tag delete ^1. -r --dry-run

# 刪除指定專案 github.com/me/my-repo 的所有 tag 及其 release
slctl s2i tag delete .+ -r --source-owner me --source-repo my-repo
```

> 請執行 `slctl s2i tag delete -h` 取得更多說明

#### tag list 

`slctl s2i tag list <TAG..>` 可列出 tag 名稱, 發佈時間及發佈人員, 範例:

```sh
# 在當前目錄的專案中, 以 regex 表示列出更多資訊
slctl s2i tag list ^1. -r
```

> 請執行 `slctl s2i tag list -h` 取得更多說明

## Example

Tag 跟 serviceID 都希望自動找到: 

```sh
slctl s2i pre/release -i
```

已有 serviceID `xxxxx`, 但 tag 希望自動找到:

```sh
slctl s2i pre/release --service-id xxxxx -i
```

已有 tag `v1.2.3` , 但 serviceID 希望協助找到:

```sh
slctl s2i pre/release v1.2.3 -i
```
