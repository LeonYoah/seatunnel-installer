export default {
  app: {
    name: 'SeaTunnel',
    subtitle: 'Enterprise Management Console',
    titleSuffix: 'SeaTunnel Console'
  },
  settings: {
    notify: {
      title: 'Notification Channels',
      dingWebhook: 'DingTalk Webhook',
      wecomWebhook: 'WeCom Webhook',
      email: 'Email Recipient',
      sendTest: 'Send Test',
      saved: 'Notification settings saved',
      testSent: 'Test notification sent'
    },
    system: {
      title: 'System Settings',
      heartbeatTimeout: 'Heartbeat Timeout',
      seconds: 'seconds',
      logRetention: 'Log Retention Days',
      days: 'days',
      autoBackup: 'Auto Backup',
      backupTime: 'Backup Time',
      pickTime: 'Select time',
      saved: 'System settings saved'
    },
    users: {
      title: 'User Management',
      add: 'Add User',
      addWip: 'Add user is under development',
      roles: {
        admin: 'Admin',
        operator: 'Operator',
        viewer: 'Viewer'
      },
      active: 'Active',
      disabled: 'Disabled',
      columns: {
        username: 'Username',
        email: 'Email',
        role: 'Role',
        status: 'Status',
        lastLogin: 'Last Login',
        createdAt: 'Created At'
      },
      resetPwd: 'Reset Password',
      edit: 'Edit user: {name}',
      confirmReset: 'Reset password for "{name}"?',
      resetOk: 'Password reset',
      confirmDelete: 'Delete user "{name}"?'
    },
    audit: {
      title: 'Audit Logs',
      exporting: 'Exporting audit logs',
      columns: {
        user: 'User',
        action: 'Action',
        resource: 'Resource',
        result: 'Result',
        ip: 'IP Address',
        timestamp: 'Time',
        details: 'Details'
      }
    }
  },
  deploy: {
    steps: {
      selectHosts: 'Select Hosts',
      start: 'Start Deploy'
    },
    select: {
      tip: 'Tip',
      desc: 'Select nodes from registered hosts. If none, go to',
      desc2: 'to register hosts and install Agent.',
      master: 'Select Master',
      masterPlaceholder: 'Choose a Master node',
      workers: 'Select Workers',
      workersPlaceholder: 'Choose Worker nodes (multi-select)',
      clusterNodes: 'Cluster Nodes',
      nodes: 'Select Nodes',
      nodesPlaceholder: 'Choose cluster nodes (multi-select)',
      selectedCount: '{count} hosts selected'
    },
    modes: {
      separated: 'Separated (Master/Worker)',
      hybrid: 'Hybrid (same role for all)'
    }
  },
	  install: {
    steps: {
      config: 'Configuration',
      precheck: 'Prechecks',
      plugins: 'Plugins',
      install: 'Install',
      complete: 'Complete'
    },
    config: {
      title: 'Install Configuration',
      basic: 'Basic Settings',
      version: 'SeaTunnel Version',
      installMode: 'Install Mode',
      online: 'Online',
      offline: 'Offline',
      baseDir: 'Install Directory',
      deploy: 'Deploy Settings',
      deployMode: 'Deploy Mode',
      separated: 'Separated',
      hybrid: 'Hybrid',
      clusterNodes: 'Cluster Nodes'
    },
    precheck: {
      allPassed: 'All Passed',
      items: {
        memory: 'Memory Check',
        cpu: 'CPU Check',
        disk: 'Disk Space Check',
        ssh: 'SSH Connectivity',
        port: 'Port Availability',
        firewall: 'Firewall Status'
      },
      msg: {
        memoryOk: 'Available memory: 16GB',
        cpuOk: 'CPU cores: 8',
        diskOk: 'Free space: 100GB',
        sshOk: 'All nodes reachable',
        portOk: 'Ports available',
        firewallOk: 'Firewall disabled'
      }
	    },
	    install: {
	      installing: 'Installing...',
	      completed: 'Installed',
	      logs: 'Install Logs',
	      scrollBottom: 'Scroll to bottom',
	      done: '[INFO] Install completed!',
	      steps: {
	        download: 'Download package',
	        unpack: 'Unpack package',
	        configure: 'Configure cluster',
	        plugins: 'Install plugins',
	        distribute: 'Distribute to nodes',
	        start: 'Start cluster'
	      }
	    },
    complete: {
      title: 'SeaTunnel cluster deployed!',
      subtitle: 'Congratulations! Your SeaTunnel cluster is up and running.',
      clusterInfo: 'Cluster Info',
      clusterName: 'Cluster Name',
      version: 'SeaTunnel Version',
      deployMode: 'Deploy Mode',
      nodeCount: 'Node Count',
      installPath: 'Install Path',
      access: 'Access',
      masterHttp: 'Master HTTP:',
      quickStart: 'Quick Start',
      stepCheck: 'Check cluster status',
      stepSubmit: 'Submit the first job',
      stepSubmitDesc: 'Go to Tasks to create and submit a job',
      stepMonitor: 'Monitor cluster',
      stepMonitorDesc: 'View node status and resources in Clusters',
      goConsole: 'Go to Console',
      viewClusters: 'View Clusters',
      downloadReport: 'Download Install Report',
      reportStart: 'Starting download of install report'
    },
    plugins: {
      selected: '{count} plugins selected',
      source: 'Source Connectors',
      sink: 'Sink Connectors'
    }
  },
  diagnostics: {
    quick: {
      title: 'Quick Diagnose',
      scope: 'Scope',
      select: 'Please select',
      target: 'Target',
      targetPlaceholder: 'e.g., worker-01 / orders_sync',
      items: 'Collect Items',
      generate: 'Generate Package'
    },
    scope: {
      cluster: 'Whole Cluster',
      node: 'Specific Node',
      task: 'Specific Task'
    },
    items: {
      logs: 'Log Files',
      config: 'Config Snapshot',
      thread: 'Thread Dump',
      heap: 'Heap Dump',
      gc: 'GC Logs'
    },
    fault: {
      title: 'Common Faults',
      pattern: 'Pattern',
      view: 'View Solution',
      autoFix: 'Auto Fix',
      viewSolution: 'View solution: {title}',
      confirmAutoFix: 'Auto fix "{title}"?',
      fixed: 'Fixed successfully'
    },
    history: {
      title: 'Diagnosis History',
      scope: 'Scope',
      target: 'Target',
      items: 'Items',
      status: 'Status',
      size: 'Size',
      createdAt: 'Created At',
      downloading: 'Downloading package #{id}',
      confirmDelete: 'Delete package #{id}?'
    },
    msg: {
      targetRequired: 'Please enter target',
      itemsRequired: 'Select at least one item',
      created: 'Diagnosis task created, collecting data...'
    }
  },
  clusters: {
    overview: {
      totalClusters: 'Clusters',
      totalNodes: 'Nodes',
      avgCpu: 'Avg CPU Usage'
    },
    listTitle: 'Cluster List',
    register: 'Register Cluster',
    version: 'Version',
    nodes: 'Nodes',
    columns: {
      nodeName: 'Node Name',
      role: 'Role',
      ip: 'IP Address',
      memory: 'Memory',
      disk: 'Disk',
      status: 'Status',
      lastHeartbeat: 'Last Heartbeat'
    },
    restart: 'Restart',
    diagnose: 'Diagnose',
    config: 'Config',
    scaleOut: 'Scale Out',
    scaleIn: 'Scale In',
    delete: 'Delete Cluster',
    msg: {
      registerWip: 'Cluster registration is under development',
      restart: 'Restart triggered: {name}',
      diagnose: 'Diagnose node: {name}',
      config: 'Config cluster: {name}',
      scaleOut: 'Scale out cluster: {name}',
      scaleIn: 'Scale in cluster: {name}',
      confirmDelete: 'Delete cluster "{name}"?'
    }
  },
  plugins: {
    refreshRepo: 'Refresh Repo',
    uploadLocal: 'Upload Local',
    searchPlaceholder: 'Search plugin name',
    available: 'Available',
    install: 'Install',
    upgrade: 'Upgrade',
    disable: 'Disable',
    columns: {
      name: 'Plugin Name',
      type: 'Type',
      version: 'Version',
      compatibility: 'Compatibility',
      desc: 'Description',
      status: 'Status'
    },
    msg: {
      uploadWip: 'Local upload is under development',
      confirmInstall: 'Install plugin "{name}"?',
      installing: 'Installing {name}...',
      upgrade: 'Upgrade plugin: {name}',
      confirmDisable: 'Disable plugin "{name}"?',
      disabled: 'Disabled successfully',
      view: 'View plugin detail: {name}'
    }
  },
  dashboard: {
    kpis: {
      health: 'Cluster Health',
      taskSuccess: 'Task Success Rate',
      latencyP95: 'Latency P95',
      activeAlerts: 'Active Alerts'
    },
    trend: {
      title: 'Task Execution Trend',
      last24h: 'Past 24 hours',
      placeholder: 'Chart placeholder — integrate ECharts later'
    },
    alerts: {
      title: 'Recent Alerts',
      empty: 'No alerts'
    },
    recentTasks: {
      title: 'Recent Tasks'
    },
    msg: {
      alertsWip: 'Alerts page is under development',
      retryTriggered: 'Retry triggered: {name}'
    }
  },
  common: {
    ok: 'OK',
    cancel: 'Cancel',
    loading: 'Loading',
    actions: 'Actions',
    confirm: 'Confirm',
    confirmDelete: 'Confirm Delete',
    deleteSuccess: 'Deleted successfully',
    cancelled: 'Cancelled',
    refreshSuccess: 'Refreshed',
    create: 'Create',
    edit: 'Edit',
    view: 'View',
    delete: 'Delete',
    run: 'Run',
    retry: 'Retry',
    refresh: 'Refresh',
    save: 'Save',
    add: 'Add',
    export: 'Export',
    download: 'Download',
    all: 'All',
    statusFilter: 'Filter by status',
    typeFilter: 'Filter by type',
    next: 'Next',
    prev: 'Previous',
    viewAll: 'View All'
  },
  status: {
    running: 'Running',
    success: 'Success',
    failed: 'Failed',
    paused: 'Paused',
    completed: 'Completed',
    processing: 'Processing',
    online: 'Online',
    offline: 'Offline',
    installed: 'Installed',
    notInstalled: 'Not Installed',
    healthy: 'Healthy',
    unhealthy: 'Unhealthy'
  },
  user: {
    profile: 'Profile',
    logout: 'Sign out',
    admin: 'Administrator'
  },
  theme: {
    toLight: 'Switch to light mode',
    toDark: 'Switch to dark mode'
  },
  lang: {
    zh: '简体中文',
    en: 'English',
    switch: 'Language'
  },
  menu: {
    dashboard: 'Overview',
    hosts: 'Hosts',
    deploy: 'Deploy',
    tasks: 'Tasks',
    clusters: 'Clusters',
    diagnostics: 'Diagnostics',
    plugins: 'Plugins',
    settings: 'Settings'
  },
  route: {
    dashboard: 'Overview',
    hosts: 'Hosts',
    deploy: 'Deploy',
    tasks: 'Tasks',
    clusters: 'Clusters',
    diagnostics: 'Diagnostics',
    plugins: 'Plugins',
    settings: 'Settings'
  },
  footer: {
    copyright: 'SeaTunnel Console | Apache SeaTunnel'
  },
  tips: {
    profileWip: 'Profile page is under development',
    loggedOut: 'Signed out'
  },
  hosts: {
    register: 'Register Host',
    edit: 'Edit Host',
    searchPlaceholder: 'Search host IP or name',
    agentNotInstalled: 'Agent not installed',
    installAgent: 'Install Agent',
    uninstallAgent: 'Uninstall Agent',
    testConnection: 'Test Connection',
    tabs: {
      list: 'Host List',
      monitor: 'Monitor Panel',
      alerts: 'Alert Rules'
    },
    actions: {
      scaleCluster: 'Scale Cluster',
      restart: 'Restart'
    },
    expandSections: {
      basicInfo: 'Basic Information',
      networkInfo: 'Network Interfaces',
      diskInfo: 'Disk Information'
    },
    columns: {
      name: 'Host Name',
      ip: 'IP Address',
      ipAddress: 'IP Address',
      port: 'SSH Port',
      user: 'SSH User',
      agentStatus: 'Agent Status',
      hostStatus: 'Host Status',
      cpu: 'CPU Usage',
      cpuUsage: 'CPU Usage',
      memory: 'Memory Usage',
      memoryUsage: 'Memory Usage',
      networkIO: 'Network IO',
      avgLoad: 'Average Load',
      lastHeartbeat: 'Last Heartbeat'
    },
    details: {
      cpuModel: 'CPU Model',
      totalMemory: 'Total Memory',
      registerTime: 'Register Time',
      os: 'Operating System',
      kernel: 'Kernel Version',
      hostname: 'Hostname'
    },
    disk: {
      name: 'Disk Name',
      used: 'Used',
      total: 'Total',
      usage: 'Usage',
      ioUtil: 'IO Util',
      mountPoint: 'Mount Point'
    },
    form: {
      name: 'Host Name',
      ip: 'IP Address',
      port: 'SSH Port',
      user: 'SSH User',
      authType: 'Auth Type',
      password: 'Password',
      key: 'Key',
      sshPassword: 'SSH Password',
      keyPath: 'Private Key Path',
      desc: 'Description',
      descPlaceholder: 'Host description',
      exampleName: 'e.g., master-01',
      exampleIp: 'e.g., 192.168.1.100',
      exampleUser: 'e.g., root',
      exampleKeyPath: 'e.g., ~/.ssh/id_rsa'
    },
    valid: {
      name: 'Please enter host name',
      ipRequired: 'Please enter IP address',
      ipFormat: 'Please enter a valid IP address',
      port: 'Please enter SSH port',
      user: 'Please enter SSH user',
      password: 'Please enter SSH password',
      keyPath: 'Please enter private key path'
    },
    msg: {
      saved: 'Host saved successfully',
      confirmInstallAgent: 'Install Agent on "{name}" ({ip})?',
      installSubmitted: 'Agent install task submitted...',
      confirmUninstallAgent: 'Uninstall Agent from "{name}" ({ip})?',
      uninstallSuccess: 'Agent uninstalled successfully',
      testing: 'Testing connection to {name}...',
      testSuccess: 'Connection successful',
      confirmDelete: 'Delete host "{name}"?'
    }
  },
  tasks: {
    searchPlaceholder: 'Search task name',
    type: { streaming: 'Streaming', batch: 'Batch' },
    columns: {
      name: 'Task Name',
      type: 'Type',
      status: 'Status',
      lastRun: 'Last Run',
      duration: 'Duration',
      version: 'Version',
      creator: 'Creator'
    },
    msg: {
      createWip: 'Task creation is under development',
      runTriggered: 'Run triggered: {name}',
      edit: 'Edit task: {name}',
      view: 'View task: {name}',
      confirmDelete: 'Delete task "{name}"?'
    }
  }
}
