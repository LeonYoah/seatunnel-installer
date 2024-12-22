#!/bin/bash

# 确保遇到错误时立即退出
set -e

# 获取脚本执行路径
EXEC_PATH=$(cd "$(dirname "$0")" && pwd)
echo "执行路径: $EXEC_PATH"

# 记录开始时间
START_TIME=$(date +%s)

# 日志文件路径
LOG_DIR="$EXEC_PATH/seatunnel-install-log-${INSTALL_USER:-$(whoami)}"
LOG_FILE="$LOG_DIR/install.log"

# 颜色输出函数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
    exit 1
}
log_success() {
  echo -e "${GREEN}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

# 最大重试次数
MAX_RETRIES=3
# SSH超时时间(秒)
SSH_TIMEOUT=10

# 安装包仓库地址映射
declare -A PACKAGE_REPOS=(
    ["apache"]="https://archive.apache.org/dist/seatunnel"
    ["aliyun"]="https://mirrors.aliyun.com/apache/seatunnel"
    ["huaweicloud"]="https://mirrors.huaweicloud.com/apache/seatunnel"
)

# 插件仓库地址映射
declare -A PLUGIN_REPOS=(
    ["apache"]="https://repo1.maven.org/maven2"
    ["aliyun"]="https://maven.aliyun.com/repository/public"
    ["huaweicloud"]="https://repo.huaweicloud.com/repository/maven"
)


# 添加错误处理函数
handle_error() {
    local exit_code=$?
    local line_number=$1
    log_error "脚本在第 $line_number 行发生错误 (退出码: $exit_code)"
}

# 设置错误处理
trap 'handle_error ${LINENO}' ERR

# 增强的清理函数
cleanup() {
    local exit_code=$?
    log_info "开始清理..."
    
    # 清理临时文件
    cleanup_temp_files
    
    # 如果安装失败,提示用户
    if [ $exit_code -ne 0 ]; then
        log_warning "安装失败。如果需要重新安装,请手动删除安装目录: $SEATUNNEL_HOME"
        log_warning "删除命令: sudo rm -rf $SEATUNNEL_HOME"
    fi
    
    exit $exit_code
}

# 设置清理trap
trap cleanup EXIT INT TERM

# 带重试的SSH命令执行
ssh_with_retry() {
    local host=$1
    local cmd=$2
    local retries=0
    
    # 确保INSTALL_USER已设置
    if [ -z "$INSTALL_USER" ]; then
        log_error "INSTALL_USER未设置"
        return 1
    fi
    
    while [ $retries -lt $MAX_RETRIES ]; do
        if timeout $SSH_TIMEOUT ssh -p $SSH_PORT "${INSTALL_USER}@${host}" "$cmd" 2>/dev/null; then
            return 0
        fi
        retries=$((retries + 1))
        log_warning "SSH到 ${INSTALL_USER}@${host} 失败，重试 $retries/$MAX_RETRIES..."
        sleep 2
    done
    
    log_error "SSH到 ${INSTALL_USER}@${host} 失败，已重试 $MAX_RETRIES 次"
    return 1
}

# 带重试的SCP命令执行
scp_with_retry() {
    local src=$1
    local host=$2
    local dest=$3
    local retries=0
    
    # 确保INSTALL_USER已设置
    if [ -z "$INSTALL_USER" ]; then
        log_error "INSTALL_USER未设置"
        return 1
    fi
    
    # 检查源文件/目录是否存在
    if [ ! -e "$src" ]; then
        log_error "源文件/目录不存在: $src"
        return 1
    fi
    
    # 测试SSH连接
    if ! timeout $SSH_TIMEOUT ssh -p $SSH_PORT -o ConnectTimeout=5 "${INSTALL_USER}@${host}" "echo >/dev/null" 2>/dev/null; then
        log_error "SSH连接失败: ${INSTALL_USER}@${host}"
        return 1
    fi
    
    while [ $retries -lt $MAX_RETRIES ]; do
        log_info "正在分发到 ${host}..."
        
        # 使用-q参数静默输出，仅显示错误
        if timeout $SSH_TIMEOUT scp -q -r -P $SSH_PORT "$src" "${INSTALL_USER}@${host}:$dest" 2>/dev/null; then
            log_info "成功分发到 ${host}"
            return 0
        else
            local exit_code=$?
            # 检查目标主机磁盘空间
            local disk_space
            disk_space=$(ssh -p $SSH_PORT "${INSTALL_USER}@${host}" "df -h $dest" 2>/dev/null | tail -n1 | awk '{print $4}')
            log_warning "分发失败，目标目录可用空间: ${disk_space:-未知}"
        fi
        
        retries=$((retries + 1))
        if [ $retries -lt $MAX_RETRIES ]; then
            log_warning "分发到 ${host} 失败，重试 $retries/$MAX_RETRIES..."
            sleep 2
        fi
    done
    
    log_error "分发到 ${host} 失败，已重试 $MAX_RETRIES 次"
    return 1
}

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
check_port() {
    local host=$1
    local port=$2
    
    # 尝试多种方式检查端口
    if command -v nc >/dev/null 2>&1; then
        # 如果有nc命令，优先使用
        if nc -z -w2 "$host" "$port" >/dev/null 2>&1; then
            log_error "端口 $port 在 $host 上已被占用"
        fi
    elif command -v telnet >/dev/null 2>&1; then
        # 如果有telnet，使用telnet
        if echo quit | timeout 2 telnet "$host" "$port" >/dev/null 2>&1; then
            log_error "端口 $port 在 $host 上已被占用"
        fi
    else
        # 最后尝试/dev/tcp
        if timeout 2 bash -c "echo >/dev/tcp/$host/$port" >/dev/null 2>&1; then
            log_error "端口 $port 在 $host 上已被占用"
        fi
    fi
}

# 替换文件内容
replace_in_file() {
    local search=$1
    local replace=$2
    local file=$3
    sed -i "s|$search|$replace|g" "$file"
}

# 读取配置文件
read_config() {
    local config_file="$EXEC_PATH/config.properties"
    check_file "$config_file"
    
    # 读取基础配置
    SEATUNNEL_VERSION=$(grep "^SEATUNNEL_VERSION=" "$config_file" | cut -d'=' -f2)
    BASE_DIR=$(grep "^BASE_DIR=" "$config_file" | cut -d'=' -f2)
    SSH_PORT=$(grep "^SSH_PORT=" "$config_file" | cut -d'=' -f2)
    DEPLOY_MODE=$(grep "^DEPLOY_MODE=" "$config_file" | cut -d'=' -f2)
    
    # 设置下载目录
    DOWNLOAD_DIR="${BASE_DIR}/downloads"
    mkdir -p "$DOWNLOAD_DIR"
    setup_permissions "$DOWNLOAD_DIR"
    
    # 读取用户配置
    INSTALL_USER=$(grep "^INSTALL_USER=" "$config_file" | cut -d'=' -f2)
    INSTALL_GROUP=$(grep "^INSTALL_GROUP=" "$config_file" | cut -d'=' -f2)
    
    # 根据部署模式读取节点配置
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：读取所有集群节点
        CLUSTER_NODES_STRING=$(grep "^CLUSTER_NODES=" "$config_file" | cut -d'=' -f2)
        [[ -z "$CLUSTER_NODES_STRING" ]] && log_error "CLUSTER_NODES 未配置"
        IFS=',' read -r -a ALL_NODES <<< "$CLUSTER_NODES_STRING"
    else
        # 分离模式：读取master和worker节点
        MASTER_IPS_STRING=$(grep "^MASTER_IP=" "$config_file" | cut -d'=' -f2)
        WORKER_IPS_STRING=$(grep "^WORKER_IPS=" "$config_file" | cut -d'=' -f2)
        [[ -z "$MASTER_IPS_STRING" ]] && log_error "MASTER_IP 未配置"
        [[ -z "$WORKER_IPS_STRING" ]] && log_error "WORKER_IPS 未配置"
        
        # 转换为数组
        IFS=',' read -r -a MASTER_IPS <<< "$MASTER_IPS_STRING"
        IFS=',' read -r -a WORKER_IPS <<< "$WORKER_IPS_STRING"
        ALL_NODES=("${MASTER_IPS[@]}" "${WORKER_IPS[@]}")
    fi
    
    # 设置SEATUNNEL_HOME
    SEATUNNEL_HOME="$BASE_DIR/apache-seatunnel-$SEATUNNEL_VERSION"
    
    # 验证必要的配置
    [[ -z "$SEATUNNEL_VERSION" ]] && log_error "SEATUNNEL_VERSION 未配置"
    [[ -z "$BASE_DIR" ]] && log_error "BASE_DIR 未配置"
    
    # 添加用户配置验证
    [[ -z "$INSTALL_USER" ]] && log_error "INSTALL_USER 未配置"
    [[ -z "$INSTALL_GROUP" ]] && log_error "INSTALL_GROUP 未配置"
    
    # 验证部署模式
    if [[ "$DEPLOY_MODE" != "hybrid" && "$DEPLOY_MODE" != "separated" ]]; then
        log_error "DEPLOY_MODE 必须是 hybrid:混合模式 或 separated:分离模式"
    fi
    
    # 读取JVM内存配置
    HYBRID_HEAP_SIZE=$(grep "^HYBRID_HEAP_SIZE=" "$config_file" | cut -d'=' -f2)
    MASTER_HEAP_SIZE=$(grep "^MASTER_HEAP_SIZE=" "$config_file" | cut -d'=' -f2)
    WORKER_HEAP_SIZE=$(grep "^WORKER_HEAP_SIZE=" "$config_file" | cut -d'=' -f2)
    
    # 验证内存配置
    [[ -z "$HYBRID_HEAP_SIZE" ]] && log_error "HYBRID_HEAP_SIZE 未配置"
    [[ -z "$MASTER_HEAP_SIZE" ]] && log_error "MASTER_HEAP_SIZE 未配置"
    [[ -z "$WORKER_HEAP_SIZE" ]] && log_error "WORKER_HEAP_SIZE 未配置"
    
    # 读取安装模式配置
    INSTALL_MODE=$(grep "^INSTALL_MODE=" "$config_file" | cut -d'=' -f2)
    [[ -z "$INSTALL_MODE" ]] && log_error "INSTALL_MODE 未配置"
    
    # 验证安装模式
    if [[ "$INSTALL_MODE" != "online" && "$INSTALL_MODE" != "offline" ]]; then
        log_error "INSTALL_MODE 必须是 online:在线安装 或 offline:离线安装"
    fi
    
    # 读取安装包相关配置
    if [[ "$INSTALL_MODE" == "offline" ]]; then
        PACKAGE_PATH=$(grep "^PACKAGE_PATH=" "$config_file" | cut -d'=' -f2)
        [[ -z "$PACKAGE_PATH" ]] && log_error "离线安装模式下 PACKAGE_PATH 未配置"
        
        # 处理版本号变量
        PACKAGE_PATH=$(echo "$PACKAGE_PATH" | sed "s/\${SEATUNNEL_VERSION}/$SEATUNNEL_VERSION/g")
        
        # 转换为绝对路径
        if [[ "$PACKAGE_PATH" != /* ]]; then
            PACKAGE_PATH="$EXEC_PATH/$PACKAGE_PATH"
        fi
    else
        # 读取安装包仓库配置
        PACKAGE_REPO=$(grep "^PACKAGE_REPO=" "$config_file" | cut -d'=' -f2)
        PACKAGE_REPO=${PACKAGE_REPO:-aliyun}  # 默认使用aliyun源
        
        # 验证仓库配置
        if [[ "$PACKAGE_REPO" == "custom" ]]; then
            CUSTOM_PACKAGE_URL=$(grep "^CUSTOM_PACKAGE_URL=" "$config_file" | cut -d'=' -f2)
            [[ -z "$CUSTOM_PACKAGE_URL" ]] && log_error "使用自定义仓库(PACKAGE_REPO=custom)时必须配置 CUSTOM_PACKAGE_URL"
        else
            [[ -z "${PACKAGE_REPOS[$PACKAGE_REPO]}" ]] && log_error "不支持的安装包仓库: $PACKAGE_REPO"
        fi
    fi

    # 读取连接器配置
    INSTALL_CONNECTORS=$(grep "^INSTALL_CONNECTORS=" "$config_file" | cut -d'=' -f2)
    INSTALL_CONNECTORS=${INSTALL_CONNECTORS:-true}  # 默认安装

    if [ "$INSTALL_CONNECTORS" = "true" ]; then
        CONNECTORS=$(grep "^CONNECTORS=" "$config_file" | cut -d'=' -f2)
        PLUGIN_REPO=$(grep "^PLUGIN_REPO=" "$config_file" | cut -d'=' -f2)
        PLUGIN_REPO=${PLUGIN_REPO:-aliyun}  # 默认使用aliyun
    fi
    
    # 读取检查点存储配置
    CHECKPOINT_STORAGE_TYPE=$(grep "^CHECKPOINT_STORAGE_TYPE=" "$config_file" | cut -d'=' -f2)
    CHECKPOINT_NAMESPACE=$(grep "^CHECKPOINT_NAMESPACE=" "$config_file" | cut -d'=' -f2)
    
    # 根据存储类型读取相应配置
    case "$CHECKPOINT_STORAGE_TYPE" in
        "HDFS")
            HDFS_NAMENODE_HOST=$(grep "^HDFS_NAMENODE_HOST=" "$config_file" | cut -d'=' -f2)
            HDFS_NAMENODE_PORT=$(grep "^HDFS_NAMENODE_PORT=" "$config_file" | cut -d'=' -f2)
            [[ -z "$HDFS_NAMENODE_HOST" ]] && log_error "HDFS模式下必须配置 HDFS_NAMENODE_HOST"
            [[ -z "$HDFS_NAMENODE_PORT" ]] && log_error "HDFS模式下必须配置 HDFS_NAMENODE_PORT"
            ;;
        "OSS"|"S3")
            STORAGE_ENDPOINT=$(grep "^STORAGE_ENDPOINT=" "$config_file" | cut -d'=' -f2)
            STORAGE_ACCESS_KEY=$(grep "^STORAGE_ACCESS_KEY=" "$config_file" | cut -d'=' -f2)
            STORAGE_SECRET_KEY=$(grep "^STORAGE_SECRET_KEY=" "$config_file" | cut -d'=' -f2)
            STORAGE_BUCKET=$(grep "^STORAGE_BUCKET=" "$config_file" | cut -d'=' -f2)
            [[ -z "$STORAGE_ENDPOINT" ]] && log_error "${CHECKPOINT_STORAGE_TYPE}模式下必须配置 STORAGE_ENDPOINT"
            [[ -z "$STORAGE_ACCESS_KEY" ]] && log_error "${CHECKPOINT_STORAGE_TYPE}模式下必须配置 STORAGE_ACCESS_KEY"
            [[ -z "$STORAGE_SECRET_KEY" ]] && log_error "${CHECKPOINT_STORAGE_TYPE}模式下必须配置 STORAGE_SECRET_KEY"
            [[ -z "$STORAGE_BUCKET" ]] && log_error "${CHECKPOINT_STORAGE_TYPE}模式下必须配置 STORAGE_BUCKET"
            ;;
        "LOCAL_FILE")
            # 本地文件模式下使用默认路径
            CHECKPOINT_NAMESPACE="$SEATUNNEL_HOME/checkpoint"
            ;;
        "")
            log_error "必须配置 CHECKPOINT_STORAGE_TYPE"
            ;;
        *)
            log_error "不支持的检查点存储类型: $CHECKPOINT_STORAGE_TYPE"
            ;;
    esac

    # 读取开机自启动配置
    ENABLE_AUTO_START=$(grep "^ENABLE_AUTO_START=" "$config_file" | cut -d'=' -f2)
    ENABLE_AUTO_START=${ENABLE_AUTO_START:-true}  

    # 读取连接器配置
    INSTALL_CONNECTORS=$(grep "^INSTALL_CONNECTORS=" "$config_file" | cut -d'=' -f2)
    INSTALL_CONNECTORS=${INSTALL_CONNECTORS:-true}  # 默认安装
    
    if [ "$INSTALL_CONNECTORS" = "true" ]; then
        CONNECTORS=$(grep "^CONNECTORS=" "$config_file" | cut -d'=' -f2)
        PLUGIN_REPO=$(grep "^PLUGIN_REPO=" "$config_file" | cut -d'=' -f2)
        PLUGIN_REPO=${PLUGIN_REPO:-aliyun}  # 默认使用aliyun
    fi
    
    # 读取端口配置
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        HYBRID_PORT=$(grep "^HYBRID_PORT=" "$config_file" | cut -d'=' -f2)
        HYBRID_PORT=${HYBRID_PORT:-5801}  # 默认端口5801
    else
        MASTER_PORT=$(grep "^MASTER_PORT=" "$config_file" | cut -d'=' -f2)
        WORKER_PORT=$(grep "^WORKER_PORT=" "$config_file" | cut -d'=' -f2)
        MASTER_PORT=${MASTER_PORT:-5801}  # 默认端口5801
        WORKER_PORT=${WORKER_PORT:-5802}  # 默认端口5802
    fi
}

# 检查用户配置
check_user() {
    log_info "检查用户配置..."
    
    # 检查sudo权限
    if ! sudo -v &>/dev/null; then
        log_error "当前用户没有sudo权限，请确保用户在sudoers中"
    fi
    
    # 检查指定用户是否存在
    if ! id "$INSTALL_USER" >/dev/null 2>&1; then
        log_error "用户 $INSTALL_USER 不存在，请先创建用户"
    fi
    
    # 检查用户组是否存在
    if ! getent group "$INSTALL_GROUP" >/dev/null; then
        log_error "用户组 $INSTALL_GROUP 不存在，请先创建用户组"
    fi
}

# 设置目录权限
setup_permissions() {
    local dir=$1
    log_info "设置目录权限: $dir"
    sudo chown -R "$INSTALL_USER:$INSTALL_GROUP" "$dir"
    sudo chmod -R 755 "$dir"
}

# 指定用户执行命令
run_as_user() {
    sudo -u "$INSTALL_USER" bash -c "$1"
}

# 创建目录
create_directory() {
    local dir=$1
    sudo mkdir -p "$dir"
}

# 临时文件管理
create_temp_file() {
    local temp_dir="$LOG_DIR/temp"
    local temp_file
    
    # 创建临时目录,如果不存在）
    if [ ! -d "$temp_dir" ]; then
        mkdir -p "$temp_dir" || log_error "无法创建临时目录: $temp_dir"
        chmod 700 "$temp_dir"
        [ -n "$INSTALL_USER" ] && [ -n "$INSTALL_GROUP" ] && chown "$INSTALL_USER:$INSTALL_GROUP" "$temp_dir"
    fi
    
    # 创建临时文件
    temp_file=$(mktemp "$temp_dir/temp.XXXXXX") || log_error "无法创建临时文件"
    chmod 600 "$temp_file"
    [ -n "$INSTALL_USER" ] && [ -n "$INSTALL_GROUP" ] && chown "$INSTALL_USER:$INSTALL_GROUP" "$temp_file"
    
    # 添加到临时文件列表
    TEMP_FILES+=("$temp_file")
    
    echo "$temp_file"
}

# 清理临时文件
cleanup_temp_files() {
    log_info "进入cleanup_temp_files函数..."
    
    # 检查数组是否已定义
    if [ -z "${TEMP_FILES+x}" ]; then
        log_info "TEMP_FILES数组未定义，退出清理"
        return 0
    fi
    
    log_info "当前TEMP_FILES数组大小: ${#TEMP_FILES[@]}"
    
    local temp_dir="$LOG_DIR/temp"
    
    # 检查是否有临时文件需要清理
    if [ ${#TEMP_FILES[@]} -eq 0 ]; then
        log_info "没有临时文件需要清理"
        return 0
    fi
    
    log_info "清理临时文件..."
    
    # 清理所有已记录的临时文件
    for temp_file in "${TEMP_FILES[@]}"; do
        if [ -f "$temp_file" ]; then
            rm -f "$temp_file"
            log_info "已删除临时文件: $temp_file"
        fi
    done
    
    # 清理临时目录（如果为空）
    if [ -d "$temp_dir" ] && [ -z "$(ls -A "$temp_dir")" ]; then
        rm -rf "$temp_dir"
        log_info "已删除空的临时目录: $temp_dir"
    fi
}

# 在文件开头添加参数处理
# 默认使用自动选择模式
FORCE_SED=false
ONLY_INSTALL_PLUGINS=false

# 解析命令行参数
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --force-sed) FORCE_SED=true ;;
        --install-plugins) ONLY_INSTALL_PLUGINS=true ;;
        *) log_error "未知参数: $1" ;;
    esac
    shift
done

# 修改check_command函数
check_command() {
    local cmd=$1
    if [ "$cmd" = "awk" ]; then
        # 检查awk是否可用
        command -v "$cmd" >/dev/null 2>&1
    else
        # 其他命令强制使用指定的模式
        if [ "$FORCE_SED" = true ]; then
            return 1
        fi
        command -v "$cmd" >/dev/null 2>&1
    fi
}

# 更新replace_yaml_section函数
replace_yaml_section() {
    local file=$1        # yaml文件路径
    local section=$2     # 要替换的部分的开始标记（如 member-list:, plugin-config:）
    local indent=$3      # 新内容的缩进空格数
    local content=$4     # 新的内容
    local temp_file
    
    log_info "修改配置文件: $file, 替换部分: $section"
    
    # 创建临时文件
    temp_file=$(create_temp_file)
    
    if [ "$FORCE_SED" = true ]; then
        log_info "强制使用sed处理文件..."
        # 获取section的缩进和完整行
        local section_line
        section_line=$(grep "$section" "$file")
        local section_indent
        section_indent=$(echo "$section_line" | sed 's/[^[:space:]].*//' | wc -c)
        
        # 预处理内容，添加缩进
        local indented_content
        indented_content=$(echo "$content" | sed "s/^/$(printf '%*s' "$indent" '')/")
        
        # 创建临时文件存储新内容
        local content_file=$(create_temp_file)
        echo "$indented_content" > "$content_file"
        
        # 第一步：找section的起始行号
        local start_line
        start_line=$(grep -n "$section" "$file" | cut -d: -f1)
        
        # 第二步：找到section的结束行号
        local end_line
        end_line=$(tail -n +$((start_line + 1)) "$file" | grep -n "^[[:space:]]\{0,$section_indent\}[^[:space:]]" | head -1 | cut -d: -f1)
        end_line=$((start_line + end_line))
        
        # 第三步：组合新文件
        # 1. 复制section之前的内容
        sed -n "1,${start_line}p" "$file" > "$temp_file"
        # 2. 添加新内容
        cat "$content_file" >> "$temp_file"
        # 3. 复制section之后的内容
        sed -n "$((end_line)),\$p" "$file" >> "$temp_file"
        
        # 清理临时内容文件
        rm -f "$content_file"
    else
        log_info "使用awk处理文件..."
        # awk版本的实现
        awk -v section="$section" -v base_indent="$indent" -v content="$content" '
        # 计算行的缩进空格数
        function get_indent(line) {
            match(line, /^[[:space:]]*/)
            return RLENGTH
        }
        
        # 为每行添加缩进的函数
        function add_indent(str, indent,    lines, i, result) {
            split(str, lines, "\n")
            result = ""
            for (i = 1; i <= length(lines); i++) {
                if (lines[i] != "") {
                    result = result sprintf("%*s%s\n", indent, "", lines[i])
                }
            }
            return result
        }
        
        BEGIN { 
            in_section = 0
            section_indent = -1
            # 预处理content，添加缩进
            indented_content = add_indent(content, base_indent)
        }
        {
            current_indent = get_indent($0)
            
            if ($0 ~ section) {
                # 找到section，记录其缩进级别
                section_indent = current_indent
                print $0
                printf "%s", indented_content
                in_section = 1
                next
            }
            
            if (in_section) {
                # 如果当前行的缩进小于等section的缩进，说明section结束
                if (current_indent <= section_indent && $0 !~ "^[[:space:]]*$") {
                    in_section = 0
                    section_indent = -1
                    print $0
                }
            } else {
                print $0
            }
        }' "$file" > "$temp_file"
    fi
    
    # 获取文件权限，使用ls -l作为备选方案
    local file_perms
    if stat --version 2>/dev/null | grep -q 'GNU coreutils'; then
        # GNU stat
        file_perms=$(stat -c %a "$file")
    else
        # 其他系统，使用ls -l解析
        file_perms=$(ls -l "$file" | cut -d ' ' -f1 | tr 'rwx-' '7500' | sed 's/^.\(.*\)/\1/' | tr -d '\n')
    fi
    
    # 复制新内容到原文件
    cp "$temp_file" "$file"
    
    # 恢复文件权限
    chmod "$file_perms" "$file"
}

# 修改hazelcast配置文件
modify_hazelcast_config() {
    local config_file=$1
    local content
    
    # 备份文件
    cp "$config_file" "${config_file}.bak"
    
    case "$config_file" in
        *"hazelcast.yaml")
            log_info "修改 hazelcast.yaml (集群通信配置)..."
            if [ "$DEPLOY_MODE" = "hybrid" ]; then
                # 生成新的member-list内容
                content=$(for node in "${ALL_NODES[@]}"; do
                    echo "- ${node}:${HYBRID_PORT:-5801}"
                done)
                replace_yaml_section "$config_file" "member-list:" 10 "$content"
                
                # 修改端口配置
                sed -i "s/port: [0-9]\+/port: ${HYBRID_PORT:-5801}/" "$config_file"
            fi
            ;;
        *"hazelcast-client.yaml")
            log_info "修改 hazelcast-client.yaml (客户端连接配置)..."
            # 生成新的cluster-members内容
            if [ "$DEPLOY_MODE" = "hybrid" ]; then
                log_info "混合模式: 客户端可连接任意节点的 ${HYBRID_PORT:-5801} 端口"
                content=$(for node in "${ALL_NODES[@]}"; do
                    echo "- ${node}:${HYBRID_PORT:-5801}"
                done)
            else
                log_info "分离模式: 客户端仅连接Master节点的 ${MASTER_PORT:-5801} 端口"
                content=$(for master in "${MASTER_IPS[@]}"; do
                    echo "- ${master}:${MASTER_PORT:-5801}"
                done)
            fi
            replace_yaml_section "$config_file" "cluster-members:" 6 "$content"
            ;;
        *"hazelcast-master.yaml")
            if [ "$DEPLOY_MODE" != "hybrid" ]; then
                log_info "修改 hazelcast-master.yaml (Master节点配置)..."
                log_info "分离模式: Master使用 ${MASTER_PORT:-5801} 端口，master使用 ${MASTER_PORT:-5801} 端口"
                # 生成新的member-list内容
                content=$(
                    for master in "${MASTER_IPS[@]}"; do
                        echo "- ${master}:${MASTER_PORT:-5801}"
                    done
                    for worker in "${WORKER_IPS[@]}"; do
                        echo "- ${worker}:${WORKER_PORT:-5802}"
                    done
                )
                replace_yaml_section "$config_file" "member-list:" 10 "$content"
                
                # 修改端口配置
                sed -i "s/port: [0-9]\+/port: ${MASTER_PORT:-5801}/" "$config_file"
            fi
            ;;
        *"hazelcast-worker.yaml")
            if [ "$DEPLOY_MODE" != "hybrid" ]; then
                log_info "修改 hazelcast-worker.yaml (Worker节点配置)..."
                log_info "分离模式: Master使用 ${MASTER_PORT:-5801} 端口，Worker使用 ${WORKER_PORT:-5802} 端口"
                # 生成新的member-list内容
                content=$(
                    for master in "${MASTER_IPS[@]}"; do
                        echo "- ${master}:${MASTER_PORT:-5801}"
                    done
                    for worker in "${WORKER_IPS[@]}"; do
                        echo "- ${worker}:${WORKER_PORT:-5802}"
                    done
                )
                replace_yaml_section "$config_file" "member-list:" 10 "$content"
                
                # 修改端口配置
                sed -i "s/port: [0-9]\+/port: ${WORKER_PORT:-5802}/" "$config_file"
            fi
            ;;
    esac
}

