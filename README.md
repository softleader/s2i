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

s2i 會優先試著透過 JIB 來建置 Image, 若發生任何問題會再試著以 `mvn package && docker build` 來建置 Image.

請執行 `slctl s2i pre -h` 取得更多說明

> 建議專案參考 [Using JIB to build image](https://github.com/softleader/softleader-microservice-wiki/wiki/Using-Dockerfile-to-build-cache-layers-image) 或 [Using Dockerfile to build cache layers image](https://github.com/softleader/softleader-microservice-wiki/wiki/Using-Dockerfile-to-build-cache-layers-image) 設定成 cache-layers 的 image

### release

![](./docs/command-release.svg)

`slctl s2i release` 的目的是加速定版程序, 如將手動去 github 下 tag 等多的步驟結合為單一 command

請執行 `slctl s2i release -h` 取得更多說明

> 上圖 2-3 的 update service 需要專案的 Jenkinsfile 配合做些調整, 請參考 [Jenkins Hook to Update Service on Deployer](https://github.com/softleader/softleader-microservice-wiki/wiki/Jenkins-Hook-to-Update-Service-on-Deployer)

### tag

`slctl s2i tag` 的目的是快速的管理某個 repo 下得 tags 及其 releases, 可控制的 sub command 有:

#### tag delete 

`tag delete <TAG..>`, `tag del <TAG..>` 或 `tag rm <TAG..>` 可以協助你刪除不必要的 tag 以及 release, 支援傳一個或多個 `TAG`, 也可以跟其他 command 輕鬆整合, 如: 

```sh
# 依所有 local tag 順序刪除 remote 的 tag 及 release
slctl s2i tag delete $(git tag -l)
```

另外也支援了透過 regex (`--regex`, `-r`) 條件來刪除符合的 tag, 但使用 regex 時需注意: 這是透過 sacn GitHub remote 上所有 tag 來過濾, 因此執行上需要較長的時間 (視 remote tag 多寡決定)

```sh
# 將所有 remote tag 及 release 刪除
slctl s2i tag delete .+ -r
```

請執行 `slctl s2i tag delete -h` 取得更多說明

#### tag list 

`tag list <TAG..>` 可列出 tag 名稱, 發佈時間及發佈人員, 一樣也支援了 regex 的過濾

請執行 `slctl s2i tag list -h` 取得更多說明

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
