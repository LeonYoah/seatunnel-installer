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

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') $1"
}

# 带重试的SSH命令执行
ssh_with_retry() {
    local host=$1
    local cmd=$2
    local retries=0
    local max_retries=3
    local timeout=10
    
    log_info "尝试SSH连接到 $host..."
    
    while [ $retries -lt $max_retries ]; do
        if timeout $timeout ssh -p "$SSH_PORT" -o StrictHostKeyChecking=no -o ConnectTimeout=5 "$host" "sudo bash -c '$cmd'" 2>ssh_error.tmp; then
            rm -f ssh_error.tmp
            return 0
        fi
        retries=$((retries + 1))
        log_warning "SSH到 $host 失败，重试 $retries/$max_retries..."
        log_warning "错误信息: $(cat ssh_error.tmp)"
        sleep 2
    done
    
    log_error "SSH到 $host 失败，已重试 $max_retries 次。错误信息: $(cat ssh_error.tmp)"
    rm -f ssh_error.tmp
    return 1
}

# 读取配置文件
read_config() {
    local config_file="config.properties"
    if [ ! -f "$config_file" ]; then
        log_error "配置文件不存在: $config_file"
    fi
    
    # 读取基础配置
    BASE_DIR=$(grep "^BASE_DIR=" "$config_file" | cut -d'=' -f2)
    INSTALL_USER=$(grep "^INSTALL_USER=" "$config_file" | cut -d'=' -f2)
    DEPLOY_MODE=$(grep "^DEPLOY_MODE=" "$config_file" | cut -d'=' -f2)
    SSH_PORT=$(grep "^SSH_PORT=" "$config_file" | cut -d'=' -f2)
    ENABLE_AUTO_START=$(grep "^ENABLE_AUTO_START=" "$config_file" | cut -d'=' -f2)
    
    # 读取节点配置
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        CLUSTER_NODES=$(grep "^CLUSTER_NODES=" "$config_file" | cut -d'=' -f2)
        IFS=',' read -r -a ALL_NODES <<< "$CLUSTER_NODES"
    else
        MASTER_IP=$(grep "^MASTER_IP=" "$config_file" | cut -d'=' -f2)
        WORKER_IPS=$(grep "^WORKER_IPS=" "$config_file" | cut -d'=' -f2)
        IFS=',' read -r -a MASTER_IPS <<< "$MASTER_IP"
        IFS=',' read -r -a WORKER_IPS <<< "$WORKER_IPS"
        # 合并所有节点到一个数组
        ALL_NODES=("${MASTER_IPS[@]}" "${WORKER_IPS[@]}")
    fi
}

