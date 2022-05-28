port=$(tty | sed -e "s:/dev/pts/::"); echo $port

./simplePBFT pbft node -id $(expr $port - 1)

kill $(sudo lsof -t -i:{8080..8100})