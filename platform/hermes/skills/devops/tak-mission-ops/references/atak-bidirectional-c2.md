# Legacy AI C2 via ATAK — Full Bidirectional Interface

## What's Happening on iTAK Right Now
You send a GeoChat message → Legacy receives it → Legacy responds with text back in chat.  
Text in, text out. You're using a tactical platform like a messaging app.

## What Opens Up on ATAK

The tak-agent on the NATS bus is built to handle far more than text. On ATAK, the full CoT conversation opens:

| Zeus Does in ATAK | Legacy Sees on NATS | Legacy Does Back |
|-------------------|---------------------|-----------------|
| Drop yellow marker with remarks | CoT event on `tak.cot.out` | Researches intel on that location, pushes a waypoint back to your map |
| Send GeoChat to Legacy room | `b-t-f` CoT message | Processes, responds in chat AND drops map markers simultaneously |
| Drop blue "INTEL" marker | Mission UID add event | Logs intel point, cross-references, acknowledges on map |
| Drop green "OBJ COMPLETE" | Mission change event | Logs completion, updates mission log on TAK server |
| Drop red "HOLD" marker | CoT event | Logs hold, maintains position watch, sends ACK |
| Hit emergency button | `a-f-G-U-C-E` CoT | Priority response, pushes support info to your map |
| Change team color Red→Green | SA update | Status signal — Legacy reads as mission state change |

## The Real Difference

**iTAK:** You talk to Legacy. Legacy talks back.

**ATAK:** You and Legacy share a live Common Operating Picture. You drop markers — Legacy reads them as commands. Legacy drops markers, waypoints, and routes back onto your map in real time. **The map IS the interface — no typing required.**

## Concrete NYC Mission Example

**iTAK right now:**
```
You:    "Legacy what's at 40.7580° N, 73.9855° W"
Legacy: "That's Times Square, here's what I know..."
```

**ATAK on S25 Ultra:**
```
You:    Drop a yellow marker on Times Square with remarks "recon this"
Legacy: Pulls intel → drops a waypoint back on your map with label
        → pushes a GeoChat summary to your Legacy room
        → logs it to your active mission on the TAK server
        All within seconds. Nothing typed.
```

## Why It's Even Better With Your OneMind Stack

Your fabric is already publishing live feeds into TAK:
- ADS-B aircraft positions on your map
- APRS ham radio operators
- ACLED/GDACS/USGS intel overlays
- MediaMTX RTSP camera streams

On ATAK you can drop a marker on any of those live feed objects and Legacy immediately pulls context on it — what is that aircraft, who is that APRS operator, what is that seismic event.

**iTAK can display those feeds. ATAK lets you interact with them and command Legacy against them.**

## Bottom Line

Right now you're using Legacy as a chatbot attached to your TAK server.  
On ATAK, the map becomes a direct command interface — your gestures on the COP are orders to Legacy, and Legacy's responses come back as map objects, not just text.  
That's what full bidirectional C2 actually means.
