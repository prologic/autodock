# autodock design

- 3 components

# autodock-server

- runs docker proxy
- runs on manager nodes

Has the following endpoints:

```
/metrics
/events
/proxy
```

## autodock-agent

- runs event collectors from docker daemon and publishes to server via msgbus
- runs on every node

Has the following endpoints:

```
/metrics
```

## autodock-<plugin>

- listens for events from server via msgbus and performs actions
- runs anywhere on any node

Has the following endpoints:

```
/metrics
```
