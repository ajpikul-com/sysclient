# sysclient

I think there is a memory leak with a write go routine that doesn't stop when it's connection is severed.

In general there is poor documentation and modularization with the modules. There is the `gitstatus` client, and API changes need to propogate through `sysclient` and `sysboss` and the `ajpikul_static2` as thats where the parser is. Would it better to supply a `MarshalText` or `MarshalHTML` function in Javascript directly in that library? And should that library contain the config structure that `sysclient` has to use?

This is a port of sysboss w/ code stripped and added from wsssh/testclient-ssh
