<p align="center">
    <font color=red size=5>cbt</font>
</p>
<p align="center">
  <b>Go</b> smart contract of card brawl for the <a href="https://neo.org">Neo</a> blockchain.
</p>
<hr />

## Features

### Main operations

- `mint`

    initial supply of tokens
    
    Params required:

    1. `to` - destination address
    2. `amount` - supply of tokens.

- `transfer`

    transfer tokens from a address to another address
    
    Params required:

    1. `from` - tokens address
    2. `to` - destination address
    3. `amount` - supply of tokens.

- `exchange`

    exchange gas to token
    
    Params required:

    1. `to` - destination address
    2. `amount` - supply of tokens.

- `sale`

    exchange token to gas
    
    Params required:

    1. `to` - destination address
    2. `amount` - supply of tokens.

- `grantNFT`

    exchange token to nft 
    
    Params required:

    1. `to` - destination address
    2. `id` - nft unqiue id.
    3. `name` - nft's name,
    4. `owner` - nft's owner
    5. `image` - nft's image
    6. `atk` - nft's attack
    7. `def` - nft's defend 
    8. `hp` - nft's hp 

### Other operations

- `setNFTPrice` - set nft price
- `changeTR` - change ratio for token/gas 
