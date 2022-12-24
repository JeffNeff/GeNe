package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// geth "github.com/ethereum/go-ethereum"

	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// GethFileLocation is the location of the geth binary
	GethFileLocation string `env:"GETH" default:"../geth"`
}

type Node struct {
	GethFileLocation string
}

type UserInputForNodeConfig struct {
	// TOMLConfig points to a TOML configuration file
	TOMLConfig string

	// RPCHTTPPort HTTP-RPC server listening port (default: 8545)
	RPCHTTPPort string
	// RPCHTTPSelectedAPIMethods HTTP-RPC API modules
	RPCHTTPSelectedAPIMethods []string

	// WSRPCInterface WS-RPC server listening interface (default: "localhost")
	WSRPCInterface string
	// WSRPCHTTPPort WS-RPC server listening port (default: 8546)
	WSRPCHTTPPort string
	// WSRPCOrigins WS-RPC allowed origins list (default: "[]")
	WSRPCOrigins string
	// WSRPCAPIs API's offered over the WS-RPC interface
	WSRPCAPIs []string

	// GraphQLEnabled Enable GraphQL on the HTTP-RPC server. Note that GraphQL can only be started if an HTTP server is started as well.
	GraphQLEnabled bool
	// GraphQLCors Comma separated list of domains from which to accept cross origin requests (browser enforced)
	GraphQLCors string
	// GraphQLVirtualHosts Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: "localhost")
	GraphQLVirtualHosts string

	// AdminAddr Listening address for authenticated APIs (default: "localhost")
	AdminAddr string
	// AdminPort Listening port for authenticated APIs (default: 8551)
	AdminPort string

	// PreloadJS Comma separated list of JavaScript files to preload into the console
	PreloadJS string
	// ExecJS  Execute JavaScript statement
	ExecJS string

	// DBEndpoint URL for remote database
	DBEndpoint string
	// TxLookupLimit Number of recent blocks to maintain transactions index for (default = about one year, 0 = entire chain) (default: 2350000)
	TxLookupLimit string
	// SyncMode Blockchain sync mode ("snap", "full" or "light") (default: snap)
	SyncMode string
	// NetworkID Explicitly set network id (integer)(For testnets: use --sepolia, --goerli instead) (default: 1)
	NetworkID string
	// P2PPort Network listening port (default: 30303)
	P2PPort string
	// DataDir Data directory for the databases and keystore (default: "~/.ethereum")
	DataDir string

	// UserAddress Public address for block mining rewards (default = first account) (default: "0")
	UserAddress string
	// MinerThreads Number of CPU threads to use for mining (default: 0 = use all available cores)
	MinerThreads string
	// NotifyURLs Comma separated HTTP URL list to notify of new work packages
	NotifyURLs string
	// MinerMinimumGasPrice Minimum gas price for mining a transaction (default: 1000000000)
	MinerMinimumGasPrice string
	// MinerGasTarget Target gas ceiling for mined blocks (default: 30000000)
	MinerGasTarget string
	// MinerExtraData Block extra data set by the miner (default = client version)
	MinerExtraData string
	// MinerRecommit Time interval to recreate the block being mined (default: 3s)
	MinerRecommit string
	// MinerDisableRemoteSealing Disable remote sealing verification
	MinerDisableRemoteSealing bool

	// DeveloperMode Flag to enable Ephemeral proof-of-authority network with a pre-funded developer account, mining enabled
	DeveloperMode bool
	// DeveloperPeriod Block period to use in developer mode (0 = mine only if transaction pending) (default: 0)
	DeveloperPeriod string
	// DeveloperGasLimit Initial block gas limit (default: 11500000)
	DeveloperGasLimit string
}

