# maelstrom-unique-ids

## design

each node has its own counter (with a lock because *goroutines~*) to ensure that when a node processes multiple requests arriving with the same timestamp, ids will stil be unique

timestamp is not really necessary for this challenge (the counter should be sufficient to guarantee uniqueness) but node id is required to ensure uniqueness across nodes even with network partitions

## format

\<timestamp\>-\<node_id\>-\<counter\>

- timestamp: unix timestamp
- node id: a string identifier
- counter: counter value for that node

all parts are joined with dashes to form the unique id