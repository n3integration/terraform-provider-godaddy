#!/bin/sh

version=1.6.0

os=$(uname -s | tr '[:upper:]' '[:lower:]')
mach=$(uname -m)
plugin="terraform-godaddy"

case "$os" in
  "linux"|"darwin")
    printf "[info] fetching latest revision..."

    if [ "$mach" == "x86_64" ]; then
      arch="amd64"
    else
      arch="386"
    fi

    file=$(printf "%s_%s_%s" $plugin $os $arch)
    archive="$file.tgz"

    cd /tmp && curl -fOsSL https://github.com/n3integration/$plugin/releases/download/v$version/$archive 2>/dev/null

    if [ ! -f /tmp/$archive ]; then
      echo "ERROR"
      echo "\t-> failed to download file"
      exit 1
    else 
      cd /tmp && gzip -dc $archive | tar xf -
      mv $file $plugin
    fi

    if [ -f $plugin ]; then
      echo "OK"
    else
      echo "ERROR"
      echo "\t-> failed to extract file contents"
      exit 1
    fi
    ;;
  *)
    echo "[error] Dagger. It doesn't look like your OS is currently supported. Please submit an issue or pull request."
    exit 1
    ;;
esac

printf "[info] installing plugin..."
mkdir -p ~/.terraform/plugins && mv $plugin ~/.terraform/plugins
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
    echo "\tgodaddy = \"$HOME/.terraform/plugins/$plugin\""
    echo "}"
    echo ""
  fi
else
  cat > ~/.terraformrc <<EOF
providers {
  godaddy = "$HOME/.terraform/plugins/$plugin"
}
EOF
  echo "[info] complete"
fi
