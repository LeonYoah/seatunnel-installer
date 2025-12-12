#!/bin/bash
# ============================================================================
# SeaTunnel Installer - 通用工具函数库
# ============================================================================
# 本文件包含可复用的工具函数，供 install_seatunnel.sh 和其他脚本使用
# 使用方式: source util.sh
# ============================================================================

# 防止重复加载
if [ -n "$_UTIL_SH_LOADED" ]; then
    return 0
fi
_UTIL_SH_LOADED=1

# ============================================================================
# 颜色定义
# ============================================================================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# ============================================================================
# 日志函数
# ============================================================================
# 注意: 在 JSON 模式下 (_JSON_MODE=true)，日志会被抑制
log_info() {
    [ "${_JSON_MODE:-false}" = "true" ] && return 0
    echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

log_warning() {
    [ "${_JSON_MODE:-false}" = "true" ] && return 0
    echo -e "${YELLOW}[WARN]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

log_error() {
    if [ "${_JSON_MODE:-false}" = "true" ]; then
        echo "{\"error\":\"$1\"}"
    else
        echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
    fi
    exit 1
}
log_debug() {
    [ "${_JSON_MODE:-false}" = "true" ] && return 0
    echo -e "${YELLOW}[DEBUG]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

log_success() {
    [ "${_JSON_MODE:-false}" = "true" ] && return 0
    echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

# ============================================================================
# SSH/SCP 工具函数
# ============================================================================

# 带重试的SSH命令执行
# 参数: $1=host, $2=cmd
ssh_with_retry() {
    local host=$1
    local cmd=$2
    local retries=0
    
    if [ -z "$INSTALL_USER" ]; then
        log_error "INSTALL_USER未设置"
        return 1
    fi
    
    while [ $retries -lt ${MAX_RETRIES:-3} ]; do
        if timeout ${SSH_TIMEOUT:-10} ssh -p ${SSH_PORT:-22} "${INSTALL_USER}@${host}" "$cmd" 2>/dev/null; then
            return 0
        fi
        retries=$((retries + 1))
        log_warning "SSH到 ${INSTALL_USER}@${host} 失败，重试 $retries/${MAX_RETRIES:-3}..."
        sleep 2
    done
    
    log_error "SSH到 ${INSTALL_USER}@${host} 失败，已重试 ${MAX_RETRIES:-3} 次"
    return 1
}

# 带重试的SCP命令执行
# 参数: $1=src, $2=host, $3=dest
scp_with_retry() {
    local src=$1
    local host=$2
    local dest=$3
    local retries=0
    
    if [ -z "$INSTALL_USER" ]; then
        log_error "INSTALL_USER未设置"
        return 1
    fi
    
    if [ ! -e "$src" ]; then
        log_error "源文件/目录不存在: $src"
        return 1
    fi
    
    # 测试SSH连接
    if ! timeout ${SSH_TIMEOUT:-10} ssh -p ${SSH_PORT:-22} -o ConnectTimeout=5 "${INSTALL_USER}@${host}" "echo >/dev/null" 2>/dev/null; then
        log_error "SSH连接失败: ${INSTALL_USER}@${host}"
        return 1
    fi
    
    while [ $retries -lt ${MAX_RETRIES:-3} ]; do
        log_info "正在分发到 ${host}..."
        
        if timeout ${SSH_TIMEOUT:-10} scp -q -r -P ${SSH_PORT:-22} "$src" "${INSTALL_USER}@${host}:$dest" 2>/dev/null; then
            log_info "成功分发到 ${host}"
            return 0
        else
            local disk_space
            disk_space=$(ssh -p ${SSH_PORT:-22} "${INSTALL_USER}@${host}" "df -h $dest" 2>/dev/null | tail -n1 | awk '{print $4}')
            log_warning "分发失败，目标目录可用空间: ${disk_space:-未知}"
        fi
        
        retries=$((retries + 1))
        if [ $retries -lt ${MAX_RETRIES:-3} ]; then
            log_warning "分发到 ${host} 失败，重试 $retries/${MAX_RETRIES:-3}..."
            sleep 2
        fi
    done
    
    log_error "分发到 ${host} 失败，已重试 ${MAX_RETRIES:-3} 次"
    return 1
}

# ============================================================================
# 文件/目录检查函数
# ============================================================================

# 检查文件是否存在
check_file() {
    if [[ ! -f "$1" ]]; then
        log_error "文件不存在: $1"
    fi
}

# 检查目录是否存在
check_dir() {
    if [[ ! -d "$1" ]]; then
        log_error "目录不存在: $1"
    fi
}

# 检查端口占用
# 参数: $1=host, $2=port
check_port() {
    local host=$1
    local port=$2
    
    if command -v nc >/dev/null 2>&1; then
        if nc -z -w2 "$host" "$port" >/dev/null 2>&1; then
            log_error "端口 $port 在 $host 上已被占用"
        fi
    elif command -v telnet >/dev/null 2>&1; then
        if echo quit | timeout 2 telnet "$host" "$port" >/dev/null 2>&1; then
            log_error "端口 $port 在 $host 上已被占用"
        fi
    else
        if timeout 2 bash -c "echo >/dev/tcp/$host/$port" >/dev/null 2>&1; then
            log_error "端口 $port 在 $host 上已被占用"
        fi
    fi
}

# 替换文件内容
# 参数: $1=search, $2=replace, $3=file
replace_in_file() {
    local search=$1
    local replace=$2
    local file=$3
    sed -i "s|$search|$replace|g" "$file"
}

# ============================================================================
# 目录权限管理
# ============================================================================

# 设置目录权限
# 参数: $1=dir
setup_permissions() {
    local dir=$1
    log_info "设置目录权限: $dir"
    
    if [ -z "$INSTALL_USER" ]; then
        log_warning "INSTALL_USER 未设置，跳过目录权限设置: $dir"
        return 0
    fi

    if [ "$INSTALL_USER" = "root" ]; then
        log_info "INSTALL_USER 为 root，跳过目录权限设置"
        return 0
    fi
    
    local user_home
    if command -v getent >/dev/null 2>&1; then
        user_home=$(getent passwd "$INSTALL_USER" | cut -d: -f6)
    else
        user_home=$(eval echo ~"$INSTALL_USER")
    fi
    
    if [[ "$dir" == "$user_home" ]]; then
        log_warning "跳过用户home目录的权限设置: $dir"
        return 0
    fi
    
    if [ ! -d "$dir" ]; then
        log_error "目录不存在: $dir"
        return 1
    fi
    
    sudo chown -R "$INSTALL_USER:$INSTALL_GROUP" "$dir"
    sudo chmod -R 755 "$dir"
}

# 指定用户执行命令
# 参数: $1=command
run_as_user() {
    sudo -u "$INSTALL_USER" bash -c "$1"
}

# 创建目录
# 参数: $1=dir
create_directory() {
    local dir=$1
    sudo mkdir -p "$dir"
}

# ============================================================================
# 临时文件管理
# ============================================================================

# 创建临时文件
create_temp_file() {
    local temp_dir="${LOG_DIR:-/tmp}/temp"
    local temp_file
    
    if [ ! -d "$temp_dir" ]; then
        mkdir -p "$temp_dir" || log_error "无法创建临时目录: $temp_dir"
        chmod 700 "$temp_dir"
        [ -n "$INSTALL_USER" ] && [ -n "$INSTALL_GROUP" ] && chown "$INSTALL_USER:$INSTALL_GROUP" "$temp_dir"
    fi
    
    temp_file=$(mktemp "$temp_dir/temp.XXXXXX") || log_error "无法创建临时文件"
    chmod 600 "$temp_file"
    [ -n "$INSTALL_USER" ] && [ -n "$INSTALL_GROUP" ] && chown "$INSTALL_USER:$INSTALL_GROUP" "$temp_file"
    
    TEMP_FILES+=("$temp_file")
    
    echo "$temp_file"
}

# 清理临时文件
cleanup_temp_files() {
    
    if [ -z "${TEMP_FILES+x}" ]; then
        return 0
    fi
    
    log_info "当前TEMP_FILES数组大小: ${#TEMP_FILES[@]}"
    
    local temp_dir="${LOG_DIR:-/tmp}/temp"
    
    if [ ${#TEMP_FILES[@]} -eq 0 ]; then
        log_info "没有临时文件需要清理"
        return 0
    fi
    
    log_info "清理临时文件..."
    
    for temp_file in "${TEMP_FILES[@]}"; do
        if [ -f "$temp_file" ]; then
            rm -f "$temp_file"
            log_info "已删除临时文件: $temp_file"
        fi
    done
    
    if [ -d "$temp_dir" ] && [ -z "$(ls -A "$temp_dir")" ]; then
        rm -rf "$temp_dir"
        log_info "已删除空的临时目录: $temp_dir"
    fi
}

# ============================================================================
# 系统架构识别
# ============================================================================

# 获取系统架构
get_arch() {
    local m
    m=$(uname -m)
    case "$m" in
        x86_64|amd64)
            echo "x86";;
        aarch64|arm64)
            echo "arm64";;
        *)
            echo "$m";;
    esac
}

# ============================================================================
# 命令安装工具
# ============================================================================

# 从本地lib安装缺失命令
# 参数: $1=cmd
install_if_missing() {
    local cmd="$1"
    local binary_path="/usr/local/bin/$cmd"
    local install_arch
    install_arch=$(get_arch)
    local src1="${COMMAND_LIB_DIR:-./lib}/${cmd}-${install_arch}"
    local src2="${COMMAND_LIB_DIR:-./lib}/${cmd}"

    if command -v "$cmd" >/dev/null 2>&1; then
        log_info "$cmd 命令已存在."
        return 0
    fi

    local src=""
    if [ -f "$src1" ]; then
        src="$src1"
    elif [ -f "$src2" ]; then
        src="$src2"
    else
        log_warning "未在本地找到 $cmd 离线包: $src1 / $src2"
        return 1
    fi

    log_info "$cmd 命令不存在. 正在安装: $src -> $binary_path"
    sudo cp "$src" "$binary_path"
    sudo chmod +x "$binary_path"
}

# 使用yq修改YAML的统一封装函数
# 参数:
#   $1: yaml文件路径
#   $2: yq表达式
#   $3: 可选的环境变量列表(格式: "VAR1=value1 VAR2=value2")
replace_yaml_with_yq() {
    local yaml_file="$1"
    local yq_expr="$2"
    local env_vars="${3:-}"
    
    if [ ! -f "$yaml_file" ]; then
        log_error "YAML文件不存在: $yaml_file"
        return 1
    fi
    
    if ! command -v yq >/dev/null 2>&1; then
        log_error "yq未启用或不可用，无法执行YAML修改"
        return 1
    fi
    
    cp "$yaml_file" "${yaml_file}.bak" 2>/dev/null || true
    
    if [ -n "$env_vars" ]; then
        eval "$env_vars yq -i '$yq_expr' '$yaml_file'"
    else
        yq -i "$yq_expr" "$yaml_file"
    fi
    
    if [ $? -ne 0 ]; then
        log_error "yq执行失败: $yq_expr"
        [ -f "${yaml_file}.bak" ] && mv "${yaml_file}.bak" "$yaml_file"
        return 1
    fi
    
    return 0
}

# ============================================================================
# 版本比较工具
# ============================================================================

# 标准化版本号
normalize_version() {
    echo "$1" | sed 's/[^0-9.].*$//'
}

# 版本比较: v1 >= v2
# 参数: $1=v1, $2=v2
# 返回: 0 表示 v1 >= v2, 1 表示 v1 < v2
version_ge() {
    local v1 v2
    v1=$(normalize_version "$1")
    v2=$(normalize_version "$2")
    if sort -V </dev/null >/dev/null 2>&1; then
        [ "$(printf '%s\n' "$v2" "$v1" | sort -V | head -n1)" = "$v2" ]
        return $?
    fi
    IFS='.' read -r a1 a2 a3 <<<"$v1"; IFS='.' read -r b1 b2 b3 <<<"$v2"
    a1=${a1:-0}; a2=${a2:-0}; a3=${a3:-0}
    b1=${b1:-0}; b2=${b2:-0}; b3=${b3:-0}
    if (( a1 > b1 )); then return 0; elif (( a1 < b1 )); then return 1; fi
    if (( a2 > b2 )); then return 0; elif (( a2 < b2 )); then return 1; fi
    if (( a3 >= b3 )); then return 0; else return 1; fi
}

# ============================================================================
# 网络工具
# ============================================================================

# 获取本机IP列表
get_local_ips() {
    hostname -I 2>/dev/null | tr ' ' '\n' | grep -v '^$' || \
    ip addr show 2>/dev/null | grep -oP 'inet \K[\d.]+' || \
    ifconfig 2>/dev/null | grep -oP 'inet \K[\d.]+'
}

# 检查是否为本地节点
# 参数: $1=node_ip
is_local_node() {
    local node=$1
    local local_ips
    local_ips=$(get_local_ips)
    
    if [[ "$node" == "localhost" ]] || [[ "$node" == "127.0.0.1" ]]; then
        return 0
    fi
    
    if echo "$local_ips" | grep -q "^${node}$"; then
        return 0
    fi
    
    return 1
}

# ============================================================================
# 字符串工具
# ============================================================================

# 去除字符串首尾空白
trim() {
    local var="$*"
    var="${var#"${var%%[![:space:]]*}"}"
    var="${var%"${var##*[![:space:]]}"}"
    echo -n "$var"
}

# 检查字符串是否为空
is_empty() {
    local var="$1"
    [ -z "$(trim "$var")" ]
}

# 检查字符串是否非空
is_not_empty() {
    local var="$1"
    [ -n "$(trim "$var")" ]
}
