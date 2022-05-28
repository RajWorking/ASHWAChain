#!/bin/bash

for i in {0..490..10}
  do
    echo $i
      ./node pbft node -id $i >> logs &
  done