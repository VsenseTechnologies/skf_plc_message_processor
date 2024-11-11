# SKF MQTT API

## 1. Real time Temperature

```json
{"rt_tp":"30.90"}
```

The value will be of floating type

## 2. Real time PID

```json
{"rt_pid":"30.90"}
```

The value will be of floating type

## 3. Blower Trip Status

```json
{"st_bl_trp":"0"}
```

The value may be 0 or 1

## 4. Elevator Trip Status

```json
{"st_el_trp":"1"}
```

The value may be 0 or 1

## 5. Rotor Trip Status

```json
{"st_rt_trp":"0"}
```

The value may be 0 or 1

## 6. Blower Run Status

```json
{"st_bl_rn":"0"}
```

The value may be 0 or 1

## 7. Elevator Run Status

```json
{"st_el_rn":"0"}
```

The value may be 0 or 1

## 8. Rotor Run Status

```json
{"st_rt_rn":"0"}
```

The value may be 0 or 1

## 9. Recipe Step

```json
{"rcp_stp":"1","tm":"40","tp":"30.90"}
```

1. “rcp_stp” → Recipe step count
2. “tm” → Recipe time
3. “tp” → Recipe temperature