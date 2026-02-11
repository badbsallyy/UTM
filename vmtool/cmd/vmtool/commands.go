package vmtool

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/utmapp/vmtool/pkg/api"
	"github.com/google/uuid"
	"github.com/utmapp/vmtool/pkg/config"
	"github.com/utmapp/vmtool/pkg/vm"
	"gopkg.in/yaml.v3"
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Creating VM: %s\n", name)

		cfg := &config.VMConfig{
			Name: name,
			UUID: uuid.New().String(),
			System: config.SystemConfig{
				Architecture: "x86_64",
				Memory:       2048,
				CPUs:         2,
			},
			Display: config.DisplayConfig{
				Enabled: true,
				VNCAddr: ":0",
			},
		}

		// Save VM to store
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := store.SaveVM(cfg); err != nil {
			fmt.Printf("Error saving VM: %v\n", err)
			return
		}
		fmt.Printf("VM %s created successfully.\n", name)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all virtual machines",
	Run: func(cmd *cobra.Command, args []string) {
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		vms := store.ListVMs()
		fmt.Printf("%-20s %-20s %s\n", "NAME", "ARCH", "STATUS")
		for _, v := range vms {
			fmt.Printf("%-20s %-20s %s\n", v.Name, v.System.Architecture, "stopped") // Manager check needed for status
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start a virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		manager := vm.NewManager(store)
		fmt.Printf("🚀 Starting VM: %s...\n", name)
		if err := manager.StartVM(cmd.Context(), name); err != nil {
			fmt.Printf("❌ Error starting VM: %v\n", err)
			return
		}
		fmt.Printf("✅ VM %s is now running.\n", name)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [name]",
	Short: "Stop a virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		manager := vm.NewManager(store)
		fmt.Printf("🛑 Stopping VM: %s...\n", name)
		if err := manager.StopVM(name); err != nil {
			fmt.Printf("❌ Error stopping VM: %v\n", err)
			return
		}
		fmt.Printf("✅ VM %s stopped.\n", name)
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause [name]",
	Short: "Pause a running virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		manager := vm.NewManager(store)
		fmt.Printf("⏸️  Pausing VM: %s...\n", name)
		if err := manager.PauseVM(name); err != nil {
			fmt.Printf("❌ Error pausing VM: %v\n", err)
			return
		}
		fmt.Printf("✅ VM %s paused.\n", name)
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume [name]",
	Short: "Resume a paused virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		manager := vm.NewManager(store)
		fmt.Printf("▶️  Resuming VM: %s...\n", name)
		if err := manager.ResumeVM(name); err != nil {
			fmt.Printf("❌ Error resuming VM: %v\n", err)
			return
		}
		fmt.Printf("✅ VM %s resumed.\n", name)
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the VM management daemon and web server",
	Run: func(cmd *cobra.Command, args []string) {
		appCfg, err := config.LoadAppConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		store, err := vm.NewStore(appCfg.Paths.VMs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		manager := vm.NewManager(store)
		server := api.NewServer(manager, appCfg)

		addr := fmt.Sprintf("%s:%d", appCfg.Server.Host, appCfg.Server.Port)
		fmt.Printf("Starting VMTool server on %s...\n", addr)
		if err := server.Run(addr); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("🗑️ Deleting VM: %s...\n", name)
		if err := store.DeleteVM(name); err != nil {
			fmt.Printf("❌ Error deleting VM: %v\n", err)
			return
		}
		fmt.Printf("✅ VM %s deleted.\n", name)
	},
}

var infoCmd = &cobra.Command{
	Use:   "info [name]",
	Short: "Show detailed information about a virtual machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dataDir := config.GetDefaultDataDir()
		store, err := vm.NewStore(filepath.Join(dataDir, "machines"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		cfg, ok := store.GetVM(name)
		if !ok {
			fmt.Printf("❌ VM %s not found.\n", name)
			return
		}

		fmt.Printf("VM Info: %s\n", cfg.Name)
		fmt.Printf("  UUID:    %s\n", cfg.UUID)
		fmt.Printf("  Arch:    %s\n", cfg.System.Architecture)
		fmt.Printf("  Memory:  %d MB\n", cfg.System.Memory)
		fmt.Printf("  CPUs:    %d\n", cfg.System.CPUs)
		fmt.Printf("  Drives:  %d\n", len(cfg.Drives))
		for _, d := range cfg.Drives {
			fmt.Printf("    - %s (%s)\n", d.ImagePath, d.Interface)
		}
	},
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage VM snapshots",
}

var snapshotCreateCmd = &cobra.Command{
	Use:   "create [vm-name] [snapshot-name]",
	Short: "Create a new snapshot",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vmName, snapName := args[0], args[1]
		dataDir := config.GetDefaultDataDir()
		store, _ := vm.NewStore(filepath.Join(dataDir, "machines"))
		manager := vm.NewManager(store)
		fmt.Printf("📸 Creating snapshot '%s' for VM '%s'...\n", snapName, vmName)
		if err := manager.CreateSnapshot(vmName, snapName); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Snapshot created.")
	},
}

var snapshotRestoreCmd = &cobra.Command{
	Use:   "restore [vm-name] [snapshot-name]",
	Short: "Restore a snapshot",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vmName, snapName := args[0], args[1]
		dataDir := config.GetDefaultDataDir()
		store, _ := vm.NewStore(filepath.Join(dataDir, "machines"))
		manager := vm.NewManager(store)
		fmt.Printf("⏪ Restoring snapshot '%s' for VM '%s'...\n", snapName, vmName)
		if err := manager.RestoreSnapshot(vmName, snapName); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Snapshot restored.")
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:   "delete [vm-name] [snapshot-name]",
	Short: "Delete a snapshot",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vmName, snapName := args[0], args[1]
		dataDir := config.GetDefaultDataDir()
		store, _ := vm.NewStore(filepath.Join(dataDir, "machines"))
		manager := vm.NewManager(store)
		fmt.Printf("🗑️ Deleting snapshot '%s' for VM '%s'...\n", snapName, vmName)
		if err := manager.DeleteSnapshot(vmName, snapName); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Snapshot deleted.")
	},
}

var importCmd = &cobra.Command{
	Use:   "import [path.utm]",
	Short: "Import a UTM bundle into vmtool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bundlePath := args[0]
		fmt.Printf("📥 Importing UTM bundle: %s...\n", bundlePath)
		cfg, warnings, err := vm.ImportUTM(bundlePath)
		if err != nil {
			fmt.Printf("❌ Error importing UTM: %v\n", err)
			return
		}

		fmt.Printf("✅ Imported: %s\n", cfg.Name)
		fmt.Println("📋 Summary:")
		fmt.Printf("   - Name: %s\n", cfg.Name)
		fmt.Printf("   - Arch: %s\n", cfg.System.Architecture)
		fmt.Printf("   - CPU:  %d cores\n", cfg.System.CPUs)
		fmt.Printf("   - RAM:  %d MB\n", cfg.System.Memory)
		fmt.Printf("   - Disks: %d\n", len(cfg.Drives))

		if len(warnings) > 0 {
			fmt.Println("\n⚠️  Warnings:")
			for _, w := range warnings {
				fmt.Printf("   - %s\n", w)
			}
		}

		dataDir := config.GetDefaultDataDir()
		store, _ := vm.NewStore(filepath.Join(dataDir, "machines"))
		store.SaveVM(cfg)

		fmt.Printf("\n💾 Saved to: %s/machines/%s.yaml\n", dataDir, cfg.Name)
		fmt.Printf("▶️  Start with: vmtool start %s\n", cfg.Name)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize vmtool with default directories and configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Initializing vmtool...")

		configDir := config.GetDefaultConfigDir()
		dataDir := config.GetDefaultDataDir()
		cacheDir := config.GetDefaultCacheDir()
		vmDir := filepath.Join(dataDir, "machines")

		dirs := []string{configDir, dataDir, cacheDir, vmDir}
		for _, dir := range dirs {
			fmt.Printf("📁 Creating directory: %s\n", dir)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("❌ Error creating directory %s: %v\n", dir, err)
				return
			}
		}

		appCfg := config.AppConfig{
			Paths: config.PathConfig{
				VMs:   vmDir,
				Cache: cacheDir,
			},
			QEMU: config.QEMUConfig{
				Binary: "auto",
			},
			Server: config.ServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
			},
			Security: config.SecurityConfig{
				APIToken: "", // Generate later
			},
		}

		cfgPath := filepath.Join(configDir, "config.yaml")
		fmt.Printf("🔧 Creating default config: %s\n", cfgPath)
		data, _ := yaml.Marshal(appCfg)
		if err := os.WriteFile(cfgPath, data, 0644); err != nil {
			fmt.Printf("❌ Error writing config: %v\n", err)
			return
		}

		fmt.Println("\n✅ vmtool initialized successfully!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(importCmd)

	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotRestoreCmd)
	snapshotCmd.AddCommand(snapshotDeleteCmd)
	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(initCmd)
}
