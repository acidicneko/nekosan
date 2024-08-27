#!/bin/bash
install() {
  go build -v
  mv nekosan "$GITMAN_BIN"/nekosan
}

uninstall() {
  rm "$GITMAN_BIN"/nekosan
}

update(){
  install
}

if [ $1 = "install" ] ; then 
	install
elif [ $1 = "uninstall" ] ; then
	uninstall
elif [ $1 = "update" ] ; then
	update
else
	echo unknown option
	exit 1
fi