# 配置混合模式
setup_hybrid_mode() {
    log_info "配置混合模式集群..."
    
    # 配置hazelcast.yaml，所有节点使用5801端口
    modify_hazelcast_config "$SEATUNNEL_HOME/config/hazelcast.yaml"
    
    # 配置client，所有点都可以作为连接点
    modify_hazelcast_config "$SEATUNNEL_HOME/config/hazelcast-client.yaml"
    
    # 配置JVM选项
    configure_jvm_options "$SEATUNNEL_HOME/config/jvm_options" "$HYBRID_HEAP_SIZE"
}

# 配置分离模式
setup_separated_mode() {
    log_info "配置分离模式集群..."
    
    # 配置master节点
    modify_hazelcast_config "$SEATUNNEL_HOME/config/hazelcast-master.yaml"
    
    # 配置worker节点
    modify_hazelcast_config "$SEATUNNEL_HOME/config/hazelcast-worker.yaml"
    
    # 配置client
    modify_hazelcast_config "$SEATUNNEL_HOME/config/hazelcast-client.yaml"
    
    # 配置JVM选项
    configure_jvm_options "$SEATUNNEL_HOME/config/jvm_master_options" "$MASTER_HEAP_SIZE"
    configure_jvm_options "$SEATUNNEL_HOME/config/jvm_worker_options" "$WORKER_HEAP_SIZE"
}

