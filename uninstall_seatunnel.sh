#!/bin/bash

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
    fi
}

# 停止服务
stop_services() {
    log_info "停止所有SeaTunnel服务..."
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        # 混合模式：停止所有节点的服务
        for node in "${ALL_NODES[@]}"; do
            log_info "停止节点 $node 的服务..."
            if [ "$node" = "$current_ip" ] || [ "$node" = "localhost" ]; then
                # 本地节点
                if systemctl is-active --quiet seatunnel; then
                    sudo systemctl stop seatunnel
                    sudo systemctl disable seatunnel
                    sudo rm -f /etc/systemd/system/seatunnel.service
                    sudo systemctl daemon-reload
                    log_success "本地节点 $node 的服务已停止"
                else
                    log_info "本地节点 $node 的服务未运行"
                fi
            else
                # 远程节点
                ssh -p "$SSH_PORT" "$node" "
                    if systemctl is-active --quiet seatunnel; then
                        sudo systemctl stop seatunnel
                        sudo systemctl disable seatunnel
                        sudo rm -f /etc/systemd/system/seatunnel.service
                        sudo systemctl daemon-reload
                        echo '节点 $node 的服务已停止'
                    else
                        echo '节点 $node 的服务未运行'
                    fi
                "
            fi
        done
    else
        # 分离模式：分别停止master和worker节点的服务
        # 停止Master节点
        for master in "${MASTER_IPS[@]}"; do
            log_info "停止Master节点 $master 的服务..."
            if [ "$master" = "$current_ip" ] || [ "$master" = "localhost" ]; then
                # 本地Master节点
                if systemctl is-active --quiet seatunnel-master; then
                    sudo systemctl stop seatunnel-master
                    sudo systemctl disable seatunnel-master
                    sudo rm -f /etc/systemd/system/seatunnel-master.service
                    sudo systemctl daemon-reload
                    log_success "本地Master节点 $master 的服务已停止"
                else
                    log_info "本地Master节点 $master 的服务未运行"
                fi
            else
                # 远程Master节点
                ssh -p "$SSH_PORT" "$master" "
                    if systemctl is-active --quiet seatunnel-master; then
                        sudo systemctl stop seatunnel-master
                        sudo systemctl disable seatunnel-master
                        sudo rm -f /etc/systemd/system/seatunnel-master.service
                        sudo systemctl daemon-reload
                        echo 'Master节点 $master 的服务已停止'
                    else
                        echo 'Master节点 $master 的服务未运行'
                    fi
                "
            fi
        done

        # 停止Worker节点
        for worker in "${WORKER_IPS[@]}"; do
            log_info "停止Worker节点 $worker 的服务..."
            if [ "$worker" = "$current_ip" ] || [ "$worker" = "localhost" ]; then
                # 本地Worker节点
                if systemctl is-active --quiet seatunnel-worker; then
                    sudo systemctl stop seatunnel-worker
                    sudo systemctl disable seatunnel-worker
                    sudo rm -f /etc/systemd/system/seatunnel-worker.service
                    sudo systemctl daemon-reload
                    log_success "本地Worker节点 $worker 的服务已停止"
                else
                    log_info "本地Worker节点 $worker 的服务未运行"
                fi
            else
                # 远程Worker节点
                ssh -p "$SSH_PORT" "$worker" "
                    if systemctl is-active --quiet seatunnel-worker; then
                        sudo systemctl stop seatunnel-worker
                        sudo systemctl disable seatunnel-worker
                        sudo rm -f /etc/systemd/system/seatunnel-worker.service
                        sudo systemctl daemon-reload
                        echo 'Worker节点 $worker 的服务已停止'
                    else
                        echo 'Worker节点 $worker 的服务未运行'
                    fi
                "
            fi
        done
    fi

    # 检查是否有遗留进程
    # local process_name="org.apache.seatunnel.core.starter.seatunnel.SeaTunnelServer"
    # for node in "${ALL_NODES[@]}"; do
    #     if [ "$node" = "$current_ip" ]; then
    #         if pgrep -f "$process_name" >/dev/null; then
    #             log_warning "发现遗留进程,强制终止..."
    #             pkill -9 -f "$process_name"
    #         fi
    #     else
    #         ssh -p "$SSH_PORT" "$node" "
    #             if pgrep -f '$process_name' >/dev/null; then
    #                 echo '发现遗留进程,强制终止...'
    #                 pkill -9 -f '$process_name'
    #             fi
    #         "
    #     fi
    # done

    log_success "所有服务已停止"
}

# 清理环境变量
clean_environment() {
    log_info "清理环境变量..."
    
    # 清理本地环境变量
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
    
    # 清理远程节点环境变量
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        for node in "${ALL_NODES[@]}"; do
            if [ "$node" != "$current_ip" ] && [ "$node" != "localhost" ]; then
                remote_home=$(ssh -p "$SSH_PORT" "$node" "echo ~$INSTALL_USER")
                ssh -p "$SSH_PORT" "$node" "sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' '$remote_home/.bashrc'"
            fi
        done
    else
        for node in "${WORKER_IPS[@]}"; do
            if [ "$node" != "$current_ip" ] && [ "$node" != "localhost" ]; then
                remote_home=$(ssh -p "$SSH_PORT" "$node" "echo ~$INSTALL_USER")
                ssh -p "$SSH_PORT" "$node" "sed -i '/# SEATUNNEL_HOME BEGIN/,/# SEATUNNEL_HOME END/d' '$remote_home/.bashrc'"
            fi
        done
    fi
}

# 删除安装目录
remove_files() {
    log_info "删除安装目录..."
    
    if [ -d "$BASE_DIR" ]; then
        sudo rm -rf "$BASE_DIR"
    fi
    
    # 删除远程节点的安装目录
    local current_ip
    current_ip=$(hostname -I | awk '{print $1}')
    
    if [ "$DEPLOY_MODE" = "hybrid" ]; then
        for node in "${ALL_NODES[@]}"; do
            if [ "$node" != "$current_ip" ] && [ "$node" != "localhost" ]; then
                ssh -p "$SSH_PORT" "$node" "sudo rm -rf '$BASE_DIR'"
            fi
        done
    else
        for node in "${WORKER_IPS[@]}"; do
            if [ "$node" != "$current_ip" ] && [ "$node" != "localhost" ]; then
                ssh -p "$SSH_PORT" "$node" "sudo rm -rf '$BASE_DIR'"
            fi
        done
    fi
}

# 主函数
main() {
    # 读取配置
    read_config

    echo -e "${YELLOW}警告: 此操作将完全删除SeaTunnel的所有组件和配置!${NC}"
    echo -e "\n${YELLOW}将执行以下操作:${NC}"
    echo "1. 停止所有SeaTunnel服务"
    echo "2. 删除所有节点上的安装目录 $BASE_DIR"
    echo "3. 清理所有节点上的环境变量配置"
    echo "4. 删除systemd服务配置(如果已配置)"
    
    echo -e "\n${YELLOW}注意事项:${NC}"
    echo "1. 环境变量将在下次登录后生效"
    echo "2. 如果有自定义的配置文件或数据，请手动备份"
    echo "3. 此操作无法撤销"
    
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
    
    # 删除安装目录
    remove_files
    
    log_info "SeaTunnel卸载完成"
}

# 执行主函数
main 