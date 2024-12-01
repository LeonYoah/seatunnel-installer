#!/bin/bash

# 确保遇到错误时立即退出
set -e

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

# 最大重试次数
MAX_RETRIES=3
# SSH超时时间(秒)
SSH_TIMEOUT=10

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
    
    # 如果安装失败，清理部分安装的文件
    if [ $exit_code -ne 0 ]; then
        log_warning "检测到安装失败，清理部分安装的文件..."
        if [ -d "$SEATUNNEL_HOME" ]; then
            log_info "除安装目录: $SEATUNNEL_HOME"
            sudo rm -rf "$SEATUNNEL_HOME"
        fi
        
        # 清理远程节点上的文件
        if [ -n "${WORKER_IPS[*]}" ]; then
            for node in "${WORKER_IPS[@]}"; do
                log_info "清理节点 $node 上的文件..."
                ssh_with_retry "$node" "sudo rm -rf $SEATUNNEL_HOME"
            done
        fi
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
    
    while [ $retries -lt $MAX_RETRIES ]; do
        if timeout $SSH_TIMEOUT scp -P $SSH_PORT "$src" "${INSTALL_USER}@${host}:$dest" 2>/dev/null; then
            return 0
        fi
        retries=$((retries + 1))
        log_warning "SCP到 ${INSTALL_USER}@${host} 失败，重试 $retries/$MAX_RETRIES..."
        sleep 2
    done
    
    log_error "SCP到 ${INSTALL_USER}@${host} 失败，已重试 $MAX_RETRIES 次"
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
    local config_file="config.properties"
    check_file "$config_file"
    
    # 读取基础配置
    SEATUNNEL_VERSION=$(grep "^SEATUNNEL_VERSION=" "$config_file" | cut -d'=' -f2)
    BASE_DIR=$(grep "^BASE_DIR=" "$config_file" | cut -d'=' -f2)
    SSH_PORT=$(grep "^SSH_PORT=" "$config_file" | cut -d'=' -f2)
    DEPLOY_MODE=$(grep "^DEPLOY_MODE=" "$config_file" | cut -d'=' -f2)
    
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
    
    # 置SEATUNNEL_HOME
    SEATUNNEL_HOME="$BASE_DIR/seatunnel-$SEATUNNEL_VERSION"
    
    # 验证必要的配置
    [[ -z "$SEATUNNEL_VERSION" ]] && log_error "SEATUNNEL_VERSION 未配置"
    [[ -z "$BASE_DIR" ]] && log_error "BASE_DIR 未配置"
    [[ -z "$MASTER_IPS_STRING" ]] && log_error "MASTER_IP 未配置"
    [[ -z "$WORKER_IPS_STRING" ]] && log_error "WORKER_IPS 未配置"
    
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
            PACKAGE_PATH="$(cd "$(dirname "$0")" && pwd)/$PACKAGE_PATH"
        fi
    else
        # 读取可选的下载URL
        DOWNLOAD_URL=$(grep "^DOWNLOAD_URL=" "$config_file" | cut -d'=' -f2)
        if [[ -z "$DOWNLOAD_URL" ]]; then
            DOWNLOAD_URL="https://archive.apache.org/dist/seatunnel/${SEATUNNEL_VERSION}/apache-seatunnel-${SEATUNNEL_VERSION}-bin.tar.gz"
        fi
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
    local temp_dir="/tmp/seatunnel-install-$INSTALL_USER"
    local temp_file
    
    # 创建临时目录（如果不存在）
    if [ ! -d "$temp_dir" ]; then
        mkdir -p "$temp_dir" || log_error "无法创建临时目录: $temp_dir"
        chmod 700 "$temp_dir"
        chown "$INSTALL_USER:$INSTALL_GROUP" "$temp_dir"
    fi
    
    # 创建临时文件
    temp_file=$(mktemp "$temp_dir/temp.XXXXXX") || log_error "无法创建临时文件"
    chmod 600 "$temp_file"
    chown "$INSTALL_USER:$INSTALL_GROUP" "$temp_file"
    
    # 添加到临时文件列表
    TEMP_FILES+=("$temp_file")
    
    echo "$temp_file"
}