# 启动集群
start_cluster() {
    log_info "启动SeaTunnel集群..."
    
    if [ "${ENABLE_AUTO_START}" = "true" ]; then
        # 使用systemd服务启动
        local current_ip
        current_ip=$(hostname -I | awk '{print $1}')
        
        if [ "$DEPLOY_MODE" = "hybrid" ]; then
            # 混合模式：在所有节点上启动服务
            for node in "${ALL_NODES[@]}"; do
                if [ "$node" = "localhost" ] || [ "$node" = "$current_ip" ]; then
                    log_info "在本地节点启动服务..."
                    sudo systemctl daemon-reload
                    sudo systemctl start seatunnel
                    continue
                fi
                
                log_info "在节点 $node 上启动服务..."
                ssh_with_retry "$node" "sudo systemctl daemon-reload && sudo systemctl start seatunnel"
            done
        else
            # 分离模式：根据节点角色启动对应服务
            # 启动Master节点
            for master in "${MASTER_IPS[@]}"; do
                if [ "$master" = "localhost" ] || [ "$master" = "$current_ip" ]; then
                    log_info "在本地Master节点启动服务..."
                    sudo systemctl daemon-reload
                    sudo systemctl start seatunnel-master
                    continue
                fi
                
                log_info "在Master节点 $master 上启动服务..."
                ssh_with_retry "$master" "sudo systemctl daemon-reload && sudo systemctl start seatunnel-master"
            done
            
            # 启动Worker节点
            for worker in "${WORKER_IPS[@]}"; do
                if [ "$worker" = "localhost" ] || [ "$worker" = "$current_ip" ]; then
                    log_info "在本地Worker节点启动服务..."
                    sudo systemctl daemon-reload
                    sudo systemctl start seatunnel-worker
                    continue
                fi
                
                log_info "在Worker节点 $worker 上启动服务..."
                ssh_with_retry "$worker" "sudo systemctl daemon-reload && sudo systemctl start seatunnel-worker"
            done
        fi
    else
        # 使用脚本启动
        sudo chmod +x "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh"
        run_as_user "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh start"
    fi
}

# 配置JVM选项
configure_jvm_options() {
    local file=$1
    local heap_size=$2
    
    log_info "配置JVM选项: $file (堆内存: ${heap_size}g)"
    
    # 备份原始文件
    cp "$file" "${file}.bak"
    
    # 修改JVM堆内存配置
    sed -i "s/-Xms[0-9]\+g/-Xms${heap_size}g/" "$file"
    sed -i "s/-Xmx[0-9]\+g/-Xmx${heap_size}g/" "$file"
}

# 检查端口占用
check_ports() {
    log_info "检查端口占用..."
    local occupied_ports=()
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        local port=${HYBRID_PORT:-5801}
        for node in "${ALL_NODES[@]}"; do
            if ! check_port "$node" "$port" 2>/dev/null; then
                occupied_ports+=("$node:$port")
            fi
        done
    else
        local master_port=${MASTER_PORT:-5801}
        local worker_port=${WORKER_PORT:-5802}
        
        for master in "${MASTER_IPS[@]}"; do
            if ! check_port "$master" "$master_port" 2>/dev/null; then
                occupied_ports+=("$master:$master_port")
            fi
        done
        for worker in "${WORKER_IPS[@]}"; do
            if ! check_port "$worker" "$worker_port" 2>/dev/null; then
                occupied_ports+=("$worker:$worker_port")
            fi
        done
    fi
    
    if [ ${#occupied_ports[@]} -gt 0 ]; then
        log_error "以下端口已被占用:\n${occupied_ports[*]}"
    fi
}

# 检查服务状态
check_services() {
    log_info "检查服务状态..."
    
    # 等待服务启动
    log_info "等待服务启动（10秒）..."
    sleep 10
    
    # 检查所有节点的服务状态
    local nodes=()
    local ports=()
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：所有节点使用相同端口
        nodes=("${ALL_NODES[@]}")
        for node in "${nodes[@]}"; do
            ports+=("${HYBRID_PORT:-5801}")
        done
    else
        # 分离模式：收集所有节点和对应端口
        for master in "${MASTER_IPS[@]}"; do
            nodes+=("$master")
            ports+=("${MASTER_PORT:-5801}")
        done
        for worker in "${WORKER_IPS[@]}"; do
            nodes+=("$worker")
            ports+=("${WORKER_PORT:-5802}")
        done
    fi
    
    # 检查每个节点的服务状态
    local success=true
    for i in "${!nodes[@]}"; do
        local node="${nodes[$i]}"
        local port="${ports[$i]}"
        
        log_info "检查节点 $node:$port 的服务状态..."
        
        # 处理localhost的情况
        if [ "$node" = "localhost" ]; then
            node="127.0.0.1"
        fi
        
        # 尝试多种方式检查端口
        local service_running=false
        
        if command -v nc >/dev/null 2>&1; then
            # 使用nc命令检查
            log_info "使用nc命令检查节点 $node:$port 的服务状态..."
            if [ "$node" = "127.0.0.1" ] || [ "$node" = "$(hostname -I | awk '{print $1}')" ]; then
                # 本地检查使用localhost
                if nc -z -w2 localhost "$port" >/dev/null 2>&1; then
                    service_running=true
                fi
            else
                # 远程节点检查
                if nc -z -w2 "$node" "$port" >/dev/null 2>&1; then
                    service_running=true
                fi
            fi
        else
            # 使用/dev/tcp
            log_info "使用/dev/tcp命令检查节点 $node:$port 的服务状态..."
            if [ "$node" = "127.0.0.1" ] || [ "$node" = "$(hostname -I | awk '{print $1}')" ]; then
                # 本地检查使用localhost
                if timeout 2 bash -c "echo >/dev/tcp/localhost/$port" >/dev/null 2>&1; then
                    service_running=true
                fi
            else
                # 远程节点检查
                if timeout 2 bash -c "echo >/dev/tcp/$node/$port" >/dev/null 2>&1; then
                    service_running=true
                fi
            fi
        fi
        
        if [ "$service_running" = true ]; then
            log_success "节点 $node:$port 服务运行正常"
        else
            log_warning "节点 $node:$port 服务未响应"
            success=false
        fi
    done
    
    if [ "$success" = true ]; then
        log_success "所有节点服务检查通过"
    else
        log_warning "部分节点服务检查未通过，请检查日志确认具体原因"
    fi
}

# 配置检查点���储
configure_checkpoint() {
    # 计算实际节点数（排除localhost）
    local actual_node_count=0
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        for node in "${ALL_NODES[@]}"; do
            if [ "$node" != "localhost" ]; then
                actual_node_count=$((actual_node_count + 1))
            fi
        done
    else
        # 分离模式：计算master和worker节点总数
        for master in "${MASTER_IPS[@]}"; do
            if [ "$master" != "localhost" ]; then
                actual_node_count=$((actual_node_count + 1))
            fi
        done
        for worker in "${WORKER_IPS[@]}"; do
            if [ "$worker" != "localhost" ]; then
                actual_node_count=$((actual_node_count + 1))
            fi
        done
    fi
    
    local content
    
    # Validate storage type
    if [[ -z "$CHECKPOINT_STORAGE_TYPE" ]]; then
        log_info "未配置检查点存储类型，使用默认配置"
        CHECKPOINT_STORAGE_TYPE="LOCAL_FILE"
    fi

    # Validate required variables based on storage type
    case "$CHECKPOINT_STORAGE_TYPE" in
        LOCAL_FILE)
            [[ -z "$CHECKPOINT_NAMESPACE" ]] && log_error "LOCAL_FILE 模式需要配置 CHECKPOINT_NAMESPACE"
            ;;
        HDFS)
            [[ -z "$CHECKPOINT_NAMESPACE" ]] && log_error "HDFS 模式需要配置 CHECKPOINT_NAMESPACE"
            [[ -z "$HDFS_NAMENODE_HOST" ]] && log_error "HDFS 模式需要配置 HDFS_NAMENODE_HOST"
            [[ -z "$HDFS_NAMENODE_PORT" ]] && log_error "HDFS 模式需要配置 HDFS_NAMENODE_PORT"
            ;;
        OSS|S3)
            [[ -z "$CHECKPOINT_NAMESPACE" ]] && log_error "${CHECKPOINT_STORAGE_TYPE} 模式需要配置 CHECKPOINT_NAMESPACE"
            [[ -z "$STORAGE_BUCKET" ]] && log_error "${CHECKPOINT_STORAGE_TYPE} 模式需要配置 STORAGE_BUCKET"
            [[ -z "$STORAGE_ENDPOINT" ]] && log_error "${CHECKPOINT_STORAGE_TYPE} 模式需要配置 STORAGE_ENDPOINT"
            [[ -z "$STORAGE_ACCESS_KEY" ]] && log_error "${CHECKPOINT_STORAGE_TYPE} 模式需要配置 STORAGE_ACCESS_KEY"
            [[ -z "$STORAGE_SECRET_KEY" ]] && log_error "${CHECKPOINT_STORAGE_TYPE} 模式需要配置 STORAGE_SECRET_KEY"
            ;;
        *)
            log_error "不支持的检查点存储类型: $CHECKPOINT_STORAGE_TYPE"
            ;;
    esac
    
    # 如果是LOCAL_FILE类型，创建本地目录
    if [[ "$CHECKPOINT_STORAGE_TYPE" == "LOCAL_FILE" ]]; then
        local checkpoint_dir="$SEATUNNEL_HOME/checkpoint"
        create_directory "$checkpoint_dir"
        setup_permissions "$checkpoint_dir"
        CHECKPOINT_NAMESPACE="$checkpoint_dir"
        
        # 在其他节点上创建目录（排除localhost）
        for node in "${ALL_NODES[@]}"; do
            # 跳过localhost和当前节点
            if [ "$node" = "localhost" ] || [ "$node" = "$(hostname -I | awk '{print $1}')" ]; then
                log_info "跳过本地节点: $node"
                continue
            fi
            
            log_info "在节点 $node 上创建检查点目录..."
            ssh_with_retry "$node" "mkdir -p $checkpoint_dir && chown $INSTALL_USER:$INSTALL_GROUP $checkpoint_dir && chmod 755 $checkpoint_dir"
        done
        
        # 只有在实际节点数大于1时才显示警告
        if [ "$actual_node_count" -gt 1 ]; then
            log_warning "检测到多节点部署，不建议使用本地文件存储作检查点。建议使用 HDFS、OSS 或 S3。"
        fi
    fi
    
    # 根据存储类型生成配置内容
    case "$CHECKPOINT_STORAGE_TYPE" in
        LOCAL_FILE)
            content="namespace: ${CHECKPOINT_NAMESPACE}
storage.type: local"
            ;;
        HDFS)
            content="namespace: ${CHECKPOINT_NAMESPACE}
