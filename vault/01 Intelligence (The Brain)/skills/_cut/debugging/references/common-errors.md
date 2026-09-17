# Common Error Reference

## Python Errors

### TypeError
```
TypeError: 'NoneType' object is not subscriptable
```
**Cause:** Trying to index something that is None
**Fix:** Check if value exists before accessing

### KeyError
```
KeyError: 'missing_key'
```
**Cause:** Dictionary key doesn't exist
**Fix:** Use `.get()` or check `if key in dict`

### ImportError
```
ImportError: No module named 'xyz'
```
**Cause:** Package not installed or wrong path
**Fix:** `pip install xyz` or check PYTHONPATH

### AttributeError
```
AttributeError: 'str' object has no attribute 'append'
```
**Cause:** Using method that doesn't exist on type
**Fix:** Check the object type, wrong variable?

## JavaScript Errors

### undefined is not a function
**Cause:** Calling something that isn't a function
**Fix:** Check variable type, async timing

### Cannot read property 'x' of undefined
**Cause:** Accessing property on undefined object
**Fix:** Optional chaining `obj?.property`

### CORS Error
**Cause:** Cross-origin request blocked
**Fix:** Server needs CORS headers

## Database Errors

### Connection refused
**Cause:** DB not running or wrong port
**Fix:** Check service status, connection string

### Deadlock
**Cause:** Two transactions waiting for each other
**Fix:** Consistent lock ordering, timeouts

## Network Errors

### ECONNREFUSED
**Cause:** Target not accepting connections
**Fix:** Check service is running, firewall

### Timeout
**Cause:** Request taking too long
**Fix:** Check network, increase timeout, optimize query


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