func main() {
	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	myApp := app.New()
	n := &Node{
		GethFileLocation: cfg.GethFileLocation,
	}

	myWindow := myApp.NewWindow("GeNe")
	myWindow.Resize(fyne.NewSize(1080, 1080))

	tomlConfigBinding := binding.NewString()
	tomlConfigInput := widget.NewForm(
		widget.NewFormItem("TOML config file location", widget.NewEntryWithData(tomlConfigBinding)),
	)

	rpcHTTPPortBinding := binding.NewString()
	rpcPortInput := widget.NewForm(
		widget.NewFormItem("HTTP-RPC server listening port (default: 8545)", widget.NewEntryWithData(rpcHTTPPortBinding)),
	)

	httpAPIMethods := []string{"eth", "net", "web3", "txpool", "debug", "admin", "miner", "shh", "clique", "les"}
	RPCHTTPSelectedAPIMethods := binding.NewStringList()
	// select from a list of methods "eth,net,web3,txpool,debug,admin,miner,shh,clique,les"
	httpAPIMethodsInput := widget.NewForm(widget.NewFormItem("API's offered over the HTTP-RPC interface", widget.NewCheckGroup(httpAPIMethods, func(s []string) {
		fmt.Println("Selected:", s)
		RPCHTTPSelectedAPIMethods.Set(s)
	})))

	httpAPIMethodsInput.Hide()

	wsRPCInterfaceBinding := binding.NewString()
	wsRPCInterfaceInput := widget.NewForm(
		widget.NewFormItem("WS-RPC server listening interface (default: localhost)", widget.NewEntryWithData(wsRPCInterfaceBinding)),
	)

	wsRPCHTTPPortBinding := binding.NewString()
	wsRPCHTTPPortInput := widget.NewForm(
		widget.NewFormItem("WS-RPC server listening port (default: 8546)", widget.NewEntryWithData(wsRPCHTTPPortBinding)),
	)

	wsRPCOriginsBinding := binding.NewString()
	wsRPCOriginsInput := widget.NewForm(
		widget.NewFormItem("WS-RPC allowed origins list (default: [])", widget.NewEntryWithData(wsRPCOriginsBinding)),
	)

	APIsToogleBinding := binding.NewBool()
	APIsToggleInput := widget.NewForm(
		widget.NewFormItem("", widget.NewCheckWithData("Show HTTP & WS API Methods ", APIsToogleBinding)),
	)

	WSRPCAPIsBinding := binding.NewStringList()
	WSRPCAPIsInput := widget.NewForm(widget.NewFormItem("API's offered over the WS-RPC interface", widget.NewCheckGroup(httpAPIMethods, func(s []string) {
		fmt.Println("Selected:", s)
		RPCHTTPSelectedAPIMethods.Set(s)
	})))

	WSRPCAPIsInput.Hide()

	graphQLEnabledBinding := binding.NewBool()
	graphQLEnabledInput := widget.NewForm(
		widget.NewFormItem("Enable GraphQL on the HTTP-RPC server. Note that GraphQL can only be started if an HTTP server is started as well.", widget.NewCheckWithData("GraphQL Enabled", graphQLEnabledBinding)),
	)

	graphQLCorsBinding := binding.NewString()
	graphQLCorsInput := widget.NewForm(
		widget.NewFormItem("Comma separated list of domains from which to accept cross origin requests (browser enforced)", widget.NewEntryWithData(graphQLCorsBinding)),
	)

	graphQLVirtualHostsBinding := binding.NewString()
	graphQLVirtualHostsInput := widget.NewForm(
		widget.NewFormItem("Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard. (default: localhost)", widget.NewEntryWithData(graphQLVirtualHostsBinding)),
	)

	adminAddrBinding := binding.NewString()
	adminAddrInput := widget.NewForm(
		widget.NewFormItem("Listening address for authenticated APIs (default: localhost)", widget.NewEntryWithData(adminAddrBinding)),
	)

	adminPortBinding := binding.NewString()
	adminPortInput := widget.NewForm(
		widget.NewFormItem("Listening port for authenticated APIs (default: 8551)", widget.NewEntryWithData(adminPortBinding)),
	)

	preloadJSBinding := binding.NewString()
	preloadJSInput := widget.NewForm(
		widget.NewFormItem("Comma separated list of JavaScript files to preload into the console", widget.NewEntryWithData(preloadJSBinding)),
	)

	execJSBinding := binding.NewString()
	execJSInput := widget.NewForm(
		widget.NewFormItem("JavaScript file to execute at the console startup (implies --preload)", widget.NewEntryWithData(execJSBinding)),
	)

	dbEndpointBinding := binding.NewString()
	dbEndpointInput := widget.NewForm(
		widget.NewFormItem("URL for remote database", widget.NewEntryWithData(dbEndpointBinding)),
	)

	txLookupLimitBinding := binding.NewString()
	txLookupLimitInput := widget.NewForm(
		widget.NewFormItem("Number of recent transactions to maintain in the local transaction history (default: 128)", widget.NewEntryWithData(txLookupLimitBinding)),
	)

	syncModeBinding := binding.NewString()
	syncModeInput := widget.NewForm(
		widget.NewFormItem("Blockchain sync mode ('fast', 'full', or 'light')", widget.NewEntryWithData(syncModeBinding)),
	)

	networkIDBinding := binding.NewString()
	networkIDInput := widget.NewForm(
		widget.NewFormItem("Network identifier (Chain ID)", widget.NewEntryWithData(networkIDBinding)),
	)

	p2pPortBinding := binding.NewString()
	p2pPortInput := widget.NewForm(
		widget.NewFormItem("Network listening port (default: 30303)", widget.NewEntryWithData(p2pPortBinding)),
	)

	dataDirBinding := binding.NewString()
	dataDirInput := widget.NewForm(
		widget.NewFormItem("Data directory for the databases and keystore", widget.NewEntryWithData(dataDirBinding)),
	)

	userAddressBinding := binding.NewString()
	userAddressInput := widget.NewForm(
		widget.NewFormItem("Public address of the signing key", widget.NewEntryWithData(userAddressBinding)),
	)

	minerThreadsBinding := binding.NewString()
	minerThreadsInput := widget.NewForm(
		widget.NewFormItem("Number of CPU threads to use for mining (default: 0)", widget.NewEntryWithData(minerThreadsBinding)),
	)

	notifyURLsBinding := binding.NewString()
	notifyURLsInput := widget.NewForm(
		widget.NewFormItem("Comma separated list of URLs to notify of new work packages (only useful in mining mode)", widget.NewEntryWithData(notifyURLsBinding)),
	)

	minerMinimumGasPriceBinding := binding.NewString()
	minerMinimumGasPriceInput := widget.NewForm(
		widget.NewFormItem("Minimum accepted gas price to allow mining a transaction (default: 18000000000)", widget.NewEntryWithData(minerMinimumGasPriceBinding)),
	)

	minerGasTargetBinding := binding.NewString()
	minerGasTargetInput := widget.NewForm(
		widget.NewFormItem("Target gas floor for mined blocks ", widget.NewEntryWithData(minerGasTargetBinding)),
	)

	minerExtraDataBinding := binding.NewString()
	minerExtraDataInput := widget.NewForm(
		widget.NewFormItem("Block extra data set by the miner", widget.NewEntryWithData(minerExtraDataBinding)),
	)

	minerRecommitBinding := binding.NewString()
	minerRecommitInput := widget.NewForm(
		widget.NewFormItem("Time interval to recreate the block mining work", widget.NewEntryWithData(minerRecommitBinding)),
	)

	minerNoverifyBinding := binding.NewBool()
	minerNoverifyInput := widget.NewForm(
		widget.NewFormItem("Disables remote agent verification", widget.NewCheckWithData("Disable", minerNoverifyBinding)),
	)

	devModeBinding := binding.NewBool()
	devModeInput := widget.NewForm(
		widget.NewFormItem("Developer mode and unsafe RPC settings", widget.NewCheckWithData("Enable", devModeBinding)),
	)

	devPeriodBinding := binding.NewString()
	devPeriodInput := widget.NewForm(
		widget.NewFormItem("Block period to use in developer mode (0 = mine only if transaction pending)", widget.NewEntryWithData(devPeriodBinding)),
	)

	devGasLimitBinding := binding.NewString()
	devGasLimitInput := widget.NewForm(
		widget.NewFormItem("Target gas limit to enforce in developer mode", widget.NewEntryWithData(devGasLimitBinding)),
	)

	// create a binding.DataListener to listen for changes to the graphQLEnabledInput
	graphQLDL := binding.NewDataListener(func() {
		r, err := graphQLEnabledBinding.Get()
		if err != nil {
			fmt.Println(err)
		}

		if r {
			graphQLCorsInput.Show()
			graphQLVirtualHostsInput.Show()
		} else {
			graphQLCorsInput.Hide()
			graphQLVirtualHostsInput.Hide()
		}
	})

	// create a binding.DataListener to listen for changes to the APIsToogleBinding
	dl := binding.NewDataListener(func() {
		r, err := APIsToogleBinding.Get()
		if err != nil {
			fmt.Println(err)
		}
		if r {
			WSRPCAPIsInput.Show()
			httpAPIMethodsInput.Show()
			wsRPCHTTPPortInput.Hide()
			wsRPCOriginsInput.Hide()
			graphQLEnabledInput.Hide()
			graphQLCorsInput.Hide()
			graphQLVirtualHostsInput.Hide()
			adminPortInput.Hide()
			adminAddrInput.Hide()
			dbEndpointInput.Hide()
			txLookupLimitInput.Hide()
			syncModeInput.Hide()
			networkIDInput.Hide()
			p2pPortInput.Hide()
			dataDirInput.Hide()

		} else {
			WSRPCAPIsInput.Hide()
			httpAPIMethodsInput.Hide()

			wsRPCHTTPPortInput.Show()
			wsRPCOriginsInput.Show()
			graphQLEnabledInput.Show()
			graphQLCorsInput.Show()
			graphQLVirtualHostsInput.Show()
			adminPortInput.Show()
			adminAddrInput.Show()
			dbEndpointInput.Show()
			txLookupLimitInput.Show()
			syncModeInput.Show()
			networkIDInput.Show()
			p2pPortInput.Show()
			dataDirInput.Show()

		}
	})
	// watch for changes to the APIsToogleBinding and enable/disable the WSAPIsInput
	APIsToogleBinding.AddListener(dl)
	graphQLEnabledBinding.AddListener(graphQLDL)

	startGethButton := widget.NewButton("Start Geth", func() {
		myWindow.SetTitle("Start Geth")
		tomlConfig, err := tomlConfigBinding.Get()
		if err != nil {
			fmt.Println("Error getting TOML config file location")
		}

		httpRPCPort, err := rpcHTTPPortBinding.Get()
		if err != nil {
			fmt.Println("Error getting HTTP-RPC server listening port")
		}

		s, err := RPCHTTPSelectedAPIMethods.Get()
		if err != nil {
			fmt.Println("Error getting API's offered over the HTTP-RPC interface")
		}

		wsRPCInterface, err := wsRPCInterfaceBinding.Get()
		if err != nil {
			fmt.Println("Error getting WS-RPC server listening interface")
		}

		wsRPCHTTPPort, err := wsRPCHTTPPortBinding.Get()
		if err != nil {
			fmt.Println("Error getting WS-RPC server listening port")
		}

		// create a toggle to hide/show the WS-RPC API's
		// toggle, err := APIsToogleBinding.Get()
		// if err != nil {
		// 	fmt.Println("Error getting WS-RPC API's toggle")
		// }

		wsRPCAPIs, err := WSRPCAPIsBinding.Get()
		if err != nil {
			fmt.Println("Error getting API's offered over the WS-RPC interface")
		}

		graphQLEnabled, err := graphQLEnabledBinding.Get()
		if err != nil {
			fmt.Println("Error getting GraphQL Enabled")
		}

		graphQLCors, err := graphQLCorsBinding.Get()
		if err != nil {
			fmt.Println("Error getting GraphQL CORS")
		}

		graphQLVirtualHosts, err := graphQLVirtualHostsBinding.Get()
		if err != nil {
			fmt.Println("Error getting GraphQL Virtual Hosts")
		}

		adminPort, err := adminPortBinding.Get()
		if err != nil {
			fmt.Println("Error getting Listening port for authenticated APIs")
		}

		adminAddr, err := adminAddrBinding.Get()
		if err != nil {
			fmt.Println("Error getting Listening address for authenticated APIs")
		}

		preloadJS, err := preloadJSBinding.Get()
		if err != nil {
			fmt.Println("Error getting Comma separated list of JavaScript files to preload into the console")
		}

		exexJS, err := execJSBinding.Get()
		if err != nil {
			fmt.Println("Error getting Comma separated list of JavaScript files to execute within the console")
		}

		dbEndpoint, err := dbEndpointBinding.Get()
		if err != nil {
			fmt.Println("Error getting Database to use for storing the blockchain")
		}

		txLookupLimit, err := txLookupLimitBinding.Get()
		if err != nil {
			fmt.Println("Error getting Number of recent transactions to maintain in the local transaction history (default: 128)")
		}

		syncMode, err := syncModeBinding.Get()
		if err != nil {
			fmt.Println("Error getting Blockchain sync mode ('fast', 'full', or 'light')")
		}

		networkID, err := networkIDBinding.Get()
		if err != nil {
			fmt.Println("Error getting Network ID")
		}

		p2pPort, err := p2pPortBinding.Get()
		if err != nil {
			fmt.Println("Error getting P2p listening port")
		}

		dataDir, err := dataDirBinding.Get()
		if err != nil {
			fmt.Println("Error getting Data directory for the databases and keystore")
		}

		userAddress, err := userAddressBinding.Get()
		if err != nil {
			fmt.Println("Error getting Ethereum address for the signing and the mining")
		}

		minerThreads, err := minerThreadsBinding.Get()
		if err != nil {
			fmt.Println("Error getting Number of CPU threads to use for mining")
		}

		notifyURLs, err := notifyURLsBinding.Get()
		if err != nil {
			fmt.Println("Error getting Comma separated list of URLs to notify of new work packages")
		}

		minerMinimumGasPrice, err := minerMinimumGasPriceBinding.Get()
		if err != nil {
			fmt.Println("Error getting Minimum accepted gas price to allow mining a transaction ")
		}

		minerGasTarget, err := minerGasTargetBinding.Get()
		if err != nil {
			fmt.Println("Error getting Target gas floor for mined blocks ")
		}

		minerExtraData, err := minerExtraDataBinding.Get()
		if err != nil {
			fmt.Println("Error getting Specify a custom extra-data for block headers")
		}

		minerRecommitInterval, err := minerRecommitBinding.Get()
		if err != nil {
			fmt.Println("Error getting Time interval to recreate the block mining template")
		}

		minerNoverify, err := minerNoverifyBinding.Get()
		if err != nil {
			fmt.Println("Error getting Disable remote mining")
		}

		devMode, err := devModeBinding.Get()
		if err != nil {
			fmt.Println("Error getting Developer mode, disables block verification")
		}

		devPeriod, err := devPeriodBinding.Get()
		if err != nil {
			fmt.Println("Error getting Block period to use in developer mode")
		}

		devGasLimit, err := devGasLimitBinding.Get()
		if err != nil {
			fmt.Println("Error getting Target gas floor for mined blocks ")
		}

		ui := UserInputForNodeConfig{
			TOMLConfig: tomlConfig,

			RPCHTTPPort:               httpRPCPort,
			RPCHTTPSelectedAPIMethods: s,

			WSRPCInterface: wsRPCInterface,
			WSRPCHTTPPort:  wsRPCHTTPPort,
			WSRPCAPIs:      wsRPCAPIs,

			GraphQLEnabled:      graphQLEnabled,
			GraphQLCors:         graphQLCors,
			GraphQLVirtualHosts: graphQLVirtualHosts,

			AdminPort: adminPort,
			AdminAddr: adminAddr,

			PreloadJS: preloadJS,
			ExecJS:    exexJS,

			DBEndpoint:    dbEndpoint,
			TxLookupLimit: txLookupLimit,
			SyncMode:      syncMode,
			NetworkID:     networkID,
			P2PPort:       p2pPort,
			DataDir:       dataDir,

			UserAddress:               userAddress,
			MinerThreads:              minerThreads,
			NotifyURLs:                notifyURLs,
			MinerMinimumGasPrice:      minerMinimumGasPrice,
			MinerGasTarget:            minerGasTarget,
			MinerExtraData:            minerExtraData,
			MinerRecommit:             minerRecommitInterval,
			MinerDisableRemoteSealing: minerNoverify,

			DeveloperMode:     devMode,
			DeveloperPeriod:   devPeriod,
			DeveloperGasLimit: devGasLimit,
		}

		n.startGeth(ui)
	})

	BasicConfigTab := container.NewVBox(
		userAddressInput,
		minerThreadsInput,
		notifyURLsInput,
		tomlConfigInput,
		startGethButton,
	)

	advancedConfigTab := container.NewVBox(
		APIsToggleInput,
		rpcPortInput,
		httpAPIMethodsInput,
		wsRPCInterfaceInput,
		wsRPCHTTPPortInput,
		wsRPCOriginsInput,
		WSRPCAPIsInput,
		graphQLEnabledInput,
		graphQLCorsInput,
		graphQLVirtualHostsInput,
		adminPortInput,
		adminAddrInput,
		dbEndpointInput,
		txLookupLimitInput,
		syncModeInput,
		networkIDInput,
		p2pPortInput,
		dataDirInput,
		startGethButton,
	)

	minerTab := container.NewVBox(
		userAddressInput,
		minerThreadsInput,
		notifyURLsInput,
		minerMinimumGasPriceInput,
		minerGasTargetInput,
		minerExtraDataInput,
		minerRecommitInput,
		minerNoverifyInput,
		startGethButton,
	)

	developerTab := container.NewVBox(
		devModeInput,
		devPeriodInput,
		devGasLimitInput,
		preloadJSInput,
		execJSInput,
		startGethButton,
	)

	tab1Container := container.New(layout.NewAdaptiveGridLayout(1), BasicConfigTab)
	tab2Container := container.New(layout.NewAdaptiveGridLayout(1), advancedConfigTab)
	tab3Container := container.New(layout.NewAdaptiveGridLayout(1), minerTab)
	tab4Container := container.New(layout.NewAdaptiveGridLayout(1), developerTab)

	tabs := container.NewAppTabs(
		container.NewTabItem("Basic Config", tab1Container),
		container.NewTabItem("Miner Config", tab3Container),
		container.NewTabItem("Advanced Config", tab2Container),
		container.NewTabItem("Developer Config", tab4Container),
	)

	final := container.New(layout.NewMaxLayout())
	final.Add(tabs)

	myWindow.SetContent(final)
	myWindow.ShowAndRun()
	n.tidyUp()
}

