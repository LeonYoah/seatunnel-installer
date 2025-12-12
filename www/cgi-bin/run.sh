#!/bin/bash
# ============================================================================
# SeaTunnel Web API - 统一入口
# ============================================================================
# 直接调用 install_seatunnel.sh 执行安装步骤
# 由 busybox httpd 作为 CGI 脚本调用
# ============================================================================

echo "Content-Type: application/json"
echo ""

# 获取项目根目录 (www/cgi-bin -> 项目根目录)
SCRIPT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
INSTALL_SCRIPT="$SCRIPT_DIR/install_seatunnel.sh"
CONFIG_FILE="$SCRIPT_DIR/config.properties"
LOG_FILE="$SCRIPT_DIR/.web_install.log"
PID_FILE="$SCRIPT_DIR/.install_pid"
STATUS_FILE="$SCRIPT_DIR/.install_status"
STEP_STATUS_FILE="$SCRIPT_DIR/.step_status"

# 解析查询参数
parse_param() {
    echo "$QUERY_STRING" | sed -n "s/.*$1=\([^&]*\).*/\1/p"
}

ACTION=$(parse_param "action")
STEP=$(parse_param "step")
STOP_AT=$(parse_param "stop_at")
START_STEP=$(parse_param "start")
END_STEP=$(parse_param "end")

