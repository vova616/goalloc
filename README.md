## goalloc
goalloc is a Go package that can allocate and free memory.
It can also convert and load types and is made safe as possible, but this package is not completly safe, its safe as long as you dont use freed memory.
I recommend not using this package and if you do, try to pass MemBlock and not the types it creates.