func (n *Node) tidyUp() {
	fmt.Println("Exited")
	n.stopGeth()
}

func (n *Node) stopGeth() {
	// kill the geth process
	cmd := exec.Command("killall", "geth")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

}

func (n *Node) startGeth(UserInputForNodeConfig UserInputForNodeConfig) {
	fmt.Println("configuring Geth with user parameters...")
	fmt.Printf("User Address: %s Data Directory: %s P2P Port: %s RPC Port: %s", UserInputForNodeConfig.UserAddress, UserInputForNodeConfig.DataDir, UserInputForNodeConfig.P2PPort, UserInputForNodeConfig.RPCHTTPPort)

	// create a base command
	cmd := exec.Command(n.GethFileLocation)

	// check each paramenter in the UserInputForNodeConfig struct for a value, if it has a value then add it to the command
	if UserInputForNodeConfig.TOMLConfig != "" {
		cmd.Args = append(cmd.Args, "--config", UserInputForNodeConfig.TOMLConfig)
	}

	if UserInputForNodeConfig.RPCHTTPPort != "" {
		cmd.Args = append(cmd.Args, "--http.port", UserInputForNodeConfig.RPCHTTPPort)
	}

	if len(UserInputForNodeConfig.RPCHTTPSelectedAPIMethods) > 0 {
		cmd.Args = append(cmd.Args, "--http.api", strings.Join(UserInputForNodeConfig.RPCHTTPSelectedAPIMethods, ","))
	}

	if UserInputForNodeConfig.WSRPCInterface != "" {
		cmd.Args = append(cmd.Args, "--ws.addr", UserInputForNodeConfig.WSRPCInterface)
	}

	if UserInputForNodeConfig.WSRPCHTTPPort != "" {
		cmd.Args = append(cmd.Args, "--ws.port", UserInputForNodeConfig.WSRPCHTTPPort)
	}

	if UserInputForNodeConfig.WSRPCOrigins != "" {
		cmd.Args = append(cmd.Args, "--ws.origins", UserInputForNodeConfig.WSRPCOrigins)
	}

	if len(UserInputForNodeConfig.WSRPCAPIs) > 0 {
		cmd.Args = append(cmd.Args, "--ws.api", strings.Join(UserInputForNodeConfig.WSRPCAPIs, ","))
	}

	if UserInputForNodeConfig.GraphQLEnabled {
		cmd.Args = append(cmd.Args, "--graphql")
	}

	if UserInputForNodeConfig.GraphQLCors != "" {
		cmd.Args = append(cmd.Args, "--graphql.corsdomain", UserInputForNodeConfig.GraphQLCors)
	}

	if UserInputForNodeConfig.GraphQLVirtualHosts != "" {
		cmd.Args = append(cmd.Args, "--graphql.vhosts", UserInputForNodeConfig.GraphQLVirtualHosts)
	}

	if UserInputForNodeConfig.AdminAddr != "" {
		cmd.Args = append(cmd.Args, "--authrpc.addr", UserInputForNodeConfig.AdminAddr)
	}

	if UserInputForNodeConfig.AdminPort != "" {
		cmd.Args = append(cmd.Args, "--authrpc.port", UserInputForNodeConfig.AdminPort)
	}

	if UserInputForNodeConfig.PreloadJS != "" {
		cmd.Args = append(cmd.Args, "--preload", UserInputForNodeConfig.PreloadJS)
	}

	if UserInputForNodeConfig.ExecJS != "" {
		cmd.Args = append(cmd.Args, "--exec", UserInputForNodeConfig.ExecJS)
	}

	if UserInputForNodeConfig.DBEndpoint != "" {
		cmd.Args = append(cmd.Args, "--db.endpoint", UserInputForNodeConfig.DBEndpoint)
	}

	if UserInputForNodeConfig.TxLookupLimit != "" {
		cmd.Args = append(cmd.Args, "--txlookuplimit", UserInputForNodeConfig.TxLookupLimit)
	}

	if UserInputForNodeConfig.SyncMode != "" {
		cmd.Args = append(cmd.Args, "--syncmode", UserInputForNodeConfig.SyncMode)
	}

	if UserInputForNodeConfig.NetworkID != "" {
		cmd.Args = append(cmd.Args, "--networkid", UserInputForNodeConfig.NetworkID)
	}

	if UserInputForNodeConfig.P2PPort != "" {
		cmd.Args = append(cmd.Args, "--port", UserInputForNodeConfig.P2PPort)
	}

	if UserInputForNodeConfig.DataDir != "" {
		cmd.Args = append(cmd.Args, "--datadir", UserInputForNodeConfig.DataDir)
	}

	if UserInputForNodeConfig.UserAddress != "" {
		cmd.Args = append(cmd.Args, "--miner.etherbase", UserInputForNodeConfig.UserAddress)
	}

	if UserInputForNodeConfig.NotifyURLs != "" {
		cmd.Args = append(cmd.Args, "--miner.notify", UserInputForNodeConfig.NotifyURLs)
	}

	if UserInputForNodeConfig.MinerMinimumGasPrice != "" {
		cmd.Args = append(cmd.Args, "--miner.gasprice", UserInputForNodeConfig.MinerMinimumGasPrice)
	}

	if UserInputForNodeConfig.MinerGasTarget != "" {
		cmd.Args = append(cmd.Args, "--miner.gastarget", UserInputForNodeConfig.MinerGasTarget)
	}

	if UserInputForNodeConfig.MinerExtraData != "" {
		cmd.Args = append(cmd.Args, "--miner.extradata", UserInputForNodeConfig.MinerExtraData)
	}

	if UserInputForNodeConfig.MinerRecommit != "" {
		cmd.Args = append(cmd.Args, "--miner.recommit", UserInputForNodeConfig.MinerRecommit)
	}

	if UserInputForNodeConfig.MinerDisableRemoteSealing {
		cmd.Args = append(cmd.Args, "--miner.noverify")
	}

	if UserInputForNodeConfig.DeveloperMode {
		cmd.Args = append(cmd.Args, "--dev")
	}

	if UserInputForNodeConfig.DeveloperPeriod != "" {
		cmd.Args = append(cmd.Args, "--dev.period", UserInputForNodeConfig.DeveloperPeriod)
	}

	if UserInputForNodeConfig.DeveloperGasLimit != "" {
		cmd.Args = append(cmd.Args, "--dev.gaslimit", UserInputForNodeConfig.DeveloperGasLimit)
	}

	fmt.Println("Starting Geth")

	// execute the comand line
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

}
