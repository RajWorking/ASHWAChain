#!/bin/bash

# arg1 = nodes, arg2 = number of ids

./generateIDs.py $1 $2 & pid=$!

wait $pid

echo "Starting nodes..."
for ((i = 1; i <= $1; i++))
do
    echo "Node $i running.."
    ./node.py $i &
done
echo "...nodes running"
