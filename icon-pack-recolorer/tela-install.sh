#!/usr/bin/env bash

ICON_DIR="$HOME/.cache/icon-pack-recolorer/tela-source"
REPO_URL="https://github.com/vinceliuice/Tela-icon-theme.git"

if [ ! -d "$ICON_DIR/.git" ]; then
    echo "Directory not found, cloning repository..."
    git clone "$REPO_URL" "$ICON_DIR"
else
    echo "Directory exists, updating repository..."
    git -C "$ICON_DIR" pull
fi

$HOME/.cache/icon-pack-recolorer/tela-source/install.sh -d /tmp/tela-out nord

echo "Updating tela-frappe"
rm -rf $HOME/.local/share/icons/Tela-frappe/
cp -r /tmp/tela-out/Tela-nord/ $HOME/.local/share/icons/Tela-frappe
rm -rf /tmp/tela-out/
