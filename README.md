# macos VPN selector

### build
```shell
make
```
### config
#### config.yaml
```yaml
accounts:
  - alias: default.bigip
    id: __VPN_ID__
    pw: __VPN_PW__
    client: bigip
  - alias: default.wg
    client: wg
    service_name: __VPN_SERVICE_NAME__
vpn_clients:
  - alias: bigip
    location: /Applications/BIG-IP Edge Client.app
```
#### config path as environment variable
```shell
VPN_AUTOMATOR_CONFIG_PATH=~/vpn-automator/doc/config.yml
```


#### ~/.zshrc
```shell
alias vpn=vpn-automator
```
### example
- `[UNKN]`: big-ip edge client is running
- `[CONN]`: connected
- `[DISC]`: disconnected
```shell
vpn-automator git:(main) ✗ vpn                                                                                                                                                                                                           [22/04/18| 9:09AM]
⚡️Select VPN to connect

[UNKN]> cst03
[UNKN]  kcloud.fin
[UNKN]  kcloud.fin-root
[CONN]  wg.cst03
[DISC]  wg.kcloud

```