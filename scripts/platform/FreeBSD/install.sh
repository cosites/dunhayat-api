#!/bin/sh

# FreeBSD install script for dunhayat
# This script installs the dunhayat binary, rc.d script, and configuration files

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

BINARY_NAME="dunhayat"
BINARY_SRC="target/api"
BINARY_DST="/usr/local/bin/${BINARY_NAME}"
RCD_SRC="scripts/platform/FreeBSD/rc.d/dunhayat"
RCD_DST="/usr/local/etc/rc.d/dunhayat"
CONFIG_DIR="/usr/local/etc/dunhayat"
CONFIG_SRC="config.yaml.example"

print_info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

print_warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
}

check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        print_error "This script must be run as root (use either sudo or doas)"
        exit 1
    fi
}

check_files() {
    if [ ! -f "$BINARY_SRC" ]; then
        print_error "Binary not found at $BINARY_SRC. Please build up first."
        exit 1
    fi
    
    if [ ! -f "$RCD_SRC" ]; then
        print_error "RC script not found at $RCD_SRC"
        exit 1
    fi
    
    if [ ! -f "$CONFIG_SRC" ]; then
        print_error "Configuration template not found at $CONFIG_SRC"
        exit 1
    fi
}

install_binary() {
    print_info "Installing binary to $BINARY_DST"
    install -m 755 "$BINARY_SRC" "$BINARY_DST"
    print_info "Binary installed successfully"
}

install_rcd_script() {
    print_info "Installing RC script to $RCD_DST"
    install -m 755 "$RCD_SRC" "$RCD_DST"
    print_info "RC script installed successfully"
}

install_config() {
    print_info "Creating configuration directory $CONFIG_DIR"
    mkdir -p "$CONFIG_DIR"
    
    print_info "Installing configuration template"
    install -m 644 "$CONFIG_SRC" "$CONFIG_DIR/config.yaml.example"
    
    if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
        print_info "Creating default configuration file"
        install -m 644 "$CONFIG_SRC" "$CONFIG_DIR/config.yaml"
        print_warn "Please edit $CONFIG_DIR/config.yaml to match your environment"
    else
        print_info "Configuration file already exists at $CONFIG_DIR/config.yaml"
    fi
}

create_log_file() {
    LOG_FILE="/var/log/dunhayat.log"
    if [ ! -f "$LOG_FILE" ]; then
        print_info "Creating log file $LOG_FILE"
        install -m 644 /dev/null "$LOG_FILE"
    fi
}

show_post_install() {
    print_info "Installation completed successfully!"
    echo
    print_info "Next steps:"
    echo "  1. Edit the configuration file: $CONFIG_DIR/config.yaml"
    echo "  2. Enable the service: service dunhayat enable"
    echo "  3. Start the service: service dunhayat start"
    echo "  4. Check service status: service dunhayat status"
    echo
    print_info "Service files installed:"
    echo "  Binary:     $BINARY_DST"
    echo "  RC script:  $RCD_DST"
    echo "  Config:     $CONFIG_DIR/config.yaml"
    echo "  Example:    $CONFIG_DIR/config.yaml.example"
    echo "  Log file:   /var/log/dunhayat.log"
}

main() {
    print_info "Starting dunhayat installation for FreeBSD"
    
    check_root
    check_files
    
    install_binary
    install_rcd_script
    install_config
    create_log_file
    
    show_post_install
}

main "$@"