# 检查节点连通性
check_nodes() {
    log_info "检查节点连通性..."
    local failed_nodes=()
    
    log_info "需要检查的节点: ${ALL_NODES[*]}"
    
    for node in "${ALL_NODES[@]}"; do
        log_info "正在检查节点: $node"
        if [ "$node" != "localhost" ] && [ "$node" != "$(hostname -I | awk '{print $1}')" ]; then
            if ! ping -c 1 -W 5 "$node" &>/dev/null; then
                log_warning "节点 $node ping测试失败"
                failed_nodes+=("$node")
                continue
            fi
            
            if ! ssh_with_retry "$node" "echo 'Connection test'" &>/dev/null; then
                log_warning "节点 $node SSH测试失败"
                failed_nodes+=("$node")
            fi
        else
            log_info "跳过本地节点检查: $node"
        fi
    done
    
    if [ ${#failed_nodes[@]} -gt 0 ]; then
        log_error "以下节点无法连接:\n${failed_nodes[*]}"
    fi
}

# 远程执行脚本函数
execute_remote_script() {
    local host=$1
    local script_content=$2
    local tmp_script="/tmp/seatunnel_tmp_$RANDOM.sh"
    
    # 创建临时脚本
    echo "#!/bin/bash" > "$tmp_script"
    echo "$script_content" >> "$tmp_script"
    chmod +x "$tmp_script"
    
    # 复制到远程主机
    if ! scp -P "$SSH_PORT" -o StrictHostKeyChecking=no "$tmp_script" "${host}:${tmp_script}" >/dev/null 2>&1; then
        rm -f "$tmp_script"
        return 1
    fi
    
    # 执行远程脚本
    if ! ssh_with_retry "$host" "bash $tmp_script"; then
        ssh_with_retry "$host" "rm -f $tmp_script" || true
        rm -f "$tmp_script"
        return 1
    fi
    
    # 清理临时文件
    ssh_with_retry "$host" "rm -f $tmp_script" || true
    rm -f "$tmp_script"
    return 0
}

# 停止服务
stop_services() {
    log_info "停止所有SeaTunnel服务..."
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    # 定义服务停止函数
    stop_service_on_node() {
        local node=$1
        local service_name=$2
        local is_local=$3
        local service_file="/etc/systemd/system/${service_name}.service"
        
        if [ "$is_local" = true ]; then
            if [ -f "$service_file" ]; then
                log_info "停止本地节点 $node 的 $service_name 服务..."
                # 如果服务在运行,先停止
                if systemctl is-active --quiet "$service_name" 2>/dev/null; then
                    sudo systemctl stop "$service_name"
                fi
                sudo systemctl disable "$service_name"
                sudo rm -f "$service_file"
                sudo systemctl daemon-reload
                log_success "本地节点 $node 的 $service_name 服务已清理"
            else
                log_info "本地节点 $node 的 $service_name 服务文件不存在,无需清理"
            fi
        else
            local script_content="
                service_file=/etc/systemd/system/${service_name}.service
                if [ -f \"\$service_file\" ]; then
                    # 如果服务在运行,先停止
                    if systemctl is-active --quiet $service_name 2>/dev/null; then
                        sudo systemctl stop $service_name
                    fi
                    sudo systemctl disable $service_name
                    sudo rm -f \"\$service_file\"
                    sudo systemctl daemon-reload
                    echo '节点 $node 的 $service_name 服务已清理'
                else
                    echo '节点 $node 的 $service_name 服务文件不存在,无需清理'
                fi"
            
            execute_remote_script "$node" "$script_content"
        fi
    }
    
    # 停止进程函数
    kill_process_on_node() {
        local node=$1
        local process_pattern=$2
        local is_local=$3
        
        if [ "$is_local" = true ]; then
            if pgrep -f "$process_pattern" >/dev/null; then
                log_warning "发现遗留进程,强制终止..."
                for pid in $(pgrep -f "$process_pattern"); do
                    sudo kill -9 "$pid" || true
                done
            fi
        else
            local script_content="
                for pid in \$(pgrep -f \"$process_pattern\"); do
                    sudo kill -9 \$pid || true
                done"
            
            execute_remote_script "$node" "$script_content"
        fi
    }
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：停止所有节点的服务
        for node in "${ALL_NODES[@]}"; do
            local is_local=false
            if [ "$node" = "$current_ip" ] || [ "$node" = "localhost" ]; then
                is_local=true
            fi
            stop_service_on_node "$node" "seatunnel" "$is_local"
            kill_process_on_node "$node" "org.apache.seatunnel" "$is_local"
        done
    else
        # 分离模式：分别停止master和worker节点的服务
        for master in "${MASTER_IPS[@]}"; do
            local is_local=false
            if [ "$master" = "$current_ip" ] || [ "$master" = "localhost" ]; then
                is_local=true
            fi
            stop_service_on_node "$master" "seatunnel-master" "$is_local"
            kill_process_on_node "$master" "org.apache.seatunnel" "$is_local"
        done
        
        for worker in "${WORKER_IPS[@]}"; do
            local is_local=false
            if [ "$worker" = "$current_ip" ] || [ "$worker" = "localhost" ]; then
                is_local=true
            fi
            stop_service_on_node "$worker" "seatunnel-worker" "$is_local"
            kill_process_on_node "$worker" "org.apache.seatunnel" "$is_local"
        done
    fi
    
    log_success "所有服务已停止"
}

# 清理环境变量
clean_environment() {
    log_info "清理环境变量..."
    
    clean_node_environment() {
        local node=$1
        local is_local=$2
        
        if [ "$is_local" = true ]; then
            local user_home
            if command -v getent >/dev/null 2>&1; then
                user_home=$(getent passwd "$INSTALL_USER" | cut -d: -f6)
            else
                user_home=$(eval echo ~"$INSTALL_USER")
            fi
            
            if [ -f "$user_home/.bashrc" ]; then
                sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' "$user_home/.bashrc"
                sed -i '/# JAVA_HOME BEGIN/,/# JAVA_HOME END/d' "$user_home/.bashrc"
            fi
        else
            local script_content="
                user_home=\$(getent passwd $INSTALL_USER | cut -d: -f6 || echo ~$INSTALL_USER)
                if [ -f \"\$user_home/.bashrc\" ]; then
                    sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' \"\$user_home/.bashrc\"
                    sed -i '/# JAVA_HOME BEGIN/,/# JAVA_HOME END/d' \"\$user_home/.bashrc\"
                fi"
            
            execute_remote_script "$node" "$script_content"
        fi
    }
    
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    for node in "${ALL_NODES[@]}"; do
        local is_local=false
        if [ "$node" = "$current_ip" ] || [ "$node" = "localhost" ]; then
            is_local=true
        fi
        clean_node_environment "$node" "$is_local"
    done
    
    log_success "环境变量已清理"
}

# 清理Java环境
clean_java() {
    log_info "清理自动安装的Java环境..."
    
    clean_node_java() {
        local node=$1
        local is_local=$2
        local java_home="${BASE_DIR}/java"
        
        if [ "$is_local" = true ]; then
            if [ -d "$java_home" ]; then
                log_info "删除本地节点Java目录: $java_home"
                sudo rm -rf "$java_home"
            fi
        else
            ssh_with_retry "$node" "
                if [ -d '$java_home' ]; then
                    sudo rm -rf '$java_home'
                fi"
        fi
    }
    
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    for node in "${ALL_NODES[@]}"; do
        local is_local=false
        if [ "$node" = "$current_ip" ] || [ "$node" = "localhost" ]; then
            is_local=true
        fi
        clean_node_java "$node" "$is_local"
    done
}

# 删除安装目录和日志
remove_files() {
    log_info "删除安装目录和日志..."
    
    remove_node_files() {
        local node=$1
        local is_local=$2
        
        if [ "$is_local" = true ]; then
            # 删除安装目录
            if [ -d "$BASE_DIR" ]; then
                log_info "删除本地安装目录: $BASE_DIR"
                sudo rm -rf "$BASE_DIR"
            fi
        else
            ssh_with_retry "$node" "
                # 删除安装目录
                if [ -d '$BASE_DIR' ]; then
                    sudo rm -rf '$BASE_DIR'
                fi"
        fi
    }
    
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    for node in "${ALL_NODES[@]}"; do
        local is_local=false
        if [ "$node" = "$current_ip" ] || [ "$node" = "localhost" ]; then
            is_local=true
        fi
        remove_node_files "$node" "$is_local"
    done
    
    log_success "文件清理完成"
}

# 主函数
main() {
    # 读取配置
    read_config
    
    # 检查节点连通性
    check_nodes
    
    echo -e "${YELLOW}警告: 此操作将完全删除SeaTunnel的所有组件和配置!${NC}"
    echo -e "\n${YELLOW}将执行以下操作:${NC}"
    echo "1. 停止所有节点的SeaTunnel服务"
    echo "2. 删除所有节点上的安装目录: $BASE_DIR"
    echo "3. 清理所有节点上的环境变量配置"
    echo "4. 删除所有节点上的systemd服务配置"
    echo "5. 清理所有节点上自动安装的Java环境"
    echo "6. 删除所有节点上的日志文件"
    
    echo -e "\n${YELLOW}注意事项:${NC}"
    echo "1. 环境变量将在下次登录后生效"
    echo "2. 如果有自定义的配置文件或数据，请手动备份"
    echo "3. 此操作无法撤销"
    echo "4. 确保所有节点都可以通过SSH访问"
    
    read -r -p "是否继续? [y/N] " response
    case "$response" in
        [yY][eE][sS]|[yY]) 
            log_info "开始卸载SeaTunnel..."
            ;;
        *)
            log_info "卸载已取消"
            exit 0
            ;;
    esac
    
    # 停止服务
    stop_services
    
    # 清理环境变量
    clean_environment
    
    # 清理Java环境
    clean_java
    
    # 删除安装目录和日志
    remove_files
    
    log_success "SeaTunnel卸载完成"
    echo -e "\n${GREEN}提示:${NC}"
    echo "1. 环境变量的修改将在下次登录后生效"
    echo "2. 如果需要立即生效,请执行: source ~/.bashrc"
}

# 执行主函数
main