# 清理临时文件
cleanup_temp_files() {
    local temp_dir="/tmp/seatunnel-install-$INSTALL_USER"
    
    # 清理所有已记录的临时文件
    for temp_file in "${TEMP_FILES[@]}"; do
        if [ -f "$temp_file" ]; then
            rm -f "$temp_file"
        fi
    done
    
    # 清理临时目录（如果为空）
    if [ -d "$temp_dir" ] && [ -z "$(ls -A "$temp_dir")" ]; then
        rm -rf "$temp_dir"
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
        
        # 第一步：找到section的起始行号
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
    
    # 备份原始文件
    cp "$config_file" "${config_file}.bak"
    
    case "$config_file" in
        *"hazelcast.yaml")
            log_info "修改 hazelcast.yaml (集群通信配置)..."
            if [ "$DEPLOY_MODE" = "hybrid" ]; then
                log_info "混合模式: 所有节点使用相端口5801"
                # 生成新的member-list内容
                content=$(for node in "${ALL_NODES[@]}"; do
                    echo "- ${node}:5801"
                done)
                replace_yaml_section "$config_file" "member-list:" 10 "$content"
            fi
            ;;
        *"hazelcast-client.yaml")
            log_info "修改 hazelcast-client.yaml (客户端连接配置)..."
            # 生成新的cluster-members内容
            if [ "$DEPLOY_MODE" = "hybrid" ]; then
                log_info "混合模式: 客户端可连接任意节点的5801端口"
                content=$(for node in "${ALL_NODES[@]}"; do
                    echo "- ${node}:5801"
                done)
            else
                log_info "分离模式: 客户端仅连接Master节点的5801端口"
                content=$(for master in "${MASTER_IPS[@]}"; do
                    echo "- ${master}:5801"
                done)
            fi
            replace_yaml_section "$config_file" "cluster-members:" 6 "$content"
            ;;
        *"hazelcast-master.yaml")
            if [ "$DEPLOY_MODE" != "hybrid" ]; then
                log_info "修改 hazelcast-master.yaml (Master节点配置)..."
                log_info "分离模式: Master使用5801端口，Worker使用5802端口"
                # 生成新的member-list内容
                content=$(
                    for master in "${MASTER_IPS[@]}"; do
                        echo "- ${master}:5801"
                    done
                    for worker in "${WORKER_IPS[@]}"; do
                        echo "- ${worker}:5802"
                    done
                )
                replace_yaml_section "$config_file" "member-list:" 10 "$content"
            fi
            ;;
        *"hazelcast-worker.yaml")
            if [ "$DEPLOY_MODE" != "hybrid" ]; then
                log_info "修改 hazelcast-worker.yaml (Worker节点配置)..."
                log_info "分离模式: Master使用5801端口，Worker使用5802端口"
                # 生成新的member-list内容
                content=$(
                    for master in "${MASTER_IPS[@]}"; do
                        echo "- ${master}:5801"
                    done
                    for worker in "${WORKER_IPS[@]}"; do
                        echo "- ${worker}:5802"
                    done
                )
                replace_yaml_section "$config_file" "member-list:" 10 "$content"
            fi
            ;;
    esac
    
    # 显示修改结果
    if [ -f "${config_file}.bak" ]; then
        log_info "配置文件已备份为: ${config_file}.bak"
        if diff -q "${config_file}.bak" "$config_file" >/dev/null; then
            log_info "配置文件未发生变化"
        else
            log_info "配置文件已更新"
            if [ "${DEBUG:-false}" = true ]; then
                echo "修改详情:"
                diff "${config_file}.bak" "$config_file" || true
            fi
        fi
    fi
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
    sudo chmod +x "$SEATUNNEL_HOME/bin/seatunnel-cluster.sh"
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        run_as_user "$SEATUNNEL_HOME/bin/seatunnel-cluster.sh -d"
    else
        # 在master节点启动master服务
        if [[ "$HOSTNAME" == "$MASTER_IP" ]]; then
            run_as_user "$SEATUNNEL_HOME/bin/seatunnel-cluster.sh -d -r master"
        else
            # 在worker节点启动worker服务
            run_as_user "$SEATUNNEL_HOME/bin/seatunnel-cluster.sh -d -r worker"
        fi
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
        for node in "${ALL_NODES[@]}"; do
            # 使用check_port并捕获其退出态
            if ! check_port "$node" 5801 2>/dev/null; then
                occupied_ports+=("$node:5801")
            fi
        done
    else
        for master in "${MASTER_IPS[@]}"; do
            if ! check_port "$master" 5801 2>/dev/null; then
                occupied_ports+=("$master:5801")
            fi
        done
        for worker in "${WORKER_IPS[@]}"; do
            if ! check_port "$worker" 5802 2>/dev/null; then
                occupied_ports+=("$worker:5802")
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
    sleep 10
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式检查所有节点的5801端口
        for node in "${ALL_NODES[@]}"; do
            if ! check_port "$node" 5801; then
                log_info "节点 $node 服务启动成功"
            else
                log_warning "节点 $node 服未响应"
            fi
        done
    else
        # 分离模式分别检查masterworker
        for master in "${MASTER_IPS[@]}"; do
            if ! check_port "$master" 5801; then
                log_info "Master节点 $master 服务启动成功"
            else
                log_warning "Master节点 $master 服务未响应"
            fi
        done
        
        for worker in "${WORKER_IPS[@]}"; do
            if ! check_port "$worker" 5802; then
                log_info "Worker节点 $worker 服务启动成功"
            else
                log_warning "Worker节点 $worker 服务未响应"
            fi
        done
    fi
}

