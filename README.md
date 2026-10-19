# Vyos Router
## Description
VyOS firewall configuration template for my homelab router

## Configuration Example
```yaml
hostname: "vyos"
mdns:
  enabled: true
  interface: ethX
interfaces:
  wan: ethX
monitoring:
  enabled: true
  exporter: prometheus   # [prometheus, telegraf] 
  endpoint: 192.168.10.1 # Recommend a private endpoint like 192.168.0.0/16
```
