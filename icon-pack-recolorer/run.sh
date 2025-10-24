#!/usr/bin/env bash

. tela-install.sh

go build
echo "Recoloring target"
./icon-pack-recolorer $HOME/.local/share/icons/Tela-frappe/

echo "Reloading theme"
. reloadtheme.sh