case "$ACTION" in
    # ========================================
    # 检测临时文件
    # ========================================
    check_temp)
        files=""
        [ -f "$LOG_FILE" ] && files="${files}log,"
        [ -f "$STATUS_FILE" ] && files="${files}status,"
        [ -f "$STEP_STATUS_FILE" ] && files="${files}step_status,"
        [ -f "$PID_FILE" ] && files="${files}pid,"
        files="${files%,}"
        
        if [ -n "$files" ]; then
            echo "{\"status\":\"found\",\"files\":\"$files\"}"
        else
            echo '{"status":"clean"}'
        fi
        ;;
    
    # ========================================
    # 清理临时文件
    # ========================================
    clean_temp)
        rm -f "$LOG_FILE" "$STATUS_FILE" "$STEP_STATUS_FILE" "$PID_FILE" 2>/dev/null
        echo '{"status":"ok","message":"临时文件已清理"}'
        ;;
    
    # ========================================
    # 获取步骤列表
    # ========================================
    steps)
        bash "$INSTALL_SCRIPT" --list-steps --json 2>/dev/null
        ;;
    
    # ========================================
    # 执行安装 (完整或指定步骤)
    # ========================================
    run)
        # 检查是否已在运行
        if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
            echo '{"status":"running","message":"安装已在进行中"}'
            exit 0
        fi
        
        # 构建命令
        cmd="bash $INSTALL_SCRIPT"
        [ -n "$STEP" ] && cmd+=" --step $STEP"
        [ -n "$STOP_AT" ] && cmd+=" --stop-at $STOP_AT"
        
        # 清空日志和状态
        > "$LOG_FILE"
        echo "running" > "$STATUS_FILE"
        
        # 后台执行
        (
            echo $$ > "$PID_FILE"
            eval "$cmd" >> "$LOG_FILE" 2>&1
            EXIT_CODE=$?
            [ $EXIT_CODE -eq 0 ] && echo "completed" > "$STATUS_FILE" || echo "failed:$EXIT_CODE" > "$STATUS_FILE"
            rm -f "$PID_FILE"
        ) &
        
        echo '{"status":"ok","message":"安装已启动"}'
        ;;
    
    # ========================================
    # 执行步骤范围 (大步骤)
    # ========================================
    run_range)
        # 检查是否已在运行
        if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
            echo '{"status":"running","message":"安装已在进行中"}'
            exit 0
        fi
        
        if [ -z "$START_STEP" ] || [ -z "$END_STEP" ]; then
            echo '{"status":"error","message":"缺少 start 或 end 参数"}'
            exit 0
        fi
        
        # 追加分隔线到日志，不清空历史
        echo "" >> "$LOG_FILE"
        echo "========== $(date '+%Y-%m-%d %H:%M:%S') 执行阶段 $START_STEP-$END_STEP ==========" >> "$LOG_FILE"
        echo "running" > "$STATUS_FILE"
        
        # 初始化步骤状态 - 先删除该阶段所有步骤，再写入
        for s in $(seq "$START_STEP" "$END_STEP"); do
            sed -i "/^$s:/d" "$STEP_STATUS_FILE" 2>/dev/null
        done
        for s in $(seq "$START_STEP" "$END_STEP"); do
            echo "$s:pending" >> "$STEP_STATUS_FILE"
        done
        
        # 后台执行步骤范围
        (
            echo $$ > "$PID_FILE"
            for step in $(seq "$START_STEP" "$END_STEP"); do
                # 更新步骤状态为 running
                sed -i "s/^$step:.*/$step:running/" "$STEP_STATUS_FILE"
                
                bash "$INSTALL_SCRIPT" --step "$step" >> "$LOG_FILE" 2>&1
                if [ $? -ne 0 ]; then
                    sed -i "s/^$step:.*/$step:failed/" "$STEP_STATUS_FILE"
                    echo "failed:step$step" > "$STATUS_FILE"
                    rm -f "$PID_FILE"
                    exit 1
                fi
                
                # 更新步骤状态为 completed
                sed -i "s/^$step:.*/$step:completed/" "$STEP_STATUS_FILE"
            done
            echo "completed" > "$STATUS_FILE"
            rm -f "$PID_FILE"
        ) &
        
        echo '{"status":"ok","message":"阶段已启动"}'
        ;;
    
    # ========================================
    # 获取安装状态
    # ========================================
    status)
        status="idle"
        progress=0
        current_step=0
        message=""
        
        # 检查进程
        [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null && status="running"
        
        # 读取状态文件
        if [ -f "$STATUS_FILE" ]; then
            case "$(cat "$STATUS_FILE")" in
                completed) status="completed"; progress=100 ;;
                failed*) status="failed" ;;
                paused) status="paused" ;;
                running) [ "$status" != "running" ] && status="pending" ;;
            esac
        fi
        
        # 从日志分析进度
        if [ -f "$LOG_FILE" ]; then
            # 匹配 [N/M] 格式 (install_seatunnel.sh 输出格式)
            last=$(grep -oE '\[[0-9]+/[0-9]+\]' "$LOG_FILE" | tail -1)
            if [ -n "$last" ]; then
                current_step=$(echo "$last" | sed 's/\[\([0-9]*\)\/.*/\1/')
                total=$(echo "$last" | sed 's/.*\/\([0-9]*\)\]/\1/')
                [ -z "$total" ] && total=16
                [ "$total" -gt 0 ] 2>/dev/null && progress=$((current_step * 100 / total))
            fi
            # 去除所有控制字符
            message=$(tail -1 "$LOG_FILE" 2>/dev/null | head -c 200 | LC_ALL=C tr -cd '[:print:]' | sed 's/"/\\"/g')
            grep -q "安装完成\|全部步骤执行完成\|SUCCESS" "$LOG_FILE" && status="completed" && progress=100
            grep -qE "\[ERROR\]|安装中止|失败$" "$LOG_FILE" && [ "$status" != "completed" ] && status="failed"
        fi
        
        # 读取步骤状态 - 每个步骤只取最后一条记录
        step_status=""
        if [ -f "$STEP_STATUS_FILE" ]; then
            # 获取所有唯一步骤号，然后取每个步骤的最后状态
            for snum in $(cut -d: -f1 "$STEP_STATUS_FILE" | sort -u); do
                sstat=$(grep "^$snum:" "$STEP_STATUS_FILE" | tail -1 | cut -d: -f2)
                [ -n "$sstat" ] && step_status="${step_status}\"$snum\":\"$sstat\","
            done
            step_status="${step_status%,}"
        fi
        
        echo "{\"status\":\"$status\",\"progress\":$progress,\"current_step\":$current_step,\"message\":\"$message\",\"steps\":{$step_status}}"
        ;;
    
    # ========================================
    # 执行单个步骤 (重试)
    # ========================================
    run_step)
        if [ -z "$STEP" ]; then
            echo '{"status":"error","message":"缺少 step 参数"}'
            exit 0
        fi
        
        # 检查是否已在运行
        if [ -f "$PID_FILE" ] && kill -0 "$(cat "$PID_FILE")" 2>/dev/null; then
            echo '{"status":"running","message":"有任务正在执行"}'
            exit 0
        fi
        
        echo "" >> "$LOG_FILE"
        echo "========== $(date '+%Y-%m-%d %H:%M:%S') 重试步骤 $STEP ==========" >> "$LOG_FILE"
        echo "running" > "$STATUS_FILE"
        
        # 更新步骤状态
        sed -i "s/^$STEP:.*/$STEP:running/" "$STEP_STATUS_FILE" 2>/dev/null || echo "$STEP:running" >> "$STEP_STATUS_FILE"
        
        # 后台执行单步骤
        (
            echo $$ > "$PID_FILE"
            bash "$INSTALL_SCRIPT" --step "$STEP" >> "$LOG_FILE" 2>&1
            if [ $? -ne 0 ]; then
                sed -i "s/^$STEP:.*/$STEP:failed/" "$STEP_STATUS_FILE"
                echo "failed:step$STEP" > "$STATUS_FILE"
            else
                sed -i "s/^$STEP:.*/$STEP:completed/" "$STEP_STATUS_FILE"
                echo "completed" > "$STATUS_FILE"
            fi
            rm -f "$PID_FILE"
        ) &
        
        echo '{"status":"ok","message":"步骤已启动"}'
        ;;
    
    # ========================================
    # 停止安装
    # ========================================
    stop)
        if [ -f "$PID_FILE" ]; then
            pid=$(cat "$PID_FILE")
            kill -TERM "$pid" 2>/dev/null
            sleep 1
            kill -KILL "$pid" 2>/dev/null || true
            rm -f "$PID_FILE"
        fi
        echo "stopped" > "$STATUS_FILE"
        echo '{"status":"ok","message":"已停止"}'
        ;;
    
    # ========================================
    # 暂停执行 (标记当前运行中的步骤为 paused)
    # ========================================
    pause)
        if [ -f "$PID_FILE" ]; then
            pid=$(cat "$PID_FILE")
            # 终止当前执行进程
            kill -TERM "$pid" 2>/dev/null
            sleep 1
            kill -KILL "$pid" 2>/dev/null || true
            rm -f "$PID_FILE"
            
            # 将所有 running 状态的步骤改为 paused
            if [ -f "$STEP_STATUS_FILE" ]; then
                sed -i 's/:running$/:paused/' "$STEP_STATUS_FILE"
            fi
            echo "paused" > "$STATUS_FILE"
            echo '{"status":"ok","message":"已暂停"}'
        else
            echo '{"status":"ok","message":"没有正在执行的任务"}'
        fi
        ;;
    
    # ========================================
    # 获取安装日志
    # ========================================
    log)
        lines=$(parse_param "lines")
        lines=${lines:-100}
        if [ -f "$LOG_FILE" ]; then
            # 过滤日志：去除 ANSI 码、控制字符、进度条行
            content=$(tail -"$lines" "$LOG_FILE" 2>/dev/null | \
                sed 's/\x1b\[[0-9;]*m//g' | \
                sed 's/\[0;[0-9]*m//g; s/\[0m//g; s/\[1;[0-9]*m//g' | \
                grep -v '^[[:space:]]*[#=O-]*[[:space:]]*[0-9.]*%*$' | \
                grep -v '^[[:space:]]*$' | \
                tr -d '\r' | \
                sed 's/\\/\\\\/g; s/"/\\"/g; s/\t/  /g' | \
                awk '{printf "%s\\n", $0}')
            echo "{\"status\":\"ok\",\"log\":\"$content\"}"
        else
            echo '{"status":"ok","log":""}'
        fi
        ;;
    
    # ========================================
    # 加载配置
    # ========================================
    config_load)
        if [ ! -f "$CONFIG_FILE" ]; then
            echo '{"status":"error","message":"配置文件不存在"}'
            exit 0
        fi
        echo -n '{"status":"ok","config":{'
        first=true
        while IFS='=' read -r key value; do
            [[ "$key" =~ ^[[:space:]]*# ]] && continue
            [[ -z "$key" ]] && continue
            key=$(echo "$key" | tr -d '[:space:]')
            value=$(echo "$value" | sed 's/"/\\"/g')
            [ "$first" = true ] && first=false || echo -n ','
            echo -n "\"$key\":\"$value\""
        done < "$CONFIG_FILE"
        echo '}}'
        ;;
    
    # ========================================
    # 保存配置 (POST)
    # ========================================
    config_save)
        read -n "${CONTENT_LENGTH:-0}" POST_DATA 2>/dev/null || POST_DATA=""
        if [ -z "$POST_DATA" ]; then
            echo '{"status":"error","message":"没有数据"}'
            exit 0
        fi
        # 使用 sed 正则解析 JSON 键值对，正确处理值中的逗号
        # 匹配 "key":"value" 或 "key":"value,with,commas"
        echo "$POST_DATA" | sed 's/^{//;s/}$//' | grep -oE '"[^"]+":"[^"]*"' | while read -r pair; do
            key=$(echo "$pair" | sed 's/^"\([^"]*\)".*/\1/')
            value=$(echo "$pair" | sed 's/^"[^"]*":"\(.*\)"$/\1/')
            [ -z "$key" ] && continue
            if grep -q "^${key}=" "$CONFIG_FILE" 2>/dev/null; then
                sed -i "s|^${key}=.*|${key}=${value}|" "$CONFIG_FILE"
            else
                echo "${key}=${value}" >> "$CONFIG_FILE"
            fi
        done
        echo '{"status":"ok","message":"保存成功"}'
        ;;
    
    # ========================================
    # 环境检查
    # ========================================
    check)
        type=$(parse_param "type")
        type=${type:-all}
        results='{'
        
        # OS 检查
        if [ "$type" = "all" ] || [ "$type" = "os" ]; then
            os_name=$(grep "^PRETTY_NAME=" /etc/os-release 2>/dev/null | cut -d'"' -f2 || uname -s)
            results+="\"os\":{\"status\":\"ok\",\"message\":\"$os_name\"},"
        fi
        
        # Java 检查
        if [ "$type" = "all" ] || [ "$type" = "java" ]; then
            if command -v java >/dev/null 2>&1; then
                java_ver=$(java -version 2>&1 | head -n1 | sed 's/.*"\([^"]*\)".*/\1/')
                if [[ "$java_ver" == 1.8* ]] || [[ "$java_ver" == 11* ]]; then
                    results+="\"java\":{\"status\":\"ok\",\"message\":\"Java $java_ver\"},"
                else
                    results+="\"java\":{\"status\":\"warning\",\"message\":\"Java $java_ver (推荐8或11)\"},"
                fi
            else
                results+="\"java\":{\"status\":\"warning\",\"message\":\"未安装Java\"},"
            fi
        fi
        
        # SSH 检查
        if [ "$type" = "all" ] || [ "$type" = "ssh" ]; then
            if systemctl is-active sshd >/dev/null 2>&1 || systemctl is-active ssh >/dev/null 2>&1; then
                results+="\"ssh\":{\"status\":\"ok\",\"message\":\"SSH运行中\"},"
            else
                results+="\"ssh\":{\"status\":\"warning\",\"message\":\"SSH未运行\"},"
            fi
        fi
        
        # 磁盘检查
        if [ "$type" = "all" ] || [ "$type" = "disk" ]; then
            avail=$(df -BG / | tail -1 | awk '{print $4}' | tr -d 'G')
            if [ "$avail" -ge 10 ] 2>/dev/null; then
                results+="\"disk\":{\"status\":\"ok\",\"message\":\"${avail}GB可用\"},"
            else
                results+="\"disk\":{\"status\":\"warning\",\"message\":\"${avail}GB可用\"},"
            fi
        fi
        
        results="${results%,}}"
        echo "{\"status\":\"ok\",\"checks\":$results}"
        ;;
    
    # ========================================
    # 未知操作
    # ========================================
    *)
        echo '{"status":"error","message":"未知action，支持: steps, run, status, stop, log, config_load, config_save, check"}'
        ;;
esac
