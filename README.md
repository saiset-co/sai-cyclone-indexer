# saiCycloneIndexer

Utility for viewing transactions of specified addresses in Cyclone SDK based blockchains.
If added address found in the transaction, this transaction will be saved to the storage and sent to notification address.

## Configurations
**config.yml** - common saiService config file.

### Common block
- `http` - http server section
  - `enabled` - enable or disable http handlers
  - `port`    - http server port

### Storage block
- `url` - sai-storage http server address
- `token` - sai-storage token
- `collection` - sai-storage collection name

### Cyclone block
- `node_address` - cyclone node address for API calls
- `start_block` - start block height
- `sleep_duration` - sleep duration between loop iteration(in seconds)

### Notifier block
- `url` - bridge service url
- `token` - bridge service token
- `sender_id` - sender name: CYCLONE

## How to run
`make build`: rebuild and start service  
`make up`: start service  
`make down`: stop service  
`make logs`: display service logs

## API
### Add address
```json lines
{
  "method": "add_address",
  "data": "$address"
}
```
#### Params
`$address` <- any wallet address to find in transaction

### Delete address
```json lines
{
  "method": "delete_address",
  "data": "$address"
}
```
#### Params
`$address` <- any wallet address to find in transaction