storage.type: hdfs
fs.defaultFS: hdfs://${HDFS_NAMENODE_HOST}:${HDFS_NAMENODE_PORT}"
            if [ ! -z "${KERBEROS_PRINCIPAL:-}" ] && [ ! -z "${KERBEROS_KEYTAB:-}" ]; then
                content+="
kerberosPrincipal: ${KERBEROS_PRINCIPAL}
kerberosKeytabFilePath: ${KERBEROS_KEYTAB}"
            fi
            ;;
        OSS)
            content="namespace: ${CHECKPOINT_NAMESPACE}
storage.type: oss
oss.bucket: ${STORAGE_BUCKET}
fs.oss.endpoint: ${STORAGE_ENDPOINT}
fs.oss.accessKeyId: ${STORAGE_ACCESS_KEY}
fs.oss.accessKeySecret: ${STORAGE_SECRET_KEY}"
            ;;
        S3)
            content="namespace: ${CHECKPOINT_NAMESPACE}
storage.type: s3
s3.bucket: ${STORAGE_BUCKET}
fs.s3a.endpoint: ${STORAGE_ENDPOINT}
fs.s3a.access.key: ${STORAGE_ACCESS_KEY}
fs.s3a.secret.key: ${STORAGE_SECRET_KEY}
fs.s3a.aws.credentials.provider: org.apache.hadoop.fs.s3a.SimpleAWSCredentialsProvider
disable.cache: true"
            ;;
    esac
    
    # 替换plugin-config部分，调整缩进值为10
    replace_yaml_section "$SEATUNNEL_HOME/config/seatunnel.yaml" "plugin-config:" 10 "$content"
}

# 检查Java环境
check_java() {
    local node=$1
    local is_remote=$2
    
    log_info "检查节点 $node 的Java环境..."
    
    # 本地节点检查
    if [ "$is_remote" = "false" ]; then
        # 检查java命令是否存在
        if ! command -v java >/dev/null 2>&1; then
            log_warning "本地节点未找到Java环境"
            
            if [ "$INSTALL_MODE" != "online" ]; then
                log_error "离线模式下无法自动安装Java，请手动安装Java 8或Java 11"
            fi
            
            # 提示用户选择安装版本
            echo -e "\n${YELLOW}请选择要安装的Java版本:${NC}"
            echo "1) Java 8 (推荐)"
            echo "2) Java 11"
            echo "3) 取消安装"
            
            read -r -p "请输入选项 [1-3]: " choice
            
            case $choice in
                1)
                    install_java "8"
                    ;;
                2)
                    install_java "11"
                    ;;
                3)
                    log_error "用户取消安装"
                    ;;
                *)
                    log_error "无效的选项"
                    ;;
            esac
        fi
        
        # 获取本地Java版本
        local java_version
        java_version=$(java -version 2>&1 | head -n 1 | awk -F '"' '{print $2}')
        if [ -z "$java_version" ]; then
            log_error "无法获取本地Java版本"
        fi
        
        # 检查Java版本是否为8或11
        if [[ $java_version == 1.8* ]]; then
            log_info "节点 $node 检测到Java 8: $java_version"
        elif [[ $java_version == 11* ]]; then
            log_info "节点 $node 检测到Java 11: $java_version"
        else
            log_error "节点 $node 不支持的Java版本: $java_version，SeaTunnel需要Java 8或Java 11"
        fi
    else
        # 远程节点检查
        # 添加超时控制
        local TIMEOUT=30
        
        # 先检查java命令是否存在
        if ! timeout $TIMEOUT ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "command -v java" >/dev/null 2>&1; then
            log_warning "节点 $node 未找到Java环境"
            if [ "$INSTALL_MODE" != "online" ]; then
                log_error "离线模式下无法自动安装Java，请在节点 $node 上手动安装Java 8或Java 11"
            fi
            # 在线模式下自动安装Java 8
            log_info "在节点 $node 上自动安装Java 8..."
            install_java "8" "$node"
            return
        fi

        # 获取远程Java版本输出
        local java_version_output
        java_version_output=$(timeout $TIMEOUT ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "java -version 2>&1")
        local exit_code=$?

        # 检查是否超时
        if [ $exit_code = 124 ]; then
            log_error "检查节点 $node 的Java环境超时(${TIMEOUT}秒)"
            return 1
        fi

        # 在本地解析Java版本
        local java_version
        java_version=$(echo "$java_version_output" | head -n 1 | awk -F '"' '{print $2}')

        if [ -z "$java_version" ]; then
            log_error "无法获取节点 $node 的Java版本"
            return 1
        fi

        # 在本地检查版本兼容性
        if [[ $java_version == 1.8* ]]; then
            log_info "节点 $node 检测到Java 8: $java_version"
        elif [[ $java_version == 11* ]]; then
            log_info "节点 $node 检测到Java 11: $java_version"
        else
            log_error "节点 $node 不支持的Java版本: $java_version，SeaTunnel需要Java 8或Java 11"
        fi
    fi
}

# 安装Java
install_java() {
    local version=$1
    local node=${2:-"localhost"}
    local java_home="${BASE_DIR}/java"
    local is_remote=false
    
    # 检查是否是远程节点
    if [ "$node" != "localhost" ] && [ "$node" != "$(hostname -I | awk '{print $1}')" ]; then
        is_remote=true
    fi
    
    # 检测系统架构
    local arch
    if [ "$is_remote" = true ]; then
        arch=$(ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "uname -m")
    else
        arch=$(uname -m)
    fi
    
    case "$arch" in
        x86_64)
            arch_suffix="x64"
            ;;
        aarch64)
            arch_suffix="aarch64"
            ;;
        *)
            log_error "不支持的系统架构: $arch"
            ;;
    esac
    
    log_info "开始在节点 $node 上安装Java $version, 系统架构: $arch"
    
    # 构建下载URL和包名
    local download_url
    local java_package
    local java_dir
    
    case $version in
        "8")
            download_url="https://repo.huaweicloud.com/java/jdk/8u202-b08/jdk-8u202-linux-${arch_suffix}.tar.gz"
            java_package="jdk-8u202-linux-${arch_suffix}.tar.gz"
            java_dir="jdk1.8.0_202"
            ;;
        "11")
            download_url="https://repo.huaweicloud.com/java/jdk/11.0.2+9/jdk-11.0.2_linux-${arch_suffix}_bin.tar.gz"
            java_package="jdk-11.0.2_linux-${arch_suffix}_bin.tar.gz"
            java_dir="jdk-11.0.2"
            ;;
        *)
            log_error "不支持的Java版本: $version"
            ;;
    esac
    
    # 使用全局下载目录
    cd "$DOWNLOAD_DIR" || log_error "无法进入下载目录"
    
    # 下载Java安装包
    log_info "下载Java安装包..."
    if ! curl -L --progress-bar -o "$java_package" "$download_url"; then
        # 如果华为云下载失败,尝试清华源
        log_warning "从华为云下载失败,尝试清华源..."
        case $version in
            "8")
                download_url="https://mirrors.tuna.tsinghua.edu.cn/Adoptium/8/jdk/${arch_suffix}/linux/OpenJDK8U-jdk_${arch_suffix}_linux_hotspot_8u432b06.tar.gz"
                java_package="OpenJDK8U-jdk_${arch_suffix}_linux_hotspot_8u432b06.tar.gz"
                java_dir="jdk8u432-b06"
                ;;
            "11")
                download_url="https://mirrors.tuna.tsinghua.edu.cn/Adoptium/11/jdk/${arch_suffix}/linux/OpenJDK11U-jdk_${arch_suffix}_linux_hotspot_11.0.25_9.tar.gz"
                java_package="OpenJDK11U-jdk_${arch_suffix}_linux_hotspot_11.0.25_9.tar.gz"
                java_dir="jdk-11.0.25+9"
                ;;
        esac
        
        if ! curl -L --progress-bar -o "$java_package" "$download_url"; then
            log_error "Java安装包下载失败"
        fi
    fi
    
    # 创建Java安装目录
    if [ "$is_remote" = true ]; then
        ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "mkdir -p $java_home"
        scp -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$DOWNLOAD_DIR/$java_package" "${INSTALL_USER}@${node}:$java_home/"
        ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "cd $java_home && tar -zxf $java_package && rm -f $java_package"
    else
        mkdir -p "$java_home"
        tar -zxf "$java_package" -C "$java_home"
    fi
    
    # 设置权限
    if [ "$is_remote" = true ]; then
        ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "chown -R $INSTALL_USER:$INSTALL_GROUP $java_home && chmod -R 755 $java_home"
    else
        chown -R "$INSTALL_USER:$INSTALL_GROUP" "$java_home"
        chmod -R 755 "$java_home"
    fi
    
    # 配置环境变量
    local bashrc_content="
# JAVA_HOME BEGIN
export JAVA_HOME=$java_home/$java_dir
export PATH=\$JAVA_HOME/bin:\$PATH
# JAVA_HOME END"
    
    if [ "$is_remote" = true ]; then
        local remote_home
        remote_home=$(ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "echo ~$INSTALL_USER")
        ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "
            sed -i '/# JAVA_HOME BEGIN/,/# JAVA_HOME END/d' $remote_home/.bashrc
            echo '$bashrc_content' >> $remote_home/.bashrc
            source $remote_home/.bashrc"
    else
        # 获取用户home目录
        local user_home
        if command -v getent >/dev/null 2>&1; then
            user_home=$(getent passwd "$INSTALL_USER" | cut -d: -f6)
        else
            user_home=$(eval echo ~"$INSTALL_USER")
        fi
        
        # 删除已存在的Java配置
        sed -i '/# JAVA_HOME BEGIN/,/# JAVA_HOME END/d' "$user_home/.bashrc"
        echo "$bashrc_content" >> "$user_home/.bashrc"
        
        # 使环境变量生效
        export JAVA_HOME="$java_home/$java_dir"
        export PATH="$JAVA_HOME/bin:$PATH"
    fi
    
    # 验证安装
    local verify_cmd="java -version"
    if [ "$is_remote" = true ]; then
        if ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "$verify_cmd" 2>&1 | grep -q "version"; then
            log_success "节点 $node 的Java $version 安装成功"
        else
            log_error "节点 $node 的Java安装失败"
        fi
    else
        if $verify_cmd 2>&1 | grep -q "version"; then
            log_success "本地节点Java $version 安装成功"
        else
            log_error "本地节点Java安装失败"
        fi
    fi
    
    # 清理安装文件
    rm -f "$DOWNLOAD_DIR/$java_package"
}


# 添加依赖检查函数
check_dependencies() {
    log_info "检查系统依赖..."
    
    # 必需的命令列表
    local required_cmds=("ssh" "scp" "tar" "grep" "sed")
    
    for cmd in "${required_cmds[@]}"; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            log_error "缺少必需的命令: $cmd"
        fi
    done
    
}

# 检查URL是否可访问
check_url() {
    local url=$1
    local timeout=5
    
    if ! command -v curl >/dev/null 2>&1; then
        log_error "未找到curl命令,请先安装curl"
    fi
    
    if curl --connect-timeout "$timeout" -sI "$url" >/dev/null 2>&1; then
        return 0
    fi
    return 1
}

# 下载安装包
download_package() {
    # 检查是否为在线模式
    if [ "$INSTALL_MODE" != "online" ]; then
        log_error "download_package函数只能在在线模式(INSTALL_MODE=online)下使用"
    fi
    
    local package_name=$1
    local version=$2
    local output_file="$DOWNLOAD_DIR/$package_name"
    local retries=3
    local retry_count=0
    
    log_info "开始下载安装包..."
    
    # 检查curl命令
    if ! command -v curl >/dev/null 2>&1; then
        log_error "未找到curl命令,请先安装curl"
    fi
    
    # 创建并进入下载目录
    mkdir -p "$DOWNLOAD_DIR"
    cd "$DOWNLOAD_DIR" || log_error "无法进入下载目录"
    
    # 获取仓库配置
    local repo=${PACKAGE_REPO:-aliyun}
    local url
    
    # 获取发布包下载地址
    url="${PACKAGE_REPOS[$repo]}"
    if [ -z "$url" ]; then
        log_error "不支持的安装包仓库: $repo"
    fi
    url="$url/${version}/apache-seatunnel-${version}-bin.tar.gz"
    
    log_info "使用下载源: $url"
    
    # 下载重试循环
    while [ $retry_count -lt $retries ]; do
        log_info "下载尝试 $((retry_count + 1))/$retries"
        
        # 检查URL是否可访问
        if ! check_url "$url"; then
            log_warning "当前下载源不可用,尝试切换到备用源..."
            if [ "$repo" = "aliyun" ]; then
                repo="apache"
                url="${PACKAGE_REPOS[$repo]}/${version}/apache-seatunnel-${version}-bin.tar.gz"
                continue
            fi
        fi
        
        # 使用curl下载，显示进度条
        if curl -L \
            --fail \
            --progress-bar \
            --connect-timeout 10 \
            --retry 3 \
            --retry-delay 2 \
            --retry-max-time 60 \
            -o "$output_file" \
            "$url" 2>&1; then
            
            # 验证下载文件
            if [ -f "$output_file" ] && [ -s "$output_file" ]; then
                log_info "下载完成: $output_file"
                echo "$output_file" > /tmp/download_path.tmp
                return 0
            else
                log_warning "下载文件为空或不存在"
            fi
        fi
        
        retry_count=$((retry_count + 1))
        [ $retry_count -lt $retries ] && log_warning "下载失败,等待重试..." && sleep 3
    done
    
    log_error "下载失败,已重试 $retries 次"
    return 1
}

