#!/bin/bash

# 根据是否有--prerelease参数来获取对应的版本信息
get_release_url() {
    if [[ $1 == "--prerelease" ]]; then
        curl --silent "https://api.github.com/repos/scriptscat/cloudcat/releases" |
            grep "browser_download_url" |
            sed -E 's/.*"([^"]+)".*/\1/' | head -n 12
    else
        curl --silent "https://api.github.com/repos/scriptscat/cloudcat/releases/latest" |
            grep "browser_download_url" |
            sed -E 's/.*"([^"]+)".*/\1/'
    fi
}

detect_os_and_arch() {
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv*) ARCH="arm" ;;
        *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    echo "Detected OS: $OS and ARCH: $ARCH"
}

download_and_extract_binary() {
    local release_url="$1"

    detect_os_and_arch

    binary_url=$(echo "$release_url" | grep "${OS}_${ARCH}.tar.gz")

    if [ -z "$binary_url" ]; then
        echo "Binary not found for ${OS}_${ARCH}"
        exit 1
    fi

    echo "Downloading $binary_url..."
    curl -L -o cloudcat.tar.gz "$binary_url"

    mkdir -p /usr/local/cloudcat
    tar xzf cloudcat.tar.gz -C /usr/local/cloudcat
    ln -sf /usr/local/cloudcat/ccatctl /usr/local/bin/ccatctl
    chmod +x /usr/local/cloudcat/cloudcat
    chmod +x /usr/local/cloudcat/ccatctl
}

install_as_service() {
    /usr/local/cloudcat/cloudcat init

    if [ -f /etc/systemd/system/cloudcat.service ]; then
        echo "CloudCat service already exists. Overwriting..."
    fi

    cat > /etc/systemd/system/cloudcat.service <<EOL
[Unit]
Description=CloudCat Service
After=network.target

[Service]
ExecStart=/usr/local/cloudcat/cloudcat server
Restart=always
User=root
Group=root
Environment=PATH=/usr/bin:/usr/local/bin

[Install]
WantedBy=multi-user.target
EOL

    systemctl daemon-reload
    systemctl enable cloudcat
    systemctl start cloudcat
    echo "CloudCat service started!"
}

main() {
    release_url=$(get_release_url $1)
    download_and_extract_binary "$release_url"
    install_as_service
}

main $1
