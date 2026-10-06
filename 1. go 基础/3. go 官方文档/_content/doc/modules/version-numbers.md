<!--{
  "Title": "模块版本号"
}-->

模块开发者通过版本号的各个部分，表达版本的稳定性和向后兼容性。每次发布时，版本号应准确反映相较上一版本的变更性质。

开发使用外部模块的代码时，可以借助版本号判断依赖的稳定性，并决定是否升级。开发自己的模块时，版本号则向其他开发者传达你的模块的稳定性和向后兼容性。

本文介绍模块版本号的含义。

**另请参阅**

* 在代码中使用外部包时，可以通过 Go 工具管理这些依赖。详情参见[管理依赖](managing-dependencies)。
* 开发供其他人使用的模块时，需要在发布时通过仓库标签指定版本号。详情参见[发布模块](publishing)。

已发布模块使用遵循语义化版本规则的版本号，如下图所示：

<img src="images/version-number.png"
     alt="语义化版本号示意图：主版本 1、次版本 4、补丁版本 0，以及预发布标识 beta 2"
     style="width: 300px;" />

下表说明版本号各部分所表达的稳定性和向后兼容性。

<table class="DocTable">
  <thead>
    <tr class="DocTable-head">
      <th class="DocTable-cell" width="20%">版本阶段</th>
      <th class="DocTable-cell">示例</th>
      <th class="DocTable-cell">传达给开发者的信息</th>
    </tr>
  </thead>
  <tbody>
    <tr class="DocTable-row">
      <td class="DocTable-cell"><a href="#in-development">开发中</a></td>
      <td class="DocTable-cell">自动生成的伪版本号<p>v<strong>0</strong>.x.x</p></td>
      <td class="DocTable-cell">表示模块仍处于<strong>开发阶段，尚不稳定</strong>，不承诺向后兼容性或稳定性。</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell"><a href="#major">主版本</a></td>
      <td class="DocTable-cell">v<strong>1</strong>.x.x</td>
      <td class="DocTable-cell">表示公共 API 存在<strong>不向后兼容的变更</strong>，不保证与之前的主版本兼容。</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell"><a href="#minor">次版本</a></td>
      <td class="DocTable-cell">vx.<strong>4</strong>.x</td>
      <td class="DocTable-cell">表示公共 API 存在<strong>向后兼容的变更</strong>，保证向后兼容性和稳定性。</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell"><a href="#patch">补丁版本</a></td>
      <td class="DocTable-cell">vx.x.<strong>1</strong></td>
      <td class="DocTable-cell">表示<strong>不影响模块公共 API</strong>及其依赖的变更，保证向后兼容性和稳定性。</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell"><a href="#pre-release">预发布版本</a></td>
      <td class="DocTable-cell">vx.x.x-<strong>beta.2</strong></td>
      <td class="DocTable-cell">表示 <strong>alpha、beta 等预发布阶段</strong>，不承诺稳定性。</td>
    </tr>
  </tbody>
</table>

<a id="in-development" ></a>
## 开发中 {#in-development}

表示模块仍处于开发阶段，**尚不稳定**，不承诺向后兼容性或稳定性。

版本号可能采用以下形式之一：

**伪版本号**

> v0.0.0-20170915032832-14c0d48ead0c

**v0 版本号**

> v0.x.x

<a id="pseudo" ></a>
### 伪版本号 {#pseudo-version-number}

如果模块在仓库中还没有版本标签，Go 工具会生成伪版本号，用于记录在依赖方代码的 go.mod 文件中。

**注意：** 最佳实践是始终让 Go 工具生成伪版本号，不要自行编造。

当模块使用者需要基于一个尚未打上语义化版本标签的提交进行开发时，伪版本号非常有用。

伪版本号由三个部分组成，以连字符分隔，形式如下：

#### 语法 {#syntax}

_baseVersionPrefix_-_timestamp_-_revisionIdentifier_

#### 各部分含义 {#parts}

* **baseVersionPrefix**（vX.0.0 或 vX.Y.Z-0）：根据该修订之前的语义化版本标签推导；没有这样的标签时，使用 vX.0.0。
* **timestamp**：修订创建时的 UTC 时间。在 Git 中使用提交时间，而不是作者时间。时间戳示例见上面的 `20170915032832`。
* **revisionIdentifier**（abcdefabcdef）：提交哈希的前 12 个字符；对于 Subversion，则是用零补齐的修订号。

<a id="v0" ></a>
### v0 版本号 {#v0-number}

以 v0 发布的模块具有正式的语义化版本号，包括主版本、次版本、补丁版本，以及可选的预发布标识。

v0 版本可以用于生产环境，但不保证稳定性或向后兼容性。另外，v1 及后续版本也允许破坏对 v0 使用者的兼容性。因此，在 v1 发布之前，使用 v0 模块的开发者需要自行适配不兼容的变更。

<a id="pre-release" ></a>
## 预发布版本 {#pre-release-version}

表示 alpha、beta 等预发布阶段，不承诺稳定性。

#### 示例 {#example}

```
vx.x.x-beta.2
```

开发者可以在任意“主版本.次版本.补丁版本”组合后添加连字符和预发布标识。

<a id="minor" ></a>
## 次版本 {#minor-version}

表示模块公共 API 的向后兼容变更，保证向后兼容性和稳定性。

#### 示例 {#example-1}

```
vx.4.x
```

次版本会改变模块的公共 API，但不会破坏调用方代码。变更可能包括调整模块自身的依赖，或增加函数、方法、结构体字段、类型。

也就是说，次版本可能通过新增函数提供其他开发者希望使用的增强功能。但如果不使用这些新功能，旧次版本的使用者无需修改自己的代码。

<a id="patch" ></a>
## 补丁版本 {#patch-version}

表示不影响模块公共 API 及其依赖的变更，保证向后兼容性和稳定性。

#### 示例 {#example-2}

```
vx.x.1
```

递增补丁版本号的更新只用于缺陷修复等小范围变更。调用方可以安全升级，无需修改代码。

<a id="major" ></a>
## 主版本 {#major-version}

表示模块公共 API 的不向后兼容变更，不保证与之前的主版本兼容。

#### 示例 {#example-3}

v1.x.x

v1 及以上的版本号表示模块已经可以稳定使用，预发布版本除外。

注意，由于 v0 不承诺稳定性或向后兼容性，从 v0 升级到 v1 的开发者需要自行适配破坏兼容性的变更。

模块开发者应仅在必要时将主版本号提升到 v1 以上，因为这种升级会给使用者带来较大影响：除了公共 API 的不兼容变更，使用者还必须更新代码中所有导入该模块内包的路径。

高于 v1 的主版本更新还会使用新的模块路径，因为模块路径末尾需要追加主版本号，如下所示：

```
module example.com/mymodule/v2 v2.0.0
```

> 译注：上例沿用原文，用于展示模块路径与版本的关系。在实际 go.mod 文件中，`module` 指令只写模块路径：`module example.com/mymodule/v2`，版本号不写在该指令中。

主版本更新会形成一个新模块，其版本历史独立于旧模块。如果你开发模块供其他人使用，参见[模块发布与版本管理流程](release-workflow)中的“发布不兼容的 API 变更”。

关于 module 指令，参见[go.mod 参考](gomod-ref)。