# 添加安装包验证函数
verify_package() {
    local package_file=$1
    
    log_info "验证安装包: $package_file"
    
    # 检查文件是否存在
    if [ ! -f "$package_file" ]; then
        log_error "安装包不存在: $package_file"
    fi
    
    # 检查文件格式
    if ! file "$package_file" | grep -q "gzip compressed data"; then
        log_error "安装包格式错误,必须是tar.gz格式"
    fi
    
    # 检查文件名是否包含版本号
    if ! echo "$package_file" | grep -q "apache-seatunnel-${SEATUNNEL_VERSION}"; then
        log_warning "安装包文件名与配置的版本号不匹配: $SEATUNNEL_VERSION"
        log_warning "安装包: $package_file"
        read -r -p "是否继续安装? [y/N] " response
        case "$response" in
            [yY][eE][sS]|[yY]) 
                log_warning "继续安装..."
                ;;
            *)
                log_error "安装已取消"
                ;;
        esac
    fi
    
    log_info "安装包验证通过"
}

# 添加集群管理脚本
setup_cluster_scripts() {
    log_info "添加集群管理脚本..."
    
    # 获取脚本所在目录的绝对路径
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
    # 创建master和workers文件
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        printf "%s\n" "${CLUSTER_NODES[@]}" > "$SEATUNNEL_HOME/bin/master"
        printf "%s\n" "${CLUSTER_NODES[@]}" > "$SEATUNNEL_HOME/bin/workers"
    else
        printf "%s\n" "${MASTER_IPS[@]}" > "$SEATUNNEL_HOME/bin/master"
        printf "%s\n" "${WORKER_IPS[@]}" > "$SEATUNNEL_HOME/bin/workers"
    fi
    
    # 创建集群启动脚本
    cat > "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh" << 'EOF'
#!/bin/bash 
  
# 定义 SeaTunnelServer 进程名称，需要根据实际情况进行修改
PROCESS_NAME="org.apache.seatunnel.core.starter.seatunnel.SeaTunnelServer"

# 获取脚本所在目录的绝对路径
bin_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_USER=root

# 定义颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
  echo -e "$(date '+%Y-%m-%d %H:%M:%S') [INFO] $1"
}

log_error() {
  echo -e "$(date '+%Y-%m-%d %H:%M:%S') [${RED}ERROR${NC}] $1"
}

log_success() {
  echo -e "$(date '+%Y-%m-%d %H:%M:%S') [${GREEN}SUCCESS${NC}] $1"
}

log_warning() {
  echo -e "$(date '+%Y-%m-%d %H:%M:%S') [${YELLOW}WARNING${NC}] $1"
}

export SEATUNNEL_HOME="$(dirname "$bin_dir")"
log_info "SEATUNNEL_HOME: ${SEATUNNEL_HOME}"
master_conf="${bin_dir}/master"
workers_conf="${bin_dir}/workers"

if [ -f "$master_conf" ]; then
    mapfile -t masters < <(sed 's/[[:space:]]*$//' "$master_conf")
else
    log_error "找不到 $master_conf 文件"
    exit 1
fi

if [ -f "$workers_conf" ]; then
    mapfile -t workers < <(sed 's/[[:space:]]*$//' "$workers_conf")
else
    log_error "找不到 $workers_conf 文件"
    exit 1
fi

mapfile -t servers < <(sort -u <(sed 's/[[:space:]]*$//' "$master_conf" "$workers_conf"))

sshPort=22
EOF

    # 继续写入脚本内容...
    cat >> "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh" << 'EOF'

start(){
    echo "-------------------------------------------------"
    for master in "${masters[@]}"; do
        if [ "$master" = "localhost" ]; then
            log_warning "检测到仅有本地进程，跳过远程执行..."
            ${bin_dir}/seatunnel-cluster.sh -d  -r master
            log_success "${master}的SeaTunnel-master启动成功"
        else
            log_info "正在 ${master} 上启动 SeaTunnelServer。"
            ssh -p $sshPort -o StrictHostKeyChecking=no "${INSTALL_USER}@${master}" "source /etc/profile && source ~/.bashrc && ${bin_dir}/seatunnel-cluster.sh -d  -r master"
            log_success "${master}的SeaTunnel-master启动成功"    
        fi
    done

    for worker in "${workers[@]}"; do
        if [ "$worker" = "localhost" ]; then
            log_warning "检测到仅有本地进程，跳过远程执行..."
            ${bin_dir}/seatunnel-cluster.sh -d  -r worker
            log_success "${worker}的SeaTunnel-worker启动成功"
        else
            log_info "正在 ${worker} 上启动 SeaTunnelServer。"
            ssh -p $sshPort -o StrictHostKeyChecking=no "${INSTALL_USER}@${worker}" "source /etc/profile && source ~/.bashrc && ${bin_dir}/seatunnel-cluster.sh -d  -r worker"
            log_success "${worker}的SeaTunnel-worker启动成功"    
        fi
    done
}

stop(){
    echo "-------------------------------------------------"
    for server in "${servers[@]}"; do
        if [ "$server" = "localhost" ]; then
            log_warning "检测到仅有本地进程，跳过远程执行..."
            ${bin_dir}/stop-seatunnel-cluster.sh
            log_success "${server}的SeaTunnel 停止成功"
        else
            log_info "正在 ${server} 上停止 SeaTunnelServer"
            ssh -p $sshPort -o StrictHostKeyChecking=no "${INSTALL_USER}@${server}" "source /etc/profile && source ~/.bashrc && ${bin_dir}/stop-seatunnel-cluster.sh"
            log_success "${server}的SeaTunnel 停止成功"
        fi
    done
}

restart(){
    stop
    sleep 2
    start
}

case "$1" in
    "start")
        start
        ;;
    "stop")
        stop
        ;;
    "restart")
        restart
        ;;
    *)
        echo "用法：$0 {start|stop|restart}"
        exit 1
esac
EOF

    # 设置脚本权限
    chmod +x "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh"
    chmod 644 "$SEATUNNEL_HOME/bin/master"
    chmod 644 "$SEATUNNEL_HOME/bin/workers"
    
    # 设置所有者
    chown "$INSTALL_USER:$INSTALL_GROUP" "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh"
    chown "$INSTALL_USER:$INSTALL_GROUP" "$SEATUNNEL_HOME/bin/master"
    chown "$INSTALL_USER:$INSTALL_GROUP" "$SEATUNNEL_HOME/bin/workers"
    
    log_info "集群管理脚本添加完成"
}

# 安装插件和依赖库
install_plugins_and_libs() {
    # 检查是否需要安装连接器
    if [ "${INSTALL_CONNECTORS}" != "true" ]; then
        log_info "跳过连接器和依赖安装"
        return 0
    fi

    # 离线模式下跳过插件下载
    if [ "$INSTALL_MODE" = "offline" ]; then
        log_warning "离线安装模式下不支持自动下载插件,如果有需要,请手动将所需插件和依赖放置到以下目录:"
        log_warning "- 插件目录: $SEATUNNEL_HOME/connectors/"
        log_warning "- 依���目录: $SEATUNNEL_HOME/lib/"
        return 0
    fi
    
    log_info "开始安装插件和依赖..."
    
    # 创建目录
    local lib_dir="$SEATUNNEL_HOME/lib"
    local connectors_dir="$SEATUNNEL_HOME/connectors"
    create_directory "$lib_dir"
    create_directory "$connectors_dir"
    setup_permissions "$lib_dir"
    setup_permissions "$connectors_dir"
    
    # 如果CONNECTORS为空，使用默认值
    if [ -z "$CONNECTORS" ]; then
        CONNECTORS="jdbc,hive"
        log_info "使用默认连接器: $CONNECTORS"
    fi
    
    # 读取启用的连接器列表
    IFS=',' read -r -a enabled_connectors <<< "${CONNECTORS}"
    
    # 使用全局变量PLUGIN_REPO，已在read_config中设置
    local retries=3
    local config_file="$EXEC_PATH/config.properties"  # 添加这行
    
    # 处理每个连接器
    for connector in "${enabled_connectors[@]}"; do
        connector=$(echo "$connector" | tr -d '[:space:]')
        log_info "处理连接器: $connector"
        
        # 下载连接器插件
        local plugin_jar="connector-${connector}-${SEATUNNEL_VERSION}.jar"
        local target_path="$connectors_dir/$plugin_jar"
        
        # 检查插件是否已存在
        if [ -f "$target_path" ]; then
            log_info "连接器插件已存在: $plugin_jar"
        else
            log_info "下载连接器插件: $plugin_jar"
            download_artifact "$target_path" "$connector" "plugin"
        fi
        
        # 读取并处理连接器的依赖库
        local libs_str
        libs_str=$(grep "^${connector}_libs=" "$config_file" | cut -d'=' -f2)
        
        if [ -n "$libs_str" ]; then
            log_info "处理 $connector 连接器的依赖库..."
            IFS=',' read -r -a libs <<< "$libs_str"
            for lib in "${libs[@]}"; do
                lib=$(echo "$lib" | tr -d '[:space:]')  # 移除空白字符
                IFS=':' read -r group_id artifact_id version <<< "$lib"
                local lib_name="${artifact_id}-${version}.jar"
                local lib_path="$lib_dir/$lib_name"
                
                # 检查依赖库是否已存在
                if [ -f "$lib_path" ]; then
                    log_info "依赖库已存在: $lib_name"
                else
                    log_info "下载依赖库: $lib_name"
                    download_artifact "$lib_path" "$lib" "lib"
                fi
            done
        else
            log_info "连接器 $connector 没有配置依赖库"
        fi
    done
    
    log_info "插件和依赖安装完成"
}

# 下载构件（插件或库）
download_artifact() {
    local target_path=$1
    local artifact=$2
    local type=$3
    local retry_count=0
    local download_success=false
    
    while [ $retry_count -lt $retries ]; do
        log_info "下载尝试 $((retry_count + 1))/$retries"
        
        # 构建下载URL
        local download_url
        local repo_url="${PLUGIN_REPOS[$PLUGIN_REPO]:-${PLUGIN_REPOS[aliyun]}}"  # 使用aliyun作为默认值
        
        if [ "$type" = "plugin" ]; then
            # 标准仓库的插件URL格式
            download_url="$repo_url/org/apache/seatunnel/connector-${artifact}/${SEATUNNEL_VERSION}/connector-${artifact}-${SEATUNNEL_VERSION}.jar"
        else
            # 处理依赖库的URL
            IFS=':' read -r group_id artifact_id version <<< "$artifact"
            group_path=$(echo "$group_id" | tr '.' '/')
            download_url="$repo_url/$group_path/$artifact_id/$version/$artifact_id-$version.jar"
        fi
        
        log_info "从 $download_url 下载..."
        
        # 检查URL是否可访问
        if ! check_url "$download_url"; then
            log_warning "当前下载源不可用，尝试切换到备用源..."
            if [ "$PLUGIN_REPO" = "aliyun" ]; then
                PLUGIN_REPO="apache"
                continue
            elif [ "$PLUGIN_REPO" = "huaweicloud" ]; then
                PLUGIN_REPO="aliyun"
                continue
            fi
        fi
        
        # 使用curl下载
        if curl -L \
            --fail \
            --progress-bar \
            --connect-timeout 10 \
            --retry 3 \
            --retry-delay 2 \
            --retry-max-time 60 \
            -o "$target_path" \
            "$download_url" 2>&1; then
            
            if [ -f "$target_path" ]; then
                chmod 644 "$target_path"
                chown "$INSTALL_USER:$INSTALL_GROUP" "$target_path"
                log_info "下载成功: $(basename "$target_path")"
                download_success=true
                break
            fi
        fi
        
        retry_count=$((retry_count + 1))
        [ $retry_count -lt $retries ] && log_warning "下载失败，等待重试..." && sleep 3
    done
    
    if [ "$download_success" != "true" ]; then
        log_error "下载失败: $download_url，已重试 $retries 次"
    fi
}

# 修改create_service_file函数，添加必要的参数
create_service_file() {
    local service_name=$1
    local role=$2
    local java_home=${3:-"$JAVA_HOME"}
    local seatunnel_home=${4:-"$SEATUNNEL_HOME"}
    local install_user=${5:-"$INSTALL_USER"}
    local install_group=${6:-"$INSTALL_GROUP"}
    
    local service_file="/etc/systemd/system/${service_name}.service"
    local description="Apache SeaTunnel ${role} Service"
    local exec_args=""
    local hazelcast_config=""
    local jvm_options_file=""
    
    # ���置配置文件和参数
    case "$role" in
        "Hybrid")
            exec_args=""
            seatunnel_logs="seatunnel-engine-server"
            hazelcast_config="${seatunnel_home}/config/hazelcast.yaml"
            jvm_options_file="${seatunnel_home}/config/jvm_options"
            ;;
        "Master")
            exec_args="-r master"
            seatunnel_logs="seatunnel-engine-master"
            hazelcast_config="${seatunnel_home}/config/hazelcast-master.yaml"
            jvm_options_file="${seatunnel_home}/config/jvm_master_options"
            ;;
        "Worker")
            exec_args="-r worker"
            seatunnel_logs="seatunnel-engine-worker"
            hazelcast_config="${seatunnel_home}/config/hazelcast-worker.yaml"
            jvm_options_file="${seatunnel_home}/config/jvm_worker_options"
            ;;
    esac

    # 从JVM配置文件读取堆内存大小
    local heap_size
    heap_size=$(grep "^-Xmx" "$jvm_options_file" | sed 's/-Xmx\([0-9]\+\)g/\1/')
    
    # 读取所有JVM配置
    local jvm_opts=""
    while IFS= read -r line; do
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "${line// }" ]] && continue
        if [[ "$line" =~ ^-XX ]]; then
            if [[ "$line" =~ HeapDumpPath ]]; then
                jvm_opts+="${line/\/tmp\/seatunnel\/dump\/zeta-server/${seatunnel_home}\/dump\/seatunnel-zeta-server} "
            else
                jvm_opts+="$line "
            fi
        fi
    done < "$jvm_options_file"

    # 创建服务文件
    cat > "$service_file" << EOF
