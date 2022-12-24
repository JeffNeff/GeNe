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

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// geth "github.com/ethereum/go-ethereum"
)

type Node struct {
	// stdout from the geth process
	Stdout string

	GethOutput *widget.Entry
}

type UserInputForNodeConfig struct {
	UserAddress        string
	DataDir            string
	P2PPort            string
	RPCPort            string
	selectedAPIMethods []string
}

func main() {
	myApp := app.New()
	n := &Node{}
	myWindow := myApp.NewWindow("TabContainer Widget")
	myWindow.Resize(fyne.NewSize(400, 400))

	userAddressBinding := binding.NewString()
	userAddressInput := widget.NewForm(
		widget.NewFormItem("User Address", widget.NewEntryWithData(userAddressBinding)),
	)

	dataDirBinding := binding.NewString()
	dataDirInput := widget.NewForm(
		widget.NewFormItem("Data Directory", widget.NewEntryWithData(dataDirBinding)),
	)

	p2pPortBinding := binding.NewString()
	p2pPortInput := widget.NewForm(
		widget.NewFormItem("P2P Port", widget.NewEntryWithData(p2pPortBinding)),
	)

	rpcPortBinding := binding.NewString()
	rpcPortInput := widget.NewForm(
		widget.NewFormItem("RPC Port", widget.NewEntryWithData(rpcPortBinding)),
	)

	httpAPIMethods := []string{"eth", "net", "web3", "txpool", "debug", "admin", "miner", "shh", "clique", "les"}
	selectedAPIMethods := binding.NewStringList()
	// select from a list of methods "eth,net,web3,txpool,debug,admin,miner,shh,clique,les"
	httpAPIMethodsInput := widget.NewForm(widget.NewFormItem("HTTP API Methods", widget.NewCheckGroup(httpAPIMethods, func(s []string) {
		fmt.Println("Selected:", s)
		selectedAPIMethods.Set(s)
	})))

	startGethButton := widget.NewButton("Start Geth", func() {
		myWindow.SetTitle("Start Geth")
		ua, err := userAddressBinding.Get()
		if err != nil {
			fmt.Println("Error getting user address")
		}

		dd, err := dataDirBinding.Get()
		if err != nil {
			fmt.Println("Error getting data directory")
		}

		pp, err := p2pPortBinding.Get()
		if err != nil {
			fmt.Println("Error getting p2p port")
		}

		rp, err := rpcPortBinding.Get()
		if err != nil {
			fmt.Println("Error getting rpc port")
		}

		s, err := selectedAPIMethods.Get()
		if err != nil {
			fmt.Println("Error getting selected API methods")
		}

		ui := UserInputForNodeConfig{
			UserAddress:        ua,
			DataDir:            dd,
			P2PPort:            pp,
			RPCPort:            rp,
			selectedAPIMethods: s,
		}
		n.startGeth(ui)
	})

	tab1Container := container.NewVBox(userAddressInput, dataDirInput, p2pPortInput, rpcPortInput, httpAPIMethodsInput, startGethButton)

	myWindow.SetContent(container.NewWithoutLayout(tab1Container))
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
	fmt.Printf("User Address: %s Data Directory: %s P2P Port: %s RPC Port: %s", UserInputForNodeConfig.UserAddress, UserInputForNodeConfig.DataDir, UserInputForNodeConfig.P2PPort, UserInputForNodeConfig.RPCPort)

	// create a base command
	cmd := exec.Command("../geth")

	// check each paramenter in the UserInputForNodeConfig struct for a value, if it has a value then add it to the command
	if UserInputForNodeConfig.UserAddress != "" {
		cmd.Args = append(cmd.Args, "--miner.etherbase", UserInputForNodeConfig.UserAddress)
	}

	if UserInputForNodeConfig.DataDir != "" {
		cmd.Args = append(cmd.Args, "--datadir", UserInputForNodeConfig.DataDir)
	}

	if UserInputForNodeConfig.P2PPort != "" {
		cmd.Args = append(cmd.Args, "--port", UserInputForNodeConfig.P2PPort)
	}

	if UserInputForNodeConfig.RPCPort != "" {
		cmd.Args = append(cmd.Args, "--http.port", UserInputForNodeConfig.RPCPort)
	}

	if len(UserInputForNodeConfig.selectedAPIMethods) > 0 {
		cmd.Args = append(cmd.Args, "--http.api", strings.Join(UserInputForNodeConfig.selectedAPIMethods, ","))
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