# 配置检查点存储
configure_checkpoint() {
    local node_count=${#ALL_NODES[@]}
    local content
    
    # Validate storage type
    if [[ -z "$CHECKPOINT_STORAGE_TYPE" ]]; then
        log_error "检查点存储类型未配置"
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
        
        # 在所有节点上创建目录
        for node in "${ALL_NODES[@]}"; do
            if [[ "$node" != "$(hostname)" ]]; then
                log_info "节点 $node 上创建检查点目录..."
                ssh_with_retry "$node" "sudo mkdir -p $checkpoint_dir && sudo chown $INSTALL_USER:$INSTALL_GROUP $checkpoint_dir && sudo chmod 755 $checkpoint_dir"
            fi
        done
        
        if [ "$node_count" -gt 1 ]; then
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
        *)
            log_error "不支持的检查点存储类型: $CHECKPOINT_STORAGE_TYPE"
            ;;
    esac
    
    # 替换plugin-config部分，调整缩进值为10
    replace_yaml_section "$SEATUNNEL_HOME/config/seatunnel.yaml" "plugin-config:" 10 "$content"
}

# 在主脚本开头添加Java版本检查函数
check_java() {
    log_info "检查Java环境..."
    
    # 检查java命令是否存在
    if ! command -v java >/dev/null 2>&1; then
        log_error "未找到Java环境，请先安装Java 8或Java 11"
    fi
    
    # 获取Java版本
    local java_version
    java_version=$(java -version 2>&1 | awk -F '"' '/version/ {print $2}')
    
    # 检查Java版本是否为8或11
    if [[ $java_version == 1.8* ]]; then
        log_info "检测到Java 8: $java_version"
    elif [[ $java_version == 11* ]]; then
        log_info "检测到Java 11: $java_version"
    else
        log_error "不支持的Java版本: $java_version，SeaTunnel需要Java 8或Java 11"
    fi
    
    # 检查JAVA_HOME是否设置
    if [[ -z "${JAVA_HOME}" ]]; then
        log_warning "JAVA_HOME环境变量未设置，建议设置JAVA_HOME"
    else
        log_info "JAVA_HOME: $JAVA_HOME"
    fi
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
    
    if command -v curl >/dev/null 2>&1; then
        if curl --connect-timeout "$timeout" -sI "$url" >/dev/null 2>&1; then
            return 0
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget --spider -q -T "$timeout" "$url" >/dev/null 2>&1; then
            return 0
        fi
    fi
    return 1
}

# 下载安装包
download_package() {
    local output_file=$1
    local version=$2
    local retries=3
    local retry_count=0
    local download_tool
    
    log_info "开始下载安装包..."
    
    # 检查下载工具
    if command -v wget >/dev/null 2>&1; then
        download_tool="wget"
    elif command -v curl >/dev/null 2>&1; then
        download_tool="curl"
    else
        log_error "未找到wget或curl,请安装后重试"
    fi
    
    # 获取仓库配置
    local repo=${PACKAGE_REPO:-aliyun}
    local url
    
    if [ "$repo" = "custom" ]; then
        if [ -z "$CUSTOM_PACKAGE_URL" ]; then
            log_error "使用自定义仓库时必须配置 CUSTOM_PACKAGE_URL"
        fi
        url=$(echo "$CUSTOM_PACKAGE_URL" | sed "s/\${version}/$version/g")
    else
        # 获取发布包下载地址
        url="${PACKAGE_REPOS[$repo]}"
        if [ -z "$url" ]; then
            log_error "不支持的安装包仓库: $repo"
        fi
        url="$url/${version}/apache-seatunnel-${version}-bin.tar.gz"
    fi
    
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
        
        if [ "$download_tool" = "wget" ]; then
            if wget -q --show-progress "$url" -O "$output_file"; then
                log_info "下载完成: $output_file"
                return 0
            fi
        else
            if curl -# -L "$url" -o "$output_file"; then
                log_info "下载完成: $output_file"
                return 0
            fi
        fi
        
        retry_count=$((retry_count + 1))
        log_warning "下载失败,等待重试..."
        sleep 3
    done
    
    log_error "下载失败,已重试 $retries 次"
}