[Unit]
Description=${description}
After=network.target

[Service]
Type=simple
User=${install_user}
Group=${install_group}
Environment="JAVA_HOME=${java_home}"
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:\${JAVA_HOME}/bin"
Environment="SEATUNNEL_HOME=${seatunnel_home}"
WorkingDirectory=${seatunnel_home}
ExecStart=${java_home}/bin/java \\
    -Dlog4j2.contextSelector=org.apache.logging.log4j.core.async.AsyncLoggerContextSelector \\
    -Dhazelcast.logging.type=log4j2 \\
    -Dlog4j2.configurationFile=${seatunnel_home}/config/log4j2.properties \\
    -Dseatunnel.logs.path=${seatunnel_home}/logs \\
    -Dseatunnel.logs.file_name=${seatunnel_logs} \\
    -Xms${heap_size}g \\
    -Xmx${heap_size}g \\
    ${jvm_opts}\\
    -Dseatunnel.config=${seatunnel_home}/config/seatunnel.yaml \\
    -Dhazelcast.config=${hazelcast_config} \\
    -cp "${seatunnel_home}/lib/*:${seatunnel_home}/starter/seatunnel-starter.jar" \\
    org.apache.seatunnel.core.starter.seatunnel.SeaTunnelServer ${exec_args}
ExecStop=/bin/kill -s TERM \$MAINPID
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    # 设置权限
    chmod 644 "$service_file"
    
    # 重新加载systemd配置
    sudo systemctl daemon-reload
    
    # 启用服务
    sudo systemctl enable "$(basename "$service_file")"
}

# 修改setup_auto_start函数
setup_auto_start() {
    if [ "${ENABLE_AUTO_START}" != "true" ]; then
        log_info "跳过开机自启动配置"
        return 0
    fi
    
    log_info "配置SeaTunnel开机自启动..."
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：在所有节点上配置相同的服务
        for node in "${ALL_NODES[@]}"; do
            # 跳过localhost和当前节点
            if [ "$node" = "localhost" ] || [ "$node" = "$current_ip" ]; then
                log_info "在本地节点配置服务..."
                create_service_file "seatunnel" "Hybrid"
                continue
            fi
            
            log_info "在节点 $node 上配置服务..."
            ssh_with_retry "$node" "$(declare -f create_service_file); create_service_file 'seatunnel' 'Hybrid' '${JAVA_HOME}' '${SEATUNNEL_HOME}' '${INSTALL_USER}' '${INSTALL_GROUP}'"
        done
    else
        # 分离模式：根据节点角色配置服务
        # 配置Master节点
        for master in "${MASTER_IPS[@]}"; do
            if [ "$master" = "localhost" ] || [ "$master" = "$current_ip" ]; then
                log_info "在本地Master节点配置服务..."
                create_service_file "seatunnel-master" "Master"
                continue
            fi
            
            log_info "在Master节点 $master 上配置服务..."
            ssh_with_retry "$master" "$(declare -f create_service_file); create_service_file 'seatunnel-master' 'Master' '${JAVA_HOME}' '${SEATUNNEL_HOME}' '${INSTALL_USER}' '${INSTALL_GROUP}'"
        done
        
        # 配置Worker节点
        for worker in "${WORKER_IPS[@]}"; do
            if [ "$worker" = "localhost" ] || [ "$worker" = "$current_ip" ]; then
                log_info "在本地Worker节点配置服务..."
                create_service_file "seatunnel-worker" "Worker"
                continue
            fi
            
            log_info "在Worker节点 $worker 上配置服务..."
            ssh_with_retry "$worker" "$(declare -f create_service_file); create_service_file 'seatunnel-worker' 'Worker' '${JAVA_HOME}' '${SEATUNNEL_HOME}' '${INSTALL_USER}' '${INSTALL_GROUP}'"
        done
    fi
}

# 检查系统内存
check_memory() {
    log_info "检查系统内存..."
    
    # 获取系统总内存(GB)
    local total_mem
    if [ -f /proc/meminfo ]; then
        total_mem=$(awk '/MemTotal/ {print int($2/1024/1024)}' /proc/meminfo)
    else
        # 对于不支持 /proc/meminfo 的系统，尝试使用其他命令
        if command -v free >/dev/null 2>&1; then
            total_mem=$(free -g | awk '/Mem:/ {print int($2)}')
        else
            log_error "无法获取系统内存信息"
        fi
    fi
    
    # 获取可用内存(GB)
    local available_mem
    if [ -f /proc/meminfo ]; then
        available_mem=$(awk '/MemAvailable/ {print int($2/1024/1024)}' /proc/meminfo)
    else
        if command -v free >/dev/null 2>&1; then
            available_mem=$(free -g | awk '/Mem:/ {print int($4)}')
        else
            log_error "无法获取系统可用内存信息"
        fi
    fi
    
    # 根据部署模式检查内存需求
    local required_mem=0
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        required_mem=$((HYBRID_HEAP_SIZE + 2)) # 额外预留2GB系统使用
        log_info "混合模式下需要 ${HYBRID_HEAP_SIZE}GB 堆内存 + 2GB 系统预留"
    else
        local is_master=false
        local is_worker=false
        
        # 检查当前机器是否为Master节点
        if [[ " ${MASTER_IPS[*]} " =~ " $current_ip " ]]; then
            is_master=true
            required_mem=$((MASTER_HEAP_SIZE))
            log_info "当前节点是Master节点，需要 ${MASTER_HEAP_SIZE}GB 堆内存"
        fi
        
        # 检查当前机器是否为Worker节点
        if [[ " ${WORKER_IPS[*]} " =~ " $current_ip " ]]; then
            is_worker=true
            if [ "$is_master" = true ]; then
                # 如果同时是Master和Worker，需要两者内存之和
                required_mem=$((required_mem + WORKER_HEAP_SIZE))
                log_info "当前节点同时是Worker节点，额外需要 ${WORKER_HEAP_SIZE}GB 堆内存"
                log_info "总共需 ${required_mem}GB 堆内存 + 2GB 系统预留"
            else
                required_mem=$((WORKER_HEAP_SIZE))
                log_info "当前节点是Worker节点，需要 ${WORKER_HEAP_SIZE}GB 堆内存"
            fi
        fi
        
        # 添加系统预留
        required_mem=$((required_mem + 2))
    fi
    
    log_info "系统总内存: ${total_mem}GB"
    log_info "系统可用内存: ${available_mem}GB"
    log_info "最小所需内存: ${required_mem}GB"
    
    # 检查总内存是否足够
    if [ $total_mem -lt $required_mem ]; then
        log_error "系统总内存不足！需要至少 ${required_mem}GB，当前只有 ${total_mem}GB"
    fi
    
    # 检查可用内存是否足够
    if [ $available_mem -lt $required_mem ]; then
        log_warning "系统可用内存不足！需要至少 ${required_mem}GB，当前只有 ${available_mem}GB"
        log_warning "建议释放一些内存后再继续安装"
        read -r -p "是否继续安装? [y/N] " response
        case "$response" in
            [yY][eE][sS]|[yY]) 
                log_warning "继续安装，但可能会影响系统性能"
                ;;
            *)
                log_error "安装已取消"
                ;;
        esac
    fi
}

# 处理安装包
handle_package() {
    if [ "$INSTALL_MODE" = "online" ]; then
        # 在线安装模式
        local package_name="apache-seatunnel-${SEATUNNEL_VERSION}-bin.tar.gz"
        local package_path
        
        # 使用全局下载目录检查文件
        if [ -f "$DOWNLOAD_DIR/$package_name" ]; then
            log_warning "发现本地已存在安装包: $DOWNLOAD_DIR/$package_name"
            read -r -p "是否重新下载? [y/N] " response
            case "$response" in
                [yY][eE][sS]|[yY])
                    rm -f "$DOWNLOAD_DIR/$package_name"
                    ;;
                *)
                    log_info "使用已存在的安装包"
                    package_path="$DOWNLOAD_DIR/$package_name"
                    ;;
            esac
        fi
        
        # 下载安装包
        if [ -z "$package_path" ]; then
            if download_package "$package_name" "$SEATUNNEL_VERSION"; then
                package_path=$(cat /tmp/download_path.tmp)
                rm -f /tmp/download_path.tmp
                
                if [ ! -f "$package_path" ]; then
                    log_error "下载的安装包不存在: $package_path"
                fi
            else
                log_error "安装包下载失败"
            fi
        fi
        
        PACKAGE_PATH="$package_path"
    fi
    
    # 验证安装包
    verify_package "$PACKAGE_PATH"
    
    # 创建安装目录
    log_info "创建安装目录..."
    create_directory "$BASE_DIR"
    setup_permissions "$BASE_DIR"
    
    # 解压安装包
    log_info "解压安装包..."
    sudo tar -zxf "$PACKAGE_PATH" -C "$BASE_DIR"
    setup_permissions "$SEATUNNEL_HOME"
}

# 分发到其他节点
distribute_to_nodes() {
    log_info "分发到其他节点..."
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：分发到所有节点
        for node in "${ALL_NODES[@]}"; do
            # 跳过localhost和当前节点
            if [ "$node" = "localhost" ] || [ "$node" = "$(hostname -I | awk '{print $1}')" ]; then
                log_info "跳过本地节点: $node"
                continue
            fi
            
            log_info "分发到 $node..."
            ssh_with_retry "$node" "mkdir -p $BASE_DIR && chown $INSTALL_USER:$INSTALL_GROUP $BASE_DIR"
            scp_with_retry "$SEATUNNEL_HOME" "$node" "$BASE_DIR/"
        done
    else
        # 分离模式：分发到master和worker节点
        # 先分发到其他master节点
        for master in "${MASTER_IPS[@]}"; do
            if [ "$master" = "localhost" ] || [ "$master" = "$(hostname -I | awk '{print $1}')" ]; then
                log_info "跳过本地Master节点: $master"
                continue
            fi
            
            log_info "分发到Master节点 $master..."
            ssh_with_retry "$master" "mkdir -p $BASE_DIR && chown $INSTALL_USER:$INSTALL_GROUP $BASE_DIR"
            scp_with_retry "$SEATUNNEL_HOME" "$master" "$BASE_DIR/"
        done
        
        # 再分发到worker节点
        for worker in "${WORKER_IPS[@]}"; do
            if [ "$worker" = "localhost" ] || [ "$worker" = "$(hostname -I | awk '{print $1}')" ]; then
                log_info "跳过本地Worker节点: $worker"
                continue
            fi
            
            log_info "分发到Worker节点 $worker..."
            ssh_with_retry "$worker" "mkdir -p $BASE_DIR && chown $INSTALL_USER:$INSTALL_GROUP $BASE_DIR"
            scp_with_retry "$SEATUNNEL_HOME" "$worker" "$BASE_DIR/"
        done
    fi
}

# 配置环境变量
setup_environment() {
    log_info "配置环境变量..."
    BASHRC_CONTENT="
# SEATUNNEL_HOME BEGIN
export SEATUNNEL_HOME=$SEATUNNEL_HOME
export PATH=\$PATH:\$SEATUNNEL_HOME/bin
# SEATUNNEL_HOME END"

    # 获取用户home目录
    USER_HOME=""
    if command -v getent >/dev/null 2>&1; then
        USER_HOME=$(getent passwd "$INSTALL_USER" | cut -d: -f6)
    else
        USER_HOME=$(eval echo ~"$INSTALL_USER")
    fi
    
    if [ -z "$USER_HOME" ]; then
        log_error "无法获取用户 $INSTALL_USER 的home目录"
    fi
    
    # 配置本地环境变量
    if grep -q "SEATUNNEL_HOME" "$USER_HOME/.bashrc"; then
        log_info "本地环境变量已存在，更新配置..."
        sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' "$USER_HOME/.bashrc"
    fi
    echo "$BASHRC_CONTENT" >> "$USER_HOME/.bashrc"
    
    # 远程节点环境变量
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：配置所有远程节点
        for node in "${ALL_NODES[@]}"; do
            # 跳过localhost和当前节点
            if [ "$node" = "localhost" ] || [ "$node" = "$(hostname -I | awk '{print $1}')" ]; then
                log_info "跳过本地节点环境变量配置: $node"
                continue
            fi
            
            log_info "配置节点 $node 的环境变量..."
            remote_home=$(ssh_with_retry "$node" "echo ~$INSTALL_USER")
            ssh_with_retry "$node" "
                if grep -q 'SEATUNNEL_HOME' '$remote_home/.bashrc'; then
                    sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' '$remote_home/.bashrc'
                fi
                echo '$BASHRC_CONTENT' >> '$remote_home/.bashrc'
            "
        done
    else
        # 分离模式：配置master和worker节点
        # 配置其他master节点
        for master in "${MASTER_IPS[@]}"; do
            if [ "$master" = "localhost" ] || [ "$master" = "$(hostname -I | awk '{print $1}')" ]; then
                log_info "跳过本地Master节点环境变量配置: $master"
                continue
            fi
            
            log_info "配置Master节点 $master 的环境变量..."
            remote_home=$(ssh_with_retry "$master" "echo ~$INSTALL_USER")
            ssh_with_retry "$master" "
                if grep -q 'SEATUNNEL_HOME' '$remote_home/.bashrc'; then
                    sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' '$remote_home/.bashrc'
                fi
                echo '$BASHRC_CONTENT' >> '$remote_home/.bashrc'
            "
        done
        
        # 配置worker节点
        for worker in "${WORKER_IPS[@]}"; do
            if [ "$worker" = "localhost" ] || [ "$worker" = "$(hostname -I | awk '{print $1}')" ]; then
                log_info "跳过本地Worker节点环境变量配置: $worker"
                continue
            fi
            
            log_info "配置Worker节点 $worker 的环境变量..."
            remote_home=$(ssh_with_retry "$worker" "echo ~$INSTALL_USER")
            ssh_with_retry "$worker" "
                if grep -q 'SEATUNNEL_HOME' '$remote_home/.bashrc'; then
                    sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' '$remote_home/.bashrc'
                fi
                echo '$BASHRC_CONTENT' >> '$remote_home/.bashrc'
            "
        done
    fi
}

