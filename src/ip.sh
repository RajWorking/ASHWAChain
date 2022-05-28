#!/bin/bash

for i in {0..500}
do
echo -n "127.0.0.1:$((7000 + i))" > "./Keys/${i}_socket"
done
