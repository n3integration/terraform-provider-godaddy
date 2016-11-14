#!/bin/sh

version=0.0.2
os=$(uname -s | tr A-Z a-z)

case "$os" in
  "linux"|"darwin")
    printf "[info] fetching latest revision..."
    curl -fsSL https://github.com/n3integration/terraform-godaddy/releases/download/v$version/terraform-godaddy-$os-$version.tgz | gzip -dc | tar xf -

    if [ -f terraform-godaddy ]; then
      echo "OK"
    else
      exit 1
    fi
    ;;
  *)
    echo "[error] Whoops. It doesn't look like your OS is currently supported."
    exit 1
    ;;
esac

printf "[info] installing plugin..."
mkdir -p ~/.terraform/plugins && mv terraform-godaddy ~/.terraform/plugins
if [ $? -eq 0 ]; then
  echo "OK"
else
  exit 1
fi

if [ -f ~/.terraformrc ]; then
  if [[ $(grep "godaddy" ~/.terraformrc) ]]; then
    echo "[info] complete"
  else
    echo "[info] append the godaddy provider to your ~/.terraformrc configuration file"
    echo ""
    echo "providers {"
    echo "\tgodaddy = \"$HOME/.terraform/plugins/terraform-godaddy\""
    echo "}"
    echo ""
  fi
else
  cat > ~/.terraformrc <<EOF
providers {
  godaddy = "$HOME/.terraform/plugins/terraform-godaddy"
}
EOF
  echo "[info] complete"
fi