function show_completion_info(){
    # 计算安装时长
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    MINUTES=$((DURATION / 60))
    SECONDS=$((DURATION % 60))
    
    # 添加安装完成后的验证提示
    echo -e "\n${GREEN}SeaTunnel安装完成!${NC}"
    echo -e "安装总耗时: ${GREEN}${MINUTES}分${SECONDS}秒${NC}"
    echo -e "\n${YELLOW}验证和使用说明:${NC}"
    echo "1. 刷新环境变量:"
    echo -e "${GREEN}source $USER_HOME/.bashrc${NC}"
    
    echo -e "\n2. 集群管理命令:"
    if [ "${ENABLE_AUTO_START}" = "true" ]; then
        if [ "$DEPLOY_MODE" = "hybrid" ]; then
            echo -e "启动服务:    ${GREEN}sudo systemctl start seatunnel${NC}"
            echo -e "停止服务:    ${GREEN}sudo systemctl stop seatunnel${NC}"
            echo -e "重启服务:    ${GREEN}sudo systemctl restart seatunnel${NC}"
            echo -e "查看状态:    ${GREEN}sudo systemctl status seatunnel${NC}"
            echo -e "查看启动日志:    ${GREEN}sudo journalctl -u seatunnel -n 100 --no-pager${NC}"
            echo -e "查看运行日志:    ${GREEN}tail -n 100 $SEATUNNEL_HOME/logs/seatunnel-engine-server.out${NC}"
        else
            echo -e "Master服务命令:"
            echo -e "启动服务:    ${GREEN}sudo systemctl start seatunnel-master${NC}"
            echo -e "停止服务:    ${GREEN}sudo systemctl stop seatunnel-master${NC}"
            echo -e "重启服务:    ${GREEN}sudo systemctl restart seatunnel-master${NC}"
            echo -e "查看状态:    ${GREEN}sudo systemctl status seatunnel-master${NC}"
            echo -e "查看启动日志:    ${GREEN}sudo journalctl -u seatunnel-master -n 100 --no-pager${NC}"
            echo -e "查看运行日志:    ${GREEN}tail -n 100 $SEATUNNEL_HOME/logs/seatunnel-engine-master.out${NC}"
            echo -e "----------------------------------------"
            echo -e "\nWorker服务命令:"
            echo -e "启动服务:    ${GREEN}sudo systemctl start seatunnel-worker${NC}"
            echo -e "停止服务:    ${GREEN}sudo systemctl stop seatunnel-worker${NC}"
            echo -e "重启服务:    ${GREEN}sudo systemctl restart seatunnel-worker${NC}"
            echo -e "查看状态:    ${GREEN}sudo systemctl status seatunnel-worker${NC}"
            echo -e "查看启动日志:    ${GREEN}sudo journalctl -u seatunnel-worker -n 100 --no-pager${NC}"
            echo -e "查看运行日志:    ${GREEN}tail -n 100 $SEATUNNEL_HOME/logs/seatunnel-engine-worker.out${NC}"
        fi
    else
        echo -e "启动集群:    ${GREEN}$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh start${NC}"
        echo -e "停止集群:    ${GREEN}$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh stop${NC}"
        echo -e "重启集群:    ${GREEN}$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh restart${NC}"
        echo -e "查看运行日志:    ${GREEN}tail -n 100 $SEATUNNEL_HOME/logs/seatunnel-engine-server.out${NC}"
    fi
    
    echo -e "\n3. 验证安装:"
    echo -e "运行示例任务: ${GREEN}$SEATUNNEL_HOME/bin/seatunnel.sh --config config/v2.batch.config.template${NC}"
    
    echo -e "\n${YELLOW}部署信息:${NC}"
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        echo "部署模式: 混合模式"
        echo -e "集群节点: ${GREEN}${CLUSTER_NODES[*]}${NC}"
    else
        echo "部署模式: 分离模式"
        echo -e "Master节点: ${GREEN}${MASTER_IPS[*]}${NC}"
        echo -e "Worker节点: ${GREEN}${WORKER_IPS[*]}${NC}"
    fi
    
    echo -e "\n${YELLOW}注意事项:${NC}"
    echo "1. 首次启动集群前，请确保所有节点的环境变量已经生效,source $USER_HOME/.bashrc"
    echo -e "2. 更多使用说明请参考：${GREEN}https://seatunnel.apache.org/docs${NC}"
}


# 在颜色定义后添加
# 错误追踪函数
trace_error() {
    local err=$?
    local line_no=$1
    local bash_command=$2
    
    echo -e "\n${RED}[ERROR TRACE]${NC} $(date '+%Y-%m-%d %H:%M:%S')"
    echo -e "${RED}错误码:${NC} $err"
    echo -e "${RED}错误行号:${NC} $line_no"
    echo -e "${RED}错误命令:${NC} $bash_command"
    
    # 输出函数调用栈
    local frame=0
    echo -e "${RED}函数调用栈:${NC}"
    while caller $frame; do
        ((frame++))
    done | awk '{printf "  %s(): 第%s行 in %s\n", $2, $1, $3}'
    
    # 如果是在函数中发生错误,显示函数名
    if [[ "${FUNCNAME[*]}" ]]; then
        echo -e "${RED}当前函数:${NC} ${FUNCNAME[1]}"
    fi
    
    # 显示最后几行日志
    echo -e "${RED}最后10行日志:${NC}"
    tail -n 10 "$LOG_DIR/install.log" 2>/dev/null
    
    handle_error $line_no  
    
    # 直接退出脚本
    exit $err
}

# 设置错误追踪
set -E           # 继承ERR trap
set -o pipefail  # 管道中的错误也会被捕获
trap 'trace_error ${LINENO} "$BASH_COMMAND"' ERR

# 创建日志目录
mkdir -p "$LOG_DIR"
chmod 755 "$LOG_DIR"
[ -n "$INSTALL_USER" ] && [ -n "$INSTALL_GROUP" ] && chown "$INSTALL_USER:$INSTALL_GROUP" "$LOG_DIR"

# 创建日志文件
touch "$LOG_FILE"
chmod 644 "$LOG_FILE"

# 如果INSTALL_USER已定义,设置目录和文件所有权
if [ -n "$INSTALL_USER" ] && [ -n "$INSTALL_GROUP" ]; then
    chown "$INSTALL_USER:$INSTALL_GROUP" "$LOG_DIR"
    chown "$INSTALL_USER:$INSTALL_GROUP" "$LOG_FILE"
fi

# 重定向输出到日志文件
exec 1> >(tee -a "$LOG_FILE")
exec 2>&1

# 记录脚本开始执行
log_info "开始执行安装脚本..."
log_info "脚本路径: $0"
log_info "脚本参数: $*"

# 添加处理SELinux的函数
handle_selinux() {
    log_info "检查SELinux状态..."
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：检查所有节点
        for node in "${ALL_NODES[@]}"; do
            if [ "$node" != "localhost" ] && [ "$node" != "$(hostname -I | awk '{print $1}')" ]; then
                check_and_handle_selinux "$node" true
            else 
                check_and_handle_selinux "localhost" false
            fi
        done
    else
        # 分离模式：检查master和worker节点
        for master in "${MASTER_IPS[@]}"; do
            if [ "$master" != "localhost" ] && [ "$master" != "$(hostname -I | awk '{print $1}')" ]; then
                check_and_handle_selinux "$master" true
            else 
                check_and_handle_selinux "localhost" false
            fi
        done
        for worker in "${WORKER_IPS[@]}"; do
            if [ "$worker" != "localhost" ] && [ "$worker" != "$(hostname -I | awk '{print $1}')" ]; then
                check_and_handle_selinux "$worker" true
            else 
                check_and_handle_selinux "localhost" false
            fi
        done
    fi
}

# 将check_and_handle_selinux提取为独立函数
check_and_handle_selinux() {
    local node=$1
    local is_remote=$2
    
    local check_script='
        if ! command -v sestatus >/dev/null 2>&1; then
            echo "NO_SELINUX"
            exit 0
        fi
    
        # 检查SELinux状态
        if sestatus | grep  "disabled"; then
            echo "SELINUX_DISABLED"
            exit 0
        fi
        
        # 检查restorecon和chcon命令
        if ! command -v restorecon >/dev/null 2>&1 || ! command -v chcon >/dev/null 2>&1; then
            echo "NO_SELINUX_TOOLS"
            exit 0
        fi
        
        # 使用restorecon恢复/data/seatunnel目录的安全上下文
        sudo restorecon -Rv /data/seatunnel >/dev/null 2>&1
        
        # 设置Java二进制文件的SELinux上下文
        if [ -d /data/seatunnel/java ]; then
            sudo chcon -R -t bin_t /data/seatunnel/java/*/bin/java >/dev/null 2>&1
            if [ $? -eq 0 ]; then
                echo "SELINUX_CONTEXT_RESTORED"
            else
                echo "JAVA_CONTEXT_FAILED"
            fi
        else
            echo "SELINUX_CONTEXT_RESTORED"
        fi'
    
    if [ "$is_remote" = true ]; then
        # 远程节点检查
        local result
        result=$(execute_remote_script "$node" "$check_script")
        
        case "$result" in
            "NO_SELINUX")
                log_info "节点 $node 未安装SELinux工具，跳过配置"
                ;;
            "SELINUX_DISABLED")
                log_info "节点 $node SELinux已禁用，跳过配置"
                ;;
            "NO_SELINUX_TOOLS")
                log_warning "节点 $node 未安装restorecon或chcon工具，跳过SELinux配置"
                ;;
            "SELINUX_CONTEXT_RESTORED")
                log_success "节点 $node 的SELinux安全上下文已恢复"
                ;;
            "JAVA_CONTEXT_FAILED")
                log_warning "节点 $node 的Java二进制文件SELinux上下文设置失败"
                ;;
            *)
                log_warning "节点 $node 的SELinux配置未完成: $result"
                ;;
        esac
    else
        # 本地节点检查
        if ! command -v sestatus >/dev/null 2>&1; then
            log_info "本地节点未安装SELinux工具，跳过配置"
            return 0
        fi
        
        if sestatus | grep "disabled"; then
            log_info "本地节点SELinux已禁用，跳过配置"
            return 0
        fi
        
        if ! command -v restorecon >/dev/null 2>&1 || ! command -v chcon >/dev/null 2>&1; then
            log_warning "本地节点未安装restorecon或chcon工具，跳过SELinux配置"
            return 0
        fi
        
        # 恢复本地安全上下文
        sudo restorecon -Rv /data/seatunnel >/dev/null 2>&1
        
        # 设置Java二进制文件的SELinux上下文
        if [ -d /data/seatunnel/java ]; then
            if sudo chcon -R -t bin_t /data/seatunnel/java/*/bin/java >/dev/null 2>&1; then
                log_success "本地节点的SELinux安全上下文已恢复，Java二进制文件上下文已设置"
            else
                log_warning "本地节点的Java二进制文件SELinux上下文设置失败"
            fi
        else
            log_success "本地节点的SELinux安全上下文已恢复"
        fi
    fi
}

# 添加execute_remote_script函数
execute_remote_script() {
    local node=$1
    local script=$2
    shift 2  # 移除前两个参数，剩余的都是要传递给脚本的参数
    local TIMEOUT=10
    
    # 创建临时脚本文件，添加参数处理
    local temp_script="/tmp/remote_script_$RANDOM.sh"
    {
        echo '#!/bin/bash'
        # 将参数数组重建为脚本变量
        echo 'SCRIPT_ARGS=('
        printf "'%s' " "$@"
        echo ')'
        echo "$script"
    } > "$temp_script"
    chmod +x "$temp_script"
    
    # 复制脚本到远程节点并执行，添加超时控制
    if ! timeout $TIMEOUT scp -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$temp_script" "${INSTALL_USER}@${node}:/tmp/" >/dev/null 2>&1; then
        log_error "向节点 $node 传输脚本失败"
        rm -f "$temp_script"
        return 1
    fi
    
    local result
    result=$(timeout $TIMEOUT ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "bash /tmp/$(basename "$temp_script")")
    local exit_code=$?
    
    # 清理临时文件
    rm -f "$temp_script"
    timeout $TIMEOUT ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "${INSTALL_USER}@${node}" "rm -f /tmp/$(basename "$temp_script")" >/dev/null 2>&1
    
    # 返回结果
    echo "$result"
    return $exit_code
}

