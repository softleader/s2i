# depl

The [slctl](https://github.com/softleader/slctl) plugin to deploy application to SoftLeader docker swarm ecosystem

> 這是一個非常針對性, 只設計給符合松凌科技開發環境 docker swarm 使用的 command

## Install

```sh
$ slctl plugin install github.com/softleader/depl
```

## Usage

### prerelease, pre

![](./docs/command-prerelease.svg)

`slctl depl prerelease` 或 `slctl depl pre` 的目的是快速的將當前修改的異動版更到開發環境 docker swarm 中, 並且自動的在 github 上將當前的 branch 保留一個版本 (pre-release),
由於這段是在 local 進行, 因此使用前請務必確保當前專案 *local 已經跟 remote 同步過程式了*

請執行 `slctl depl pre -h` 取得更多說明

> 此 command 僅適用於已經依照[此篇](https://github.com/softleader/softleader-microservice-wiki/wiki/Using-JIB-to-build-image)步驟調整成 jib 包版的專案

### release

![](./docs/command-release.svg)

`slctl depl release` 的目的是快速化標準的定版程序, 如將手動去 github 下 tag 等多的步驟結合為單一 command

請執行 `slctl depl release -h` 取得更多說明
