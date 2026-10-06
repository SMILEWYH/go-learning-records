---
title: FIPS 140-3 合规支持
layout: article
---

从 Go 1.24 开始，Go 二进制程序可以原生运行在有助于满足 FIPS 140-3 合规要求的模式下。此外，工具链还可以使用冻结版本的密码学包进行构建，这些包共同组成 Go 密码模块（Go Cryptographic Module）。

## FIPS 140-3 {#fips-140-3}

NIST FIPS 140-3 是美国政府针对密码学应用制定的合规体系。其要求包括使用获准的算法，以及使用经过 [CMVP](https://csrc.nist.gov/projects/cryptographic-module-validation-program) 验证、并在目标运行环境中测试过的密码模块。

本页介绍的机制可帮助 Go 应用满足这些合规要求。

不需要满足 FIPS 140-3 合规要求的应用，可以忽略这些机制，也不应启用 FIPS 140-3 模式。

**注意：**仅使用符合 FIPS 140-3 且已获验证的密码模块，本身未必就能满足全部相关监管要求。对于具体用户使用所提供的 FIPS 140-3 模式是否满足特定监管要求，Go 团队无法提供保证或支持。请仔细判断该模块是否满足你的具体要求。

## Go 密码模块 {#the-go-cryptographic-module}

Go 密码模块由 `crypto/internal/fips140/...` 下的一组 Go 标准库包组成，实现了 FIPS 140-3 批准的算法。

`crypto/ecdsa`、`crypto/rand` 等公共 API 包会透明地使用 Go 密码模块来实现 FIPS 140-3 算法。

## FIPS 140-3 模式 {#fips-140-3-mode}

在 FIPS 140-3 模式下运行时：

- Go 密码模块会在 `init` 阶段自动进行完整性自检，将构建时计算的模块目标文件校验和与加载到内存中的符号进行比较。
- 所有算法都会按照相关 FIPS 140-3 实施指南，在 `init` 阶段或首次使用时执行已知答案自检。
- 对生成的密码密钥执行成对一致性测试。注意，某些密钥类型的生成耗时可能因此达到原来的两倍，这对临时密钥尤其有影响。
- [`crypto/rand.Reader`](/pkg/crypto/rand/#Reader) 使用符合 NIST SP 800-90A 的确定性随机比特生成器（DRBG）实现。为确保与未运行在 FIPS 140-3 模式下的程序具有相同安全级别，每次 `Read` 还会从平台的密码学安全伪随机数生成器（CSPRNG）获取随机字节，作为不计入熵额度的附加数据混入输出。
- [`crypto/tls`](/pkg/crypto/tls/) 包会忽略所有未获 FIPS 140-3 批准的协议版本、密码套件、签名算法和密钥交换机制，不与对端协商使用它们。这等价于旧版 Go+BoringCrypto 中需主动启用的 `crypto/tls/fipsonly` 机制。
- [`crypto/rsa.SignPSS`](/pkg/crypto/rsa/#SignPSS) 使用 [`PSSSaltLengthAuto`](/pkg/crypto/rsa/#PSSSaltLengthAuto) 时，盐值长度不会超过哈希值长度。

OpenBSD、Wasm、AIX 和 32 位 Windows 不支持 FIPS 140-3 模式。

## `crypto/fips140` 包 {#the-cryptofips140-package}

[`crypto/fips140.Enabled`](/pkg/crypto/fips140/#Enabled) 函数报告当前是否启用了 FIPS 140-3 模式。

[`crypto/fips140.Version`](/pkg/crypto/fips140/#Version) 函数返回正在使用的 Go 密码模块版本。

## `GOFIPS140` 环境变量 {#the-gofips140-environment-variable}

`GOFIPS140` 环境变量可与 `go build`、`go install` 和 `go test` 一起使用，选择链接到可执行程序中的 Go 密码模块版本，并默认启用 FIPS 140-3 模式。

- `off` 是默认值，使用当前标准库目录树中的 `crypto/internal/fips140/...` 包。
- `latest` 与 `off` 类似，但会默认启用 FIPS 140-3 模式。
- `v1.0.0` 和 `v1.26.0` 分别选择对应版本的 Go 密码模块，并默认启用 FIPS 140-3 模式。
- `inprocess` 和 `certified` 分别等价于：选择已进入 [CMVP Modules In Process List][] 的最新版本，以及已获得 [CMVP validation certificate][] 的最新版本。

[CMVP Modules In Process List]: https://csrc.nist.gov/Projects/cryptographic-module-validation-program/modules-in-process/modules-in-process-list
[CMVP validation certificate]: https://csrc.nist.gov/projects/cryptographic-module-validation-program/validated-modules/search?SearchMode=Basic&ModuleName=Go+Cryptographic+Module&CertificateStatus=Active&ValidationYear=0

## `fips140` GODEBUG 选项 {#the-fips140-godebug-option}

运行时的 `fips140` [GODEBUG](/doc/godebug) 选项控制 Go 密码模块是否以 FIPS 140-3 模式运行。程序启动后不能再修改此选项。

除非构建时设置了 `GOFIPS140`，否则默认值为 `off`。

设为 `on` 时启用 FIPS 140-3 模式，即使构建时没有设置 `GOFIPS140` 也可以这样做。

设为 `only` 时，不符合 FIPS 140-3 的密码算法会返回错误或触发 panic。注意，这是一种尽力检查的模式，用于测试、评估和调试。*它不适合用于生产环境*，也不是安全策略所要求的模式。它按设计会引入崩溃和可能未被处理的错误，并且可能出现误报或漏报。

大多数程序不应直接设置此选项，而应在构建时使用 `GOFIPS140`。

## 模块版本、验证与兼容性 {#module-versions-validations-and-compatibility}

Google 目前与 [Geomys](https://geomys.org/) 签有合作协议，推动 Go 密码模块至少每年进行一次 CMVP 验证。提交验证时，我们会冻结 Go 密码模块，并创建一个新模块版本用于提交。

验证会在全面的运行环境集合中进行测试，覆盖许多常见的操作系统与硬件平台组合。

只要更新的模块版本尚未取得 CMVP 验证证书，旧版 Go 密码模块就会继续受到支持并保持可用。更新版本一旦取得 CMVP 验证证书，旧版本就会被移除。

如果使用从较早 Go 版本冻结而来的 Go 密码模块，部分标准库功能可能不可用并返回错误。

### Go 密码模块 v1.26.0 {#go-cryptographic-module-v1260}

Go 密码模块 v1.26.0 于 2026 年初从 Go 1.26 冻结而来。

它可用于 Go 1.26 及更新版本。

截至 2026-04-28，它在 CMVP Modules In Process List 中处于 Pending Review（待审查）状态，并受 [CAVP 证书 A8028](https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?validation=40638) 覆盖。

#### 相对于 v1.0.0 的变化 {#changes-from-v100}

- 实现了 ML-DSA。
- 开始支持 [testing/cryptotest.SetGlobalRandom](/pkg/testing/cryptotest#SetGlobalRandom)。
- 引入新的 AES-GCM 合规 API，供 `crypto/hpke` 和未来公开的 API 使用。
- Go 密码模块现使用 CPU 抖动熵源，具有 [ESV 证书 #E318](https://csrc.nist.gov/projects/cryptographic-module-validation-program/entropy-validations/certificate/318) 和 [CAVP 证书 A7715](https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?product=20498)。（平台 CSPRNG 仍会作为不计入熵额度的附加数据源，用于所有随机字节的生成。）
- 多项安全性和性能改进。

### Go 密码模块 v1.0.0 {#go-cryptographic-module-v100}

原文记载，Go 密码模块 v1.0.0 于 2024 年初从 Go 1.24 冻结而来。

它可用于 Go 1.24 及更新版本。

它受 [CMVP 证书 #5247](https://csrc.nist.gov/projects/cryptographic-module-validation-program/certificate/5247) 和 [CAVP 证书 A6650](https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/details?product=19371) 覆盖。

## Go+BoringCrypto {#goboringcrypto}

过去通过 BoringCrypto 模块使用某些 FIPS 140-3 批准算法的机制目前仍然可用，但并不受支持。未来版本计划移除该机制，并以本页介绍的机制取代。

Go+BoringCrypto 与原生 FIPS 140-3 模式不兼容。