# 检查防火墙状态
check_firewall() {
    log_info "检查防火墙状态..."
    
    check_node_firewall() {
        local node=$1
        local is_remote=$2
        local ports=()
        local firewall_status_shown=false
        local ports_shown=false
        
        # 根据部署模式确定需要检查的端口
        if [ "$DEPLOY_MODE" = "hybrid" ]; then
            ports+=("${HYBRID_PORT:-5801}")
        else
            if [[ " ${MASTER_IPS[*]} " =~ " $node " ]]; then
                ports+=("${MASTER_PORT:-5801}")
            fi
            if [[ " ${WORKER_IPS[*]} " =~ " $node " ]]; then
                ports+=("${WORKER_PORT:-5802}")
            fi
        fi
        
        # 添加SSH端口
        ports+=("$SSH_PORT")
        
        # 优化端口显示,只在第一次显示
        if [ ! "$ports_shown" ]; then
            log_info "检查节点 $node 的端口: ${ports[*]}/tcp"
            ports_shown=true
        fi

        local check_script='
            # 获取传入的端口参数
            ports=("${SCRIPT_ARGS[@]}")
            
            # 检查防火墙状态和端口
            if command -v systemctl >/dev/null 2>&1; then
                if systemctl is-active --quiet firewalld; then
                    # 检查端口是否开放
                    ports_status=""
                    for port in "${ports[@]}"; do
                        if firewall-cmd --list-ports | grep -w "$port/tcp" >/dev/null 2>&1; then
                            ports_status="${ports_status}${port}:open,"
                        else
                            ports_status="${ports_status}${port}:closed,"
                        fi
                    done
                    echo "FIREWALLD_ACTIVE:$ports_status"
                elif systemctl is-active --quiet ufw; then
                    # 检查UFW端口状态
                    ports_status=""
                    for port in "${ports[@]}"; do
                        if ufw status | grep -w "$port/tcp.*ALLOW" >/dev/null 2>&1; then
                            ports_status="${ports_status}${port}:open,"
                        else
                            ports_status="${ports_status}${port}:closed,"
                        fi
                    done
                    echo "UFW_ACTIVE:$ports_status"
                else
                    echo "FIREWALL_INACTIVE"
                fi
            elif command -v service >/dev/null 2>&1; then
                if service iptables status >/dev/null 2>&1; then
                    # 检查iptables端口状态
                    ports_status=""
                    for port in "${ports[@]}"; do
                        if iptables -L -n | grep -w "tcp dpt:$port.*ACCEPT" >/dev/null 2>&1; then
                            ports_status="${ports_status}${port}:open,"
                        else
                            ports_status="${ports_status}${port}:closed,"
                        fi
                    done
                    echo "IPTABLES_ACTIVE:$ports_status"
                else
                    echo "FIREWALL_INACTIVE"
                fi
            else
                echo "NO_FIREWALL"
            fi'
        
        if [ "$is_remote" = true ]; then
            # 远程节点检查，传递端口参数
            local result
            result=$(execute_remote_script "$node" "$check_script" "${ports[@]}")
            
            case "${result%%:*}" in
                "FIREWALLD_ACTIVE"|"UFW_ACTIVE"|"IPTABLES_ACTIVE")
                    local firewall_type="${result%%:*}"
                    local ports_status="${result#*:}"
                    local closed_ports=()
                    
                    # 解析端口状态
                    IFS=',' read -r -a port_array <<< "$ports_status"
                    for port_status in "${port_array[@]}"; do
                        if [[ "$port_status" =~ ([0-9]+):closed ]]; then
                            closed_ports+=("${BASH_REMATCH[1]}")
                        fi
                    done
                    
                    if [ ${#closed_ports[@]} -gt 0 ]; then
                        log_warning "节点 $node 的防火墙已启用，以下端口未开放: ${closed_ports[*]}"
                        echo -e "\n${YELLOW}请在节点 $node 上执行以下操作之一：${NC}"
                        echo "1. 关闭防火墙（推荐用于测试环境）："
                        case "$firewall_type" in
                            "FIREWALLD_ACTIVE")
                                echo -e "${GREEN}sudo systemctl stop firewalld && sudo systemctl disable firewalld${NC}"
                                echo "2. 开放必要端口（推荐用于生产环境）："
                                for port in "${closed_ports[@]}"; do
                                    echo -e "${GREEN}sudo firewall-cmd --permanent --add-port=$port/tcp${NC}"
                                done
                                echo -e "${GREEN}sudo firewall-cmd --reload${NC}"
                                ;;
                            "UFW_ACTIVE")
                                echo -e "${GREEN}sudo ufw disable${NC}"
                                echo "2. 开放必要端口（推荐用于生产环境）："
                                for port in "${closed_ports[@]}"; do
                                    echo -e "${GREEN}sudo ufw allow $port/tcp${NC}"
                                done
                                ;;
                            "IPTABLES_ACTIVE")
                                echo -e "${GREEN}sudo service iptables stop && sudo chkconfig iptables off${NC}"
                                echo "2. 开放必要端口（推荐用于生产环境）："
                                for port in "${closed_ports[@]}"; do
                                    echo -e "${GREEN}sudo iptables -A INPUT -p tcp --dport $port -j ACCEPT${NC}"
                                done
                                echo -e "${GREEN}sudo service iptables save${NC}"
                                ;;
                        esac
                        return 1
                    else
                        log_info "节点 $node 的防火墙已启用，所需端口已开放"
                    fi
                    ;;
                "FIREWALL_INACTIVE"|"NO_FIREWALL")
                    log_info "节点 $node 的防火墙未启用"
                    ;;
                *)
                    log_warning "节点 $node 的防火墙状态检查失败: $result"
                    return 1
                    ;;
            esac
        else
            # 本地节点检查
            if command -v systemctl >/dev/null 2>&1; then
                if systemctl is-active --quiet firewalld; then
                    if [ ! "$firewall_status_shown" ]; then
                        log_info "检测到防火墙(firewalld)已启用"
                        firewall_status_shown=true
                    fi
                    local closed_ports=()
                    for port in "${ports[@]}"; do
                        if ! firewall-cmd --list-ports | grep -w "$port/tcp" >/dev/null 2>&1; then
                            closed_ports+=("$port")
                        fi
                    done
                    
                    if [ ${#closed_ports[@]} -gt 0 ]; then
                        log_warning "本地节点防火墙(firewalld)已启用，以下端口未开放: ${closed_ports[*]}"
                        echo -e "\n${YELLOW}请执行以下操作之一：${NC}"
                        echo "1. 关闭防火墙（推荐用于测试环境）："
                        echo -e "${GREEN}sudo systemctl stop firewalld && sudo systemctl disable firewalld${NC}"
                        echo "2. 开放必要端口（推荐用于生产环境）："
                        for port in "${closed_ports[@]}"; do
                            echo -e "${GREEN}sudo firewall-cmd --permanent --add-port=$port/tcp${NC}"
                        done
                        echo -e "${GREEN}sudo firewall-cmd --reload${NC}"
                        return 1
                    else
                        log_info "本地节点防火墙已启用，所需端口已开放"
                    fi
                elif systemctl is-active --quiet ufw; then
                    if [ ! "$firewall_status_shown" ]; then
                        log_info "检测到防火墙(ufw)已启用"
                        firewall_status_shown=true
                    fi
                    local closed_ports=()
                    for port in "${ports[@]}"; do
                        if ! ufw status | grep -w "$port/tcp.*ALLOW"; then
                            closed_ports+=("$port")
                        fi
                    done
                    
                    if [ ${#closed_ports[@]} -gt 0 ]; then
                        log_warning "本地节点防火墙(ufw)已启用，以下端口未开放: ${closed_ports[*]}"
                        echo -e "\n${YELLOW}请执行以下操作之一：${NC}"
                        echo "1. 关闭防火墙（推荐用于测试环境）："
                        echo -e "${GREEN}sudo ufw disable${NC}"
                        echo "2. 开放必要端口（推荐用于生产环境）："
                        for port in "${closed_ports[@]}"; do
                            echo -e "${GREEN}sudo ufw allow $port/tcp${NC}"
                        done
                        return 1
                    else
                        log_info "本地节点防火墙已启用，所需端口已开放"
                    fi
                fi
            elif command -v service >/dev/null 2>&1; then
                if service iptables status >/dev/null 2>&1; then
                    local closed_ports=()
                    for port in "${ports[@]}"; do
                        if ! iptables -L -n | grep -q "tcp dpt:$port.*ACCEPT"; then
                            closed_ports+=("$port")
                        fi
                    done
                    
                    if [ ${#closed_ports[@]} -gt 0 ]; then
                        log_warning "本地节点防火墙(iptables)已启用，以下端口未开放: ${closed_ports[*]}"
                        echo -e "\n${YELLOW}请执行以下操作之一：${NC}"
                        echo "1. 关闭防火墙（推荐用于测试环境）："
                        echo -e "${GREEN}sudo service iptables stop && sudo chkconfig iptables off${NC}"
                        echo "2. 开放必要端口（推荐用于生产环境）："
                        for port in "${closed_ports[@]}"; do
                            echo -e "${GREEN}sudo iptables -A INPUT -p tcp --dport $port -j ACCEPT${NC}"
                        done
                        echo -e "${GREEN}sudo service iptables save${NC}"
                        return 1
                    else
                        log_info "本地节点防火墙已启用，所需端口已开放"
                    fi
                fi
            fi
            # 移除重复的日志
            if [ ! "$firewall_status_shown" ]; then
                log_info "本地节点防火墙未启用或所需端口已开放"
            fi
        fi
        return 0
    }
    
    local firewall_check_failed=false
    
    # 检查所有节点的防火墙状态
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        for node in "${ALL_NODES[@]}"; do
            if [ "$node" != "localhost" ] && [ "$node" != "$(hostname -I | awk '{print $1}')" ]; then
                if ! check_node_firewall "$node" true; then
                    firewall_check_failed=true
                fi
            else
                if ! check_node_firewall "$node" false; then
                    firewall_check_failed=true
                fi
            fi
        done
    else
        for master in "${MASTER_IPS[@]}"; do
            if [ "$master" != "localhost" ] && [ "$master" != "$(hostname -I | awk '{print $1}')" ]; then
                if ! check_node_firewall "$master" true; then
                    firewall_check_failed=true
                fi
            else
                if ! check_node_firewall "$master" false; then
                    firewall_check_failed=true
                fi
            fi
        done
        for worker in "${WORKER_IPS[@]}"; do
            if [ "$worker" != "localhost" ] && [ "$worker" != "$(hostname -I | awk '{print $1}')" ]; then
                if ! check_node_firewall "$worker" true; then
                    firewall_check_failed=true
                fi
            else
                if ! check_node_firewall "$worker" false; then
                    firewall_check_failed=true
                fi
            fi
        done
    fi
    
    if [ "$firewall_check_failed" = true ]; then
        echo -e "\n${RED}检测到防火墙问题，请按上述提示处理后再次运行安装脚本${NC}"
        exit 1
    fi
}

# 主函数
main() {
    # 读取配置
    log_info "开始读取配置文件..."
    read_config
    
    # 检查防火墙状态
    check_firewall
    
    # 检查Java环境
    log_info "开始检查所有节点的Java环境..."
    
    # 检查节点java环境
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：检查所有节点
        for node in "${ALL_NODES[@]}"; do
            if [ "$node" != "localhost" ] && [ "$node" != "$(hostname -I | awk '{print $1}')" ]; then
                check_java "$node" "true"
            else 
                check_java "localhost" "false"
            fi
        done
    else
        # 分离模式：检查master和worker节点
        for master in "${MASTER_IPS[@]}"; do
            if [ "$master" != "localhost" ] && [ "$master" != "$(hostname -I | awk '{print $1}')" ]; then
                check_java "$master" "true"
            else 
                check_java "localhost" "false"
            fi
        done
        for worker in "${WORKER_IPS[@]}"; do
            if [ "$worker" != "localhost" ] && [ "$worker" != "$(hostname -I | awk '{print $1}')" ]; then
                check_java "$worker" "true"
            else 
                check_java "localhost" "false"
            fi
        done
    fi
    
    # 检查系统依赖
    log_info "开始检查系统依赖..."
    check_dependencies
    
    # 处理SELinux
    log_info "开始处理SELinux..."
    handle_selinux
    
    # 初始化临时文件列表
    log_info "初始化临时文件列表..."
    declare -a TEMP_FILES=()
    
    # 设置退出时清理
    log_info "设置退出清理..."
    trap cleanup_temp_files EXIT INT TERM
    
    # 检查用户配置
    log_info "检查用户配置..."
    check_user
    
    # 检查系统内存
    log_info "检查系统内存..."
    check_memory
    
    # 检查端口占用
    log_info "检查端口占用..."
    check_ports
    
    # 处理安装包
    handle_package
    
    # 配置文件修改
    log_info "修改配置文件..."
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        setup_hybrid_mode
    else
        setup_separated_mode
    fi
    
    # 配置检查点存储
    log_info "配置检查点存储..."
    configure_checkpoint
    
    # 安装插件和依赖
    log_info "安装插件和依赖..."
    install_plugins_and_libs
    
    # 添加集群管理脚本
    setup_cluster_scripts
    sed -i "s/root/$INSTALL_USER/g" "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh"
    
    # 分发到其他节点
    distribute_to_nodes
    
    # 配置环境变量
    setup_environment
    
    # 配置开机自启动
    setup_auto_start
    
    # 启动集群
    start_cluster
    
    # 检查服务状态
    check_services
    
    # 显示安装完成信息
    show_completion_info
}

# 执行主函数
main