# 添加安装包验证函数
verify_package() {
    local package_file=$1
    
    log_info "验证安装包: $package_file"
    
    # 检查文件是否存在
    if [ ! -f "$package_file" ]; then
        log_error "安装包不存在: $package_file"
    fi
    
    # 检文件格式
    if ! file "$package_file" | grep -q "gzip compressed data"; then
        log_error "安装包格式错误,必须是tar.gz格式"
    fi
    
    # 检查是否包含版本号
    if ! tar -tzf "$package_file" 2>/dev/null | grep -q "seatunnel-$SEATUNNEL_VERSION"; then
        log_warning "安装包可能与配置的版本号不匹配: $SEATUNNEL_VERSION"
    fi
    
    log_info "安装包验证通过"
}

# 添加集群管理脚本
setup_cluster_scripts() {
    log_info "添加集群管理脚本..."
    
    # 创建master和workers文件
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        printf "%s\n" "${CLUSTER_NODES[@]}" > "$SEATUNNEL_HOME/bin/master"
        printf "%s\n" "${CLUSTER_NODES[@]}" > "$SEATUNNEL_HOME/bin/workers"
    else
        printf "%s\n" "${MASTER_IPS[@]}" > "$SEATUNNEL_HOME/bin/master"
        printf "%s\n" "${WORKER_IPS[@]}" > "$SEATUNNEL_HOME/bin/workers"
    fi
    
    # 创建集群启动脚本
    cat > "$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh" << EOF
#!/bin/bash 
  
# 定义 SeaTunnelServer 进程名称，需要根据实际情况进行修改
PROCESS_NAME="org.apache.seatunnel.core.starter.seatunnel.SeaTunnelServer"
# 定义颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 安装用户
INSTALL_USER="$INSTALL_USER"

# 日志函数
log_info() {
  echo -e "\$(date '+%Y-%m-%d %H:%M:%S') [INFO] \$1"
}

log_error() {
  echo -e "\$(date '+%Y-%m-%d %H:%M:%S') [${RED}ERROR${NC}] \$1"
}

log_success() {
  echo -e "\$(date '+%Y-%m-%d %H:%M:%S') [${GREEN}SUCCESS${NC}] \$1"
}

log_warning() {
  echo -e "\$(date '+%Y-%m-%d %H:%M:%S') [${YELLOW}WARNING${NC}] \$1"
}

bin_dir="$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
export SEATUNNEL_HOME="\$(dirname \$bin_dir)"
log_info "SEATUNNEL_HOME: \${SEATUNNEL_HOME}"
master_conf="\${bin_dir}/master"
workers_conf="\${bin_dir}/workers"

if [ -f "\$master_conf" ]; then
    mapfile -t masters < <(sed 's/[[:space:]]*$//' "\$master_conf")
else
    log_error "找不到 \$master_conf 文件"
    exit 1
fi

if [ -f "\$workers_conf" ]; then
    mapfile -t workers < <(sed 's/[[:space:]]*$//' "\$workers_conf")
else
    log_error "找不到 \$workers_conf 文件"
    exit 1
fi

mapfile -t servers < <(sort -u <(sed 's/[[:space:]]*$//' "\$master_conf" "\$workers_conf"))

sshPort=22

start(){
    echo "-------------------------------------------------"
    for master in "\${masters[@]}"; do
        if [ "\$master" = "localhost" ]; then
            log_warning "检测到仅有本地程跳过远程执行..."
            \${bin_dir}/seatunnel-cluster.sh -d  -r master
            log_success "\${master}的SeaTunnel-master启动成功"
        else
            log_info "正在 \${master} 上启动 SeaTunnelServer。"
            ssh -p \$sshPort -o StrictHostKeyChecking=no "\${INSTALL_USER}@\${master}" "source /etc/profile && source ~/.bashrc && \${bin_dir}/seatunnel-cluster.sh -d  -r master"
            log_success "\${master}的SeaTunnel-master启动成功"    
        fi
    done

    for worker in "\${workers[@]}"; do
        if [ "\$worker" = "localhost" ]; then
            log_warning "检测到仅有本地进程，跳过远程执行..."
            \${bin_dir}/seatunnel-cluster.sh -d  -r worker
            log_success "\${worker}的SeaTunnel-worker启动成功"
        else
            log_info "正在 \${worker} 上启动 SeaTunnelServer。"
            ssh -p \$sshPort -o StrictHostKeyChecking=no "\${INSTALL_USER}@\${worker}" "source /etc/profile && source ~/.bashrc && \${bin_dir}/seatunnel-cluster.sh -d  -r worker"
            log_success "\${worker}的SeaTunnel-worker启动成功"    
        fi
    done
}

stop(){
    echo "-------------------------------------------------"
    for server in "\${servers[@]}"; do
        if [ "\$server" = "localhost" ]; then
            log_warning "检测到仅有本地进程，跳过远程执行..."
            \${bin_dir}/stop-seatunnel-cluster.sh
            log_success "\${server}的SeaTunnel 停止成功"
        else
            log_info "正在 \${server} 上停止 SeaTunnelServer。"
            ssh -p \$sshPort -o StrictHostKeyChecking=no "\${INSTALL_USER}@\${server}" "source /etc/profile && source ~/.bashrc && \${bin_dir}/stop-seatunnel-cluster.sh"
            log_success "\${server}的SeaTunnel 停止成功"
        fi
    done
}

restart(){
    stop
    sleep 2
    start
}

case "\$1" in
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
        echo "用法：\$0 {start|stop|restart}"
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

# 安装插件和JDBC驱动
install_plugins_and_drivers() {
    # 检查是否需要安装连接器
    if [ "${INSTALL_CONNECTORS}" != "true" ]; then
        log_info "跳过连接器和依赖安装"
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
    
    # 获取插件仓库配置
    local repo=${PLUGIN_REPO:-aliyun}
    local repo_url
    
    if [ "$repo" = "custom" ]; then
        if [ -z "$CUSTOM_PLUGIN_URL" ]; then
            log_error "使用自定义插件仓库时必须配置 CUSTOM_PLUGIN_URL"
        fi
        repo_url="$CUSTOM_PLUGIN_URL"
    else
        repo_url="${PLUGIN_REPOS[$repo]}"
        if [ -z "$repo_url" ]; then
            log_error "不支持的插件仓库: $repo"
        fi
    fi
    
    # 下载连接器
    for connector in "${enabled_connectors[@]}"; do
        connector=$(echo "$connector" | tr -d '[:space:]')
        
        # 构建下载URL
        local download_url=$(echo "$PLUGIN_DOWNLOAD_URL" | \
            sed "s/\${version}/$SEATUNNEL_VERSION/g" | \
            sed "s/\${connector}/$connector/g" | \
            sed "s|\${repo}|$repo_url|g")
            
        local connector_jar="$connectors_dir/connector-${connector}-${SEATUNNEL_VERSION}.jar"
        
        log_info "下载连接器: $connector (从 $repo 仓库)"
        if command -v wget >/dev/null 2>&1; then
            if ! wget -q --show-progress "$download_url" -O "$connector_jar"; then
                # 如果当前仓库下载失败,尝试Apache仓库
                if [ "$repo" != "apache" ]; then
                    log_warning "从 $repo 下载失败,尝试Apache仓库..."
                    download_url=$(echo "$PLUGIN_DOWNLOAD_URL" | \
                        sed "s/\${version}/$SEATUNNEL_VERSION/g" | \
                        sed "s/\${connector}/$connector/g" | \
                        sed "s|\${repo}|${PLUGIN_REPOS[apache]}|g")
                    if ! wget -q --show-progress "$download_url" -O "$connector_jar"; then
                        log_error "下载连接器失败: $connector"
                    fi
                else
                    log_error "下载连接器失败: $connector"
                fi
            fi
        else
            if ! curl -# -L "$download_url" -o "$connector_jar"; then
                # 如果当前仓库下载失败,尝试Apache仓库
                if [ "$repo" != "apache" ]; then
                    log_warning "从 $repo 下载失败,尝试Apache仓库..."
                    download_url=$(echo "$PLUGIN_DOWNLOAD_URL" | \
                        sed "s/\${version}/$SEATUNNEL_VERSION/g" | \
                        sed "s/\${connector}/$connector/g" | \
                        sed "s|\${repo}|${PLUGIN_REPOS[apache]}|g")
                    if ! curl -# -L "$download_url" -o "$connector_jar"; then
                        log_error "下载连接器失败: $connector"
                    fi
                else
                    log_error "下载连接器失败: $connector"
                fi
            fi
        fi
        
        # 安装连接器依赖
        local libs_var="${connector}_libs[@]"
        local libs=("${!libs_var}")
        
        if [ ${#libs[@]} -gt 0 ]; then
            log_info "安装 $connector 连接器的依赖..."
            
            # 处理依赖列表
            for dep in "${libs[@]}"; do
                dep=$(echo "$dep" | tr -d '[:space:]' | tr -d '"')  # 去除空格和引号
                
                if [[ "$dep" == *":"* ]]; then
                    # Maven依赖
                    IFS=':' read -r groupId artifactId version <<< "$dep"
                    local jar_name="${artifactId}-${version}.jar"
                    local jar_path="$lib_dir/$jar_name"
                    
                    # 构建下载URL (使用阿里云Maven仓库)
                    local group_path=$(echo "$groupId" | tr '.' '/')
                    local download_url="https://maven.aliyun.com/repository/public/${group_path}/${artifactId}/${version}/${jar_name}"
                    
                    log_info "下载依赖: $jar_name"
                    if command -v wget >/dev/null 2>&1; then
                        if ! wget -q --show-progress "$download_url" -O "$jar_path"; then
                            # 如果阿里云下载失败，尝试中央仓库
                            download_url="https://repo1.maven.org/maven2/${group_path}/${artifactId}/${version}/${jar_name}"
                            if ! wget -q --show-progress "$download_url" -O "$jar_path"; then
                                log_error "下载依赖失败: $jar_name"
                            fi
                        fi
                    else
                        if ! curl -# -L "$download_url" -o "$jar_path"; then
                            # 如果阿里云下载失败尝试中央仓库
                            download_url="https://repo1.maven.org/maven2/${group_path}/${artifactId}/${version}/${jar_name}"
                            if ! curl -# -L "$download_url" -o "$jar_path"; then
                                log_error "下载依赖失败: $jar_name"
                            fi
                        fi
                    fi
                fi
            done
        fi
    done
    
    # 设置权限
    setup_permissions "$connectors_dir"
    setup_permissions "$lib_dir"
    
    log_info "插件和依赖安装完成"
}

# 配置开机自启动
setup_auto_start() {
    if [ "${ENABLE_AUTO_START}" != "true" ]; then
        log_info "跳过开机自启动配置"
        return 0
    fi
    
    log_info "配置SeaTunnel开机自启动..."
    
    # 创建systemd服务文件
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式服务文件
        local service_file="/etc/systemd/system/seatunnel.service"
        cat > "$service_file" << EOF
[Unit]
Description=Apache SeaTunnel Service (Hybrid Mode)
After=network.target

[Service]
Type=simple
User=${INSTALL_USER}
Group=${INSTALL_GROUP}
Environment="JAVA_HOME=${JAVA_HOME}"
Environment="PATH=${PATH}:${JAVA_HOME}/bin"
WorkingDirectory=${SEATUNNEL_HOME}
ExecStart=${SEATUNNEL_HOME}/bin/seatunnel-cluster.sh -d
ExecStop=/bin/bash -c 'pkill -f "org.apache.seatunnel.core.starter.seatunnel.SeaTunnelServer -d"'
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
    else
        # 分离模式 - 创建master服务
        if [[ " ${MASTER_IPS[*]} " =~ " $(hostname) " ]]; then
            local service_file="/etc/systemd/system/seatunnel-master.service"
            cat > "$service_file" << EOF
[Unit]
Description=Apache SeaTunnel Master Service
After=network.target

[Service]
Type=simple
User=${INSTALL_USER}
Group=${INSTALL_GROUP}
Environment="JAVA_HOME=${JAVA_HOME}"
Environment="PATH=${PATH}:${JAVA_HOME}/bin"
WorkingDirectory=${SEATUNNEL_HOME}
ExecStart=${SEATUNNEL_HOME}/bin/seatunnel-cluster.sh -d -r master
ExecStop=/bin/bash -c 'pkill -f "org.apache.seatunnel.core.starter.seatunnel.SeaTunnelServer -d -r master"'
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
        else
            # 分离模式 - 创建worker服务
            local service_file="/etc/systemd/system/seatunnel-worker.service"
            cat > "$service_file" << EOF
[Unit]
Description=Apache SeaTunnel Worker Service
After=network.target

[Service]
Type=simple
User=${INSTALL_USER}
Group=${INSTALL_GROUP}
Environment="JAVA_HOME=${JAVA_HOME}"
Environment="PATH=${PATH}:${JAVA_HOME}/bin"
WorkingDirectory=${SEATUNNEL_HOME}
ExecStart=${SEATUNNEL_HOME}/bin/seatunnel-cluster.sh -d -r worker
ExecStop=/bin/bash -c 'pkill -f "org.apache.seatunnel.core.starter.seatunnel.SeaTunnelServer -d -r worker"'
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
        fi
    fi

    # 设置权限
    chmod 644 "$service_file"
    
    # 重新加载systemd配置
    sudo systemctl daemon-reload
    
    # 启用服务
    sudo systemctl enable "$(basename "$service_file")"
    
    log_info "开机自启动配置完成"
    
    # 显示服务状态
    echo -e "\n当前服务状态:"
    sudo systemctl status "$(basename "$service_file")"
    
    echo -e "\n可用的systemd命令:"
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        echo "启动服务:    sudo systemctl start seatunnel"
        echo "停止服务:    sudo systemctl stop seatunnel"
        echo "重启服务:    sudo systemctl restart seatunnel"
        echo "查看状态:    sudo systemctl status seatunnel"
        echo "禁用自启动:  sudo systemctl disable seatunnel"
    else
        if [[ " ${MASTER_IPS[*]} " =~ " $(hostname) " ]]; then
            echo "启动服务:    sudo systemctl start seatunnel-master"
            echo "停止服务:    sudo systemctl stop seatunnel-master"
            echo "重启服务:    sudo systemctl restart seatunnel-master"
            echo "查看状态:    sudo systemctl status seatunnel-master"
            echo "禁用自启动:  sudo systemctl disable seatunnel-master"
        else
            echo "启动服务:    sudo systemctl start seatunnel-worker"
            echo "停止服务:    sudo systemctl stop seatunnel-worker"
            echo "重启服务:    sudo systemctl restart seatunnel-worker"
            echo "查看状态:    sudo systemctl status seatunnel-worker"
            echo "禁用自启动:  sudo systemctl disable seatunnel-worker"
        fi
    fi
}

# 主函数
main() {
    log_info "使用模式: \$([ "\$FORCE_SED" = true ] && echo "强制sed" || echo "自动选择")"
    
    # 检查系统依赖
    check_dependencies
    
    # 检查Java环境
    check_java
    
    # 初始化临时文件列表
    declare -a TEMP_FILES=()
    
    # 设置退出时清理
    trap cleanup_temp_files EXIT INT TERM
    
    # 读取配置
    read_config
    
    # 检查用户配置
    check_user
    
    # 处理安装包
    if [ "\$INSTALL_MODE" = "online" ]; then
        # 在线安装模式
        local package_name="apache-seatunnel-\${SEATUNNEL_VERSION}-bin.tar.gz"
        local package_path="\$(cd "\$(dirname "\$0")" && pwd)/\$package_name"
        
        # 如果本地已存在,询问是否重新下载
        if [ -f "\$package_path" ]; then
            log_warning "发现本地已存在安装包: \$package_path"
            read -r -p "是否重新下载? [y/N] " response
            case "\$response" in
                [yY][eE][sS]|[yY])
                    rm -f "\$package_path"
                    ;;
                *)
                    log_info "使用已存在的安装包"
                    ;;
            esac
        fi
        
        # 下载安装包
        if [ ! -f "\$package_path" ]; then
            download_package "\$package_path" "\$SEATUNNEL_VERSION"
        fi
        
        PACKAGE_PATH="\$package_path"
    fi
    
    # 验证安装包
    verify_package "\$PACKAGE_PATH"
    
    # 创建安装目录
    log_info "创建安装目录..."
    create_directory "\$BASE_DIR"
    setup_permissions "\$BASE_DIR"
    
    # 解压安装包
    log_info "解压安装包..."
    sudo tar -zxf "\$PACKAGE_PATH" -C "\$BASE_DIR"
    setup_permissions "\$SEATUNNEL_HOME"
    
    # 配置文件修改
    log_info "修改配置文件..."
    
    # 根据部署模式配置集群
    if [ "\$DEPLOY_MODE" = "hybrid" ]; then
        setup_hybrid_mode
    else
        setup_separated_mode
    fi

    # 添加集群管理脚本
    setup_cluster_scripts

    # 分发到其他节点
    log_info "分发其他节点..."
    for node in "\${WORKER_IPS[@]}"; do
        log_info "分发到 \$node..."
        ssh_with_retry "\$node" "sudo mkdir -p \$BASE_DIR && sudo chown \$INSTALL_USER:\$INSTALL_GROUP \$BASE_DIR"
        scp_with_retry "\$SEATUNNEL_HOME" "\$node" "\$BASE_DIR/"
    done

    # 配置环境变量
    log_info "配置环境变量..."
    BASHRC_CONTENT="
# SEATUNNEL_HOME BEGIN
export SEATUNNEL_HOME=\$SEATUNNEL_HOME
export PATH=\$PATH:\$SEATUNNEL_HOME/bin
# SEATUNNEL_HOME END"

    # 获取用户home目录，使用多种方式
    USER_HOME=""
    if command -v getent >/dev/null 2>&1; then
        # 使用getent
        USER_HOME=\$(getent passwd "\$INSTALL_USER" | cut -d: -f6)
    else
        # 使用eval
        USER_HOME=\$(eval echo ~"\$INSTALL_USER")
    fi
    
    if [ -z "\$USER_HOME" ]; then
        log_error "无法获取用户 \$INSTALL_USER 的home目录"
    fi
    
    # 为指定用户添加环境变量
    run_as_user "echo '\$BASHRC_CONTENT' >> \$USER_HOME/.bashrc"
    
    # 远程节点环境变量
    for node in "\${WORKER_IPS[@]}"; do
        ssh -p \$SSH_PORT "\$node" "sudo -u \$INSTALL_USER bash -c \"echo '\$BASHRC_CONTENT' >> \$(eval echo ~\$INSTALL_USER)/.bashrc\""
    done

    # 启动集群
    start_cluster

    # 检查服务状态
    check_services

    # 配置检查点存储
    log_info "配置检查点存储..."
    configure_checkpoint

    # 添加安装完成后的验证提示
    echo -e "\n\${GREEN}SeaTunnel安装完成!\$NC"
    echo -e "\n\${YELLOW}验证和使用说明:${NC}"
    echo "1. 刷新环境变量:"
    echo -e "\${GREEN}source \$USER_HOME/.bashrc\$NC"
    
    echo -e "\n2. 集群管理命令:"
    echo -e "启动集群:  \${GREEN}\$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh start\$NC"
    echo -e "停止集群:  \${GREEN}\$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh stop\$NC"
    echo -e "重启集群:  \${GREEN}\$SEATUNNEL_HOME/bin/seatunnel-start-cluster.sh restart\$NC"
    
    echo -e "\n3. 验证安装:"
    echo -e "运行示例任务: \${GREEN}\$SEATUNNEL_HOME/bin/seatunnel.sh --config config/v2.batch.config.template\$NC"
    
    echo -e "\n\${YELLOW}部信息:${NC}"
    if [ "\$DEPLOY_MODE" = "hybrid" ]; then
        echo "部署模式: 混合模式"
        echo -e "集群节: \${GREEN}\${CLUSTER_NODES[*]}\$NC"
    else
        echo "部署模式: 分离模式"
        echo -e "Master节点: \${GREEN}\${MASTER_IPS[*]}\$NC"
        echo -e "Worker节点: \${GREEN}\${WORKER_IPS[*]}\$NC"
    fi
    
    echo -e "\n\${YELLOW}注意事项:${NC}"
    echo "1. 首次启动集群前，请确保所有节点的环境变量已经生效"
    echo "2. 建议先执行 start 命令启动集群，验证所有节点是否正常启动"
    echo "3. 如果需要停止集群，请使用 stop 命令，确保所有进程正常关闭"
    echo -e "4. 更多使用说明请参考：\${GREEN}https://seatunnel.apache.org/docs\$NC"

    # 配置开机自启动
    setup_auto_start
}

# 执行主函数
main 