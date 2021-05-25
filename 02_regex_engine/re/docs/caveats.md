## Worse case
- The most expensive operation is **concatenating sets**, the more characters in the set, the most expensive it gets. This doesn't usually pose a problem, but concatenating '\S' is faster than '[a-zA-Z0-9\_]', even though the first set is bigger.
- During compiling time, the following pattern takes exponential time: "a\*A"^n, with "a\*A\*"^14 taking 1 minute to compile on my machine. Meaning this would be an attack vector if it were in production.

## Comparison with other engines
In terms of **runtime** performance, it's a little slower in concatenation than the standard Go regexr engine, but faster in sets and closure. Overall it's as fast as Go's regexr, but there are differences in some patterns.

In terms of **compile** performance, it's quite slow, sometimes by orders of magnitude. But that is expected, it assumes the Machine is going to compile only once in runtime. Besides, it's under 200 microseconds for small regexes in most machines.
