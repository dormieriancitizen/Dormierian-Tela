#!/usr/bin/env bash

go build
rm -rf $HOME/.local/share/icons/Tela-frappe/
cp -r $HOME/.local/share/icons/Tela-nord/ $HOME/.local/share/icons/Tela-frappe
./icon-pack-recolorer $HOME/.local/share/icons/Tela-frappe/
