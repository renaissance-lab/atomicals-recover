# atomicals recover工具介绍如下：

## 背景
在铸造铭文的时候，有时候会出现reveal交易没有发出来，资金被锁死到中间地址的情况。
还有些情况铸造过程不小心用燃烧染色币染色的UTXO来进行铸币了，这个时候我们就需要尽量的保留染色币，同时使用新一个utxo作为手续费，
发起一笔reveal交易，同时保留出原来的染色币和铸造出新的染色币。

## 操作示例
一个找回例子如下：
- commit 交易：https://mempool.space/zh/tx/0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f
- reveal 交易：https://mempool.space/zh/tx/ba4b74892aa0a835520d0db041ac97fca457b58271b936e608ac8f35c00fd7c0

## 操作步骤
1. 找到铸币的commit交易信息，例如https://mempool.space/zh/tx/0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f

- 根据交易信息填写信息
- lockedUtxo:
  - hash: "0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f"
  - index: 0
  - amount: 3700

2. 找到铸币交易的funding地址的私钥 （例子中bc1p04kwemnk0ht38csrc2j7qtx0dqv8kzl0u4vxzsnqg9cdl0rs5aaq3pc8qr地址的私钥，公钥）
- 填写如下信息
- locked: 
    - address: "bc1pameu4wg32ud0u46fxnrzg9d2gm3mhldlcz79gh94xq5jtslvccdq904zhx" // 铸币中间地址
    - WIF: // 铸币交易的funding地址的私钥
    - pubkey: "a61e002383194e6a336aeba0189d10ec93e8e14d3e313cc5e4860c2c10b3d46c" // 铸币交易的funding地址公钥

3. 找到当时铸币time 和nonce信息 填写
- lockedScript:
    - time: 1709005980
    - nonce : 1

4. 从funding地址找一个未染色的utxo作为手续地址，填写
- funding: // funding 地址 
    - address: "bc1p04kwemnk0ht38csrc2j7qtx0dqv8kzl0u4vxzsnqg9cdl0rs5aaq3pc8qr" 
    - WIF: ""
- fundingUtxo: // 支付收付费的utxo信息
    - hash: "88888f7614273607e234787b53a14d15465476a3c7f36870365cacbcba4761eb" 
    - index: 1
    - amount: 14957

5. 填写其他信息
- redeemAddress: "bc1prrkv0qy075asknchdtk8zlnw4le0nc4uw3nssxsmnaf62r654g9qhy2y76" // 找回资金地址
- feeRate : 20 // 手续费聪
- network: 1 // 主网
- token: // 铸币的币对信息
    - protocol: atom
    - opType: dmt
    - bitworkc: "0000"
    - bitworkr: ""
    - mint_ticker: "photon"
    - mint_need_amount: 1000


6. 执行程序：
   - go build
   - ./atomicals-recover
   - 拷贝 rawtx： 之后原始交易数据到链上广播

打赏的话： bc1q3fq7zfzgu98nf6h9cupffcr58a4xpzjq9snshy
