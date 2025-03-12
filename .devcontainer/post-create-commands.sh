#!/bin/bash -e

echo "Installing inputrc..."
cat << EOF >> "$HOME/.inputrc"
# Make history search (up and down arrow) take into account what you've already typed.
"\e[A": history-search-backward
"\e[B": history-search-forward
EOF

echo "Installing bash completions and setting path..."
cat << EOF >> "$HOME/.bashrc"
# Bash completions
source <(kubectl completion bash)
source <(k3d completion bash)
EOF
