#!/bin/bash
# ============================================================================
# SeaTunnel Web 安装向导启动脚本
# ============================================================================
# 使用 lib 目录下的 busybox httpd 离线启动 Web 服务器
# API 脚本: www/cgi-bin/run.sh
# ============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WWW_DIR="$SCRIPT_DIR/www"
LIB_DIR="$SCRIPT_DIR/lib"
INSTALL_SCRIPT="$SCRIPT_DIR/install_seatunnel.sh"
PID_FILE="$SCRIPT_DIR/.web_server.pid"
HTTPD_CONF="$SCRIPT_DIR/.httpd.conf"
PORT=${1:-8888}

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# 获取架构
get_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "x86";;
        aarch64|arm64) echo "arm64";;
        *) echo "$(uname -m)";;
    esac
}

# 获取 busybox
get_busybox() {
    local arch=$(get_arch)
    local busybox_path="$LIB_DIR/busybox-$arch"
    
    if [ -f "$busybox_path" ]; then
        chmod +x "$busybox_path" 2>/dev/null || true
        echo "$busybox_path"
        return 0
    fi
    command -v busybox >/dev/null 2>&1 && echo "busybox" && return 0
    return 1
}

# 检查端口
check_port() {
    if command -v ss >/dev/null 2>&1; then
        ss -tuln | grep -q ":$PORT " && return 1
    elif command -v netstat >/dev/null 2>&1; then
        netstat -tuln | grep -q ":$PORT " && return 1
    fi
    return 0
}

# 获取本机 IP
get_local_ip() {
    hostname -I 2>/dev/null | awk '{print $1}' || echo "127.0.0.1"
}

# 创建 httpd 配置
# busybox httpd 配置格式:
# A:IP  - 允许访问的IP
# /cgi-bin:user - CGI目录和用户
# *.ext:mime/type - MIME类型
create_httpd_conf() {
    cat > "$HTTPD_CONF" << EOF
A:*
/cgi-bin:root
*.html:text/html
*.css:text/css
*.js:application/javascript
*.json:application/json
*.png:image/png
*.jpg:image/jpeg
*.svg:image/svg+xml
*.ico:image/x-icon
EOF
}

# 启动服务器
start_server() {
    local busybox_cmd
    busybox_cmd=$(get_busybox) || log_error "找不到 busybox，请确保 lib/busybox-x86 或 lib/busybox-arm64 存在"
    
    if ! "$busybox_cmd" --list 2>/dev/null | grep -q httpd; then
        log_error "busybox 不支持 httpd"
    fi
    
    # 检查 CGI 脚本
    if [ ! -f "$WWW_DIR/cgi-bin/run.sh" ]; then
        log_error "CGI 脚本不存在: $WWW_DIR/cgi-bin/run.sh"
    fi
    chmod +x "$WWW_DIR/cgi-bin/run.sh"
    
    log_info "使用 busybox httpd: $busybox_cmd"
    
    # busybox httpd 启动
    # -f 前台运行, -p 端口, -h Web根目录
    "$busybox_cmd" httpd \
        -f \
        -p "0.0.0.0:$PORT" \
        -h "$WWW_DIR" &
    echo $! > "$PID_FILE"
    
    sleep 1
    kill -0 "$(cat "$PID_FILE")" 2>/dev/null || log_error "httpd 启动失败"
    log_info "httpd 已启动 (PID: $(cat "$PID_FILE"))"
}

# 停止服务器
stop_server() {
    log_info "停止 Web 服务器..."
    [ -f "$PID_FILE" ] && {
        kill -TERM "$(cat "$PID_FILE")" 2>/dev/null
        rm -f "$PID_FILE"
    }
    pkill -f "httpd.*$PORT" 2>/dev/null || true
    rm -f "$HTTPD_CONF"
    log_info "已停止"
}

# 帮助
show_help() {
    cat << EOF
SeaTunnel Web 安装向导 (busybox httpd)

用法: $0 [命令] [端口]

命令:
    start [端口]    启动 (默认: 8888)
    stop            停止
    status          状态
    help            帮助

API (GET /api/run.sh?action=xxx):
    action=steps              获取步骤列表
    action=run                启动完整安装
    action=run&step=N         执行指定步骤
    action=status             获取状态
    action=stop               停止安装
    action=log&lines=100      获取日志
    action=config_load        加载配置
    action=config_save (POST) 保存配置
    action=check&type=all     环境检查

示例:
    $0              # 启动
    $0 9000         # 端口9000启动
    $0 stop         # 停止
EOF
}

# 主函数
main() {
    local action="${1:-start}"
    
    case "$action" in
        start|[0-9]*)
            [[ "$action" =~ ^[0-9]+$ ]] && PORT="$action" || [ -n "$2" ] && PORT="$2"
            check_port || log_error "端口 $PORT 已被占用"
            [ -d "$WWW_DIR" ] || log_error "Web 目录不存在: $WWW_DIR"
            [ -f "$INSTALL_SCRIPT" ] || log_error "安装脚本不存在: $INSTALL_SCRIPT"
            
            start_server
            
            echo ""
            echo -e "${BLUE}============================================${NC}"
            echo -e "${GREEN}SeaTunnel Web 安装向导已启动!${NC}"
            echo -e "${BLUE}============================================${NC}"
            echo ""
            echo -e "访问: ${GREEN}http://$(get_local_ip):$PORT${NC}"
            echo -e "CLI:  ${YELLOW}./install_seatunnel.sh --help${NC}"
            echo ""
            echo -e "按 ${YELLOW}Ctrl+C${NC} 停止"
            echo ""
            
            trap stop_server EXIT INT TERM
            wait
            ;;
        stop) stop_server ;;
        status)
            [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null && \
                log_info "运行中 (PID: $(cat "$PID_FILE"), 端口: $PORT)" || echo "未运行"
            ;;
        help|--help|-h) show_help ;;
        *) log_error "未知命令: $action" ;;
    esac
}

main "$@"
