export default {
  app: {
    name: 'SeaTunnel',
    subtitle: '企业级管理平台',
    titleSuffix: 'SeaTunnel 企业级管理平台'
  },
  settings: {
    notify: {
      title: '通知渠道',
      dingWebhook: '钉钉 Webhook',
      wecomWebhook: '企微 Webhook',
      email: '邮件接收人',
      sendTest: '发送测试',
      saved: '通知配置已保存',
      testSent: '测试通知已发送'
    },
    system: {
      title: '系统配置',
      heartbeatTimeout: '心跳超时',
      seconds: '秒',
      logRetention: '日志保留天数',
      days: '天',
      autoBackup: '自动备份',
      backupTime: '备份时间',
      pickTime: '选择时间',
      saved: '系统配置已保存'
    },
    users: {
      title: '用户管理',
      add: '添加用户',
      addWip: '添加用户功能开发中',
      roles: {
        admin: '管理员',
        operator: '操作员',
        viewer: '查看者'
      },
      active: '激活',
      disabled: '禁用',
      columns: {
        username: '用户名',
        email: '邮箱',
        role: '角色',
        status: '状态',
        lastLogin: '最后登录',
        createdAt: '创建时间'
      },
      resetPwd: '重置密码',
      edit: '编辑用户：{name}',
      confirmReset: '确定要重置用户 "{name}" 的密码吗？',
      resetOk: '密码已重置',
      confirmDelete: '确定要删除用户 "{name}" 吗？'
    },
    audit: {
      title: '审计日志',
      exporting: '开始导出审计日志',
      columns: {
        user: '操作者',
        action: '操作',
        resource: '资源',
        result: '结果',
        ip: 'IP 地址',
        timestamp: '时间',
        details: '详情'
      }
    }
  },
  deploy: {
    steps: {
      selectHosts: '选择主机',
      start: '开始部署'
    },
    select: {
      tip: '提示',
      desc: '请从已注册的主机中选择要部署 SeaTunnel 集群的节点。如果没有可用主机，请先前往',
      desc2: '注册主机并安装 Agent。',
      master: '选择 Master',
      masterPlaceholder: '请选择 Master 节点',
      workers: '选择 Workers',
      workersPlaceholder: '请选择 Worker 节点（可多选）',
      clusterNodes: '集群节点',
      nodes: '选择节点',
      nodesPlaceholder: '请选择集群节点（可多选）',
      selectedCount: '已选主机（{count} 台）'
    },
    modes: {
      separated: '分离模式（Master/Worker）',
      hybrid: '混合模式（所有节点角色相同）'
    }
  },
  install: {
    steps: {
      config: '配置参数',
      precheck: '环境检查',
      plugins: '插件选择',
      install: '安装部署',
      complete: '完成'
    },
    config: {
      title: '安装配置',
      basic: '基础配置',
      version: 'SeaTunnel 版本',
      installMode: '安装模式',
      online: '在线安装',
      offline: '离线安装',
      baseDir: '安装目录',
      deploy: '部署配置',
      deployMode: '部署模式',
      separated: '分离模式',
      hybrid: '混合模式',
      clusterNodes: '集群节点'
    }
    ,
    precheck: {
      allPassed: '全部通过',
      items: {
        memory: '内存检查',
        cpu: 'CPU 检查',
        disk: '磁盘空间检查',
        ssh: 'SSH 连通性检查',
        port: '端口占用检查',
        firewall: '防火墙状态检查'
      },
      msg: {
        memoryOk: '可用内存: 16GB',
        cpuOk: 'CPU 核心数: 8',
        diskOk: '可用空间: 100GB',
        sshOk: '所有节点连接正常',
        portOk: '端口可用',
        firewallOk: '防火墙已关闭'
      }
    }
    ,
    plugins: {
      selected: '已选 {count} 个插件',
      source: 'Source 连接器',
      sink: 'Sink 连接器'
    }
    ,
    install: {
      installing: '安装中...',
      completed: '安装完成',
      logs: '安装日志',
      scrollBottom: '滚动到底部',
      done: '[INFO] 安装完成！',
      steps: {
        download: '下载安装包',
        unpack: '解压安装包',
        configure: '配置集群',
        plugins: '安装插件',
        distribute: '分发到节点',
        start: '启动集群'
      }
    }
  },
  diagnostics: {
    quick: {
      title: '一键诊断',
      scope: '诊断范围',
      select: '请选择',
      target: '目标对象',
      targetPlaceholder: '如 worker-01 / orders_sync',
      items: '收集内容',
      generate: '生成诊断包'
    },
    scope: {
      cluster: '整个集群',
      node: '指定节点',
      task: '指定任务'
    },
    items: {
      logs: '日志文件',
      config: '配置快照',
      thread: '线程 Dump',
      heap: '堆内存 Dump',
      gc: 'GC 日志'
    },
    fault: {
      title: '常见故障库',
      pattern: '特征',
      view: '查看方案',
      autoFix: '一键修复',
      viewSolution: '查看故障解决方案：{title}',
      confirmAutoFix: '确定要自动修复 "{title}" 吗？',
      fixed: '修复成功'
    },
    history: {
      title: '诊断历史',
      scope: '范围',
      target: '目标',
      items: '收集内容',
      status: '状态',
      size: '大小',
      createdAt: '创建时间',
      downloading: '开始下载诊断包 #{id}',
      confirmDelete: '确定要删除诊断包 #{id} 吗？'
    },
    msg: {
      targetRequired: '请输入目标对象',
      itemsRequired: '请选择至少一项收集内容',
      created: '诊断任务已创建，正在收集数据...'
    }
  },
  // 合并完整安装完成文案
  // 为避免重复层级，这里放在 install.complete 下
  //（前面的 install 已定义 steps/config/precheck/install）
  'install.complete': {
    title: 'SeaTunnel 集群部署成功！',
    subtitle: '恭喜！您的 SeaTunnel 集群已成功部署并启动。',
    clusterInfo: '集群信息',
    clusterName: '集群名称',
    version: 'SeaTunnel 版本',
    deployMode: '部署模式',
    nodeCount: '节点数量',
    installPath: '安装路径',
    access: '访问地址',
    masterHttp: 'Master HTTP 端口:',
    quickStart: '快速开始',
    stepCheck: '查看集群状态',
    stepSubmit: '提交第一个任务',
    stepSubmitDesc: '前往任务管理页面创建和提交数据集成任务',
    stepMonitor: '监控集群',
    stepMonitorDesc: '在集群管理页面查看节点状态和资源使用情况',
    goConsole: '进入控制台',
    viewClusters: '查看集群',
    downloadReport: '下载安装报告',
    reportStart: '开始下载安装报告'
  },
  clusters: {
    overview: {
      totalClusters: '集群总数',
      totalNodes: '节点总数',
      avgCpu: '平均 CPU 使用率'
    },
    listTitle: '集群列表',
    register: '注册集群',
    version: '版本',
    nodes: '节点',
    columns: {
      nodeName: '节点名称',
      role: '角色',
      ip: 'IP 地址',
      memory: '内存',
      disk: '磁盘',
      status: '状态',
      lastHeartbeat: '最后心跳'
    },
    restart: '重启',
    diagnose: '诊断',
    config: '配置管理',
    scaleOut: '扩容',
    scaleIn: '缩容',
    delete: '删除集群',
    msg: {
      registerWip: '注册集群功能开发中',
      restart: '已触发重启：{name}',
      diagnose: '诊断节点：{name}',
      config: '配置集群：{name}',
      scaleOut: '扩容集群：{name}',
      scaleIn: '缩容集群：{name}',
      confirmDelete: '确定要删除集群 "{name}" 吗？'
    }
  },
  plugins: {
    refreshRepo: '刷新仓库',
    uploadLocal: '本地上传',
    searchPlaceholder: '搜索插件名称',
    available: '可用',
    install: '安装',
    upgrade: '升级',
    disable: '禁用',
    columns: {
      name: '插件名称',
      type: '类型',
      version: '版本',
      compatibility: '兼容性',
      desc: '描述',
      status: '状态'
    },
    msg: {
      uploadWip: '本地上传功能开发中',
      confirmInstall: '确定要安装插件 "{name}" 吗？',
      installing: '正在安装 {name}...',
      upgrade: '升级插件：{name}',
      confirmDisable: '确定要禁用插件 "{name}" 吗？',
      disabled: '禁用成功',
      view: '查看插件详情：{name}'
    }
  },
  dashboard: {
    kpis: {
      health: '集群健康',
      taskSuccess: '任务成功率',
      latencyP95: '延迟 P95',
      activeAlerts: '活跃告警'
    },
    trend: {
      title: '任务执行趋势',
      last24h: '过去 24 小时',
      placeholder: '图表占位 - 后续接入 ECharts 展示趋势数据'
    },
    alerts: {
      title: '最近告警',
      empty: '暂无告警'
    },
    recentTasks: {
      title: '最近任务'
    },
    msg: {
      alertsWip: '告警功能开发中',
      retryTriggered: '已触发重试：{name}'
    }
  },
  common: {
    ok: '确定',
    cancel: '取消',
    loading: '加载中',
    actions: '操作',
    confirm: '确定',
    confirmDelete: '确认删除',
    deleteSuccess: '删除成功',
    cancelled: '已取消',
    refreshSuccess: '刷新成功',
    create: '新建',
    edit: '编辑',
    view: '详情',
    delete: '删除',
    run: '运行',
    retry: '重试',
    refresh: '刷新',
    save: '保存',
    add: '添加',
    export: '导出',
    download: '下载',
    all: '全部',
    statusFilter: '状态筛选',
    typeFilter: '类型筛选',
    next: '下一步',
    prev: '上一步',
    viewAll: '查看全部'
  },
  status: {
    running: '运行中',
    success: '成功',
    failed: '失败',
    paused: '暂停',
    completed: '已完成',
    processing: '进行中',
    online: '在线',
    offline: '离线',
    installed: '已安装',
    notInstalled: '未安装',
    healthy: '健康',
    unhealthy: '异常'
  },
  user: {
    profile: '个人信息',
    logout: '退出登录',
    admin: '管理员'
  },
  theme: {
    toLight: '切换到明亮模式',
    toDark: '切换到暗色模式'
  },
  lang: {
    zh: '简体中文',
    en: 'English',
    switch: '语言'
  },
  menu: {
    dashboard: '总览',
    hosts: '主机管理',
    deploy: '部署集群',
    tasks: '任务管理',
    clusters: '集群管理',
    diagnostics: '诊断中心',
    plugins: '插件市场',
    settings: '设置'
  },
  route: {
    dashboard: '总览',
    hosts: '主机管理',
    deploy: '部署集群',
    tasks: '任务管理',
    clusters: '集群管理',
    diagnostics: '诊断中心',
    plugins: '插件市场',
    settings: '设置'
  },
  footer: {
    copyright: 'SeaTunnel 企业级管理平台 | Apache SeaTunnel'
  },
  tips: {
    profileWip: '个人信息功能开发中',
    loggedOut: '已退出登录'
  },
  hosts: {
    register: '注册主机',
    edit: '编辑主机',
    searchPlaceholder: '搜索主机 IP 或名称',
    agentNotInstalled: 'Agent 未安装',
    installAgent: '安装 Agent',
    uninstallAgent: '卸载 Agent',
    testConnection: '测试连接',
    columns: {
      name: '主机名称',
      ip: 'IP 地址',
      port: 'SSH 端口',
      user: 'SSH 用户',
      agentStatus: 'Agent 状态',
      hostStatus: '主机状态',
      cpu: 'CPU 使用率',
      memory: '内存使用率',
      lastHeartbeat: '最后心跳'
    },
    form: {
      name: '主机名称',
      ip: 'IP 地址',
      port: 'SSH 端口',
      user: 'SSH 用户',
      authType: '认证方式',
      password: '密码',
      key: '密钥',
      sshPassword: 'SSH 密码',
      keyPath: '私钥路径',
      desc: '描述',
      descPlaceholder: '主机描述信息',
      exampleName: '例如：master-01',
      exampleIp: '例如：192.168.1.100',
      exampleUser: '例如：root',
      exampleKeyPath: '例如：~/.ssh/id_rsa'
    },
    valid: {
      name: '请输入主机名称',
      ipRequired: '请输入 IP 地址',
      ipFormat: '请输入有效的 IP 地址',
      port: '请输入 SSH 端口',
      user: '请输入 SSH 用户',
      password: '请输入 SSH 密码',
      keyPath: '请输入私钥路径'
    },
    msg: {
      saved: '主机注册成功',
      confirmInstallAgent: '确定要在主机 "{name}" ({ip}) 上安装 Agent 吗？',
      installSubmitted: 'Agent 安装任务已提交，请稍候...',
      confirmUninstallAgent: '确定要卸载主机 "{name}" ({ip}) 上的 Agent 吗？',
      uninstallSuccess: 'Agent 卸载成功',
      testing: '正在测试与主机 {name} 的连接...',
      testSuccess: '连接测试成功',
      confirmDelete: '确定要删除主机 "{name}" 吗？'
    }
  },
  tasks: {
    searchPlaceholder: '搜索任务名称',
    type: { streaming: '实时', batch: '批处理' },
    columns: {
      name: '任务名称',
      type: '类型',
      status: '状态',
      lastRun: '最近运行',
      duration: '耗时',
      version: '版本',
      creator: '创建者'
    },
    msg: {
      createWip: '新建任务功能开发中',
      runTriggered: '已触发运行：{name}',
      edit: '编辑任务：{name}',
      view: '查看任务详情：{name}',
      confirmDelete: '确定要删除任务 "{name}" 吗？'
    }
  }
}
