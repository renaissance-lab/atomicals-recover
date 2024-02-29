# Introduction to the Atomicals Recover Tool:
## Background
Sometimes the reveal transaction is not sent out and the funds are locked to an intermediate address when minting inscriptions. In other cases, during the minting process, we may accidentally UTXOs  with colored coins for minting. At this time, we need to try to preserve the colored coins, use a new UTXO as for tx fee, send a reveal transaction to retain the original colored coins while minting new colored coins.

## Acknowledge 
  Thanks to [@atomicals](https://github.com/atomicals) technical support and [lyluckyJJ's](https://twitter.com/lyluckyJJ) sponsor.

## Operation Example
An example of recovery is as follows:
- Commit transaction: [0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f](https://mempool.space/zh/tx/0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f)
- Reveal transaction: [ba4b74892aa0a835520d0db041ac97fca457b58271b936e608ac8f35c00fd7c0](https://mempool.space/zh/tx/ba4b74892aa0a835520d0db041ac97fca457b58271b936e608ac8f35c00fd7c0)

## Steps
1. Find the commit transaction information for minting: 
[0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f](https://mempool.space/zh/tx/0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f)

- Fill in config.yml based on the transaction information
    - lockedUtxo:
        - hash: "0000c81399da5e448a72b07d5b0358c84d595fb7819bb438ba9b588bf3c3471f"
        - index: 0
        - amount: 3700

2. Find the private key of the funding address for the coin minting transaction (the private key and public key of the bc1p04kwemnk0ht38csrc2j7qtx0dqv8kzl0u4vxzsnqg9cdl0rs5aaq3pc8qr address in the example).

- Fill in the following information in the config.yml file:

    - locked:
        - address: "bc1pameu4wg32ud0u46fxnrzg9d2gm3mhldlcz79gh94xq5jtslvccdq904zhx" // intermediate address for coin minting
        - WIF: // private key of the funding address for coin minting transaction
        - pubkey: "a61e002383194e6a336aeba0189d10ec93e8e14d3e313cc5e4860c2c10b3d46c" // public key of the funding address for coin minting transaction

3. Find the time and nonce information for the coin minting, and fill in config.yml
    - lockedScript:
        - time: 1709005980
        - nonce : 1

4. Find an uncolored UTXO from the funding address or other address as the tx fee and fill in the config.yml file.

- funding: // funding address
- address: "bc1p04kwemnk0ht38csrc2j7qtx0dqv8kzl0u4vxzsnqg9cdl0rs5aaq3pc8qr"
- WIF: ""
- fundingUtxo: // UTXO information for payment and receipt of fees
- hash: "88888f7614273607e234787b53a14d15465476a3c7f36870365cacbcba4761eb"
- index: 1
- amount: 14957

5. Fill in other information in config.yml file
- redeemAddress: "bc1prrkv0qy075asknchdtk8zlnw4le0nc4uw3nssxsmnaf62r654g9qhy2y76" // address to recover - funds
- feeRate : 20 // transaction fee in satoshis
- network: 1 // mainnet
- token: // token protocol info
    - protocol: atom
    - opType: dmt
    - bitworkc: "0000"
    - bitworkr: ""
    - mint_ticker: "photon"
    - mint_need_amount: 1000

6. Execute the program:
- go build
- ./atomicals-recover
- Copy  output after "rawtx:" ,then broadcast the raw transaction data to the chain.


7. Others :
    How to setup golang env Pls refer to [golang docs](https://go.dev/doc/tutorial/getting-started)


Buy me coffee : bc1pucgnh3j6sy2r9dk9978rcc4pp0p0n0tuft0c0vmjwn8wamkrtzms8gx43u